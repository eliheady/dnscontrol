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
		size:    14003,
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
