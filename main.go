package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"toon-showcase/internal/adapter/codec"
	"toon-showcase/internal/adapter/presenter"
	"toon-showcase/internal/core/usecase"
	"toon-showcase/internal/infrastructure/cli"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr, time.Now))
}

func run(args []string, stdout, stderr io.Writer, now func() time.Time) int {
	opts, err := cli.Parse(args)
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		return 2
	}

	runner := usecase.Runner{
		Codecs:      codec.Build(opts.IndentSize),
		Now:         now,
		NewProgress: presenter.NewProgressBar,
	}

	out, err := runner.Run(opts)
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		return 1
	}

	presenter.PrintBenchmarkReport(stdout, out.Results, out.Baseline, opts.Records, opts.Iters, opts.Warmup, out.ResolvedSeed, opts.IndentSize)
	presenter.PrintCharsetMatrix(stdout, out.Charset)
	presenter.PrintShapeMatrix(stdout, out.Shape)
	if opts.OutputExample {
		fmt.Fprintf(stdout, "Encoded examples output dir: %s\n", out.ResolvedOutputExampleDir)
		if err := presenter.WriteShapeEncodedExamples(out.ResolvedOutputExampleDir, out.Shape.EncodedArtifacts); err != nil {
			fmt.Fprintln(stderr, fmt.Sprintf("failed to write encoded examples: %v", err))
			return 1
		}
	}

	if opts.CSVMode == "off" {
		return 0
	}

	benchmarkCSV, err := presenter.BuildBenchmarkCSV(out.Results, out.Baseline, opts.Records, opts.Iters, opts.Warmup, out.ResolvedSeed, opts.IndentSize)
	if err != nil {
		fmt.Fprintln(stderr, fmt.Sprintf("failed to build benchmark csv: %v", err))
		return 1
	}
	charsetCSV, err := presenter.BuildCharsetCSV(out.Charset)
	if err != nil {
		fmt.Fprintln(stderr, fmt.Sprintf("failed to build charset csv: %v", err))
		return 1
	}
	shapeCSV, err := presenter.BuildShapeCSV(out.Shape)
	if err != nil {
		fmt.Fprintln(stderr, fmt.Sprintf("failed to build shape csv: %v", err))
		return 1
	}
	if err := presenter.EmitCSVs(stdout, opts.CSVMode, opts.CSVFile, benchmarkCSV, charsetCSV, shapeCSV); err != nil {
		fmt.Fprintln(stderr, err.Error())
		return 1
	}

	return 0
}
