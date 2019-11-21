package middleware

import (
	"net/http"

	"github.com/DoubleChuang/Drl/api"
	"github.com/DoubleChuang/Drl/pkg/ratelimit"
	"github.com/julienschmidt/httprouter"
)

type middleWareHandler struct {
	router         *httprouter.Router
	connectLimiter *ratelimit.ConnectLimiter
}

func NewMiddleWareHandler(r *httprouter.Router) http.Handler {
	m := middleWareHandler{}
	m.router = r
	m.connectLimiter = ratelimit.GetConnectLimiter()
	return m
}

func (m middleWareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	if _, err = m.connectLimiter.Take(r); err != nil {
		api.SendErrorResponse(w, http.StatusTooManyRequests, "Error")
		return
	}
	m.router.ServeHTTP(w, r)
}
