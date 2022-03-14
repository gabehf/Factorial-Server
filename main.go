package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// type ContextKeyString is a type that is a string meant for use as the key in an https context value
type ContextKeyString string
type args struct {
	A int `json:"a"`
	B int `json:"b"`
}
type requestError struct {
	Err string `json:"error"`
}

func factorial(a int, ch chan int) {
	r := 1
	for f := 1; f <= a; f++ {
		r *= f
	}
	ch <- r
}

func postError(w http.ResponseWriter, r *http.Request, err string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(requestError{Err: err})
}

func middleware(f httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var v args
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		err := dec.Decode(&v)
		if err != nil || v.A < 0 || v.B < 0 {
			postError(w, r, "Incorrect input")
		} else {
			ctx := r.Context()
			// using a new type definition to avoid compiler warnings
			ctx = context.WithValue(ctx, ContextKeyString("a"), v.A)
			ctx = context.WithValue(ctx, ContextKeyString("b"), v.B)
			r = r.WithContext(ctx)
			f(w, r, p)
		}
	}
}

func calculate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var v args
	v.A = r.Context().Value(ContextKeyString("a")).(int)
	v.B = r.Context().Value(ContextKeyString("b")).(int)
	chanA := make(chan int)
	chanB := make(chan int)
	go factorial(v.A, chanA)
	go factorial(v.B, chanB)
	v.B = <-chanB
	v.A = <-chanA
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func main() {
	router := httprouter.New()
	router.POST("/calculate", middleware(calculate))

	log.Fatal(http.ListenAndServe(":8989", router))
}
