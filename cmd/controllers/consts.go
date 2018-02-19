package controllers

import "fmt"

const (
	queryOrderParam = "ordering"
	queryPageParam = "page"
)

var errPageParamMsg = fmt.Sprintf("\"%s\" should be an integer", queryPageParam)
