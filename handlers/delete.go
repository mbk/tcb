package handlers

import (
	"fmt"
	"github.com/mbk/tcb/handlers/mux"
	"github.com/mbk/tcb/store"
	"net/http"
)

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	store := store.EnsureStore()
	path, err := mux.DeletePath(r)
	if err != nil {
		panic(err)
	}

	obfName, err := store.Get(path, "obfuscatedName")
	backend, err := store.Get(path, "backend")
	storageBackend := GetBackend(backend)

	err = storageBackend.DeleteObject(obfName)
	err_metadelete := store.DeleteName(path)

	if err != nil || err_metadelete != nil {
		http.Error(w, "Failed to delete file", 404)
	} else {
		fmt.Fprint(w, "ok")
	}

	return

}
