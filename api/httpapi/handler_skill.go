package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListSkillsHandler 返回 Skill 列表
func (h *Handler) ListSkillsHandler(c *gin.Context) {
	if h.skillService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "skill service is not configured",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"skills": h.skillService.ListSkills(),
	})
}
