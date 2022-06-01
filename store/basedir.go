package store

import (
	"os"
	"path/filepath"
)

func exeDir() string {
	exe, err := os.Executable()
	if err != nil {
		return "."
	}
	return filepath.Dir(exe)
}

func BaseDir(addr string) string {
	return filepath.Join(exeDir(), addr)
}
