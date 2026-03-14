package registry

import "github.com/gin-gonic/gin"

// key验证
func (r *Registry) auth(ctx *gin.Context) {
	if ctx.GetHeader("rkey") != r.Key {
		ctx.Status(400)
		return
	}
	ctx.Next()
}
