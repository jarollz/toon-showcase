package presenter

import (
	"fmt"
	"os"
	"path/filepath"

	"toon-showcase/internal/core/entity"
)

func WriteShapeEncodedExamples(baseDir string, artifacts []entity.ShapeEncodedArtifact) error {
	root := filepath.Join(baseDir, "shape-encoded")
	for _, artifact := range artifacts {
		formatDir := filepath.Join(root, artifact.FormatKey)
		if err := os.MkdirAll(formatDir, 0o755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}

		if artifact.Result.EncodeOK {
			ext := ".txt"
			if artifact.Binary {
				ext = ".bin"
			}
			path := filepath.Join(formatDir, artifact.CaseName+ext)
			if err := os.WriteFile(path, artifact.Encoded, 0o644); err != nil {
				return fmt.Errorf("failed to write encoded example %s/%s: %w", artifact.FormatKey, artifact.CaseName, err)
			}
			continue
		}

		path := filepath.Join(formatDir, artifact.CaseName+".error.txt")
		if err := os.WriteFile(path, []byte(artifact.Result.Detail+"\n"), 0o644); err != nil {
			return fmt.Errorf("failed to write encoded example error %s/%s: %w", artifact.FormatKey, artifact.CaseName, err)
		}
	}
	return nil
}
