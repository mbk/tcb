package handlers

import (
	"github.com/mbk/tcb/handlers/mux"
	"github.com/mbk/tcb/store"
	"net/http"
)

func ExistsHandler(w http.ResponseWriter, r *http.Request) {
	store := store.EnsureStore()
	path, err := mux.ExistsPath(r)
	if err != nil {
		panic(err)
	}

	if store.Exists(path, "length") {
		//Note that this is an indication as it is compressed/encrypted
		//Also, the data is returned with chunked encoding
		length, _ := store.Get(path, "length")
		w.Header().Add("Content-Length", length)
	} else {
		http.Error(w, "Failed to download file, probably doesn't exist.", 500)
	}
	return
}
