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

	// Group
	rt.router.POST("/api/users/:userId/groups", rt.wrap(rt.createGroup))
	rt.router.GET("/api/users/:userId/groups/:groupId", rt.wrap(rt.getGroupInfo))
	rt.router.DELETE("/api/users/:userId/groups/:groupId", rt.wrap(rt.leaveGroup))
	rt.router.POST("/api/groups/:groupId/users", rt.wrap(rt.addToGroup))
	rt.router.PATCH("/api/groups/:groupId/name", rt.wrap(rt.setGroupName))
	rt.router.PUT("/api/groups/:groupId/photo", rt.wrap(rt.setGroupPhoto))
	rt.router.DELETE("/api/users/:userId}/groups/:groupId/member/:userName", rt.wrap(rt.kickUserFromGroup))
	rt.router.DELETE("/api/groups/:groupId}/users/:userId", rt.wrap(rt.removeGroup))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
