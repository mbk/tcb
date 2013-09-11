package handlers

import (
	"../store"
	"./mux"
	"fmt"
	"net/http"
)

var myStore store.Store = store.EnsureStore()

func MetadataSetHandler(w http.ResponseWriter, r *http.Request) {

	path, key, value, err := mux.SetMetadataPath(r)
	if err != nil {
		panic(err)
	}
	myStore.Put(path, key, value)
	fmt.Fprint(w, "POSTed metadata for "+path)
	return
}

func MetadataGetHandler(w http.ResponseWriter, r *http.Request) {

	path, key, err := mux.GetMetadataPath(r)
	if err != nil {
		panic(err)
	}
	value, err := myStore.Get(path, key)
	if err != nil {
		fmt.Fprint(w, "GET metadata failed for "+path+" with key "+key+" with a timeout.")
	} else {
		if value == "" {
			fmt.Fprint(w, "GET metadata failed for "+path+" with key "+key+" - not found.")
		} else {
			fmt.Fprint(w, value)
		}
	}
	return
}

func MetadataDeleteHandler(w http.ResponseWriter, r *http.Request) {

	path, key, err := mux.DeleteMetadataPath(r)
	if err != nil {
		panic(err)
	}
	err = myStore.Delete(path, key)
	if err != nil {
		fmt.Fprint(w, "false")
	} else {
		fmt.Fprint(w, "true")
	}
	return
}

func MetadataExistsHandler(w http.ResponseWriter, r *http.Request) {
	path, key, err := mux.ExistsMetadataPath(r)
	if err != nil {
		panic(err)
	}

	if myStore.Exists(path, key) {
		fmt.Fprint(w, "")
	} else {
		http.Error(w, "Not found", 404)
	}
	return
}

func MetadataHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT", "POST":
		MetadataSetHandler(w, r)
	case "GET":
		MetadataGetHandler(w, r)
	case "HEAD":
		MetadataExistsHandler(w, r)
	case "DELETE":
		MetadataDeleteHandler(w, r)
	}

}
