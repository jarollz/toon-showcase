package presenter

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"toon-showcase/internal/core/entity"
)

func TestPrintReports(t *testing.T) {
	results := []entity.BenchmarkResult{
		{
			Name:             "TOON",
			MarshalTotal:     2 * time.Millisecond,
			MarshalAvg:       200 * time.Nanosecond,
			UnmarshalTotal:   3 * time.Millisecond,
			UnmarshalAvg:     300 * time.Nanosecond,
			EncodedBytes:     10,
			EncodedChars:     9,
			RoundTripDetails: "ok",
			Sample:           "abc",
		},
	}
	baseline := results[0]

	buf := &bytes.Buffer{}
	PrintBenchmarkReport(buf, results, baseline, 1, 2, 3, 4, 2)
	text := buf.String()
	if !strings.Contains(text, "Multi-Format Marshal/Unmarshal Showcase") {
		t.Fatalf("missing benchmark header")
	}
	if !strings.Contains(text, "TOON") {
		t.Fatalf("missing row")
	}

	charset := entity.CharsetMatrix{
		FormatNames: []string{"A"},
		Tests:       []entity.CharsetCase{{Name: "ASCII", Value: "x"}},
		ResultsByFormat: [][]entity.CharsetResult{{
			{MarshalOK: false, RoundTripOK: false, Detail: "bad"},
		}},
		PassCounts: []int{0},
	}
	buf.Reset()
	PrintCharsetMatrix(buf, charset)
	if !strings.Contains(buf.String(), "Character Set Support Comparison") {
		t.Fatalf("missing charset header")
	}
	if !strings.Contains(buf.String(), "bad") {
		t.Fatalf("missing charset detail")
	}

	shape := entity.ShapeMatrix{
		FormatNames: []string{"A"},
		Tests:       []string{"shape-a"},
		ResultsByFormat: [][]entity.ShapeResult{{
			{EncodeOK: false, DecodeOK: false, RoundTripOK: false, Detail: "bad"},
		}},
		PassCounts: []int{0},
	}
	buf.Reset()
	PrintShapeMatrix(buf, shape)
	if !strings.Contains(buf.String(), "Shape Compatibility Matrix") {
		t.Fatalf("missing shape header")
	}
	if !strings.Contains(buf.String(), "shape-a") {
		t.Fatalf("missing shape case")
	}
}

func TestTextHelpers(t *testing.T) {
	if got := ratio(10, 0); got != "n/a" {
		t.Fatalf("ratio base zero mismatch: %s", got)
	}
	if got := ratioInt(-1, 10); got != "n/a" {
		t.Fatalf("ratioInt negative mismatch: %s", got)
	}

	if got := deltaDurationText(10, 0); got != "n/a" {
		t.Fatalf("deltaDuration base zero mismatch: %s", got)
	}
	if got := deltaDurationText(10, 10); got != "same as baseline" {
		t.Fatalf("deltaDuration equal mismatch: %s", got)
	}
	if !strings.Contains(deltaDurationText(5, 10), "faster") {
		t.Fatalf("deltaDuration faster branch missing")
	}
	if !strings.Contains(deltaDurationText(15, 10), "slower") {
		t.Fatalf("deltaDuration slower branch missing")
	}

	if got := deltaIntText(-1, 10); got != "n/a" {
		t.Fatalf("deltaInt negative mismatch: %s", got)
	}
	if got := deltaIntText(10, 10); got != "same as baseline" {
		t.Fatalf("deltaInt equal mismatch: %s", got)
	}
	if !strings.Contains(deltaIntText(5, 10), "smaller") {
		t.Fatalf("deltaInt smaller branch missing")
	}
	if !strings.Contains(deltaIntText(15, 10), "larger") {
		t.Fatalf("deltaInt larger branch missing")
	}
}
