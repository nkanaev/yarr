"""Provider abstraction for embeddings and LLM calls.

Supports: Ollama (local), Google Gemini, xAI Grok.
"""

import json
import logging
import time
from collections.abc import AsyncIterator
from typing import Protocol, runtime_checkable

import httpx

log = logging.getLogger(__name__)


# --- Protocols ---

@runtime_checkable
class EmbedProvider(Protocol):
    def embed(self, texts: list[str]) -> list[list[float]]: ...
    def dimension(self) -> int: ...


@runtime_checkable
class LLMProvider(Protocol):
    def model_name(self) -> str: ...
    async def chat(
        self, messages: list[dict], stream: bool = True, temperature: float = 0.3
    ) -> str | AsyncIterator[str]: ...


# --- Ollama Providers ---

class OllamaEmbed:
    BATCH_SIZE = 64

    def __init__(self, model: str, ollama_url: str):
        self.model = model
        self.url = ollama_url.rstrip("/")

    def dimension(self) -> int:
        return 768

    def embed(self, texts: list[str]) -> list[list[float]]:
        # Filter empty strings
        cleaned = []
        empty_indices = set()
        for idx, text in enumerate(texts):
            if text and text.strip():
                cleaned.append(text)
            else:
                empty_indices.add(idx)

        if not cleaned:
            return [[0.0] * 768 for _ in texts]

        all_embeddings: list[list[float]] = []
        for i in range(0, len(cleaned), self.BATCH_SIZE):
            batch = cleaned[i : i + self.BATCH_SIZE]
            resp = httpx.post(
                f"{self.url}/api/embed",
                json={"model": self.model, "input": batch},
                timeout=120.0,
            )
            if resp.status_code != 200:
                log.error("Ollama embed %d: %s", resp.status_code, resp.text[:500])
            resp.raise_for_status()
            all_embeddings.extend(resp.json()["embeddings"])

        if empty_indices:
            dim = len(all_embeddings[0]) if all_embeddings else 768
            result = []
            embed_idx = 0
            for idx in range(len(texts)):
                if idx in empty_indices:
                    result.append([0.0] * dim)
                else:
                    result.append(all_embeddings[embed_idx])
                    embed_idx += 1
            return result

        return all_embeddings


class OllamaLLM:
    def __init__(self, model: str, ollama_url: str):
        self.model = model
        self.url = ollama_url.rstrip("/")

    def model_name(self) -> str:
        return f"ollama/{self.model}"

    async def chat(
        self, messages: list[dict], stream: bool = True, temperature: float = 0.3
    ) -> str | AsyncIterator[str]:
        if stream:
            return self._stream(messages, temperature)
        else:
            return await self._generate(messages, temperature)

    async def _generate(self, messages: list[dict], temperature: float) -> str:
        async with httpx.AsyncClient(timeout=120.0) as client:
            resp = await client.post(
                f"{self.url}/api/chat",
                json={
                    "model": self.model,
                    "messages": messages,
                    "stream": False,
                    "options": {"temperature": temperature},
                },
            )
            resp.raise_for_status()
            content = resp.json()["message"]["content"]
            # Strip <think> tags for deepseek-r1 models
            import re
            content = re.sub(r"<think>.*?</think>", "", content, flags=re.DOTALL).strip()
            return content

    async def _stream(self, messages: list[dict], temperature: float) -> AsyncIterator[str]:
        in_think = False
        buffer = ""

        async with httpx.AsyncClient(timeout=300.0) as client:
            async with client.stream(
                "POST",
                f"{self.url}/api/chat",
                json={
                    "model": self.model,
                    "messages": messages,
                    "stream": True,
                    "options": {"temperature": temperature},
                },
            ) as resp:
                resp.raise_for_status()
                async for line in resp.aiter_lines():
                    if not line.strip():
                        continue
                    try:
                        data = json.loads(line)
                    except json.JSONDecodeError:
                        continue

                    token = data.get("message", {}).get("content", "")
                    if not token:
                        if data.get("done"):
                            break
                        continue

                    # Think-tag filtering for deepseek-r1
                    buffer += token
                    while buffer:
                        if in_think:
                            end = buffer.find("</think>")
                            if end != -1:
                                buffer = buffer[end + 8:]
                                in_think = False
                            else:
                                buffer = ""
                                break
                        else:
                            start = buffer.find("<think>")
                            if start != -1:
                                if start > 0:
                                    yield buffer[:start]
                                buffer = buffer[start + 7:]
                                in_think = True
                            else:
                                if "<" in buffer and buffer.rstrip().endswith("<"):
                                    break
                                yield buffer
                                buffer = ""

        if buffer and not in_think:
            yield buffer


# --- Gemini Providers ---

GEMINI_BASE = "https://generativelanguage.googleapis.com/v1beta"


class GeminiEmbed:
    BATCH_SIZE = 50  # Tier 1 allows 10,000 RPM

    def __init__(self, model: str, api_key: str):
        self.model = model
        self.api_key = api_key

    def dimension(self) -> int:
        return 3072  # gemini-embedding-001 outputs 3072-dim vectors

    def embed(self, texts: list[str]) -> list[list[float]]:
        # Filter empty strings
        cleaned = []
        empty_indices = set()
        for idx, text in enumerate(texts):
            if text and text.strip():
                cleaned.append(text)
            else:
                empty_indices.add(idx)

        if not cleaned:
            return [[0.0] * self.dimension() for _ in texts]

        all_embeddings: list[list[float]] = []
        for i in range(0, len(cleaned), self.BATCH_SIZE):
            batch = cleaned[i : i + self.BATCH_SIZE]
            requests_body = [
                {"model": f"models/{self.model}", "content": {"parts": [{"text": t}]}}
                for t in batch
            ]

            # Proactive throttle: ~25 RPS (1500 RPM) for free tier, higher for Tier 1
            if i > 0:
                time.sleep(0.04)

            for attempt in range(5):
                resp = httpx.post(
                    f"{GEMINI_BASE}/models/{self.model}:batchEmbedContents?key={self.api_key}",
                    json={"requests": requests_body},
                    timeout=120.0,
                )
                if resp.status_code == 429:
                    # Respect Retry-After header if present, else exponential backoff
                    retry_after = resp.headers.get("Retry-After")
                    if retry_after:
                        wait = int(retry_after)
                    else:
                        wait = 2 * (2 ** attempt)  # 2, 4, 8, 16, 32
                    log.warning("Gemini embed rate limited, retrying in %ds... (attempt %d/5)", wait, attempt + 1)
                    time.sleep(wait)
                    continue
                if resp.status_code != 200:
                    log.error("Gemini embed %d: %s", resp.status_code, resp.text[:500])
                resp.raise_for_status()
                break

            data = resp.json()
            for emb in data["embeddings"]:
                all_embeddings.append(emb["values"])

        # Re-insert zero vectors at empty positions
        if empty_indices:
            dim = len(all_embeddings[0]) if all_embeddings else 768
            result = []
            embed_idx = 0
            for idx in range(len(texts)):
                if idx in empty_indices:
                    result.append([0.0] * dim)
                else:
                    result.append(all_embeddings[embed_idx])
                    embed_idx += 1
            return result

        return all_embeddings


class GeminiLLM:
    def __init__(self, model: str, api_key: str):
        self.model = model
        self.api_key = api_key

    def model_name(self) -> str:
        return f"gemini/{self.model}"

    async def chat(
        self, messages: list[dict], stream: bool = True, temperature: float = 0.3
    ) -> str | AsyncIterator[str]:
        if stream:
            return self._stream(messages, temperature)
        else:
            return await self._generate(messages, temperature)

    def _convert_messages(self, messages: list[dict]) -> tuple[str | None, list[dict]]:
        """Convert OpenAI-style messages to Gemini format."""
        system = None
        contents = []
        for msg in messages:
            role = msg["role"]
            if role == "system":
                system = msg["content"]
            else:
                gemini_role = "user" if role == "user" else "model"
                contents.append({
                    "role": gemini_role,
                    "parts": [{"text": msg["content"]}],
                })
        return system, contents

    async def _generate(self, messages: list[dict], temperature: float) -> str:
        system, contents = self._convert_messages(messages)
        body: dict = {
            "contents": contents,
            "generationConfig": {"temperature": temperature},
        }
        if system:
            body["systemInstruction"] = {"parts": [{"text": system}]}

        async with httpx.AsyncClient(timeout=120.0) as client:
            for attempt in range(6):
                resp = await client.post(
                    f"{GEMINI_BASE}/models/{self.model}:generateContent?key={self.api_key}",
                    json=body,
                )
                if resp.status_code == 429:
                    retry_after = resp.headers.get("Retry-After")
                    wait = int(retry_after) if retry_after else 10 * (2 ** attempt)  # 10, 20, 40, 80...
                    log.warning("Gemini generate rate limited, retrying in %ds (attempt %d/6)...", wait, attempt + 1)
                    import asyncio
                    await asyncio.sleep(wait)
                    continue
                resp.raise_for_status()
                break

            data = resp.json()
            candidates = data.get("candidates", [])
            if not candidates:
                log.warning("Gemini returned no candidates: %s", str(data)[:300])
                return ""
            return candidates[0]["content"]["parts"][0]["text"]

    async def _stream(self, messages: list[dict], temperature: float) -> AsyncIterator[str]:
        system, contents = self._convert_messages(messages)
        body: dict = {
            "contents": contents,
            "generationConfig": {"temperature": temperature},
        }
        if system:
            body["systemInstruction"] = {"parts": [{"text": system}]}

        async with httpx.AsyncClient(timeout=300.0) as client:
            async with client.stream(
                "POST",
                f"{GEMINI_BASE}/models/{self.model}:streamGenerateContent?alt=sse&key={self.api_key}",
                json=body,
            ) as resp:
                resp.raise_for_status()
                async for line in resp.aiter_lines():
                    if not line.startswith("data: "):
                        continue
                    raw = line[6:]
                    if raw == "[DONE]":
                        break
                    try:
                        data = json.loads(raw)
                        parts = data.get("candidates", [{}])[0].get("content", {}).get("parts", [])
                        for part in parts:
                            text = part.get("text", "")
                            if text:
                                yield text
                    except (json.JSONDecodeError, IndexError, KeyError):
                        continue


# --- Grok (xAI) Provider ---

GROK_BASE = "https://api.x.ai/v1"


class GrokLLM:
    def __init__(self, model: str, api_key: str):
        self.model = model
        self.api_key = api_key

    def model_name(self) -> str:
        return f"grok/{self.model}"

    async def chat(
        self, messages: list[dict], stream: bool = True, temperature: float = 0.3
    ) -> str | AsyncIterator[str]:
        if stream:
            return self._stream(messages, temperature)
        else:
            return await self._generate(messages, temperature)

    async def _generate(self, messages: list[dict], temperature: float) -> str:
        async with httpx.AsyncClient(timeout=120.0) as client:
            resp = await client.post(
                f"{GROK_BASE}/chat/completions",
                headers={"Authorization": f"Bearer {self.api_key}"},
                json={
                    "model": self.model,
                    "messages": messages,
                    "stream": False,
                    "temperature": temperature,
                },
            )
            resp.raise_for_status()
            return resp.json()["choices"][0]["message"]["content"]

    async def _stream(self, messages: list[dict], temperature: float) -> AsyncIterator[str]:
        async with httpx.AsyncClient(timeout=300.0) as client:
            async with client.stream(
                "POST",
                f"{GROK_BASE}/chat/completions",
                headers={"Authorization": f"Bearer {self.api_key}"},
                json={
                    "model": self.model,
                    "messages": messages,
                    "stream": True,
                    "temperature": temperature,
                },
            ) as resp:
                resp.raise_for_status()
                async for line in resp.aiter_lines():
                    if not line.startswith("data: "):
                        continue
                    raw = line[6:]
                    if raw.strip() == "[DONE]":
                        break
                    try:
                        data = json.loads(raw)
                        delta = data["choices"][0].get("delta", {})
                        text = delta.get("content", "")
                        if text:
                            yield text
                    except (json.JSONDecodeError, IndexError, KeyError):
                        continue


# --- Fallback Wrapper ---

class FallbackEmbed:
    """Tries primary embed provider, falls back to secondary on error."""

    def __init__(self, primary: EmbedProvider, fallback: EmbedProvider):
        self.primary = primary
        self.fallback = fallback

    def dimension(self) -> int:
        return self.primary.dimension()

    def embed(self, texts: list[str]) -> list[list[float]]:
        try:
            return self.primary.embed(texts)
        except Exception as e:
            log.warning("Primary embed failed (%s), falling back to Ollama: %s", type(self.primary).__name__, e)
            return self.fallback.embed(texts)


class FallbackLLM:
    """Tries primary LLM provider, falls back to secondary on error."""

    def __init__(self, primary: LLMProvider, fallback: LLMProvider):
        self.primary = primary
        self.fallback = fallback

    def model_name(self) -> str:
        return self.primary.model_name()

    async def chat(
        self, messages: list[dict], stream: bool = True, temperature: float = 0.3
    ) -> str | AsyncIterator[str]:
        if not stream:
            # Non-streaming: simple try/except
            try:
                return await self.primary.chat(messages, stream=False, temperature=temperature)
            except Exception as e:
                log.warning("Primary LLM failed (%s), falling back: %s", self.primary.model_name(), e)
                return await self.fallback.chat(messages, stream=False, temperature=temperature)
        else:
            # Streaming: catch errors during iteration and fall back
            async def fallback_stream() -> AsyncIterator[str]:
                try:
                    result = await self.primary.chat(messages, stream=True, temperature=temperature)
                    async for token in result:
                        yield token
                except Exception as e:
                    log.warning("Primary LLM streaming failed (%s), falling back: %s", self.primary.model_name(), e)
                    try:
                        result = await self.fallback.chat(messages, stream=True, temperature=temperature)
                        async for token in result:
                            yield token
                    except Exception as e2:
                        log.error("Fallback LLM also failed: %s", e2)
                        yield f"[Error: AI service unavailable. Primary: {e}, Fallback: {e2}]"
            return fallback_stream()


# --- Factory Functions ---

def create_embed_provider(config) -> EmbedProvider:
    """Create embedding provider based on config."""
    provider = config.embed_provider

    # Auto-detect
    if provider == "auto":
        provider = "gemini" if config.gemini_api_key else "ollama"

    if provider == "gemini" and config.gemini_api_key:
        primary = GeminiEmbed(config.gemini_embed_model, config.gemini_api_key)
        fallback = OllamaEmbed(config.embed_model, config.ollama_url)
        log.info("Embed provider: Gemini %s (fallback: Ollama %s)", config.gemini_embed_model, config.embed_model)
        return FallbackEmbed(primary, fallback)

    log.info("Embed provider: Ollama %s", config.embed_model)
    return OllamaEmbed(config.embed_model, config.ollama_url)


def create_llm_provider(config, purpose: str = "chat") -> LLMProvider:
    """Create LLM provider based on config and purpose.

    purpose: "chat" (interactive streaming) or "label" (structured, non-streaming)
    """
    provider = config.llm_provider

    # Auto-detect
    if provider == "auto":
        if purpose == "chat" and config.grok_api_key:
            provider = "grok"
        elif config.gemini_api_key:
            provider = "gemini"
        else:
            provider = "ollama"

    ollama_fallback = OllamaLLM(
        config.chat_model if purpose == "chat" else config.label_model,
        config.ollama_url,
    )

    if provider == "grok" and config.grok_api_key:
        primary = GrokLLM(config.grok_chat_model, config.grok_api_key)
        log.info("LLM provider (%s): Grok %s (fallback: Ollama %s)", purpose, config.grok_chat_model, ollama_fallback.model)
        return FallbackLLM(primary, ollama_fallback)

    if provider == "gemini" and config.gemini_api_key:
        primary = GeminiLLM(config.gemini_chat_model, config.gemini_api_key)
        log.info("LLM provider (%s): Gemini %s (fallback: Ollama %s)", purpose, config.gemini_chat_model, ollama_fallback.model)
        return FallbackLLM(primary, ollama_fallback)

    log.info("LLM provider (%s): Ollama %s", purpose, ollama_fallback.model)
    return ollama_fallback
