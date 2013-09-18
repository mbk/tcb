package handlers

import (
	"crypto/aes"
	"crypto/cipher"
	rand "crypto/rand"
	"fmt"
	mux "github.com/mbk/tcb/handlers/mux"
	"github.com/mbk/tcb/store"
	"github.com/satori/uuid"
	"io"
	"net/http"
	"os"
	"strconv"
)

func storeUploadTemp(in io.Reader) (map[string]string, *os.File, error) {
	key := (uuid.NewV4())
	obfuscatedName := (uuid.NewV4().String())

	block, err := aes.NewCipher(key[0:])
	if err != nil {
		panic(err)
	}

	tmpFileName := os.TempDir() + uuid.NewV4().String()

	// If the key is unique for each ciphertext, then it's ok to use a zero
	// IV.

	//var iv [aes.BlockSize]byte

	iv := make([]byte, aes.BlockSize)
	n, err := io.ReadFull(rand.Reader, iv)
	if n != len(iv) || err != nil {
		panic(err)
	}

	stream := cipher.NewOFB(block, iv[:])

	//outFile, err := os.OpenFile("/tmp/"+obfuscatedName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	outFile, err := os.OpenFile(tmpFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	writer := &cipher.StreamWriter{S: stream, W: outFile}
	defer writer.Close()
	// Copy the input file to the output file, encrypting as we go.
	_, err = io.Copy(writer, in)
	if err != nil {
		panic(err)
	}
	//We need to flush and close before we can read  back
	writer.Close()
	outFile.Close()

	retFile, errz := os.Open(tmpFileName)
	if errz != nil {
		panic(errz)
	}
	//Stat the temp file, so we have the real length
	stat, err := os.Stat(tmpFileName)
	if err != nil {
		panic(err)
	}
	length := stat.Size()
	metadata := make(map[string]string)
	metadata["encr"] = string(key[0:])
	metadata["length"] = strconv.FormatInt(length, 10)
	metadata["obfuscatedName"] = obfuscatedName

	return metadata, retFile, errz
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	m := r.Method

	path, backend, err := mux.UploadPath(r)

	if err != nil {
		panic(err)
	}

	defer func() {
		if err := recover(); err != nil {
			//http.Error(w, "Failed to store file", 500)
			e := err.(error)
			http.Error(w, e.Error(), 500)
		}
	}()

	switch m {

	case "POST", "PUT":

		metadata, tmpFile, err := storeUploadTemp(r.Body)
		tmpPath := tmpFile.Name()
		defer tmpFile.Close()
		//Delete the temp file
		defer os.Remove(tmpPath)
		if err != nil {
			panic(err)
		}
		store := store.EnsureStore()
		storageBackend := GetBackend(backend)
		errz := storageBackend.StoreObject(metadata["obfuscatedName"], tmpFile, metadata)
		if errz != nil {
			panic(errz)
		}

		store.Put(path, "encr", metadata["encr"])
		store.Put(path, "length", metadata["length"])
		store.Put(path, "obfuscatedName", metadata["obfuscatedName"])
		store.Put(path, "backend", backend)

		fmt.Fprint(w, "Uploaded "+path+" for "+backend)

	default:
		http.Error(w, "Unknown error", 500)
	}
}
