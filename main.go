package main

import (
	"github.com/nkanaev/yarr/server"
)

func main() {
	srv := server.New()
	srv.ListenAndServe()
}
