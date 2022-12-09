package file

import (
	"net/http"
	"server/internal/api"

	"github.com/gin-gonic/gin"
)

func (h *Handler) uploadHandle(gin *gin.Context) {
	g := api.GetG(gin)

	id := g.Param("id")

	h.RLock()
	cache := h.uploadCache[id]
	h.RUnlock()
	if cache == nil {
		g.HTTPError(http.StatusNotFound, "not found")
		return
	}
	if g.GetHeader("X-Token") != cache.token {
		g.HTTPError(http.StatusForbidden, "access denied")
		return
	}
	g.File(cache.dir)
	h.removeUploadCache(id)
}
