package presenter

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"toon-showcase/internal/core/entity"
)

func BuildBenchmarkCSV(results []entity.BenchmarkResult, baseline entity.BenchmarkResult, records, iters, warmup int, seed int64, indentSize int) (string, error) {
	if len(results) == 0 {
		return "", nil
	}

	buf := &bytes.Buffer{}
	w := csv.NewWriter(buf)
	err := w.Write([]string{
		"format", "baseline_format", "records", "iters", "warmup", "seed", "indent_size", "binary",
		"marshal_total_ms", "marshal_avg_ns", "unmarshal_total_ms", "unmarshal_avg_ns",
		"encoded_bytes", "encoded_chars", "roundtrip_ok", "roundtrip_details",
		"marshal_ratio_vs_baseline", "unmarshal_ratio_vs_baseline",
		"byte_ratio_vs_baseline", "char_ratio_vs_baseline",
	})
	if err != nil {
		return "", err
	}

	for _, res := range results {
		chars := ""
		if res.EncodedChars >= 0 {
			chars = fmt.Sprintf("%d", res.EncodedChars)
		}
		err = w.Write([]string{
			res.Name,
			baseline.Name,
			fmt.Sprintf("%d", records),
			fmt.Sprintf("%d", iters),
			fmt.Sprintf("%d", warmup),
			fmt.Sprintf("%d", seed),
			fmt.Sprintf("%d", indentSize),
			fmt.Sprintf("%t", res.Binary),
			fmt.Sprintf("%.6f", float64(res.MarshalTotal.Microseconds())/1000.0),
			fmt.Sprintf("%d", res.MarshalAvg.Nanoseconds()),
			fmt.Sprintf("%.6f", float64(res.UnmarshalTotal.Microseconds())/1000.0),
			fmt.Sprintf("%d", res.UnmarshalAvg.Nanoseconds()),
			fmt.Sprintf("%d", res.EncodedBytes),
			chars,
			fmt.Sprintf("%t", res.RoundTripOK),
			res.RoundTripDetails,
			ratioFloat(res.MarshalAvg, baseline.MarshalAvg),
			ratioFloat(res.UnmarshalAvg, baseline.UnmarshalAvg),
			ratioIntFloat(res.EncodedBytes, baseline.EncodedBytes),
			ratioIntFloat(res.EncodedChars, baseline.EncodedChars),
		})
		if err != nil {
			return "", err
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func BuildCharsetCSV(report entity.CharsetMatrix) (string, error) {
	buf := &bytes.Buffer{}
	w := csv.NewWriter(buf)

	err := w.Write([]string{"format", "case", "marshal_ok", "roundtrip_ok", "status", "detail"})
	if err != nil {
		return "", err
	}

	for fi, name := range report.FormatNames {
		for ti, tc := range report.Tests {
			res := report.ResultsByFormat[fi][ti]
			status := "FAIL"
			if res.MarshalOK && res.RoundTripOK {
				status = "OK"
			}
			err = w.Write([]string{name, tc.Name, fmt.Sprintf("%t", res.MarshalOK), fmt.Sprintf("%t", res.RoundTripOK), status, res.Detail})
			if err != nil {
				return "", err
			}
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func BuildShapeCSV(report entity.ShapeMatrix) (string, error) {
	buf := &bytes.Buffer{}
	w := csv.NewWriter(buf)

	err := w.Write([]string{"format", "case", "encode_ok", "decode_ok", "roundtrip_ok", "status", "detail"})
	if err != nil {
		return "", err
	}

	for fi, name := range report.FormatNames {
		for ti, testName := range report.Tests {
			res := report.ResultsByFormat[fi][ti]
			status := "FAIL"
			if res.EncodeOK && res.DecodeOK && res.RoundTripOK {
				status = "OK"
			}
			err = w.Write([]string{
				name,
				testName,
				fmt.Sprintf("%t", res.EncodeOK),
				fmt.Sprintf("%t", res.DecodeOK),
				fmt.Sprintf("%t", res.RoundTripOK),
				status,
				res.Detail,
			})
			if err != nil {
				return "", err
			}
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func EmitCSVs(w io.Writer, mode, benchmarkPath, benchmarkCSV, charsetCSV, shapeCSV string) error {
	if mode == "stdout" {
		fmt.Fprintln(w, "=== CSV Benchmark Output ===")
		fmt.Fprint(w, benchmarkCSV)
		fmt.Fprintln(w, "=== CSV Charset Output ===")
		fmt.Fprint(w, charsetCSV)
		fmt.Fprintln(w, "=== CSV Shape Compatibility Output ===")
		fmt.Fprint(w, shapeCSV)
		return nil
	}

	if err := os.WriteFile(benchmarkPath, []byte(benchmarkCSV), 0644); err != nil {
		return fmt.Errorf("failed to write benchmark csv: %w", err)
	}
	charsetPath := deriveCSVPath(benchmarkPath, ".charset.csv")
	if err := os.WriteFile(charsetPath, []byte(charsetCSV), 0644); err != nil {
		return fmt.Errorf("failed to write charset csv: %w", err)
	}
	shapePath := deriveCSVPath(benchmarkPath, ".shape.csv")
	if err := os.WriteFile(shapePath, []byte(shapeCSV), 0644); err != nil {
		return fmt.Errorf("failed to write shape csv: %w", err)
	}
	fmt.Fprintf(w, "CSV written: %s\n", benchmarkPath)
	fmt.Fprintf(w, "CSV written: %s\n", charsetPath)
	fmt.Fprintf(w, "CSV written: %s\n", shapePath)
	return nil
}

func deriveCSVPath(path, suffix string) string {
	if strings.HasSuffix(path, ".csv") {
		return strings.TrimSuffix(path, ".csv") + suffix
	}
	return path + suffix
}

func ratioFloat(v, base time.Duration) string {
	if base == 0 {
		return ""
	}
	return fmt.Sprintf("%.6f", float64(v)/float64(base))
}

func ratioIntFloat(v, base int) string {
	if base <= 0 || v < 0 {
		return ""
	}
	return fmt.Sprintf("%.6f", float64(v)/float64(base))
}
