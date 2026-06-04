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

func sampleBenchmarkResults() ([]entity.BenchmarkResult, entity.BenchmarkResult) {
	results := []entity.BenchmarkResult{
		{
			Name:             "TOON",
			MarshalTotal:     10 * time.Millisecond,
			MarshalAvg:       100 * time.Nanosecond,
			UnmarshalTotal:   20 * time.Millisecond,
			UnmarshalAvg:     200 * time.Nanosecond,
			EncodedBytes:     100,
			EncodedChars:     90,
			RoundTripOK:      true,
			RoundTripDetails: "ok",
		},
		{
			Name:             "MessagePack",
			Binary:           true,
			MarshalTotal:     5 * time.Millisecond,
			MarshalAvg:       50 * time.Nanosecond,
			UnmarshalTotal:   6 * time.Millisecond,
			UnmarshalAvg:     60 * time.Nanosecond,
			EncodedBytes:     80,
			EncodedChars:     -1,
			RoundTripOK:      true,
			RoundTripDetails: "ok",
		},
	}
	return results, results[0]
}

func TestBuildBenchmarkCSV(t *testing.T) {
	results, baseline := sampleBenchmarkResults()

	csvText, err := BuildBenchmarkCSV(results, baseline, 10, 20, 5, 42, 2)
	if err != nil {
		t.Fatalf("BuildBenchmarkCSV error: %v", err)
	}
	if !strings.Contains(csvText, "format,baseline_format") {
		t.Fatalf("missing header")
	}
	if !strings.Contains(csvText, "MessagePack") {
		t.Fatalf("missing row")
	}

	empty, err := BuildBenchmarkCSV(nil, baseline, 1, 1, 0, 1, 2)
	if err != nil {
		t.Fatalf("BuildBenchmarkCSV empty error: %v", err)
	}
	if empty != "" {
		t.Fatalf("expected empty csv for empty results")
	}
}

func TestBuildCharsetAndShapeCSV(t *testing.T) {
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

	c1, err := BuildCharsetCSV(charset)
	if err != nil {
		t.Fatalf("BuildCharsetCSV error: %v", err)
	}
	if !strings.Contains(c1, "ASCII") {
		t.Fatalf("missing charset row")
	}

	c2, err := BuildShapeCSV(shape)
	if err != nil {
		t.Fatalf("BuildShapeCSV error: %v", err)
	}
	if !strings.Contains(c2, "shape") {
		t.Fatalf("missing shape row")
	}
}

func TestEmitCSVs(t *testing.T) {
	out := &bytes.Buffer{}
	if err := EmitCSVs(out, "stdout", "x.csv", "a", "b", "c"); err != nil {
		t.Fatalf("EmitCSVs stdout error: %v", err)
	}
	if !strings.Contains(out.String(), "CSV Benchmark Output") {
		t.Fatalf("missing stdout section")
	}

	dir := t.TempDir()
	path := filepath.Join(dir, "report.csv")
	out.Reset()
	if err := EmitCSVs(out, "file", path, "a", "b", "c"); err != nil {
		t.Fatalf("EmitCSVs file error: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("benchmark file missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "report.charset.csv")); err != nil {
		t.Fatalf("charset file missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "report.shape.csv")); err != nil {
		t.Fatalf("shape file missing: %v", err)
	}

	if err := EmitCSVs(out, "file", "/not/a/real/path/report.csv", "a", "b", "c"); err == nil {
		t.Fatalf("expected write error")
	}
}

func TestCSVHelpers(t *testing.T) {
	if got := deriveCSVPath("abc.csv", ".x.csv"); got != "abc.x.csv" {
		t.Fatalf("derive csv path mismatch: %s", got)
	}
	if got := deriveCSVPath("abc", ".x.csv"); got != "abc.x.csv" {
		t.Fatalf("derive csv path mismatch: %s", got)
	}

	if got := ratioFloat(10, 0); got != "" {
		t.Fatalf("ratioFloat base zero should be empty")
	}
	if got := ratioIntFloat(-1, 10); got != "" {
		t.Fatalf("ratioIntFloat negative should be empty")
	}
}
