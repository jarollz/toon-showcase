package usecase

import (
	"fmt"

	"toon-showcase/internal/core/entity"
)

func (r Runner) Run(opts Options) (RunOutput, error) {
	if err := opts.Validate(); err != nil {
		return RunOutput{}, err
	}
	if len(r.Codecs) == 0 {
		return RunOutput{}, fmt.Errorf("no codecs configured")
	}

	baselineKey, err := NormalizeBaselineKey(opts.Baseline)
	if err != nil {
		return RunOutput{}, err
	}

	now := r.Now
	if now == nil {
		return RunOutput{}, fmt.Errorf("clock is not configured")
	}

	seed := opts.Seed
	if seed == -1 {
		seed = now().Unix()
	}

	dataset := GenerateDataset(opts.Records, seed)

	var progress Progress
	if !opts.NoProgress && r.NewProgress != nil {
		progress = r.NewProgress(totalBenchmarkSteps(len(r.Codecs), opts.Iters, opts.Warmup))
		progress.SetStage("initializing benchmark")
		progress.Start()
	}

	results := make([]entity.BenchmarkResult, 0, len(r.Codecs))
	for _, c := range r.Codecs {
		if progress != nil {
			progress.SetStage("benchmark " + c.Name())
		}
		res, benchErr := runBenchmark(c, dataset, opts.Iters, opts.Warmup, opts.SampleLimit, progress)
		if benchErr != nil {
			if progress != nil {
				progress.Finish()
			}
			return RunOutput{}, fmt.Errorf("%s failed: %w", c.Name(), benchErr)
		}
		results = append(results, res)
	}
	if progress != nil {
		progress.Finish()
	}

	baseline, err := findBaselineResult(results, baselineKey)
	if err != nil {
		return RunOutput{}, err
	}

	return RunOutput{
		Results:      results,
		Baseline:     baseline,
		Charset:      runCharsetMatrix(r.Codecs),
		Shape:        runShapeMatrix(r.Codecs),
		ResolvedSeed: seed,
		BaselineKey:  baselineKey,
	}, nil
}

func totalBenchmarkSteps(formatCount, iters, warmup int) int64 {
	if formatCount <= 0 {
		return 1
	}
	perFormat := int64(warmup*2 + iters*2 + 1)
	total := int64(formatCount) * perFormat
	if total <= 0 {
		return 1
	}
	return total
}

func findBaselineResult(results []entity.BenchmarkResult, baselineKey string) (entity.BenchmarkResult, error) {
	for _, res := range results {
		if res.Key == baselineKey {
			return res, nil
		}
	}
	return entity.BenchmarkResult{}, fmt.Errorf("baseline key %q not found in benchmark results", baselineKey)
}
