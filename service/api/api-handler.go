package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/context", rt.wrap(rt.getContextReply))

	rt.router.POST("/session", rt.loginUser) //
	rt.router.GET("/users/:user_id", rt.getUser)
	rt.router.PUT("/settings/name", rt.setUsername) //
	rt.router.PUT("/settings/photo", rt.setPhoto)   //
	rt.router.GET("/users", rt.getAllUsers)

	rt.router.GET("/chats/:chat_id", rt.getChat)                       //
	rt.router.GET("/chats", rt.getAllChats)                            //
	rt.router.POST("/chat", rt.createChat)                             //
	rt.router.PUT("/chats/:chat_id/settings/name", rt.setChatName)     //
	rt.router.PUT("/chats/:chat_id/settings/photo", rt.setChatPhoto)   //
	rt.router.POST("/chats/:chat_id/settings/users", rt.addUserToChat) //
	rt.router.DELETE("/chats/:chat_id/settings", rt.leaveChat)         //

	rt.router.POST("/chats/:chat_id", rt.sendMessage)                                                  //
	rt.router.POST("/chats/:chat_id/messages/:message_id", rt.forwardMessage)                          //
	rt.router.DELETE("/chats/:chat_id/messages/:message_id", rt.deleteMessage)                         //
	rt.router.POST("/chats/:chat_id/messages/:message_id/reactions", rt.addReaction)                   //
	rt.router.DELETE("/chats/:chat_id/messages/:message_id/reactions/:reaction_id", rt.deleteReaction) //

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
