package main

import (
	"crypto/rand"
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

const (
	Vowels = "aeioAEIO"
	Leet   = "43104310"
)

func randomReplaceVowels(s string) (string, error) {
	for pos, runeValue := range s {
		if index := strings.IndexRune(Vowels, runeValue); index > -1 {
			randInt, err := rand.Int(rand.Reader, big.NewInt(2))
			if err != nil {
				return "", nil
			}
			if randInt.Cmp(big.NewInt(1)) == 0 {
				s = s[0:pos] + string(Leet[index]) + s[pos+1:len(s)]
			}
		}
	}
	return s, nil
}

func queryInt(r *http.Request, key string) (int, error) {
	value := r.URL.Query().Get(key)
	if value == "" {
		return 0, nil
	} else {
		return strconv.Atoi(value)
	}
}

func main() {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		var invalidParams []string
		length, err := queryInt(r, "length")
		if err != nil {
			invalidParams = append(invalidParams, "length")
		}
		numDigits, err := queryInt(r, "numDigits")
		if err != nil {
			fmt.Println(err.Error())
			invalidParams = append(invalidParams, "numDigits")
		}
		numSymbols, err := queryInt(r, "numSymbols")
		if err != nil {
			invalidParams = append(invalidParams, "numSymbols")
		}
		if len(invalidParams) > 0 {
			message := fmt.Sprintf("invalid parameter(s): %s", strings.Join(invalidParams, ","))
			http.Error(w, message, 422)
			return
		}
		res, err := password.Generate(length, numDigits, numSymbols, false, true)
		if err != nil {
			http.Error(w, err.Error(), 422)
			return
		}
		res, err = randomReplaceVowels(res)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Write([]byte(res))
	})

	err := http.ListenAndServe(":3334", r)
	if err != nil {
		log.Print(err)
	}
}
