package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))

	// Auth (no authorization required)
	rt.router.POST("/api/auth", rt.wrap(rt.doLogin))

	// User
	rt.router.PATCH("/api/users/:userId/me/name", rt.wrap(rt.setMyUserName))
	rt.router.PUT("/api/users/:userId/me/photo", rt.wrap(rt.setMyPhoto))
	rt.router.GET("/api/users/:userId/other/:userName", rt.wrap(rt.findUser))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
