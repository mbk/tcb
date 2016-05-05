package handlers

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	hmac "crypto/hmac"
	hsh "crypto/sha256"
	"errors"
	"github.com/mbk/tcb/handlers/mux"
	"github.com/mbk/tcb/store"
	"io"
	"log"
	"net/http"
	"os"
)

var logger = log.New(os.Stdout, "******", log.LstdFlags)

func handleDownload(to io.Writer, name string) {

	store := store.EnsureStore()
	key, err := store.Get(name, "encr")
	obname, err := store.Get(name, "obfuscatedName")
	backend, err := store.Get(name, "backend")
	hmacKey, err := store.Get(name, "hmacKey")
	computedHash, err := store.Get(name, "hmacSha256")
	ivString, err := store.Get(name, "iv")

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
	hmacReadKey := []byte(hmacKey)
	hashType := hmac.New(hsh.New, hmacReadKey)
	hashWrapper := newNopWriter(hashType)

	// If the key is unique for each ciphertext, then it's ok to use a zero
	// IV.
	iv := []byte(ivString)
	stream := cipher.NewCFBDecrypter(block, iv)

	reader := &cipher.StreamReader{S: stream, R: inFile}
	toWrapped := newNopWriter(to)

	multiWriter := newMultiWriterCloser(toWrapped, hashWrapper)

	// Copy the input file to the output file, decrypting as we go.
	if _, err := io.Copy(multiWriter, reader); err != nil {
		panic(err)
	}

	originalHash := []byte(computedHash)
	readHash := make([]byte, (hashType.Size()))
	readHash = hashType.Sum(readHash)

	if !bytes.Equal(readHash, originalHash) {
		logger.Println("hashes do not match")
		panic(errors.New("Hashes do not match"))
		return 
	}
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {

	defer func() {
		if err := recover(); err != nil {
			//http.Error(w, "Failed to download file, probably doesn't exist.", 500)
			//This can be commented in in stead of the above for debuging purposes
			e := err.(error)
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
