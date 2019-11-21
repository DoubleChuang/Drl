package api

import (
	"io"
	"net/http"
	"strconv"

	"github.com/DoubleChuang/Drl/pkg/ratelimit"
	"github.com/julienschmidt/httprouter"
)

func helloHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	cl := ratelimit.GetConnectLimiter()
	v, err := cl.Get(r)
	if err != nil {
		SendErrorResponse(w, http.StatusTooManyRequests, "Error")
	}
	io.WriteString(w, strconv.FormatInt(v, 10))
}
