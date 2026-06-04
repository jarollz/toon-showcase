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

	if !presenter.IsMarkdownStdout(opts.MarkdownFile) {
		presenter.PrintBenchmarkReport(stdout, out.Results, out.Baseline, opts.Records, opts.Iters, opts.Warmup, out.ResolvedSeed, opts.IndentSize)
		presenter.PrintCharsetMatrix(stdout, out.Charset)
		presenter.PrintShapeMatrix(stdout, out.Shape)
	}

	if err := presenter.EmitEncodedExamplesOutput(stdout, opts.EncodedExamplesDir, out.ResolvedOutputExampleDir, out.Shape.EncodedArtifacts); err != nil {
		fmt.Fprintln(stderr, err.Error())
		return 1
	}

	if err := presenter.EmitMarkdownOutput(stdout, opts.MarkdownFile, presenter.BuildFullMarkdownReport(out.Results, out.Baseline, out.Charset, out.Shape, opts.Records, opts.Iters, opts.Warmup, out.ResolvedSeed, opts.IndentSize)); err != nil {
		fmt.Fprintln(stderr, err.Error())
		return 1
	}

	if err := presenter.EmitCSVOutput(stdout, opts.CSVDir, out.Results, out.Baseline, out.Charset, out.Shape, opts.Records, opts.Iters, opts.Warmup, out.ResolvedSeed, opts.IndentSize); err != nil {
		fmt.Fprintln(stderr, err.Error())
		return 1
	}

	return 0
}
