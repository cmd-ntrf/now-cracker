package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type thing struct {
	time time.Time
	pass string
	hash string
}

const PIN_LENGTH = 8
const PIN_PADDING = "%08d"
const MAX_PIN = 100000000

var things = make([]thing, 0)

func filterOutOldThings(t time.Time) {
	new_things := make([]thing, 0)
	for _, thing := range things {
		if thing.time.Add(time.Second * 2).After(t) {
			new_things = append(new_things, thing)
		}
	}
	things = new_things
	fmt.Printf("%v\n", things)
}

func addNewThing(t time.Time) string {
	pass := fmt.Sprintf(PIN_PADDING, rand.Intn(MAX_PIN))
	hash := sha256.Sum256([]byte(pass))
	hexDigest := hex.EncodeToString(hash[:])
	things = append(things, thing{t, pass, hexDigest})
	return hexDigest
}

func Handler(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	filterOutOldThings(now)
	switch r.Method {
	case "POST":
		givenHash := r.URL.Path[1:]
		body, _ := ioutil.ReadAll(r.Body)
		givenPass := string(body)
		for _, thing := range things {
			if thing.hash == givenHash && thing.pass == givenPass {
				fmt.Fprintln(w, os.Getenv("FLAG"))
			}
		}
	case "GET":
		hash := addNewThing(now)
		fmt.Fprintln(w, hash)
	}
}
