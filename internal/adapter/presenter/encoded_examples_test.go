package presenter

import (
	"os"
	"path/filepath"
	"testing"

	"toon-showcase/internal/core/entity"
)

func TestWriteShapeEncodedExamples(t *testing.T) {
	dir := t.TempDir()
	err := WriteShapeEncodedExamples(dir, []entity.ShapeEncodedArtifact{
		{
			FormatKey: "json-compact",
			CaseName:  "slice_struct",
			Binary:    false,
			Result:    entity.ShapeResult{EncodeOK: true, DecodeOK: true, RoundTripOK: true, Detail: "ok"},
			Encoded:   []byte("{\"a\":1}"),
		},
		{
			FormatKey: "messagepack",
			CaseName:  "slice_struct",
			Binary:    true,
			Result:    entity.ShapeResult{EncodeOK: true, DecodeOK: true, RoundTripOK: true, Detail: "ok"},
			Encoded:   []byte{0x81, 0xa1, 0x61, 0x01},
		},
		{
			FormatKey: "toml",
			CaseName:  "slice_any_mixed",
			Binary:    false,
			Result:    entity.ShapeResult{EncodeOK: false, DecodeOK: false, RoundTripOK: false, Detail: "unsupported type"},
		},
	})
	if err != nil {
		t.Fatalf("WriteShapeEncodedExamples error: %v", err)
	}

	textPath := filepath.Join(dir, "shape-encoded", "json-compact", "slice_struct.txt")
	if _, err := os.Stat(textPath); err != nil {
		t.Fatalf("missing text output: %v", err)
	}

	binPath := filepath.Join(dir, "shape-encoded", "messagepack", "slice_struct.bin")
	if _, err := os.Stat(binPath); err != nil {
		t.Fatalf("missing binary output: %v", err)
	}

	errPath := filepath.Join(dir, "shape-encoded", "toml", "slice_any_mixed.error.txt")
	content, err := os.ReadFile(errPath)
	if err != nil {
		t.Fatalf("missing error output: %v", err)
	}
	if string(content) != "unsupported type\n" {
		t.Fatalf("error file mismatch: %q", string(content))
	}
}

func TestWriteShapeEncodedExamplesError(t *testing.T) {
	err := WriteShapeEncodedExamples("/dev/null/bad", []entity.ShapeEncodedArtifact{{
		FormatKey: "json-compact",
		CaseName:  "slice_struct",
		Result:    entity.ShapeResult{EncodeOK: true},
		Encoded:   []byte("{}"),
	}})
	if err == nil {
		t.Fatalf("expected write error")
	}
}
