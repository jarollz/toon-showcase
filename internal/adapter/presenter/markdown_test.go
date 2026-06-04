package presenter

import (
	"strings"
	"testing"
	"time"

	"toon-showcase/internal/core/entity"
)

func TestBuildBenchmarkMarkdownReport(t *testing.T) {
	results := []entity.BenchmarkResult{
		{
			Name:             "TOON",
			MarshalAvg:       100 * time.Nanosecond,
			UnmarshalAvg:     200 * time.Nanosecond,
			EncodedBytes:     100,
			EncodedChars:     90,
			RoundTripDetails: "ok",
			Sample:           "toon sample",
		},
		{
			Name:             "JSON compact",
			MarshalAvg:       90 * time.Nanosecond,
			UnmarshalAvg:     300 * time.Nanosecond,
			EncodedBytes:     100,
			EncodedChars:     90,
			RoundTripDetails: "ok",
			Sample:           "json sample",
		},
	}
	charset := entity.CharsetMatrix{
		FormatNames: []string{"TOON", "JSON compact"},
		Tests: []entity.CharsetCase{
			{Name: "ASCII", Value: "x"},
		},
		ResultsByFormat: [][]entity.CharsetResult{
			{{MarshalOK: true, RoundTripOK: true, Detail: "ok"}},
			{{MarshalOK: false, RoundTripOK: false, Detail: "bad | detail"}},
		},
		PassCounts: []int{1, 0},
	}
	shape := entity.ShapeMatrix{
		FormatNames: []string{"TOON", "JSON compact"},
		Tests:       []string{"shape-a"},
		ResultsByFormat: [][]entity.ShapeResult{
			{{EncodeOK: true, DecodeOK: true, RoundTripOK: true, Detail: "ok"}},
			{{EncodeOK: false, DecodeOK: false, RoundTripOK: false, Detail: "shape bad"}},
		},
		PassCounts: []int{1, 0},
	}
	md := BuildFullMarkdownReport(results, results[0], charset, shape, 20, 50, 5, 42, 2)

	if !strings.Contains(md, "| Format | Marshal (ns/op) | Unmarshal (ns/op) | Bytes | Roundtrip |") {
		t.Fatalf("missing benchmark markdown table header: %s", md)
	}
	if !strings.Contains(md, "Emoji legend: `✅` better, `❌` worse, `➖` same.") {
		t.Fatalf("missing emoji legend: %s", md)
	}
	if !strings.Contains(md, "JSON compact") || !strings.Contains(md, "✅") || !strings.Contains(md, "❌") || !strings.Contains(md, "➖") {
		t.Fatalf("missing expected delta marker rows: %s", md)
	}
	if !strings.Contains(md, "## Character Set Support Comparison") {
		t.Fatalf("missing charset section: %s", md)
	}
	if !strings.Contains(md, "## Shape Compatibility Matrix") {
		t.Fatalf("missing shape section: %s", md)
	}
	if !strings.Contains(md, "bad \\| detail") {
		t.Fatalf("missing escaped pipe in markdown cell: %s", md)
	}
}

func TestBuildBenchmarkMarkdownReportWrapperAndEmptyMatrices(t *testing.T) {
	results := []entity.BenchmarkResult{
		{
			Name:             "TOON",
			MarshalAvg:       100 * time.Nanosecond,
			UnmarshalAvg:     200 * time.Nanosecond,
			EncodedBytes:     100,
			EncodedChars:     90,
			RoundTripDetails: "ok",
			Sample:           "line1|line2\nline3",
		},
	}
	md := BuildBenchmarkMarkdownReport(results, results[0], 20, 50, 5, 42, 2)

	if !strings.Contains(md, "No character set report data.") {
		t.Fatalf("missing empty charset fallback: %s", md)
	}
	if !strings.Contains(md, "No shape report data.") {
		t.Fatalf("missing empty shape fallback: %s", md)
	}
	if !strings.Contains(md, "line1\\|line2<br>line3") {
		t.Fatalf("missing escaped sample markdown: %s", md)
	}
}

func TestMarkdownDeltaHelpersAndEscape(t *testing.T) {
	if got := deltaCellDuration(1, 0); got != "n/a" {
		t.Fatalf("deltaCellDuration base zero mismatch: %s", got)
	}
	if got := deltaCellDuration(10, 10); got != "0.0% ➖" {
		t.Fatalf("deltaCellDuration equal mismatch: %s", got)
	}
	if got := deltaCellInt(-1, 10); got != "n/a" {
		t.Fatalf("deltaCellInt invalid mismatch: %s", got)
	}
	if got := deltaCellInt(10, 10); got != "0.0% ➖" {
		t.Fatalf("deltaCellInt equal mismatch: %s", got)
	}

	if got := mdEscapeCell("a|b\r\nc\nd\r"); got != "a\\|b<br>c<br>d<br>" {
		t.Fatalf("mdEscapeCell mismatch: %q", got)
	}
}
