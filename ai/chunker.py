"""Sentence-aware text chunking with overlap."""

import re


def estimate_tokens(text: str) -> int:
    """Rough token estimate via word count."""
    return len(text.split())


def chunk_text(text: str, chunk_size: int = 500, overlap: int = 50) -> list[str]:
    """Split text into chunks respecting sentence boundaries.

    Strategy:
    1. Split by paragraphs
    2. Split paragraphs into sentences
    3. Accumulate sentences until chunk_size tokens
    4. Overlap by reusing last N words from previous chunk
    """
    if not text or not text.strip():
        return []

    paragraphs = re.split(r"\n\s*\n", text)
    sentences: list[str] = []
    for para in paragraphs:
        para = para.strip()
        if not para:
            continue
        parts = re.split(r"(?<=[.!?])\s+", para)
        sentences.extend(p.strip() for p in parts if p.strip())

    if not sentences:
        return [text.strip()] if text.strip() else []

    chunks: list[str] = []
    current: list[str] = []
    current_tokens = 0

    for sentence in sentences:
        stokens = estimate_tokens(sentence)
        if current_tokens + stokens > chunk_size and current:
            chunk_text_str = " ".join(current)
            chunks.append(chunk_text_str)
            # Overlap: take last `overlap` words
            words = chunk_text_str.split()
            overlap_words = words[-overlap:] if len(words) > overlap else words
            current = [" ".join(overlap_words)]
            current_tokens = len(overlap_words)
        current.append(sentence)
        current_tokens += stokens

    if current:
        chunks.append(" ".join(current))

    return chunks
