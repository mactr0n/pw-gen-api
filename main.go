package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sethvargo/go-password/password"
)

func queryInt(r *http.Request, key string) (int, error) {
	value := r.URL.Query().Get(key)
	return strconv.Atoi(value)
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
		w.Write([]byte(res))
	})

	err := http.ListenAndServe(":3334", r)
	if err != nil {
		log.Print(err)
	}
}
