"""AI tag generation using LLM."""

import asyncio
import logging
import re

log = logging.getLogger(__name__)


def generate_tags(
    text: str,
    title: str,
    llm_provider=None,
    model: str = "llama3.2:latest",
    ollama_url: str = "http://localhost:11434",
) -> str:
    """Generate topic tags for an article using LLM.

    Returns comma-separated lowercase tags (e.g., "inflation,argentina,economy").
    """
    text_excerpt = text[:1000]

    prompt = f"""Article Title: {title}

Article Excerpt:
{text_excerpt}

Extract 3-5 topic tags that describe the main themes of this article.
Return ONLY comma-separated lowercase tags with no spaces after commas.
Example: "inflation,argentina,economy"

Tags:"""

    try:
        if llm_provider:
            import concurrent.futures
            with concurrent.futures.ThreadPoolExecutor() as pool:
                content = pool.submit(
                    lambda: asyncio.run(llm_provider.chat(
                        [{"role": "user", "content": prompt}],
                        stream=False, temperature=0.3
                    ))
                ).result()
        else:
            import httpx
            resp = httpx.post(
                f"{ollama_url.rstrip('/')}/api/chat",
                json={
                    "model": model,
                    "messages": [{"role": "user", "content": prompt}],
                    "stream": False,
                },
                timeout=60.0,
            )
            resp.raise_for_status()
            content = resp.json()["message"]["content"]
    except Exception as e:
        log.warning("Tag generation failed: %s", e)
        return ""

    # Strip <think> tags
    content = re.sub(r"<think>.*?</think>", "", content, flags=re.DOTALL).strip()

    # Parse tags
    tags = [t.strip().lower() for t in content.split(",") if t.strip()]
    tags = [t for t in tags if t][:5]

    return ",".join(tags)
