package rpa

import (
	"io"
	"server/internal/api"

	"github.com/gin-gonic/gin"
)

func (h *Handler) file(gin *gin.Context) {
	g := api.GetG(gin)

	id := gin.Param("id")
	h.RLock()
	c, ok := h.files[id]
	h.RUnlock()
	if !ok {
		g.Notfound("file")
		return
	}

	g.Header("Content-Type", "image/jpeg")
	io.Copy(g.Writer, c.cache)
}
