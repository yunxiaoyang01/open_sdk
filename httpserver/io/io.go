package io

import (
	"github.com/gin-gonic/gin"
)

type Binder interface {
	Bind(ctx *gin.Context, method string, obj interface{}) error
}

type Render interface {
	Rend(ctx *gin.Context, method string, response interface{})
}
