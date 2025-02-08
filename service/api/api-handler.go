package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	//rt.router.GET("/", rt.wrap(rt.insertUser))
	rt.router.GET("/context", rt.wrap(rt.getContextReply))

	rt.router.POST("/session", rt.loginUser)
	rt.router.GET("/users/:user_id", rt.getUser)
	rt.router.PUT("/settings/changename", rt.setUsername)
	rt.router.PUT("/settings/changephoto", rt.setPhoto)
	rt.router.GET("/users", rt.getAllUsers)

	rt.router.GET("/chats/:chat_id", rt.getChat)
	rt.router.GET("/chats", rt.getAllChats)
	rt.router.POST("/createchat", rt.createChat)
	rt.router.PUT("/chats/:chat_id/settings/changename", rt.setChatName)
	rt.router.PUT("/chats/:chat_id/settings/changephoto", rt.setChatPhoto)
	rt.router.POST("/chats/:chat_id/settings/add", rt.addUserToChat)
	rt.router.DELETE("/chats/:chat_id/settings/leave", rt.leaveChat)

	rt.router.POST("/chats/:chat_id", rt.sendMessage)
	rt.router.POST("/chats/:chat_id/messages/:message_id", rt.forwardMessage)
	rt.router.DELETE("/chats/:chat_id/messages/:message_id/delete", rt.deleteMessage)

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
