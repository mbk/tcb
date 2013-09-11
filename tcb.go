package main

import (
	"flag"
	"github.com/mbk/tcb/config"
	"github.com/mbk/tcb/handlers"
	"github.com/mbk/tcb/handlers/mux"
	"net/http"
	"strconv"
)

func main() {
	var port = flag.Int("port", 8080, "Specify the port to listen on")
	var config_file = flag.String("config", "./tcb.ini", "specify the location of the config file")
	var ssl = flag.Bool("usessl", false, "Enables ssl. Make sure to have the certificate and key specified in the config file.")
	flag.Parse()
	config.EnsureConfiguration(*config_file)

	http.HandleFunc("/exists/", mux.ForVerbs([]string{"HEAD"}, handlers.ExistsHandler))
	http.HandleFunc("/metadata/", mux.ForVerbs([]string{"POST", "PUT", "GET", "HEAD", "DELETE"}, handlers.MetadataHandler))
	http.HandleFunc("/upload/", mux.ForVerbs([]string{"POST", "PUT"}, handlers.UploadHandler))
	http.HandleFunc("/download/", mux.ForVerbs([]string{"GET"}, handlers.DownloadHandler))
	http.HandleFunc("/delete/", mux.ForVerbs([]string{"DELETE"}, handlers.DeleteHandler))

	if *ssl {
		certFile := config.GetOrElse("cert_file", "cert.pem")
		keyFile := config.GetOrElse("key_file", "key.pem")
		http.ListenAndServeTLS(":"+strconv.Itoa(*port), certFile, keyFile, nil)

	} else {
		http.ListenAndServe(":"+strconv.Itoa(*port), nil)
	}

}
