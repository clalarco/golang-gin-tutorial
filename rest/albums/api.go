package albums

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddRoutes(context *gin.Engine, groupPath string) {
	// Create a handler (which will later have injected services)
	h := GetHandler()

	group := context.Group(groupPath)
	group.GET("", h.GetAlbums)
	group.GET(":ID", h.GetAlbum)
	group.POST("", h.AddAlbum)
	group.DELETE(":ID", h.DeleteAlbum)
}

// getAlbums responds with the list of all albums as JSON.
func (h *Handler) GetAlbums(c *gin.Context) {
	albums, error := h.DB.GetAlbums()
	if error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": error.Error()})
		return
	}
	if len(albums) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No albums found"})
		return
	}
	c.JSON(http.StatusOK, albums)
}

// getAlbums responds with the list of all albums as JSON.
func (h *Handler) GetAlbum(c *gin.Context) {
	id, isParamFound := c.Params.Get("ID")
	if !isParamFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parameter ID is required"})
		return
	}
	album, error := h.DB.GetAlbum(id)
	if error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": error.Error()})
		return
	}
	c.JSON(http.StatusOK, album)
}

// AddAlbum adds an album from JSON received in the request body.
func (h *Handler) AddAlbum(c *gin.Context) {
	var newAlbum album

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}
	error := h.DB.AddAlbum(newAlbum)
	if error != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "ID already exists"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Album added successfully"})
}

func (h *Handler) DeleteAlbum(c *gin.Context) {
	id, isParamFound := c.Params.Get("ID")
	if !isParamFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parameter ID is required"})
		return
	}
	error := h.DB.DeleteAlbum(id)
	if error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Album deleted successfully"})
}
