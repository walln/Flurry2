package main

import (
	"github.com/walln/flurry2/flurry"
)

func main() {

	server := flurry.Initialize()
	server.ListenAndServe()

}
