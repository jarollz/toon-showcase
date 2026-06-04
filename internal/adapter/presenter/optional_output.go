package presenter

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"toon-showcase/internal/core/entity"
)

func IsMarkdownStdout(mdFile string) bool {
	return strings.TrimSpace(mdFile) == "stdout"
}

func EmitEncodedExamplesOutput(w io.Writer, rawDir, resolvedDir string, artifacts []entity.ShapeEncodedArtifact) error {
	trimmed := strings.TrimSpace(rawDir)
	if trimmed == "" {
		return nil
	}
	if err := ensureDirectory(trimmed); err != nil {
		return fmt.Errorf("invalid encoded examples dir: %w", err)
	}
	fmt.Fprintf(w, "Encoded examples output dir: %s\n", resolvedDir)
	if err := WriteShapeEncodedExamples(resolvedDir, artifacts); err != nil {
		return fmt.Errorf("failed to write encoded examples: %w", err)
	}
	return nil
}

func EmitMarkdownOutput(w io.Writer, rawPath, markdown string) error {
	trimmed := strings.TrimSpace(rawPath)
	if trimmed == "" {
		return nil
	}
	if trimmed == "stdout" {
		rendered, err := RenderMarkdownDark(markdown)
		if err != nil {
			return fmt.Errorf("failed to render markdown: %w", err)
		}
		fmt.Fprint(w, rendered)
		return nil
	}
	if err := ensureMarkdownParentDir(trimmed); err != nil {
		return fmt.Errorf("invalid md-file path: %w", err)
	}
	if err := os.WriteFile(trimmed, []byte(markdown), 0644); err != nil {
		return fmt.Errorf("failed to write markdown file: %w", err)
	}
	fmt.Fprintf(w, "Markdown written: %s\n", trimmed)
	return nil
}

func EmitCSVOutput(w io.Writer, rawCSVDir string, results []entity.BenchmarkResult, baseline entity.BenchmarkResult, charset entity.CharsetMatrix, shape entity.ShapeMatrix, records, iters, warmup int, seed int64, indentSize int) error {
	trimmed := strings.TrimSpace(rawCSVDir)
	if trimmed == "" {
		return nil
	}
	if trimmed != "stdout" {
		if err := ensureDirectory(trimmed); err != nil {
			return fmt.Errorf("invalid csv dir: %w", err)
		}
	}
	benchmarkCSV, err := BuildBenchmarkCSV(results, baseline, records, iters, warmup, seed, indentSize)
	if err != nil {
		return fmt.Errorf("failed to build benchmark csv: %w", err)
	}
	charsetCSV, err := BuildCharsetCSV(charset)
	if err != nil {
		return fmt.Errorf("failed to build charset csv: %w", err)
	}
	shapeCSV, err := BuildShapeCSV(shape)
	if err != nil {
		return fmt.Errorf("failed to build shape csv: %w", err)
	}
	if err := EmitCSVs(w, trimmed, benchmarkCSV, charsetCSV, shapeCSV); err != nil {
		return err
	}
	return nil
}

func ensureMarkdownParentDir(path string) error {
	parent := filepath.Dir(path)
	if parent == "." {
		return nil
	}
	return os.MkdirAll(parent, 0755)
}
