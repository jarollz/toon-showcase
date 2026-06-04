package presenter

import (
	"fmt"
	"io"
	"time"

	"toon-showcase/internal/core/entity"
)

func PrintBenchmarkReport(w io.Writer, results []entity.BenchmarkResult, baseline entity.BenchmarkResult, records, iters, warmup int, seed int64, indentSize int) {
	fmt.Fprintln(w)
	fmt.Fprintln(w, "=== Multi-Format Marshal/Unmarshal Showcase ===")
	fmt.Fprintf(w, "Dataset records: %d\n", records)
	fmt.Fprintf(w, "Iterations: %d (warmup: %d)\n", iters, warmup)
	fmt.Fprintf(w, "Seed: %d\n", seed)
	fmt.Fprintf(w, "Indent length used for JSON pretty + TOON + YAML + TOML + XML pretty: %d spaces\n", indentSize)
	fmt.Fprintln(w)

	fmt.Fprintf(w, "%-13s  %-13s  %-13s  %-13s  %-13s  %-8s  %-8s  %-10s\n",
		"Format", "Marshal(ms)", "Marshal(ns)", "Unmarshal(ms)", "Unmarshal(ns)", "Bytes", "Chars", "RoundTrip")
	for _, res := range results {
		chars := fmt.Sprintf("%d", res.EncodedChars)
		if res.EncodedChars < 0 {
			chars = "n/a"
		}
		fmt.Fprintf(w, "%-13s  %-13.2f  %-13d  %-13.2f  %-13d  %-8d  %-8s  %-10s\n",
			res.Name,
			float64(res.MarshalTotal.Microseconds())/1000.0,
			res.MarshalAvg.Nanoseconds(),
			float64(res.UnmarshalTotal.Microseconds())/1000.0,
			res.UnmarshalAvg.Nanoseconds(),
			res.EncodedBytes,
			chars,
			res.RoundTripDetails,
		)
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w, "Ratios vs selected baseline (lower is better):")
	fmt.Fprintf(w, "Baseline format: %s\n", baseline.Name)
	fmt.Fprintf(w, "%-13s  %-16s  %-16s  %-16s  %-16s\n", "Format", "MarshalRatio", "UnmarshalRatio", "ByteRatio", "CharRatio")
	for _, res := range results {
		fmt.Fprintf(w, "%-13s  %-16s  %-16s  %-16s  %-16s\n",
			res.Name,
			ratio(res.MarshalAvg, baseline.MarshalAvg),
			ratio(res.UnmarshalAvg, baseline.UnmarshalAvg),
			ratioInt(res.EncodedBytes, baseline.EncodedBytes),
			ratioInt(res.EncodedChars, baseline.EncodedChars),
		)
	}

	fmt.Fprintln(w)
	fmt.Fprintf(w, "Human-friendly delta vs %s baseline:\n", baseline.Name)
	fmt.Fprintln(w, "(negative direction for time is faster; negative direction for size is smaller)")
	fmt.Fprintf(w, "%-13s  %-26s  %-26s  %-24s  %-24s\n", "Format", "Marshal", "Unmarshal", "Bytes", "Chars")
	for _, res := range results {
		fmt.Fprintf(w, "%-13s  %-26s  %-26s  %-24s  %-24s\n",
			res.Name,
			deltaDurationText(res.MarshalAvg, baseline.MarshalAvg),
			deltaDurationText(res.UnmarshalAvg, baseline.UnmarshalAvg),
			deltaIntText(res.EncodedBytes, baseline.EncodedBytes),
			deltaIntText(res.EncodedChars, baseline.EncodedChars),
		)
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w, "Encoded output sample (truncated):")
	for _, res := range results {
		fmt.Fprintf(w, "- %s: %s\n", res.Name, res.Sample)
	}
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Notes:")
	fmt.Fprintln(w, "- MessagePack is binary; char metrics shown as n/a.")
	fmt.Fprintln(w, "- XML is included in two lanes: compact and pretty.")
	fmt.Fprintln(w)
}

func PrintCharsetMatrix(w io.Writer, report entity.CharsetMatrix) {
	fmt.Fprintln(w, "=== Character Set Support Comparison ===")
	fmt.Fprintln(w, "Scope: practical behavior in selected Go libraries")
	fmt.Fprintln(w)

	fmt.Fprintf(w, "%-17s", "Case")
	for _, name := range report.FormatNames {
		fmt.Fprintf(w, "  %-16s", name)
	}
	fmt.Fprintln(w)

	for i, tc := range report.Tests {
		fmt.Fprintf(w, "%-17s", tc.Name)
		for fi := range report.FormatNames {
			cell := "FAIL"
			if report.ResultsByFormat[fi][i].MarshalOK && report.ResultsByFormat[fi][i].RoundTripOK {
				cell = "OK"
			}
			fmt.Fprintf(w, "  %-16s", cell)
		}
		fmt.Fprintln(w)
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w, "Support summary:")
	for fi, name := range report.FormatNames {
		fmt.Fprintf(w, "- %s: %d/%d cases roundtrip\n", name, report.PassCounts[fi], len(report.Tests))
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w, "Notable details:")
	for fi, name := range report.FormatNames {
		for ti, tc := range report.Tests {
			res := report.ResultsByFormat[fi][ti]
			if res.MarshalOK && res.RoundTripOK {
				continue
			}
			fmt.Fprintf(w, "- %s / %s: %s\n", name, tc.Name, res.Detail)
		}
	}
	fmt.Fprintln(w)
}

func PrintShapeMatrix(w io.Writer, report entity.ShapeMatrix) {
	fmt.Fprintln(w, "=== Shape Compatibility Matrix ===")
	fmt.Fprintln(w, "Scope: structural compatibility of complex object shapes")
	fmt.Fprintln(w)

	fmt.Fprintf(w, "%-27s", "Shape case")
	for _, name := range report.FormatNames {
		fmt.Fprintf(w, "  %-16s", name)
	}
	fmt.Fprintln(w)

	for i, testName := range report.Tests {
		fmt.Fprintf(w, "%-27s", testName)
		for fi := range report.FormatNames {
			res := report.ResultsByFormat[fi][i]
			cell := "FAIL"
			if res.EncodeOK && res.DecodeOK && res.RoundTripOK {
				cell = "OK"
			}
			fmt.Fprintf(w, "  %-16s", cell)
		}
		fmt.Fprintln(w)
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w, "Support summary:")
	for fi, name := range report.FormatNames {
		fmt.Fprintf(w, "- %s: %d/%d shape cases roundtrip\n", name, report.PassCounts[fi], len(report.Tests))
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w, "Notable details:")
	for fi, name := range report.FormatNames {
		for ti, testName := range report.Tests {
			res := report.ResultsByFormat[fi][ti]
			if res.EncodeOK && res.DecodeOK && res.RoundTripOK {
				continue
			}
			fmt.Fprintf(w, "- %s / %s: %s\n", name, testName, res.Detail)
		}
	}
	fmt.Fprintln(w)
}

func ratio(v, base time.Duration) string {
	if base == 0 {
		return "n/a"
	}
	return fmt.Sprintf("%.3fx", float64(v)/float64(base))
}

func ratioInt(v, base int) string {
	if base <= 0 || v < 0 {
		return "n/a"
	}
	return fmt.Sprintf("%.3fx", float64(v)/float64(base))
}

func deltaDurationText(v, base time.Duration) string {
	if base == 0 {
		return "n/a"
	}
	if v == base {
		return "same as baseline"
	}
	if v < base {
		pct := (float64(base-v) / float64(base)) * 100
		return fmt.Sprintf("faster %.1f%%", pct)
	}
	pct := (float64(v-base) / float64(base)) * 100
	return fmt.Sprintf("slower %.1f%%", pct)
}

func deltaIntText(v, base int) string {
	if base <= 0 || v < 0 {
		return "n/a"
	}
	if v == base {
		return "same as baseline"
	}
	if v < base {
		pct := (float64(base-v) / float64(base)) * 100
		return fmt.Sprintf("smaller %.1f%%", pct)
	}
	pct := (float64(v-base) / float64(base)) * 100
	return fmt.Sprintf("larger %.1f%%", pct)
}
