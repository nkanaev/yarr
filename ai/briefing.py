"""On-demand briefing/digest generation."""

import logging
from collections.abc import AsyncIterator

from .chunker import estimate_tokens
from .providers import LLMProvider
from .store import list_articles, get_article_chunks

log = logging.getLogger(__name__)

SUMMARY_SYSTEM_PROMPT = """You are a news analyst summarizing RSS articles for a reader.

The knowledge base contains articles in multiple languages. Always respond in the same language the user is likely using -- infer it from the article titles.

Your task is to produce a thematic news digest -- organized by THEME, not by article. Group related information together and identify the main developments, trends, and patterns across all articles.

CRITICAL: Cite sources inline using [1], [2], etc. within your sentences -- not as section headers or at the end. Every factual claim must have a citation next to it.

Good example: "El dolar blue cerro la semana en baja [3], mientras el tipo de cambio oficial se mantuvo estable [1][5]."
Bad example: "Article 3 discusses the dollar. Article 1 and 5 cover exchange rates."

Do not make up information -- only summarize what is in the provided articles."""


def prepare_briefing(
    collection,
    topic: str | None = None,
    tag: str | None = None,
    since_ts: float | None = None,
    context_window: int = 4096,
) -> tuple[list[dict], str] | None:
    """Build article list and context string for briefing."""
    articles = list_articles(collection, folder=topic, tag=tag)

    if since_ts is not None:
        articles = [a for a in articles if a.get("published_ts", 0) >= since_ts]

    if not articles:
        return None

    titles_parts = []
    for i, a in enumerate(articles):
        titles_parts.append(f"[{i + 1}] {a['title']} ({a.get('published', 'n/a')})")
    titles_section = "\n".join(titles_parts)

    response_reserve = 600
    framing_overhead = estimate_tokens(SUMMARY_SYSTEM_PROMPT) + estimate_tokens(titles_section) + 300
    budget = context_window - framing_overhead - response_reserve

    detail_parts = []
    tokens_used = 0
    for i, a in enumerate(articles):
        chunks = get_article_chunks(collection, a["url"])
        if not chunks:
            continue
        first_chunk = chunks[0]["text"]
        entry = f"[{i + 1}] **{a['title']}**\n{first_chunk}"
        entry_tokens = estimate_tokens(entry)
        if tokens_used + entry_tokens > budget:
            break
        detail_parts.append(entry)
        tokens_used += entry_tokens

    scope_parts = []
    if topic:
        scope_parts.append(f"topic={topic}")
    if tag:
        scope_parts.append(f"tag={tag}")
    scope_desc = ", ".join(scope_parts) if scope_parts else "all articles"

    detail_section = "\n\n".join(detail_parts)

    user_prompt = f"""Scope: {scope_desc}
Total articles found: {len(articles)} ({len(detail_parts)} with full excerpts below)

Article list:
{titles_section}

Article excerpts:

{detail_section}

Write a thematic summary -- group by topic, NOT one section per article. Use inline citations [1], [2], etc. within sentences next to each specific claim."""

    return articles, user_prompt


async def generate_briefing(
    llm_provider: LLMProvider,
    collection,
    context_window: int = 4096,
    temperature: float = 0.3,
    topic: str | None = None,
    tag: str | None = None,
    since_ts: float | None = None,
) -> AsyncIterator[str]:
    """Stream a briefing digest via LLM provider."""
    result = prepare_briefing(
        collection, topic=topic, tag=tag,
        since_ts=since_ts, context_window=context_window,
    )

    if result is None:
        yield "No articles found for the requested scope and time window."
        return

    articles, user_prompt = result

    messages = [
        {"role": "system", "content": SUMMARY_SYSTEM_PROMPT},
        {"role": "user", "content": user_prompt},
    ]

    stream = await llm_provider.chat(messages, stream=True, temperature=temperature)
    async for token in stream:
        yield token
