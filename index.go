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

	"github.com/patrickmn/go-cache"
)

const PIN_LENGTH = 8
const PIN_PADDING = "%08d"
const MAX_PIN = 100000000
const TIME_LIMIT = 2

var FLAG = os.Getenv("FLAG")

var pin_cache = cache.New(TIME_LIMIT*time.Second, 2*TIME_LIMIT*time.Second)

func generatePIN() string {
	pass := fmt.Sprintf(PIN_PADDING, rand.Intn(MAX_PIN))
	hash := sha256.Sum256([]byte(pass))
	hexDigest := hex.EncodeToString(hash[:])
	pin_cache.Set(hexDigest, pass, cache.DefaultExpiration)
	fmt.Println(pass, hexDigest)
	return hexDigest
}

func Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		givenHash := r.URL.Path[1:]
		body, _ := ioutil.ReadAll(r.Body)
		givenPass := string(body)
		pin, found := pin_cache.Get(givenHash)
		if found && pin.(string) == givenPass {
			fmt.Fprintln(w, FLAG)
		}
	case "GET":
		hash := generatePIN()
		fmt.Fprintln(w, hash)
	}
}
