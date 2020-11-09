package file

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/stack-labs/stack-rpc/pkg/config/storage"
)

const _perm = 0644

type file struct {
	file string
}

func NewStorage(f string) storage.Storage {
	return &file{file: f}
}

func (f *file) Exist() bool {
	exit, _ := existFile(f.file)
	return exit
}

func (f *file) FileName() string {
	return f.file
}

func (f *file) Write(content []byte) (err error) {
	if err := os.MkdirAll(filepath.Dir(f.file), _perm); err != nil {
		return err
	}

	if err := ioutil.WriteFile(f.file, content, _perm); err != nil {
		return err
	}

	return nil
}

func (f *file) Load() (config []byte, err error) {
	return ioutil.ReadFile(f.file)
}

func existFile(file string) (bool, error) {
	_, err := os.Stat(file)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
