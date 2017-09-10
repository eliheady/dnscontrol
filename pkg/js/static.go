package js

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

type _escLocalFS struct{}

var _escLocal _escLocalFS

type _escStaticFS struct{}

var _escStatic _escStaticFS

type _escDirectory struct {
	fs   http.FileSystem
	name string
}

type _escFile struct {
	compressed string
	size       int64
	modtime    int64
	local      string
	isDir      bool

	once sync.Once
	data []byte
	name string
}

func (_escLocalFS) Open(name string) (http.File, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	return os.Open(f.local)
}

func (_escStaticFS) prepare(name string) (*_escFile, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	var err error
	f.once.Do(func() {
		f.name = path.Base(name)
		if f.size == 0 {
			return
		}
		var gr *gzip.Reader
		b64 := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(f.compressed))
		gr, err = gzip.NewReader(b64)
		if err != nil {
			return
		}
		f.data, err = ioutil.ReadAll(gr)
	})
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fs _escStaticFS) Open(name string) (http.File, error) {
	f, err := fs.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.File()
}

func (dir _escDirectory) Open(name string) (http.File, error) {
	return dir.fs.Open(dir.name + name)
}

func (f *_escFile) File() (http.File, error) {
	type httpFile struct {
		*bytes.Reader
		*_escFile
	}
	return &httpFile{
		Reader:   bytes.NewReader(f.data),
		_escFile: f,
	}, nil
}

func (f *_escFile) Close() error {
	return nil
}

func (f *_escFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func (f *_escFile) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f *_escFile) Name() string {
	return f.name
}

func (f *_escFile) Size() int64 {
	return f.size
}

func (f *_escFile) Mode() os.FileMode {
	return 0
}

func (f *_escFile) ModTime() time.Time {
	return time.Unix(f.modtime, 0)
}

func (f *_escFile) IsDir() bool {
	return f.isDir
}

func (f *_escFile) Sys() interface{} {
	return f
}

// _escFS returns a http.Filesystem for the embedded assets. If useLocal is true,
// the filesystem's contents are instead used.
func _escFS(useLocal bool) http.FileSystem {
	if useLocal {
		return _escLocal
	}
	return _escStatic
}

// _escDir returns a http.Filesystem for the embedded assets on a given prefix dir.
// If useLocal is true, the filesystem's contents are instead used.
func _escDir(useLocal bool, name string) http.FileSystem {
	if useLocal {
		return _escDirectory{fs: _escLocal, name: name}
	}
	return _escDirectory{fs: _escStatic, name: name}
}

// _escFSByte returns the named file from the embedded assets. If useLocal is
// true, the filesystem's contents are instead used.
func _escFSByte(useLocal bool, name string) ([]byte, error) {
	if useLocal {
		f, err := _escLocal.Open(name)
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(f)
		f.Close()
		return b, err
	}
	f, err := _escStatic.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.data, nil
}

// _escFSMustByte is the same as _escFSByte, but panics if name is not present.
func _escFSMustByte(useLocal bool, name string) []byte {
	b, err := _escFSByte(useLocal, name)
	if err != nil {
		panic(err)
	}
	return b
}

// _escFSString is the string version of _escFSByte.
func _escFSString(useLocal bool, name string) (string, error) {
	b, err := _escFSByte(useLocal, name)
	return string(b), err
}

// _escFSMustString is the string version of _escFSMustByte.
func _escFSMustString(useLocal bool, name string) string {
	return string(_escFSMustByte(useLocal, name))
}

var _escData = map[string]*_escFile{

	"/helpers.js": {
		local:   "pkg/js/helpers.js",
		size:    14054,
		modtime: 0,
		compressed: `
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
`,
	},

	"/": {
		isDir: true,
		local: "pkg/js",
	},
}
