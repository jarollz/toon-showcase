package cli

import (
	"flag"

	"toon-showcase/internal/core/usecase"
)

func Parse(args []string) (usecase.Options, error) {
	fs := flag.NewFlagSet("toon-showcase", flag.ContinueOnError)

	records := fs.Int("records", 120, "number of orders in dataset")
	iters := fs.Int("iters", 500, "benchmark iterations for each marshal and unmarshal")
	warmup := fs.Int("warmup", 50, "warmup iterations before measuring")
	seed := fs.Int64("seed", -1, "seed for dataset generation; use -1 for current unix timestamp")
	indentSize := fs.Int("indent", 2, "explicit indent size for JSON pretty, TOON, YAML, and TOML")
	sampleLimit := fs.Int("sample-limit", 180, "max characters shown for encoded sample")
	noProgress := fs.Bool("no-progress", false, "disable animated progress bar")
	baseline := fs.String("baseline", "toon", "baseline format for ratio/delta: toon|json-compact|json-pretty|yaml|toml|xml-compact|xml-pretty|messagepack")
	csvMode := fs.String("csv", "off", "csv output mode: off|stdout|file")
	csvFile := fs.String("csv-file", "benchmark_report.csv", "csv output path when -csv=file")
	outputExample := fs.Bool("output-example", false, "write full encoded examples by format and shape case")
	outputExampleDir := fs.String("output-example-dir", "", "base directory for encoded examples; empty means auto default")

	if err := fs.Parse(args); err != nil {
		return usecase.Options{}, err
	}

	return usecase.Options{
		Records:          *records,
		Iters:            *iters,
		Warmup:           *warmup,
		Seed:             *seed,
		IndentSize:       *indentSize,
		SampleLimit:      *sampleLimit,
		NoProgress:       *noProgress,
		Baseline:         *baseline,
		CSVMode:          *csvMode,
		CSVFile:          *csvFile,
		OutputExample:    *outputExample,
		OutputExampleDir: *outputExampleDir,
	}, nil
}
