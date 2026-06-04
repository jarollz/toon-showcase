package cli

import (
	"strings"
	"testing"
	"time"

	"toon-showcase/internal/adapter/presenter"
	"toon-showcase/internal/core/entity"
)

func TestCoverageMarkdownSmoke(t *testing.T) {
	results := []entity.BenchmarkResult{{
		Name:             "TOON",
		MarshalAvg:       1 * time.Nanosecond,
		UnmarshalAvg:     2 * time.Nanosecond,
		EncodedBytes:     1,
		EncodedChars:     1,
		RoundTripDetails: "ok",
		Sample:           "x",
	}}
	charset := entity.CharsetMatrix{
		FormatNames: []string{"TOON"},
		Tests:       []entity.CharsetCase{{Name: "ASCII", Value: "x"}},
		ResultsByFormat: [][]entity.CharsetResult{{
			{MarshalOK: true, RoundTripOK: true, Detail: "ok"},
		}},
		PassCounts: []int{1},
	}
	shape := entity.ShapeMatrix{
		FormatNames: []string{"TOON"},
		Tests:       []string{"shape-a"},
		ResultsByFormat: [][]entity.ShapeResult{{
			{EncodeOK: true, DecodeOK: true, RoundTripOK: true, Detail: "ok"},
		}},
		PassCounts: []int{1},
	}
	md := presenter.BuildFullMarkdownReport(results, results[0], charset, shape, 1, 1, 0, 1, 2)
	if !strings.Contains(md, "## Character Set Support Comparison") {
		t.Fatalf("markdown smoke output missing charset section")
	}
}
