package presenter

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"toon-showcase/internal/core/entity"
)

func TestIsMarkdownStdout(t *testing.T) {
	if !IsMarkdownStdout(" stdout ") {
		t.Fatalf("expected stdout sentinel detection")
	}
	if IsMarkdownStdout("./stdout") {
		t.Fatalf("./stdout must not be treated as sentinel")
	}
}

func TestEmitEncodedExamplesOutput(t *testing.T) {
	buf := &bytes.Buffer{}
	if err := EmitEncodedExamplesOutput(buf, "", "", nil); err != nil {
		t.Fatalf("empty encoded examples dir should be ignored: %v", err)
	}

	dir := t.TempDir()
	artifacts := []entity.ShapeEncodedArtifact{{
		FormatKey: "toon",
		CaseName:  "slice_struct",
		Result: entity.ShapeResult{
			EncodeOK: true,
		},
		Encoded: []byte("x"),
	}}
	if err := EmitEncodedExamplesOutput(buf, dir, dir, artifacts); err != nil {
		t.Fatalf("emit encoded examples error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "shape-encoded", "toon", "slice_struct.txt")); err != nil {
		t.Fatalf("expected encoded artifact file: %v", err)
	}

	filePath := filepath.Join(t.TempDir(), "file.txt")
	if err := os.WriteFile(filePath, []byte("x"), 0644); err != nil {
		t.Fatalf("setup file failed: %v", err)
	}
	if err := EmitEncodedExamplesOutput(buf, filePath, filePath, artifacts); err == nil {
		t.Fatalf("expected directory validation error")
	}
}

func TestEmitMarkdownOutput(t *testing.T) {
	buf := &bytes.Buffer{}
	if err := EmitMarkdownOutput(buf, "", "# x"); err != nil {
		t.Fatalf("empty markdown path should be ignored: %v", err)
	}

	mdFile := filepath.Join(t.TempDir(), "a", "report.md")
	if err := EmitMarkdownOutput(buf, mdFile, "# heading"); err != nil {
		t.Fatalf("emit markdown file error: %v", err)
	}
	b, err := os.ReadFile(mdFile)
	if err != nil {
		t.Fatalf("failed read markdown file: %v", err)
	}
	if !strings.Contains(string(b), "heading") {
		t.Fatalf("unexpected markdown file content: %s", string(b))
	}

	buf.Reset()
	if err := EmitMarkdownOutput(buf, "stdout", "# heading"); err != nil {
		t.Fatalf("emit markdown stdout error: %v", err)
	}
	if !strings.Contains(buf.String(), "\x1b[") {
		t.Fatalf("expected ansi rendering for markdown stdout")
	}
}

func TestEmitCSVOutput(t *testing.T) {
	results, baseline := sampleBenchmarkResults()
	charset := entity.CharsetMatrix{
		FormatNames: []string{"A"},
		Tests:       []entity.CharsetCase{{Name: "ASCII", Value: "x"}},
		ResultsByFormat: [][]entity.CharsetResult{{
			{MarshalOK: true, RoundTripOK: true, Detail: "ok"},
		}},
	}
	shape := entity.ShapeMatrix{
		FormatNames: []string{"A"},
		Tests:       []string{"shape"},
		ResultsByFormat: [][]entity.ShapeResult{{
			{EncodeOK: true, DecodeOK: true, RoundTripOK: true, Detail: "ok"},
		}},
	}
	buf := &bytes.Buffer{}
	if err := EmitCSVOutput(buf, "", results, baseline, charset, shape, 1, 1, 0, time.Now().Unix(), 2); err != nil {
		t.Fatalf("empty csv dir should be ignored: %v", err)
	}
	if err := EmitCSVOutput(buf, "stdout", results, baseline, charset, shape, 1, 1, 0, time.Now().Unix(), 2); err != nil {
		t.Fatalf("stdout csv output error: %v", err)
	}

	csvDir := filepath.Join(t.TempDir(), "csv")
	if err := EmitCSVOutput(buf, csvDir, results, baseline, charset, shape, 1, 1, 0, time.Now().Unix(), 2); err != nil {
		t.Fatalf("csv dir output error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(csvDir, "benchmark_report.csv")); err != nil {
		t.Fatalf("missing benchmark report csv: %v", err)
	}
}
