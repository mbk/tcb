package handlers

import (
	"crypto/aes"
	"crypto/cipher"
	"github.com/mbk/tcb/handlers/mux"
	"github.com/mbk/tcb/store"
	"io"
	"net/http"
)

func handleDownload(to io.Writer, name string) {

	store := store.EnsureStore()
	key, err := store.Get(name, "encr")
	obname, err := store.Get(name, "obfuscatedName")
	backend, err := store.Get(name, "backend")

	//Switch on backend type here to call GetReader
	storageBackend := GetBackend(backend)
	inFile, err := storageBackend.GetReader(obname)
	defer inFile.Close()
	if err != nil {
		panic(err)
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}

	// If the key is unique for each ciphertext, then it's ok to use a zero
	// IV.
	var iv [aes.BlockSize]byte
	stream := cipher.NewOFB(block, iv[:])

	reader := &cipher.StreamReader{S: stream, R: inFile}
	// Copy the input file to the output file, decrypting as we go.
	if _, err := io.Copy(to, reader); err != nil {
		panic(err)
	}
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {

	defer func() {
		if err := recover(); err != nil {
			e := err.(error)
			//http.Error(w, "Failed to download file, probably doesn't exist.", 500)
			http.Error(w, e.Error(), 500)
		}
	}()
	//We will use this later on to get the filename etc.
	path, err := mux.DownloadPath(r)
	if err != nil {
		panic(err)
	}

	switch r.Method {
	case "GET":
		handleDownload(w, path)
	default:
		http.Error(w, "Unknown HTTP verb for download, error", 500)
	}
}
