package gorag

import (
	"strings"
	"testing"
)

func TestCleanNormalizesWhitespace(t *testing.T) {
	in := "Hello\r\n\r\n\r\n   world  \t here​.\n\n\n\nBye"
	got := Clean(in)
	if strings.Contains(got, "\r") || strings.Contains(got, "​") || strings.Contains(got, "\t") {
		t.Fatalf("control/zero-width/tab survived: %q", got)
	}
	if strings.Contains(got, "   ") || strings.Contains(got, "\n\n\n") {
		t.Fatalf("runs not collapsed: %q", got)
	}
	if !strings.Contains(got, "world here") {
		t.Fatalf("words/spacing mangled: %q", got)
	}
}

func TestParagraphChunker(t *testing.T) {
	text := "Para one is short.\n\nPara two is also short.\n\n" + strings.Repeat("x", 700)
	chunks := ParagraphChunker{MaxChars: 100}.Chunk(text)
	if len(chunks) < 2 {
		t.Fatalf("expected multiple paragraph chunks, got %d", len(chunks))
	}
	if last := chunks[len(chunks)-1]; len(last) < 500 {
		t.Fatalf("oversized paragraph should stand alone, got len %d", len(last))
	}
}

func TestFixedChunkerOverlap(t *testing.T) {
	text := strings.Repeat("abcdefghij", 60) // 600 chars
	chunks := FixedChunker{Size: 200, Overlap: 50}.Chunk(text)
	if len(chunks) < 3 {
		t.Fatalf("expected >=3 windows, got %d", len(chunks))
	}
	for _, c := range chunks {
		if len([]rune(c)) > 200 {
			t.Fatalf("chunk exceeds size: %d", len([]rune(c)))
		}
	}
}

func TestChunkerByNameTunable(t *testing.T) {
	if ChunkerByName("paragraph", 0, 0).Name() != "paragraph" ||
		ChunkerByName("fixed", 0, 0).Name() != "fixed" ||
		ChunkerByName("nonsense", 0, 0).Name() != "whole" {
		t.Fatal("ChunkerByName resolution wrong")
	}
	// size/overlap flow through to the fixed strategy.
	f, ok := ChunkerByName("fixed", 300, 40).(FixedChunker)
	if !ok || f.Size != 300 || f.Overlap != 40 {
		t.Fatalf("fixed params not applied: %+v", f)
	}
	if p, ok := ChunkerByName("paragraph", 250, 0).(ParagraphChunker); !ok || p.MaxChars != 250 {
		t.Fatalf("paragraph size not applied: %+v", p)
	}
}