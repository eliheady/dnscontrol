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
		size:    14101,
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
