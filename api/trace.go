package api

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	//pb "github.com/inkah-trace/daemon/protobuf"
)

func TraceList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	bdbc := getBoltDB(r)

	traces, _ := bdbc.TraceService().TraceList()

	jsonResponse(w, traces)
}

func TraceGet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	bdbc := getBoltDB(r)

	id := ps[0].Value

	events, _ := bdbc.TraceService().Trace(id)

	jsonResponse(w, events)
}