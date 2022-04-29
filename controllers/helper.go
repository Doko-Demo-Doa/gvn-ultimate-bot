package controllers

import "github.com/gin-gonic/gin"

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func HTTPRes(c *gin.Context, httpCode int, msg string, data interface{}) {
	// c.JSON()
}
