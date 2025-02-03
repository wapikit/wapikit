package main

import (
	"os"
	"path"
	"strings"

	"github.com/knadh/stuffbin"
)

func joinFSPaths(root string, paths []string) []string {
	out := make([]string, 0, len(paths))
	for _, p := range paths {
		// real_path:stuffbin_alias
		f := strings.Split(p, ":")
		out = append(out, path.Join(root, f[0])+":"+f[1])
	}
	return out
}

// initFileSystem initializes the stuffbin FileSystem to provide
// access to bundled static assets to the app.
func initFS(appDir, frontendDir string) stuffbin.FileSystem {
	var (
		// These paths are joined with "." which is appDir.
		appFiles = []string{
			"./config.toml.sample:config.toml.sample",
			"./internal/database/migrations/:/migrations/",
		}

		// These path are joined with frontend/out dir
		frontendFiles = []string{
			// frontend/out files should be available on the root path following the file path .
			"./:/",
		}

		// ! TODO: add a static dir path if somebody mounts any other static directory here
	)

	// Get the executable's execPath.
	execPath, err := os.Executable()
	if err != nil {
		logger.Error("error getting executable path: %v", err)
	}

	// Load embedded files in the executable.
	hasEmbed := true
	fs, err := stuffbin.UnStuff(execPath)
	if err != nil {
		hasEmbed = false
		// Running in local mode. Load local assets into
		// the in-memory stuffbin.FileSystem.
		logger.Info("unable to initialize embedded filesystem (%v). Using local filesystem", err)
		fs, err = stuffbin.NewLocalFS("/")
		if err != nil {
			logger.Error("failed to initialize local file for assets: %v", err)
		}
	}

	// If the embed failed, load app and frontend files from the compile-time paths.
	files := []string{}
	if !hasEmbed {
		files = append(files, joinFSPaths(appDir, appFiles)...)
		if frontendDir != "" {
			files = append(files, joinFSPaths(frontendDir, frontendFiles)...)
		}
	}

	// No additional files to load.
	if len(files) == 0 {
		return fs
	}

	// Load files from disk and overlay into the FS.
	fStatic, err := stuffbin.NewLocalFS("/", files...)
	if err != nil {
		logger.Error("failed reading static files from disk: '%s': %v", err)
	}

	if err := fs.Merge(fStatic); err != nil {
		logger.Error("error merging static files: '%s': %v", err)
	}

	return fs
}
