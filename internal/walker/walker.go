package walker

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/Runic-Studios/Palimpsest/internal/merger"
)

// Walk merges files from multiple overlayDirs and writes results to outputDir.
// The overlayDirs are processed in order: the earliest overlay that has a file
// is used as that file's "base," and any later overlays with the same file
// override/merge accordingly. A file present only in later overlays still gets written.
func Walk(overlayDirs []string, outputDir string) error {
	if len(overlayDirs) == 0 {
		return fmt.Errorf("at least one overlay dir is required")
	}

	// Gather *all* files from *all* overlays into a map:
	pathMap, err := gatherAllPaths(overlayDirs)
	if err != nil {
		return err
	}

	for relPath, configFile := range pathMap {
		indices := configFile.overlays
		loader := configFile.loader

		var config []merger.Config
		for _, idx := range indices {
			overlayPath := filepath.Join(overlayDirs[idx], relPath)
			overlayData, err := loader.Load(overlayPath)
			if err != nil {
				return err
			}
			config = append(config, overlayData)
		}
		data := merger.Merge(config)

		outPath := filepath.Join(outputDir, relPath)
		if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
			return err
		}

		if err := loader.Write(outPath, data); err != nil {
			return err
		}
	}

	return nil
}

type ConfigFile struct {
	overlays []int
	loader   merger.ConfigLoader
}

// gatherAllPaths walks all overlays in order and accumulates a map:
//
//	relativePath -> list of overlay indices that contain that file.
func gatherAllPaths(overlayDirs []string) (map[string]*ConfigFile, error) {
	pathMap := make(map[string]*ConfigFile)
	for i, dir := range overlayDirs {
		err := filepath.Walk(dir, func(fullPath string, info fs.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return err // skip dirs or errors
			}
			// relative path is consistent with how we write to outputDir
			rel, err := filepath.Rel(dir, fullPath)
			if err != nil {
				return err
			}

			configFile, ok := pathMap[rel]
			if !ok {
				configFile = &ConfigFile{}
				pathMap[rel] = configFile
			}

			if configFile.loader == nil {
				ext := filepath.Ext(rel)
				loader, err := merger.ForExtension(ext)
				if err != nil {
					delete(pathMap, rel)
					return nil
				}
				configFile.loader = loader
			}

			configFile.overlays = append(configFile.overlays, i)
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return pathMap, nil
}
