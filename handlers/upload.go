package handlers

import (
	"crypto/aes"
	"crypto/cipher"
	hmac "crypto/hmac"
	rand "crypto/rand"
	hsh "crypto/sha256"
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
	hmacKey := make([]byte, 32)
	n, err := io.ReadFull(rand.Reader, hmacKey)
	if n != len(hmacKey) || err != nil {
		panic(err)
	}
	hashType := hmac.New(hsh.New, hmacKey)
	hashWrapper := newNopWriter(hashType)

	iv := make([]byte, aes.BlockSize)
	n, err = io.ReadFull(rand.Reader, iv)
	if n != len(iv) || err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv[:])

	//outFile, err := os.OpenFile("/tmp/"+obfuscatedName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	outFile, err := os.OpenFile(tmpFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	multiWriter := newMultiWriterCloser(outFile, hashWrapper)

	writer := &cipher.StreamWriter{S: stream, W: multiWriter}
	defer writer.Close()
	// Copy the input file to the output file, encrypting as we go.
	_, err = io.Copy(writer, in)
	if err != nil {
		panic(err)
	}
	//We need to flush and close before we can read  back
	writer.Close()
	multiWriter.Close()

	retFile, errz := os.Open(tmpFileName)
	if errz != nil {
		panic(errz)
	}

	computedHash := make([]byte, hashType.Size())
	computedHash = hashType.Sum(computedHash)
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
	metadata["hmacSha256"] = string(computedHash)
	metadata["hmacKey"] = string(hmacKey)
	metadata["iv"] = string(iv)

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
		store.Put(path, "hmacSha256", metadata["hmacSha256"])
		store.Put(path, "hmacKey", metadata["hmacKey"])
		store.Put(path, "iv", metadata["iv"])

		fmt.Fprint(w, "Uploaded "+path+" for "+backend)

	default:
		http.Error(w, "Unknown error", 500)
	}
}
