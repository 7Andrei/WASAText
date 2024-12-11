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
	rt.router.PUT("/settings/changename", rt.setUsername)
	rt.router.GET("/getuser", rt.getUser)
	rt.router.PUT("/settings/changephoto", rt.setPhoto)

	rt.router.GET("/getchat", rt.getChat)
	rt.router.POST("/createchat", rt.createChat)
	//rt.router.GET("/insert", rt.insertUser)
	//rt.router.GET("/", rt.getUserName)

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
