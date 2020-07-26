package server

import (
	"encoding/json"
	"net/http"
	"log"
)

func writeJSON(rw http.ResponseWriter, data interface{}) {
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	reply, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	rw.Write(reply)
	rw.Write([]byte("\n"))
}
