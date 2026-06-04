package presenter

import (
	"fmt"
	"strings"
	"time"

	"charm.land/glamour/v2"

	"toon-showcase/internal/core/entity"
)

func BuildFullMarkdownReport(results []entity.BenchmarkResult, baseline entity.BenchmarkResult, charset entity.CharsetMatrix, shape entity.ShapeMatrix, records, iters, warmup int, seed int64, indentSize int) string {
	b := &strings.Builder{}
	fmt.Fprintln(b, "# Multi-Format Marshal/Unmarshal Showcase")
	fmt.Fprintln(b)
	fmt.Fprintf(b, "- Dataset records: %d\n- Iterations: %d (warmup: %d)\n- Seed: %d\n- Indent length: %d\n\n", records, iters, warmup, seed, indentSize)

	fmt.Fprintln(b, "## Benchmark")
	fmt.Fprintln(b)
	fmt.Fprintln(b, "| Format | Marshal (ns/op) | Unmarshal (ns/op) | Bytes | Roundtrip |")
	fmt.Fprintln(b, "| --- | ---: | ---: | ---: | :---: |")
	for _, r := range results {
		fmt.Fprintf(b, "| %s | %d | %d | %d | %s |\n", mdEscapeCell(r.Name), r.MarshalAvg.Nanoseconds(), r.UnmarshalAvg.Nanoseconds(), r.EncodedBytes, mdEscapeCell(r.RoundTripDetails))
	}

	fmt.Fprintf(b, "\nRatios vs selected baseline (lower is better):\nBaseline format: %s\n\n", mdEscapeCell(baseline.Name))
	fmt.Fprintln(b, "| Format | MarshalRatio | UnmarshalRatio | ByteRatio | CharRatio |")
	fmt.Fprintln(b, "| --- | ---: | ---: | ---: | ---: |")
	for _, r := range results {
		fmt.Fprintf(b, "| %s | %s | %s | %s | %s |\n", mdEscapeCell(r.Name), ratio(r.MarshalAvg, baseline.MarshalAvg), ratio(r.UnmarshalAvg, baseline.UnmarshalAvg), ratioInt(r.EncodedBytes, baseline.EncodedBytes), ratioInt(r.EncodedChars, baseline.EncodedChars))
	}

	fmt.Fprintf(b, "\nHuman-friendly delta vs %s baseline:\n\n", mdEscapeCell(baseline.Name))
	fmt.Fprintln(b, "(negative direction for time is faster; negative direction for size is smaller)")
	fmt.Fprintln(b)
	fmt.Fprintln(b, "| Format | Marshal | Unmarshal | Bytes | Chars |")
	fmt.Fprintln(b, "| --- | --- | --- | --- | --- |")
	for _, r := range results {
		fmt.Fprintf(b, "| %s | %s | %s | %s | %s |\n", mdEscapeCell(r.Name), mdEscapeCell(deltaDurationText(r.MarshalAvg, baseline.MarshalAvg)), mdEscapeCell(deltaDurationText(r.UnmarshalAvg, baseline.UnmarshalAvg)), mdEscapeCell(deltaIntText(r.EncodedBytes, baseline.EncodedBytes)), mdEscapeCell(deltaIntText(r.EncodedChars, baseline.EncodedChars)))
	}

	fmt.Fprintf(b, "\nPercent delta vs %s baseline (`-` faster/smaller, `+` slower/larger):\n\n", mdEscapeCell(baseline.Name))
	fmt.Fprintln(b, "Emoji legend: `✅` better, `❌` worse, `➖` same.")
	fmt.Fprintln(b)
	fmt.Fprintln(b, "| Format | Marshal delta | Unmarshal delta | Bytes delta |")
	fmt.Fprintln(b, "| --- | ---: | ---: | ---: |")
	for _, r := range results {
		fmt.Fprintf(b, "| %s | %s | %s | %s |\n", mdEscapeCell(r.Name), deltaCellDuration(r.MarshalAvg, baseline.MarshalAvg), deltaCellDuration(r.UnmarshalAvg, baseline.UnmarshalAvg), deltaCellInt(r.EncodedBytes, baseline.EncodedBytes))
	}

	fmt.Fprintln(b, "\nEncoded output sample (truncated):")
	for _, r := range results {
		fmt.Fprintf(b, "- %s: %s\n", mdEscapeCell(r.Name), mdEscapeCell(r.Sample))
	}
	fmt.Fprintln(b)
	fmt.Fprintln(b, "Notes:")
	fmt.Fprintln(b, "- MessagePack is binary; char metrics shown as n/a.")
	fmt.Fprintln(b, "- XML is included in two lanes: compact and pretty.")
	fmt.Fprintln(b)

	if len(charset.Tests) == 0 || len(charset.FormatNames) == 0 {
		fmt.Fprintln(b, "## Character Set Support Comparison")
		fmt.Fprintln(b)
		fmt.Fprintln(b, "No character set report data.")
		fmt.Fprintln(b)
	} else {
		tests := make([]string, 0, len(charset.Tests))
		for _, tc := range charset.Tests {
			tests = append(tests, tc.Name)
		}
		writeMatrixSection(b, "Character Set Support Comparison", "Scope: practical behavior in selected Go libraries", "Case", tests, charset.FormatNames, charset.PassCounts, len(charset.Tests), func(fi, ti int) (string, bool, string) {
			r := charset.ResultsByFormat[fi][ti]
			ok := r.MarshalOK && r.RoundTripOK
			if ok {
				return "OK", true, ""
			}
			return "FAIL", false, r.Detail
		})
	}

	if len(shape.Tests) == 0 || len(shape.FormatNames) == 0 {
		fmt.Fprintln(b, "## Shape Compatibility Matrix")
		fmt.Fprintln(b)
		fmt.Fprintln(b, "No shape report data.")
		fmt.Fprintln(b)
	} else {
		writeMatrixSection(b, "Shape Compatibility Matrix", "Scope: structural compatibility of complex object shapes", "Shape case", shape.Tests, shape.FormatNames, shape.PassCounts, len(shape.Tests), func(fi, ti int) (string, bool, string) {
			r := shape.ResultsByFormat[fi][ti]
			ok := r.EncodeOK && r.DecodeOK && r.RoundTripOK
			if ok {
				return "OK", true, ""
			}
			return "FAIL", false, r.Detail
		})
	}

	return b.String()
}

func writeMatrixSection(b *strings.Builder, title, scope, leadCol string, tests, formats []string, passCounts []int, total int, cell func(fi, ti int) (status string, pass bool, detail string)) {
	fmt.Fprintf(b, "## %s\n\n%s\n\n", title, scope)
	fmt.Fprintf(b, "| %s |", leadCol)
	for _, name := range formats {
		fmt.Fprintf(b, " %s |", mdEscapeCell(name))
	}
	fmt.Fprintln(b)
	fmt.Fprint(b, "| --- |")
	for range formats {
		fmt.Fprint(b, " :---: |")
	}
	fmt.Fprintln(b)

	details := make([]string, 0)
	for ti, test := range tests {
		fmt.Fprintf(b, "| %s |", mdEscapeCell(test))
		for fi := range formats {
			status, pass, detail := cell(fi, ti)
			fmt.Fprintf(b, " %s |", status)
			if !pass && detail != "" {
				details = append(details, fmt.Sprintf("- %s / %s: %s", mdEscapeCell(formats[fi]), mdEscapeCell(test), mdEscapeCell(detail)))
			}
		}
		fmt.Fprintln(b)
	}

	fmt.Fprintln(b)
	fmt.Fprintln(b, "Support summary:")
	for fi, name := range formats {
		fmt.Fprintf(b, "- %s: %d/%d %s roundtrip\n", mdEscapeCell(name), passCounts[fi], total, sectionNoun(title))
	}
	fmt.Fprintln(b)
	fmt.Fprintln(b, "Notable details:")
	for _, line := range details {
		fmt.Fprintln(b, line)
	}
	fmt.Fprintln(b)
}

func sectionNoun(title string) string {
	if strings.Contains(title, "Shape") {
		return "shape cases"
	}
	return "cases"
}

func BuildBenchmarkMarkdownReport(results []entity.BenchmarkResult, baseline entity.BenchmarkResult, records, iters, warmup int, seed int64, indentSize int) string {
	return BuildFullMarkdownReport(results, baseline, entity.CharsetMatrix{}, entity.ShapeMatrix{}, records, iters, warmup, seed, indentSize)
}

func RenderMarkdownDark(md string) (string, error) { return glamour.Render(md, "dark") }

func deltaCellDuration(v, base time.Duration) string {
	if base == 0 {
		return "n/a"
	}
	if v == base {
		return "0.0% ➖"
	}
	pct := (float64(v-base) / float64(base)) * 100
	if pct < 0 {
		return fmt.Sprintf("%.1f%% ✅", pct)
	}
	return fmt.Sprintf("+%.1f%% ❌", pct)
}

func deltaCellInt(v, base int) string {
	if base <= 0 || v < 0 {
		return "n/a"
	}
	if v == base {
		return "0.0% ➖"
	}
	pct := (float64(v-base) / float64(base)) * 100
	if pct < 0 {
		return fmt.Sprintf("%.1f%% ✅", pct)
	}
	return fmt.Sprintf("+%.1f%% ❌", pct)
}

func mdEscapeCell(v string) string {
	v = strings.ReplaceAll(v, "|", "\\|")
	v = strings.ReplaceAll(v, "\r\n", "<br>")
	v = strings.ReplaceAll(v, "\n", "<br>")
	v = strings.ReplaceAll(v, "\r", "<br>")
	return v
}
