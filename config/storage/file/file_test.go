package file

import (
	"bytes"
	"testing"
)

func TestFile_Write_Load(t *testing.T) {
	data := []byte("123456")
	storage := NewStorage("/tmp/test.file")
	if err := storage.Write(data); err != nil {
		t.Fatal(err)
	}

	d, err := storage.Load()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(data, d) {
		t.Fatal(data, d)
	}
}
