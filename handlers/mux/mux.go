package mux

import (
	cfg "../../config"
	"errors"
	"net/http"
	str "strings"
)

func ForVerbs(verbs []string, handler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		var verbFound = false
		for _, verb := range verbs {
			if verb == r.Method {
				verbFound = true
				break
			}
		}

		if verbFound {
			handler(w, r)
		} else {
			http.Error(w, "Invalid HTTP verb", 500)
		}

	}
}

func DownloadPath(r *http.Request) (path string, err error) {
	if !str.HasPrefix(r.URL.Path, "/download/") {
		err = errors.New("Incorrect path prefix")
		return
	} else {
		path = r.URL.Path[9:]
	}
	return
}

func DeletePath(r *http.Request) (path string, err error) {
	if !str.HasPrefix(r.URL.Path, "/delete/") {
		err = errors.New("Incorrect path prefix")
		return
	} else {
		path = r.URL.Path[7:]
	}
	return
}

func UploadPath(r *http.Request) (path, backend string, err error) {

	if !str.HasPrefix(r.URL.Path, "/upload/") {
		err = errors.New("Incorrect path prefix")
		return
	} else {
		pathComponents := str.Split(r.URL.Path[8:], "/")
		length := len(pathComponents)
		if (length >= 3) && (pathComponents[length-2] == "to") {
			backend = pathComponents[length-1]
			path = "/" + str.Join(pathComponents[:(length-2)], "/")
		} else {
			backend = cfg.GetOrElse("default_cloud_storage_type", "local")
			path = "/" + str.Join(pathComponents, "/")
		}
	}
	return
}

func ExistsPath(r *http.Request) (path string, err error) {
	if !str.HasPrefix(r.URL.Path, "/exists/") {
		err = errors.New("Incorrect path prefix")
		return
	} else {
		path = r.URL.Path[7:]
	}
	return
}

func SetMetadataPath(r *http.Request) (path, key, value string, err error) {
	if !str.HasPrefix(r.URL.Path, "/metadata/") {
		err = errors.New("Incorrect path prefix")
		return
	} else {
		pathComponents := str.Split(r.URL.Path[10:], "/")
		length := len(pathComponents)
		if !(length >= 5) && (pathComponents[length-4] == "key") && (pathComponents[length-2] == "value") {
			err = errors.New("metadata does not have a key or value component")
		} else {
			key = pathComponents[length-3]
			value = pathComponents[length-1]
			path = "/" + str.Join(pathComponents[:(length-4)], "/")
		}
	}
	return
}

func GetMetadataPath(r *http.Request) (path, key string, err error) {
	if !str.HasPrefix(r.URL.Path, "/metadata/") {
		err = errors.New("Incorrect path prefix")
		return
	} else {
		pathComponents := str.Split(r.URL.Path[10:], "/")
		length := len(pathComponents)
		if !(length >= 3) && (pathComponents[length-2] == "key") {
			err = errors.New("metadata does not have a key component")
		} else {
			key = pathComponents[length-1]
			path = "/" + str.Join(pathComponents[:(length-2)], "/")
		}
	}
	return
}

func ExistsMetadataPath(r *http.Request) (path, key string, err error) {
	return GetMetadataPath(r)
}

func DeleteMetadataPath(r *http.Request) (path, key string, err error) {
	return GetMetadataPath(r)
}
