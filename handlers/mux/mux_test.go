package mux

import (
	mux "."
	config "../../config"
	http "net/http"
	url "net/url"
	test "testing"
)

var r = new(http.Request)

func TestDownloadPrefix(t *test.T) {

	u := new(url.URL)
	u.Path = "/download/some/path/as/file"
	r.URL = u
	path, err := mux.DownloadPath(r)
	t.Log("path = " + path)
	if err != nil || path != "/some/path/as/file" {
		t.Log("Download path not extracted correctly")
		t.Fail()
	}
}

func TestUploadPrefix(t *test.T) {

	u := new(url.URL)
	u.Path = "/upload/some/path/as/file/to/s3"
	r.URL = u
	path, backend, err := mux.UploadPath(r)
	t.Log("path = " + path)
	t.Log("backend = " + backend)
	if err != nil || path != "/some/path/as/file" || backend != "s3" {
		t.Log("Upload path or backend not extracted correctly")
		t.Fail()
	}
}

func TestDeletePrefix(t *test.T) {

	u := new(url.URL)
	u.Path = "/delete/some/path/as/file"
	r.URL = u
	path, err := mux.DeletePath(r)
	t.Log("path = " + path)
	if err != nil || path != "/some/path/as/file" {
		t.Log("Delete path not extracted correctly")
		t.Fail()
	}
}

func TestUploadPrefixDefaultBackend(t *test.T) {
	config.EnsureConfiguration("../../tcb.ini")
	default_backend := config.GetOrElse("default_cloud_storage_type", "local")
	u := new(url.URL)
	u.Path = "/upload/some/path/as/file"
	r.URL = u
	path, backend, err := mux.UploadPath(r)
	t.Log("path = " + path)
	t.Log("backend = " + backend)
	if err != nil || path != "/some/path/as/file" || backend != default_backend {
		t.Log("Upload path or backend not extracted correctly")
		t.Fail()
	}
}

func TestExistsPrefix(t *test.T) {

	u := new(url.URL)
	u.Path = "/exists/some/path/as/file"
	r.URL = u
	path, err := mux.ExistsPath(r)
	t.Log("path = " + path)
	if err != nil || path != "/some/path/as/file" {
		t.Log("Download path not extracted correctly")
		t.Fail()
	}
}

func TestSetMetadataPath(t *test.T) {
	u := new(url.URL)
	u.Path = "/metadata/some/path/as/file/key/sleutel/value/waarde"
	r.URL = u
	path, key, value, err := mux.SetMetadataPath(r)
	t.Log("path = " + path)
	t.Log("key = " + key)
	t.Log("value = " + value)
	if err != nil || path != "/some/path/as/file" || key != "sleutel" || value != "waarde" {
		t.Log("Set metadata path failed")
		t.Fail()
	}
}

func TestGetMetadataPath(t *test.T) {
	u := new(url.URL)
	u.Path = "/metadata/some/path/as/file/key/sleutel"
	r.URL = u
	path, key, err := mux.GetMetadataPath(r)
	t.Log("path = " + path)
	t.Log("key = " + key)
	if err != nil || path != "/some/path/as/file" || key != "sleutel" {
		t.Log("Set metadata path failed")
		t.Fail()
	}
}

func TestExistsMetadataPath(t *test.T) {
	u := new(url.URL)
	u.Path = "/metadata/some/path/as/file/key/sleutel"
	r.URL = u
	path, key, err := mux.ExistsMetadataPath(r)
	t.Log("path = " + path)
	t.Log("key = " + key)
	if err != nil || path != "/some/path/as/file" || key != "sleutel" {
		t.Log("Set metadata path failed")
		t.Fail()
	}
}

func TestDeleteMetadataPath(t *test.T) {
	u := new(url.URL)
	u.Path = "/metadata/some/path/as/file/key/sleutel"
	r.URL = u
	path, key, err := mux.DeleteMetadataPath(r)
	t.Log("path = " + path)
	t.Log("key = " + key)
	if err != nil || path != "/some/path/as/file" || key != "sleutel" {
		t.Log("Set metadata path failed")
		t.Fail()
	}
}
