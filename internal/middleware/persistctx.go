package middleware

import (
	"codeberg.org/gruf/go-pctx"
	"github.com/gin-gonic/gin"
)

// PersistContext is a middleware handler that replaces an existing
// request context with one capable of storing persisting key-values.
// Ideally this should be the first middleware in the handler stack.
func PersistContext() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get base request context.
		base := ctx.Request.Context()

		// Wrap within persistent ctx.
		pCtx := pctx.Get(base)

		// Replace existing request context.
		r := ctx.Request.WithContext(pCtx)
		ctx.Request = r

		// Next handler.
		ctx.Next()
	}
}
