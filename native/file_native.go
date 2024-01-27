//go:build exclude

package main

import (
	"os"
	"path/filepath"
)

type OsFS[File os.File, FileMode os.FileMode, FileInfo os.FileInfo] struct {
}

func (osFileSystem *OsFS[File, FileMode, FileInfo]) Mkdir(name string, perm os.FileMode) error {
	return os.Mkdir(name, perm)
}
func (osFileSystem *OsFS[File, FileMode, FileInfo]) Create(name string) (*os.File, error) {
	return os.Create(name)
}
func (osFileSystem *OsFS[File, FileMode, FileInfo]) Open(name string) (*os.File, error) {
	return os.Open(name)
}
func (osFileSystem *OsFS[File, FileMode, FileInfo]) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}
func (osFileSystem *OsFS[File, FileMode, FileInfo]) WriteFile(name string, data []byte, perm os.FileMode) error {
	return os.WriteFile(name, data, perm)
}
func (osFileSystem *OsFS[File, FileMode, FileInfo]) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}
func (osFileSystem *OsFS[File, FileMode, FileInfo]) IsNotExist(err error) bool {
	return os.IsNotExist(err)
}

func (osFileSystem *OsFS[File, FileMode, FileInfo]) Walk(root string, fn WalkFunc[os.FileInfo, os.FileMode]) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		return fn(path, info, err)
	})
}
func (osFileSystem *OsFS[File, FileMode, FileInfo]) Glob(pattern string) (matches []string, err error) {
	return filepath.Glob(pattern)
}
