package api

import (
	"log"
	"net/http"
	"encoding/json"
	"github.com/inkah-trace/server/bolt"
)

func decodeJson(r *http.Request, v interface{}) interface{} {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(v)

	if err != nil {
		log.Printf("Error decoding profile: %s", err)
	}
	return v
}

func jsonResponse(w http.ResponseWriter, v interface{}) {
	resp, err := json.Marshal(v)
	if err != nil {
		w.Write([]byte("Error!"))
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func getBoltDB(r *http.Request) *bolt.Client {
	ctx := r.Context()
	bdbc := ctx.Value("boltdb").(*bolt.Client)
	return bdbc
}