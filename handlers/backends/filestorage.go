package backends

import (
	"github.com/mbk/tcb/config"
	"github.com/mbk/tcb/store"
	"io"
	"os"
)

type FileStore struct {
	destDir string
}

func (store *FileStore) StoreObject(name string, source *os.File, path string, metadata store.Store) error {

	dest, err := os.OpenFile(store.destDir+"/"+name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(dest, source)
	return err
}

func (store *FileStore) GetReader(name string) (io.ReadCloser, error) {

	inFile, err := os.Open(store.destDir + "/" + name)
	if err != nil {
		panic(err)
	}

	return inFile, err

}

func (store *FileStore) DeleteObject(name string) error {
	return os.Remove(store.destDir + "/" + name)
}

func GetFileStore() *FileStore {
	store := new(FileStore)
	store.destDir = config.GetOrElse("file_store_path", "/tmp")
	return store

}
