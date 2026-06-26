package manifest

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type Manifest struct {
	SchemaVersion int     `json:"schema_version"`
	GeneratedAt   string  `json:"generated_at"`
	Target        string  `json:"target"`
	OutputPath    string  `json:"output_path"`
	Chunks        []Chunk `json:"chunks"`
}

type Chunk struct {
	SourceFile         string `json:"source_file"`
	Section            string `json:"section"`
	AssembledLineStart int    `json:"assembled_line_start"`
	AssembledLineEnd   int    `json:"assembled_line_end"`
	SourceLineStart    int    `json:"source_line_start"`
	SourceLineEnd      int    `json:"source_line_end"`
}

func sidecarPath(root, outputPath string) string {
	base := filepath.Base(outputPath)
	return filepath.Join(root, ".jctx", ".manifest", base+".json")
}

func Write(root, outputPath string, m *Manifest) error {
	p := sidecarPath(root, outputPath)
	if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0644)
}

func Read(root, outputPath string) (*Manifest, error) {
	p := sidecarPath(root, outputPath)
	data, err := os.ReadFile(p)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

func Delete(root, outputPath string) error {
	p := sidecarPath(root, outputPath)
	err := os.Remove(p)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

func ListAll(root string) ([]string, error) {
	dir := filepath.Join(root, ".jctx", ".manifest")
	entries, err := os.ReadDir(dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}

	var paths []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return nil, err
		}
		var m Manifest
		if err := json.Unmarshal(data, &m); err != nil {
			return nil, err
		}
		paths = append(paths, m.OutputPath)
	}
	return paths, nil
}
