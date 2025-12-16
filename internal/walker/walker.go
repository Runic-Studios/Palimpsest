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
func Walk(overlayDirs []string, outputDir string, verbose bool) error {
	if len(overlayDirs) == 0 {
		return fmt.Errorf("at least one overlay dir is required")
	}

	// Gather *all* files from *all* overlays into a map:
	pathMap, err := gatherAllPaths(overlayDirs, verbose)
	if err != nil {
		return err
	}

	for relPath, configFile := range pathMap {
		indices := configFile.overlays
		loader := configFile.loader

		var config []merger.Config
		for _, idx := range indices {
			overlayPath := filepath.Join(overlayDirs[idx], relPath)
			if verbose {
				fmt.Printf("Loading %s from overlay %d: %s\n", relPath, idx, overlayPath)
			}
			overlayData, err := loader.Load(overlayPath)
			if err != nil {
				return err
			}
			config = append(config, overlayData)
		}
		if len(config) > 0 {
			if len(config) == 1 {
				idx := indices[0]
				absOverlay, err1 := filepath.Abs(overlayDirs[idx])
				absOutput, err2 := filepath.Abs(outputDir)
				if err1 == nil && err2 == nil && absOverlay == absOutput {
					if verbose {
						fmt.Printf("Skipping %s (only present in output dir)\n", relPath)
					}
					continue
				}
			}

			if verbose {
				fmt.Printf("Merging %d configs for %s\n", len(config), relPath)
			}
			data := merger.Merge(config)

			outPath := filepath.Join(outputDir, relPath)
			if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
				return err
			}

			if err := loader.Write(outPath, data); err != nil {
				return err
			}
			if verbose {
				fmt.Printf("Written %s to %s\n", relPath, outPath)
			}
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
func gatherAllPaths(overlayDirs []string, verbose bool) (map[string]*ConfigFile, error) {
	pathMap := make(map[string]*ConfigFile)
	for i, dir := range overlayDirs {
		if verbose {
			fmt.Printf("Scanning overlay %d: %s\n", i, dir)
		}
		err := filepath.Walk(dir, func(fullPath string, info fs.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return err // skip dirs or errors
			}
			// relative path is consistent with how we write to outputDir
			rel, err := filepath.Rel(dir, fullPath)
			if err != nil {
				return err
			}

			if verbose {
				fmt.Printf("Found file %s in overlay %d\n", rel, i)
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
					if verbose {
						fmt.Printf("Skipping %s (unsupported extension: %s)\n", rel, ext)
					}
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
