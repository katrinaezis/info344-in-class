package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

type zip struct {
	// Allows you to assign your variable to the json value
	// single quote used for runes, back quote for any string where you
	// want to perserve line breaks as well as fields
	// use `json:"-"` to ignore field (never exports into output json)
	Zip   string `json:"zip"`
	City  string `json:"city"`
	State string `json:"state"`
}

type zipSlice []*zip
type zipIndex map[string]zipSlice

// * passes by reference, no star passes by value
func helloHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	// must add headers before writing body
	w.Header().Add("Content-Type", "text/plain")
	w.Header().Add("Access-Control-Allow-Origin", "*")

	w.Write([]byte("hello " + name))
}

func (zi zipIndex) zipsForCityHandler(w http.ResponseWriter, r *http.Request) {
	// want to request /zips/city/seattle
	_, city := path.Split(r.URL.Path)
	lcity := strings.ToLower(city)

	w.Header().Add("Content-Type", "application/json; charset=utf-8")

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(zi[lcity]); err != nil {
		http.Error(w, "error encoding json: "+err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		log.Fatal("please set ADDR enviro var")
	}

	f, err := os.Open("../data/zips.json")
	if err != nil {
		log.Fatal("error opening zips file: " + err.Error())
	}

	// declared slice of pointers to zip struct
	zips := make(zipSlice, 0, 43000)
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&zips); err != nil {
		log.Fatal("error decoding zips json:" + err.Error())
	}
	fmt.Printf("loaded %d zips\n", len(zips))

	zi := make(zipIndex)

	for _, z := range zips {
		lower := strings.ToLower(z.City)
		zi[lower] = append(zi[lower], z)

	}

	fmt.Printf("there are %d zips in seattle\n", len(zi["seattle"]))

	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/zips/city/", zi.zipsForCityHandler)

	fmt.Printf("server is listening at %s...\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
