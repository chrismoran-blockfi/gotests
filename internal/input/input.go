package input

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/cweill/gotests/internal/models"
)

// Files returns all the Golang files for the given path. Ignores hidden files. Conditionally ignores generated files.
func Files(srcPath string, ignoreGenerated bool) ([]models.Path, error) {
	srcPath, err := filepath.Abs(srcPath)
	if err != nil {
		return nil, fmt.Errorf("filepath.Abs: %v\n", err)
	}
	var fi os.FileInfo
	if fi, err = os.Stat(srcPath); err != nil {
		return nil, fmt.Errorf("os.Stat: %v\n", err)
	}
	if fi.IsDir() {
		return dirFiles(srcPath, ignoreGenerated)
	}
	return file(srcPath, ignoreGenerated)
}

func dirFiles(srcPath string, ignoreGenerated bool) ([]models.Path, error) {
	ps, err := filepath.Glob(path.Join(srcPath, "*.go"))
	if err != nil {
		return nil, fmt.Errorf("filepath.Glob: %v\n", err)
	}
	var srcPaths []models.Path
	for _, p := range ps {
		src := models.Path(p)
		if isHiddenFile(p) || src.IsTestPath() || (ignoreGenerated && src.IsGenPath()) {
			continue
		}
		srcPaths = append(srcPaths, src)
	}
	return srcPaths, nil
}

func file(srcPath string, ignoreGenerated bool) ([]models.Path, error) {
	src := models.Path(srcPath)
	if ignoreGenerated && src.IsGenPath() {
		return nil, fmt.Errorf("generated source files are ignored: %v", srcPath)
	}
	if filepath.Ext(srcPath) != ".go" || isHiddenFile(srcPath) {
		return nil, fmt.Errorf("no Go source files found at %v", srcPath)
	}
	return []models.Path{src}, nil
}

func isHiddenFile(path string) bool {
	return []rune(filepath.Base(path))[0] == '.'
}
