package master

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type HTTP struct {
	u Usecase
}

func NewHTTP(
	r *gin.Engine,
	u Usecase,
) {
	h := &HTTP{u: u}
	r.POST("/enqueue", h.enqueue)
}

type EnqueueRequest struct {
	EndNode         int    `json:"end_node"`
	DelayTimeInSecs int    `json:"delay_time_in_secs"`
	Message         string `json:"message"`
}

func (h *HTTP) enqueue(ctx *gin.Context) {
	var r EnqueueRequest
	if err := ctx.BindJSON(&r); err != nil {
		return
	}
	err := h.u.Enqueue(ctx.Request.Context(), EnqueueCommand{
		EndNode:         r.EndNode,
		DelayTimeInSecs: r.DelayTimeInSecs,
		Message:         r.Message,
	})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
