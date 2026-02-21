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
	rt.router.PATCH("/api/users/:userID/me/name", rt.wrap(rt.setMyUserName))
	rt.router.PUT("/api/users/:userID/me/photo", rt.wrap(rt.setMyPhoto))
	rt.router.GET("/api/users/:userID/other/:userName", rt.wrap(rt.findUser))

	// Group
	rt.router.POST("/api/users/:userID/groups", rt.wrap(rt.createGroup))
	rt.router.GET("/api/users/:userID/groups/:groupId", rt.wrap(rt.getGroupInfo))
	rt.router.DELETE("/api/users/:userID/groups/:groupId", rt.wrap(rt.leaveGroup))
	rt.router.POST("/api/groups/:groupId/users", rt.wrap(rt.addToGroup))
	rt.router.PATCH("/api/groups/:groupId/name", rt.wrap(rt.setGroupName))
	rt.router.PUT("/api/groups/:groupId/photo", rt.wrap(rt.setGroupPhoto))
	rt.router.DELETE("/api/users/:userID/groups/:groupId/member/:userName", rt.wrap(rt.kickUserFromGroup))
	rt.router.DELETE("/api/groups/:groupId/users/:userID", rt.wrap(rt.removeGroup))

	//Conversation
	rt.router.GET("/api/users/:userID/conversations", rt.wrap(rt.getMyConversations))
	rt.router.GET(" /api/users/:userID/conversations/users/:conversationID/messages", rt.wrap(rt.getConversation))
	rt.router.GET(" /api/users/:userID/conversations/groups/:conversationID/messages", rt.wrap(rt.getConversationGroups))

	// Messages
	rt.router.POST("/api/users/:userID/conversations/users/:conversationID/messages", rt.wrap(rt.sendMessage))
	rt.router.POST("/api/users/:userID/conversations/users/:conversationID/messages/:messageID", rt.wrap(rt.forwardMessage))
	rt.router.DELETE("/api/users/:userID/conversations/users/:conversationID/messages/:messageID", rt.wrap(rt.deleteMessage))
	rt.router.POST("/api/users/:userID/conversations/groups/:conversationID/messages", rt.wrap(rt.sendMessageGroup))
	rt.router.POST("/api/users/:userID/conversations/groups/:conversationID/messages/:messageID", rt.wrap(rt.forwardMessageGroup))
	rt.router.DELETE("/api/users/:userID/conversations/groups/:conversationID/messages/:messageID", rt.wrap(rt.deleteMessageGroup))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}