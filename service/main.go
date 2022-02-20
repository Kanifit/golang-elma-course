package main

import (
	"golang-elma-course/service/http/router"
	"net/http"
)

func main() {
	err := http.ListenAndServe(":8080", router.Router())
	if err != nil {
		panic(err)
	}
}
