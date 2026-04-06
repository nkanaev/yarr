"""RAG chat engine with hybrid search and streaming LLM responses."""

import logging
from collections.abc import AsyncIterator

from .chunker import estimate_tokens
from .providers import EmbedProvider, LLMProvider
from .search import bm25_search, reciprocal_rank_fusion
from .store import query

log = logging.getLogger(__name__)

SYSTEM_PROMPT = """You are a helpful assistant that answers questions based on a knowledge base of RSS articles.

{topic_summary}

The knowledge base contains articles in multiple languages. Always respond in the same language the user writes in.

Each article chunk in the context includes:
- Article title (in bold)
- Metadata: folder/topic, feed source, and publication date
- The actual text content

Ground your answers in the provided context. Extract and synthesize information from the articles even if they are in a different language than the question -- translate and summarize as needed.
Cite articles using [1], [2], etc. when referencing information.

If you need to think through a problem, you may use internal reasoning, but present your final answer clearly and concisely."""

USER_PROMPT_TEMPLATE = """Context from knowledge base:

{context}

Question: {query}

Please answer the question based on the context provided. Cite sources using [1], [2], etc."""


class ChatEngine:
    def __init__(self, config, collection, bm25_index, bm25_docs,
                 embed_provider: EmbedProvider = None,
                 llm_provider: LLMProvider = None):
        self.config = config
        self.collection = collection
        self.bm25_index = bm25_index
        self.bm25_docs = bm25_docs
        self.embed_provider = embed_provider
        self.llm_provider = llm_provider
        self.topic_summary = self._build_topic_summary()

    def _build_topic_summary(self) -> str:
        from .store import list_topics
        topics = list_topics(self.collection)
        if not topics:
            return "The knowledge base is currently empty."
        top = topics[:10]
        lines = [f"- {t['folder']}: {t['article_count']} articles" for t in top]
        return "Current knowledge base topics:\n" + "\n".join(lines)

    def search(
        self,
        user_query: str,
        topic: str | None = None,
        tag: str | None = None,
        since_ts: float | None = None,
    ) -> list[dict]:
        """Hybrid search: vector + BM25 + RRF."""
        filters = []
        if topic:
            filters.append({"folder": {"$eq": topic}})
        if tag:
            filters.append({"tags": {"$eq": tag}})

        where = None
        if len(filters) == 1:
            where = filters[0]
        elif len(filters) > 1:
            where = {"$and": filters}

        vector_results = query(
            self.collection,
            user_query,
            n_results=self.config.n_results,
            where=where,
            distance_threshold=self.config.distance_threshold,
            since_ts=since_ts,
            embed_fn=self.embed_provider,
        )

        if self.bm25_index is not None:
            bm25_results = bm25_search(
                self.bm25_index,
                self.bm25_docs,
                user_query,
                n_results=self.config.n_results,
                folder=topic,
                tag=tag,
                since_ts=since_ts,
            )
            merged = reciprocal_rank_fusion(vector_results, bm25_results)
            return merged[: self.config.n_results]

        return vector_results

    def _format_context(self, results: list[dict]) -> tuple[str, list[dict]]:
        """Format search results into context string with budget management."""
        system_tokens = estimate_tokens(
            SYSTEM_PROMPT.format(topic_summary=self.topic_summary)
        )
        response_reserve = 500
        budget = self.config.context_window - system_tokens - response_reserve - 200

        parts = []
        used = []
        tokens_used = 0

        for i, r in enumerate(results):
            chunk = (
                f"[{i + 1}] **{r.get('title', 'Untitled')}**\n"
                f"    Folder: {r.get('folder', '')} | "
                f"Feed: {r.get('feed_name', '')} | "
                f"Published: {r.get('published', '')}\n"
                f"    {r.get('text', '')}"
            )
            chunk_tokens = estimate_tokens(chunk)
            if tokens_used + chunk_tokens > budget:
                break
            parts.append(chunk)
            used.append(r)
            tokens_used += chunk_tokens

        return "\n\n".join(parts), used

    async def generate_response(
        self,
        user_query: str,
        results: list[dict],
        history: list[dict] | None = None,
    ) -> AsyncIterator[str]:
        """Stream LLM response with context."""
        context, used_results = self._format_context(results)

        system_msg = SYSTEM_PROMPT.format(topic_summary=self.topic_summary)
        user_msg = USER_PROMPT_TEMPLATE.format(context=context, query=user_query)

        messages = [{"role": "system", "content": system_msg}]

        if history:
            max_msgs = self.config.max_history * 2
            messages.extend(history[-max_msgs:])

        messages.append({"role": "user", "content": user_msg})

        if self.llm_provider:
            result = await self.llm_provider.chat(
                messages, stream=True, temperature=self.config.temperature
            )
            async for token in result:
                yield token
        else:
            yield "LLM provider not configured."

    def rebuild_index(self, bm25_index, bm25_docs):
        """Update BM25 index after new articles are indexed."""
        self.bm25_index = bm25_index
        self.bm25_docs = bm25_docs
        self.topic_summary = self._build_topic_summary()
