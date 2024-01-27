//go:build exclude

package main

import (
	"time"
)

type IFile interface {
	Readdirnames(n int) (names []string, err error)
	Seek(offset int64, whence int) (ret int64, err error)
	Read(b []byte) (n int, err error)
	Close()
}

type IFileInfo[FileMode IFileMode] interface {
	IsDir() bool
	ModTime() time.Time
	Mode() FileMode
	Name() string
	Size() int64
}

type IFileMode uint32

type WalkFunc[FileInfo IFileInfo[FileMode], FileMode IFileMode] func(path string, info FileInfo, err error) error

type AbstractFileSystem[File IFile, FileMode IFileMode, FileInfo IFileInfo[FileMode]] interface {
	Mkdir(name string, perm FileMode) error
	Create(name string) (*IFile, error)
	Open(name string) (*IFile, error)
	ReadFile(name string) ([]byte, error)
	WriteFile(name string, data []byte, perm FileMode) error
	Stat(name string) (FileInfo, error)
	IsNotExist(err error) bool

	Walk(root string, fn WalkFunc[FileInfo, FileMode]) error
	Glob(pattern string) (matches []string, err error)
}
