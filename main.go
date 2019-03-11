package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sethvargo/go-password/password"
)

// characters in "vowels" are replaced as they
// appear in "leet" (e.g. a -> 4, e -> 3, ...)
const (
	vowels = "aeioAEIO"
	leet   = "43104310"
)

func main() {
	// use chi router for routing
	r := chi.NewRouter()

	// proposed base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// there is only one route which returns password candidates
	r.Get("/passwords", getPasswords)

	err := http.ListenAndServe(":3334", r)
	if err != nil {
		log.Print(err)
	}
}

// response holding password candidates
type passwordsResponse struct {
	Candidates []string `json:"candidates"`
}

// get password candidates
func getPasswords(w http.ResponseWriter, r *http.Request) {
	// holds all invalid query param names
	var invalidParams []string

	// query param: length of the password
	length, err := queryInt(r, "length", 8)
	if err != nil {
		invalidParams = append(invalidParams, "length")
	}
	// query param: number of digits of the password
	numDigits, err := queryInt(r, "numDigits", 0)
	if err != nil {
		invalidParams = append(invalidParams, "numDigits")
	}
	// query param: number of symbols of the password
	numSymbols, err := queryInt(r, "numSymbols", 0)
	if err != nil {
		invalidParams = append(invalidParams, "numSymbols")
	}
	// query param: number of password candidates to generate
	numCandidates, err := queryInt(r, "numCandidates", 4)
	if err != nil {
		invalidParams = append(invalidParams, "numCandidates")
	}
	// query param: enabled replacing vowels randomly
	replaceVowels, err := queryBool(r, "replaceVowels", false)
	if err != nil {
		invalidParams = append(invalidParams, "replaceVowels")
	}
	// if any invalid query params exists return error
	if len(invalidParams) > 0 {
		message := fmt.Sprintf("invalid parameter(s): %s", strings.Join(invalidParams, ","))
		http.Error(w, message, 422)
		return
	}
	var passwordsResponse passwordsResponse
	// generate the desired number of password candidates and add them to the response
	for i := 0; i < numCandidates; i++ {
		candidate, err := generatePassword(length, numDigits, numSymbols, replaceVowels)
		// if there is an error it most likely happend during password generation
		// -> status code 422
		if err != nil {
			http.Error(w, err.Error(), 422)
			return
		}
		passwordsResponse.Candidates = append(passwordsResponse.Candidates, candidate)
	}
	// serialize response (to json) and send it
	response, _ := json.Marshal(passwordsResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(response)
}

// replaces vowels in given string randomly and returns result
func randomReplaceVowels(s string) (string, error) {
	// iterate over string rune by rune
	for pos, runeValue := range s {
		// check if rune at current position is vowel
		if index := strings.IndexRune(vowels, runeValue); index > -1 {
			// rand int from [0,1]
			randInt, err := rand.Int(rand.Reader, big.NewInt(2))
			if err != nil {
				return "", nil
			}
			// 0 -> do not replace
			// 1 -> replace
			if randInt.Cmp(big.NewInt(1)) == 0 {
				// replace vowel with number assuming "rune equals byte"
				s = s[0:pos] + string(leet[index]) + s[pos+1:len(s)]
			}
		}
	}
	return s, nil
}

// returns int query param from request or default value if query param does not exists
func queryInt(r *http.Request, key string, defaultValue int) (int, error) {
	value := r.URL.Query().Get(key)
	if value == "" {
		return defaultValue, nil
	}
	return strconv.Atoi(value)
}

// returns bool query param from request or default value if query param does not exists
func queryBool(r *http.Request, key string, defaultValue bool) (bool, error) {
	value := r.URL.Query().Get(key)
	if value == "" {
		return defaultValue, nil
	}
	return strconv.ParseBool(value)

}

// generates and returns password
func generatePassword(length int, numDigits int, numSymbols int, replaceVowles bool) (string, error) {
	// generate password according to parameters
	res, err := password.Generate(length, numDigits, numSymbols, false, true)
	if err != nil {
		return "", err
	}
	// replace vowels (optionally)
	if replaceVowles {
		return randomReplaceVowels(res)
	}
	return res, nil
}
