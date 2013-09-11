package handlers

import (
	"../store"
	"./mux"
	"fmt"
	"net/http"
)

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	store := store.EnsureStore()
	path, err := mux.DeletePath(r)
	fmt.Println("DELETE PATH = " + path)
	if err != nil {
		panic(err)
	}

	obfName, err := store.Get(path, "obfuscatedName")
	fmt.Println("DELETE OBFNAME = " + obfName)
	backend, err := store.Get(path, "backend")
	storageBackend := GetBackend(backend)

	err = storageBackend.DeleteObject(obfName)

	if err != nil {
		http.Error(w, "Failed to delete file", 404)
	} else {
		fmt.Fprint(w, "ok")
	}

	return

}
