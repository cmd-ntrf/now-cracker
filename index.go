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

type tuple struct {
	time time.Time
	pass string
}

const PIN_LENGTH = 8
const PIN_PADDING = "%08d"
const MAX_PIN = 100000000

var tuples = make(map[string]tuple)

func filterOutOldTuples(t time.Time) {
	new_tuples := make(map[string]tuple)
	for key, tuple := range tuples {
		if tuple.time.Add(time.Second * 2).After(t) {
			new_tuples[key] = tuple
		}
	}
	tuples = new_tuples
	fmt.Printf("%v\n", tuples)
}

func addNewThing(t time.Time) string {
	pass := fmt.Sprintf(PIN_PADDING, rand.Intn(MAX_PIN))
	hash := sha256.Sum256([]byte(pass))
	hexDigest := hex.EncodeToString(hash[:])
	tuples[hexDigest] = tuple{t, pass}
	return hexDigest
}

func Handler(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	filterOutOldTuples(now)
	switch r.Method {
	case "POST":
		givenHash := r.URL.Path[1:]
		body, _ := ioutil.ReadAll(r.Body)
		givenPass := string(body)
		tuple, ok := tuples[givenHash]
		if ok && tuple.pass == givenPass {
			fmt.Fprintln(w, os.Getenv("FLAG"))
		}
	case "GET":
		hash := addNewThing(now)
		fmt.Fprintln(w, hash)
	}
}
