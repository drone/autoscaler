// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

// +build ignore

package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

func main() {
	addr := ":3333"

	// serve templates with dummy data
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.FormValue("data")
		if path == "" {
			http.Error(w, "missing data parameter", 500)
			return
		}

		tmpl := r.FormValue("template")
		if path == "" {
			http.Error(w, "missing template parameter", 500)
			return
		}

		// read the json data from file.
		rawjson, err := ioutil.ReadFile(filepath.Join("testdata", path))
		if err != nil {
			http.Error(w, "cannot open json file", 500)
			return
		}

		// unmarshal the json data
		data := map[string]interface{}{}
		err = json.Unmarshal(rawjson, &data)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		// load the templates
		T := template.New("_").Funcs(funcMap)
		matches, _ := filepath.Glob("files/*.tmpl")
		for _, match := range matches {
			raw, _ := ioutil.ReadFile(match)
			base := filepath.Base(match)
			T = template.Must(
				T.New(base).Parse(string(raw)),
			)
		}

		// render the template
		w.Header().Set("Content-Type", "text/html")
		err = T.ExecuteTemplate(w, tmpl, data)
		if err != nil {
			log.Println(err)
		}
	})

	// serve static content.
	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(
				http.Dir("../static/files"),
			),
		),
	)

	log.Printf("listening at %s", addr)
	log.Fatalln(http.ListenAndServe(addr, nil))
}

// mirros the func map in template.go
var funcMap = map[string]interface{}{
	"substr": func(v string, i int) string {
		return v[0:i]
	},
	"timestamp": func(v float64) string {
		return time.Unix(int64(v), 0).UTC().Format("2006-01-02T15:04:05Z")
	},
}
