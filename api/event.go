package api

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	pb "github.com/inkah-trace/daemon/protobuf"
)

func EventCreate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	bdbc := getBoltDB(r)

	var e pb.ForwardedEvent
	decodeJson(r, &e)
	defer r.Body.Close()

	bdbc.TraceService().CreateEvent(&e)

	jsonResponse(w, e)
}