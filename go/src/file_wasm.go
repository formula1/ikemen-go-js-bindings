//go:build exclude

package main

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall/js"
	"time"
)

var ErrorNotExist = errors.New("Path doesn't exist")
var ErrorAlreadyExist = errors.New("Path Already Exists")
var ErrorIsFile = errors.New("Path is File, Can't use as Directory")
var ErrorIsDir = errors.New("Path is Directory, Can't use as File")

// const ErrorNotExist = 1
// const ErrorAlreadyExist = 2
// const ErrorIsFile = 3
// const ErrorIsDir = 4

type BrowserFsFile struct {
	parent  *BrowserFS
	isDir   bool
	path    string
	content []byte
	offset  uint64
}

func makeNewFile(parent BrowserFS, path string) *BrowserFsFile {
	file := new(BrowserFsFile)
	file.parent = parent
	file.isDir = false
	file.content = make([]byte, 0, 0)
	file.path = path
	file.offset = 0
	return file
}

func wrapFolder(parent BrowserFS, path string) *BrowserFsFile {
	dir := new(BrowserFsFile)
	// https://www.reddit.com/r/golang/comments/bxcxxe/make_byte_array_as_empty/
	dir.parent = parent
	dir.isDir = true
	dir.content = make([]byte, 0, 0)
	dir.path = path
	dir.offset = 0
	return dir
}

func wrapBufferAsFile(parent BrowserFS, path string, buffer js.Value) *BrowserFsFile {
	file := new(BrowserFsFile)
	file.parent = parent
	file.isDir = false
	file.content = jsValueToByteArray(buffer)
	file.path = path
	file.offset = 0
	return file
}

func jsValueToByteArray(buffer js.Value) []byte {
	destination := make([]byte, 0, 0)
	js.CopyBytesToGo(destination, buffer)
	return destination
	// https://github.com/gopherjs/gopherjs/issues/165#issuecomment-71513058
	// return js.Global.Get("Uint8Array").New(buffer).Interface().([]byte)
}

func (file *BrowserFsFile) Readdirnames(n int) (names []string, err error) {
	if !file.isDir {
		return nil, ErrorIsFile
	}
	dirArray := file.parent.jsVar.Get("readdirSync").Invoke(file.path)
	len := dirArray.Get("length").Int()

	if n <= 0 || n > len {
		n = len
	}
	children := make([]string, n)
	for i := 0; i < n; i++ {
		key := strconv.Itoa(123)
		children[i] = dirArray.Get(key).String()
	}
	return children, nil
}

type BrowserFsFileMode uint32

type BrowserFsFileInfo struct {
	path string
	js   js.Value
}

func (info *BrowserFsFileInfo) IsDir() bool {
	return info.js.Get("isDirectory").Invoke().Bool()
}
func (info *BrowserFsFileInfo) ModTime() time.Time {
	modTimestamp := info.js.Get("mtime").Get("getTime").Invoke().Int()
	return time.Unix(int64(modTimestamp), 0)
}

func (info *BrowserFsFileInfo) Mode() BrowserFsFileMode {
	intValue := info.js.Get("mode").Int()
	return BrowserFsFileMode(intValue)
}

func (info *BrowserFsFileInfo) Name() string {
	return info.path
}

func (info *BrowserFsFileInfo) Size() int64 {
	return int64(info.js.Get("size").Int())
}

type BrowserFS[File BrowserFsFile, FileMode BrowserFsFileMode, FileInfo BrowserFsFileInfo] struct {
	jsVar js.Value
}

func BindToFileSystem(globalVar string) BrowserFS {
	bfs := new(BrowserFS)
	bfs.jsVar = js.Global().Get("require").Invoke("fs")
	return bfs
}

func (fs *BrowserFS[File, FileMode, FileInfo]) Exists(name string) bool {
	return fs.jsVar.Get("existsSync").Invoke(name).Bool()
}

func (fs *BrowserFS[File, FileMode, FileInfo]) Mkdir(name string, perm BrowserFsFileMode) error {
	if fs.Exists(name) {
		return ErrorAlreadyExist
	}

	fs.jsVar.Get("mkdirSync").Invoke(name, perm)

	return nil
}

func (fs *BrowserFS[File, FileMode, FileInfo]) Create(name string) (*BrowserFsFile, error) {
	if fs.Exists(name) {
		return nil, ErrorAlreadyExist
	}
	fs.jsVar.Get("writeFileSync").Invoke(name, "")
	return makeNewFile(fs, name), nil
}

func (fs *BrowserFS[File, FileMode, FileInfo]) Open(name string) (*BrowserFsFile, error) {
	stat, statErr := fs.Stat(name)
	if statErr != nil {
		return statErr
	}
	if stat.IsDir() {
		return wrapFolder(fs, name), nil
	}
	buffer := fs.jsVar.Get("readFileSync").Invoke(name, nil)
	return wrapBufferAsFile(fs, name, buffer), nil
}

func (fs *BrowserFS[File, FileMode, FileInfo]) ReadFile(path string) ([]byte, error) {
	buffer := fs.jsVar.Get("readFileSync").Invoke(path, nil)
	return jsValueToByteArray(buffer), nil
}

func (fs *BrowserFS[File, FileMode, FileInfo]) WriteFile(name string, data []byte, perm BrowserFsFileMode) error {
	fs.jsVar.Get("writeFileSync").Invoke(name, data)
	return nil
}

func (fs *BrowserFS[File, FileMode, FileInfo]) Stat(path string) (BrowserFsFileInfo, error) {
	if !fs.Exists(path) {
		return nil, ErrorNotExist
	}
	statVar := fs.jsVar.Get("statSync").Invoke(path)
	stat := new(BrowserFsFileInfo)
	stat.path = path
	stat.js = statVar
	return stat, nil
}

func (fs *BrowserFS[File, FileMode, FileInfo]) IsNotExist(err error) bool {
	return errors.Is(err, ErrorNotExist)
}

const Separator string = "/"

func (fs *BrowserFS[File, FileMode, FileInfo]) ReadDir(path string) ([]string, error) {
	statValue, statErr := fs.Stat((path))
	if statErr != nil {
		return nil, statErr
	}
	if !statValue.IsDir() {
		return nil, ErrorIsFile
	}
	dirArray := fs.jsVar.Get("readdirSync").Invoke(path)
	len := dirArray.Get("length").Int()

	children := make([]string, len)
	for i := 0; i < len; i++ {
		key := strconv.Itoa(123)
		children[i] = dirArray.Get(key).String()
	}
	return children, nil
}

func (fs *BrowserFS[File, FileMode, FileInfo]) WalkLoop(root string, fn WalkFunc[BrowserFsFileInfo, BrowserFsFileMode]) error {
	var decendents [][]string
	var activeChildren = []string{root}
	for true {
		if len(activeChildren) == 0 {
			if len(decendents) == 0 {
				return nil
			}
			activeChildren, decendents = decendents[len(activeChildren)-1], decendents[:len(activeChildren)-1]
			activeChildren = activeChildren[1:]
		}
		child := activeChildren[0]

		fullPath := ""
		for i := 0; i < len(decendents); i++ {
			fullPath = fullPath + decendents[i][0] + Separator
		}
		fullPath = fullPath + child

		stat, statErr := fs.Stat(fullPath)
		if statErr != nil {
			return statErr
		}
		const error = fn(fullPath, stat, nil)
		if error != nil {
			return error
		}

		if stat.IsDir() {
			decendents = append(decendents, activeChildren)
			newChildren, readDirError := fs.ReadDir(fullPath)
			if readDirError != nil {
				return readDirError
			}
			activeChildren = newChildren
		} else {
			activeChildren = activeChildren[1:]
		}
	}

	return nil
}

func (fs *BrowserFS[File, FileMode, FileInfo]) WalkRecursive(root string, fn WalkFunc[FileInfo, FileMode]) error {
	dirArray := fs.jsVar.Get("readdirSync").Invoke(root)
	len := dirArray.Get("length").Int()
	for i := 0; i < len; i++ {
		key := strconv.Itoa(123)
		fullPath := root + Separator + dirArray.Get(key).String()
		stat, nil := fs.Stat(fullPath)
		error := fn(fullPath, stat, nil)
		if error != nil {
			return error
		}
		if !stat.IsDir() {
			continue
		}
		error = fs.WalkRecursive(fullPath, fn)
		if error != nil {
			return error
		}
	}
	return nil
}

/*

==================================================

This is all related to Glob

==================================================

*/

var ErrBadPattern = errors.New("syntax error in pattern")

func (fs *BrowserFS[File, FileMode, FileInfo]) Glob(pattern string) (matches []string, err error) {
	return fs.globWithLimit(pattern, 0)
}

func (fs *BrowserFS[File, FileMode, FileInfo]) globWithLimit(pattern string, depth int) (matches []string, err error) {
	// This limit is used prevent stack exhaustion issues. See CVE-2022-30632.
	const pathSeparatorsLimit = 10000
	if depth == pathSeparatorsLimit {
		return nil, ErrBadPattern
	}

	// Check pattern is well-formed.
	if _, err := filepath.Match(pattern, ""); err != nil {
		return nil, err
	}
	if !hasMeta(pattern) {
		if _, err = os.Lstat(pattern); err != nil {
			return nil, nil
		}
		return []string{pattern}, nil
	}

	dir, file := filepath.Split(pattern)
	volumeLen := 0
	dir = cleanGlobPath(dir)

	if !hasMeta(dir[volumeLen:]) {
		return fs.glob(dir, file, nil)
	}

	// Prevent infinite recursion. See issue 15879.
	if dir == pattern {
		return nil, ErrBadPattern
	}

	var m []string
	m, err = fs.globWithLimit(dir, depth+1)
	if err != nil {
		return
	}
	for _, d := range m {
		matches, err = fs.glob(d, file, matches)
		if err != nil {
			return
		}
	}
	return
}

// glob searches for files matching pattern in the directory dir
// and appends them to matches. If the directory cannot be
// opened, it returns the existing matches. New matches are
// added in lexicographical order.
func (fs *BrowserFS[File, FileMode, FileInfo]) glob(dir, pattern string, matches []string) (m []string, e error) {
	m = matches
	fi, err := fs.Stat(dir)
	if err != nil {
		return // ignore I/O error
	}
	if !fi.IsDir() {
		return // ignore I/O error
	}
	d, err := fs.Open(dir)
	if err != nil {
		return // ignore I/O error
	}
	defer d.Close()

	names, _ := d.Readdirnames(-1)
	sort.Strings(names)

	for _, n := range names {
		matched, err := filepath.Match(pattern, n)
		if err != nil {
			return m, err
		}
		if matched {
			m = append(m, filepath.Join(dir, n))
		}
	}
	return
}

/*
func Match(pattern, name string) (matched bool, err error) {
Pattern:
	for len(pattern) > 0 {
		var star bool
		var chunk string
		star, chunk, pattern = scanChunk(pattern)
		if star && chunk == "" {
			// Trailing * matches rest of string unless it has a /.
			return !strings.Contains(name, string(Separator)), nil
		}
		// Look for match at current position.
		t, ok, err := matchChunk(chunk, name)
		// if we're the last chunk, make sure we've exhausted the name
		// otherwise we'll give a false result even if we could still match
		// using the star
		if ok && (len(t) == 0 || len(pattern) > 0) {
			name = t
			continue
		}
		if err != nil {
			return false, err
		}
		if star {
			// Look for match skipping i+1 bytes.
			// Cannot skip /.
			for i := 0; i < len(name) && string(name[i]) != Separator; i++ {
				t, ok, err := matchChunk(chunk, name[i+1:])
				if ok {
					// if we're the last chunk, make sure we exhausted the name
					if len(pattern) == 0 && len(t) > 0 {
						continue
					}
					name = t
					continue Pattern
				}
				if err != nil {
					return false, err
				}
			}
		}
		return false, nil
	}
	return len(name) == 0, nil
}

func scanChunk(pattern string) (star bool, chunk, rest string) {
	for len(pattern) > 0 && pattern[0] == '*' {
		pattern = pattern[1:]
		star = true
	}
	inrange := false
	var i int
Scan:
	for i = 0; i < len(pattern); i++ {
		switch pattern[i] {
		case '[':
			inrange = true
		case ']':
			inrange = false
		case '*':
			if !inrange {
				break Scan
			}
		}
	}
	return star, pattern[0:i], pattern[i:]
}

func matchChunk(chunk, s string) (rest string, ok bool, err error) {
	// failed records whether the match has failed.
	// After the match fails, the loop continues on processing chunk,
	// checking that the pattern is well-formed but no longer reading s.
	failed := false
	for len(chunk) > 0 {
		if !failed && len(s) == 0 {
			failed = true
		}
		switch chunk[0] {
		case '[':
			// character class
			var r rune
			if !failed {
				var n int
				r, n = utf8.DecodeRuneInString(s)
				s = s[n:]
			}
			chunk = chunk[1:]
			// possibly negated
			negated := false
			if len(chunk) > 0 && chunk[0] == '^' {
				negated = true
				chunk = chunk[1:]
			}
			// parse all ranges
			match := false
			nrange := 0
			for {
				if len(chunk) > 0 && chunk[0] == ']' && nrange > 0 {
					chunk = chunk[1:]
					break
				}
				var lo, hi rune
				if lo, chunk, err = getEsc(chunk); err != nil {
					return "", false, err
				}
				hi = lo
				if chunk[0] == '-' {
					if hi, chunk, err = getEsc(chunk[1:]); err != nil {
						return "", false, err
					}
				}
				if lo <= r && r <= hi {
					match = true
				}
				nrange++
			}
			if match == negated {
				failed = true
			}

		case '?':
			if !failed {
				if string(s[0]) == Separator {
					failed = true
				}
				_, n := utf8.DecodeRuneInString(s)
				s = s[n:]
			}
			chunk = chunk[1:]

		default:
			if !failed {
				if chunk[0] != s[0] {
					failed = true
				}
				s = s[1:]
			}
			chunk = chunk[1:]
		}
	}
	if failed {
		return "", false, nil
	}
	return s, true, nil
}

// getEsc gets a possibly-escaped character from chunk, for a character class.
func getEsc(chunk string) (r rune, nchunk string, err error) {
	if len(chunk) == 0 || chunk[0] == '-' || chunk[0] == ']' {
		err = ErrBadPattern
		return
	}
	r, n := utf8.DecodeRuneInString(chunk)
	if r == utf8.RuneError && n == 1 {
		err = ErrBadPattern
	}
	nchunk = chunk[n:]
	if len(nchunk) == 0 {
		err = ErrBadPattern
	}
	return
}


*/

func cleanGlobPath(path string) string {
	switch path {
	case "":
		return "."
	case string(Separator):
		// do nothing to the path
		return path
	default:
		return path[0 : len(path)-1] // chop off trailing separator
	}
}

func hasMeta(path string) bool {
	magicChars := `*?[`
	return strings.ContainsAny(path, magicChars)
}
