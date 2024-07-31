package common

import "github.com/gin-gonic/gin"

// Handler struct holds required services for handler to function
type Handler struct{}

func AddRoutes(context *gin.Engine, groupPath string) {
	// Create a handler (which will later have injected services)
	h := &Handler{} // currently has no properties

	group := context.Group(groupPath)
	group.GET("", h.Ping)

}

func (h *Handler) Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
