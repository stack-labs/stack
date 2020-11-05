package file

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/stack-labs/stack-rpc/config/storage"
)

type file struct {
	file string
}

func NewStorage(f string) storage.Storage {
	return &file{file: f}
}

func (f *file) Write(content []byte) (err error) {
	if err := os.MkdirAll(filepath.Dir(f.file), 0755); err != nil {
		return err
	}
	// create backup file
	exist, err := existFile(f.file)
	if err != nil {
		return err
	}
	if exist {
		if err := copyFile(f.file, fmt.Sprintf("%s_backup", f.file)); err != nil {
			return err
		}
	}

	tmpFile := fmt.Sprintf("%s_tmp", f.file)
	if err := ioutil.WriteFile(tmpFile, content, 0755); err != nil {
		return err
	}

	if err := os.Rename(tmpFile, f.file); err != nil {
		return err
	}

	return nil
}

func (f *file) Load() (config []byte, err error) {
	return ioutil.ReadFile(f.file)
}

func copyFile(source string, dest string) error {
	src, err := os.Open(source)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	info, err := os.Stat(source)
	if err != nil {
		err = os.Chmod(dest, info.Mode())
		if err != nil {
			return err
		}
	}

	return nil
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
