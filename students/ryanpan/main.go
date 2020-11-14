package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"urlshort/students/ryanpan/urlshort"
)

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	path := flag.String("file", "", "file path for parsing path to url mappings")
	flag.Parse()

	extension := filepath.Ext(*path)
	fmt.Println(extension)
	data, err := ioutil.ReadFile(*path)
	checkErr(err)

	var handler http.HandlerFunc
	switch extension {
	case ".yml":
		handler, err = urlshort.YAMLHandler(data, mapHandler)
		break
	case ".json":
		handler, err = urlshort.JSONHandler(data, mapHandler)
		break
	default:
		fmt.Println("The filepath provided is not valid")
		os.Exit(1)
	}

	checkErr(err)

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", handler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func checkErr(err error) {
	if err != nil {
		os.Exit(1)
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
