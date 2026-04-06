"""Extract plain text from HTML content."""

import logging

from bs4 import BeautifulSoup

log = logging.getLogger(__name__)

try:
    import trafilatura
    HAS_TRAFILATURA = True
except ImportError:
    HAS_TRAFILATURA = False


def html_to_text(html: str) -> str:
    """Convert HTML to plain text.

    Uses trafilatura as primary extractor, falls back to BeautifulSoup.
    """
    if not html or not html.strip():
        return ""

    # Primary: trafilatura
    if HAS_TRAFILATURA:
        try:
            result = trafilatura.extract(
                html, include_comments=False, include_tables=True
            )
            if result and result.strip():
                return result.strip()
        except Exception as e:
            log.debug("trafilatura failed: %s", e)

    # Fallback: BeautifulSoup
    try:
        soup = BeautifulSoup(html, "html.parser")
        for tag in soup(["script", "style", "nav", "footer"]):
            tag.decompose()
        text = soup.get_text(separator="\n", strip=True)
        lines = [line.strip() for line in text.splitlines() if line.strip()]
        return "\n".join(lines)
    except Exception as e:
        log.warning("BeautifulSoup failed: %s", e)
        return ""
