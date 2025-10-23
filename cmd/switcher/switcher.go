package main

import (
	"net/http"
	"strconv"

	muxswitcher "github.com/AlexChukhleb/muxswitcher"
)

func main() {
	sw := muxswitcher.New(nil, nil)
	if err := sw.NewHandler(newMux(sw, 0)); err != nil {
		panic(err)
	}

	http.ListenAndServe(":8087", sw)
}

func newMux(sw muxswitcher.Switcher, index int) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(strconv.Itoa(index) + " ok"))
	})
	mux.HandleFunc("/new", func(w http.ResponseWriter, r *http.Request) {

		m := newMux(sw, index+1)

		if err := sw.NewHandler(m); err != nil {
			http.Error(w, "failed to switch handler: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte("switched to new handler"))
	})
	return mux
}
