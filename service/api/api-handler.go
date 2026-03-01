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
	rt.router.PATCH("/api/users/:userID/me/name", rt.wrapAuth(rt.setMyUserName))
	rt.router.PUT("/api/users/:userID/me/photo", rt.wrapAuth(rt.setMyPhoto))
	rt.router.GET("/api/users/:userID/other/:userName", rt.wrapAuth(rt.findUser))

	// Group
	rt.router.POST("/api/users/:userID/groups", rt.wrapAuth(rt.createGroup))
	rt.router.GET("/api/users/:userID/groups/:groupId", rt.wrapAuth(rt.getGroupInfo))
	rt.router.DELETE("/api/users/:userID/groups/:groupId", rt.wrapAuth(rt.leaveGroup))
	rt.router.POST("/api/groups/:groupId/users", rt.wrapAuth(rt.addToGroup))
	rt.router.PATCH("/api/groups/:groupId/name", rt.wrapAuth(rt.setGroupName))
	rt.router.PUT("/api/groups/:groupId/photo", rt.wrapAuth(rt.setGroupPhoto))
	rt.router.DELETE("/api/users/:userID/groups/:groupId/member/:userName", rt.wrapAuth(rt.kickUserFromGroup))
	rt.router.DELETE("/api/groups/:groupId/users/:userID", rt.wrapAuth(rt.removeGroup))

	// Conversation
	rt.router.GET("/api/users/:userID/conversations", rt.wrapAuth(rt.getMyConversations))
	rt.router.GET("/api/users/:userID/conversations/users/:conversationID/messages", rt.wrapAuth(rt.getConversation))
	rt.router.GET("/api/users/:userID/conversations/groups/:conversationID/messages", rt.wrapAuth(rt.getConversationGroups))

	// Messages
	rt.router.POST("/api/users/:userID/conversations/users/:conversationID/messages", rt.wrapAuth(rt.sendMessage))
	rt.router.POST("/api/users/:userID/conversations/users/:conversationID/messages/:messageID", rt.wrapAuth(rt.forwardMessage))
	rt.router.DELETE("/api/users/:userID/conversations/users/:conversationID/messages/:messageID", rt.wrapAuth(rt.deleteMessage))
	rt.router.POST("/api/users/:userID/conversations/groups/:conversationID/messages", rt.wrapAuth(rt.sendMessageGroup))
	rt.router.POST("/api/users/:userID/conversations/groups/:conversationID/messages/:messageID", rt.wrapAuth(rt.forwardMessageGroup))
	rt.router.DELETE("/api/users/:userID/conversations/groups/:conversationID/messages/:messageID", rt.wrapAuth(rt.deleteMessageGroup))

	// Comment (moved to avoid httprouter wildcard conflict)
	rt.router.POST("/api/comments/users/:userID/messages/:messageID", rt.wrapAuth(rt.commentMessage))
	rt.router.DELETE("/api/comments/users/:userID/messages/:messageID/:commentID", rt.wrapAuth(rt.deleteComment))
	rt.router.POST("/api/comments/groups/:groupID/messages/:messageID", rt.wrapAuth(rt.commentMessageGroup))
	rt.router.DELETE("/api/comments/groups/:groupID/messages/:messageID/:commentID", rt.wrapAuth(rt.deleteCommentGroup))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
