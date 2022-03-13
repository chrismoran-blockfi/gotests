package input

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"

	"github.com/cweill/gotests/internal/models"
)

// Files returns all the Golang files for the given path. Ignores hidden files. Conditionally ignores generated files.
func Files(srcPath string, ignore *regexp.Regexp) ([]models.Path, error) {
	srcPath, err := filepath.Abs(srcPath)
	if err != nil {
		return nil, fmt.Errorf("filepath.Abs: %v\n", err)
	}
	var fi os.FileInfo
	if fi, err = os.Stat(srcPath); err != nil {
		return nil, fmt.Errorf("os.Stat: %v\n", err)
	}
	if fi.IsDir() {
		return dirFiles(srcPath, ignore)
	}
	return file(srcPath, ignore)
}

func dirFiles(srcPath string, ignore *regexp.Regexp) ([]models.Path, error) {
	ps, err := filepath.Glob(path.Join(srcPath, "*.go"))
	if err != nil {
		return nil, fmt.Errorf("filepath.Glob: %v\n", err)
	}
	var srcPaths []models.Path
	for _, p := range ps {
		src := models.Path(p)
		if isHiddenFile(p) || src.IsTestPath() || (isIgnored(src.FilePart(), ignore)) {
			continue
		}
		srcPaths = append(srcPaths, src)
	}
	return srcPaths, nil
}

func isIgnored(path string, excl *regexp.Regexp) bool {
	return excl != nil && excl.MatchString(path)
}

func file(srcPath string, ignore *regexp.Regexp) ([]models.Path, error) {
	src := models.Path(srcPath)
	if isIgnored(src.FilePart(), ignore) {
		return nil, nil
	}
	if filepath.Ext(srcPath) != ".go" || isHiddenFile(srcPath) {
		return nil, fmt.Errorf("no Go source files found at %v", srcPath)
	}
	return []models.Path{src}, nil
}

func isHiddenFile(path string) bool {
	return []rune(filepath.Base(path))[0] == '.'
}
