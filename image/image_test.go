package image

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestListFiles(t *testing.T) {
	imageFilename := filepath.Join("testdata", "esp.img")
	files, err := ListFiles(imageFilename)
	require.Nil(t, err)
	require.IsType(t, []string{}, files)
	for _, file := range files {
		fmt.Println(file)
	}
}

func TestExtractFiles(t *testing.T) {
	imageFilename := filepath.Join("testdata", "esp.img")
	targetDir := filepath.Join("testdata", "files")
	stat, err := os.Stat(targetDir)
	if err == nil {
		if stat.IsDir() {
			err := os.RemoveAll(targetDir)
			require.Nil(t, err)
		}
	}
	err = ExtractFiles(imageFilename, targetDir)
	require.Nil(t, err)
}
