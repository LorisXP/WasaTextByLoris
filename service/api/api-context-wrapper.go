package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/LorisXP/WasaTextByLoris/service/api/reqcontext"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

// httpRouterHandler is the signature for functions that accepts a reqcontext.RequestContext in addition to those
// required by the httprouter package.
type httpRouterHandler func(http.ResponseWriter, *http.Request, httprouter.Params, reqcontext.RequestContext)

// wrap parses the request and adds a reqcontext.RequestContext instance related to the request.
func (rt *_router) wrap(fn httpRouterHandler) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		reqUUID, err := uuid.NewV4()
		if err != nil {
			rt.baseLogger.WithError(err).Error("can't generate a request UUID")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var ctx = reqcontext.RequestContext{
			ReqUUID: reqUUID,
		}

		// Create a request-specific logger
		ctx.Logger = rt.baseLogger.WithFields(logrus.Fields{
			"reqid":     ctx.ReqUUID.String(),
			"remote-ip": r.RemoteAddr,
		})

		// Call the next handler in chain (usually, the handler function for the path)
		fn(w, r, ps, ctx)
	}
}

// wrapAuth è come wrap, ma richiede un header Authorization: Bearer <userID> valido.
// Se il path contiene il parametro :userID, verifica che corrisponda al token Bearer.
func (rt *_router) wrapAuth(fn httpRouterHandler) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		reqUUID, err := uuid.NewV4()
		if err != nil {
			rt.baseLogger.WithError(err).Error("can't generate a request UUID")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var ctx = reqcontext.RequestContext{
			ReqUUID: reqUUID,
		}

		// Create a request-specific logger
		ctx.Logger = rt.baseLogger.WithFields(logrus.Fields{
			"reqid":     ctx.ReqUUID.String(),
			"remote-ip": r.RemoteAddr,
		})

		// Estrai e valida il token Bearer dall'header Authorization
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.Logger.Warn("missing or invalid Authorization header")
			http.Error(w, "missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		bearerUserID, errToken := strconv.Atoi(tokenStr)
		if errToken != nil || bearerUserID <= 0 {
			ctx.Logger.Warn("invalid Bearer token: not a valid user identifier")
			http.Error(w, "invalid Bearer token", http.StatusUnauthorized)
			return
		}
		ctx.BearerUserID = bearerUserID

		// Se la rotta ha :userID nel path, verifica che corrisponda al Bearer token
		if pathUserIDStr := ps.ByName("userID"); pathUserIDStr != "" {
			pathUserID, errPath := strconv.Atoi(pathUserIDStr)
			if errPath == nil && pathUserID != bearerUserID {
				ctx.Logger.Warnf("Bearer userID %d does not match path userID %d", bearerUserID, pathUserID)
				http.Error(w, "forbidden: token does not match the requested user", http.StatusForbidden)
				return
			}
		}

		// Call the next handler in chain
		fn(w, r, ps, ctx)
	}
}
