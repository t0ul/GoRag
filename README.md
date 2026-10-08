# GoRag

For data things — the fleet's RAG ingestion toolkit: dependency-free building blocks
that turn raw documents into clean, chunked units ready to index.

- `Clean(text)` — normalize whitespace, strip control/zero-width characters.
- `Chunker` — pluggable, tunable chunking strategies:
  - `WholeChunker` — whole document (baseline).
  - `ParagraphChunker{MaxChars}` — blank-line split, greedily packed.
  - `FixedChunker{Size, Overlap}` — fixed rune window with overlap.
- `ChunkerByName(name, size, overlap)` — resolve a strategy by operator choice.

Stdlib only. The retrieval store, embeddings, and app-specific controls live with the
consumer, not here.
