package controllers

import "fmt"

const (
	queryOrderParam = "ordering"
	queryPageParam  = "page"
)

var errPageParamMsg = fmt.Sprintf("\"%s\" should be a positive integer\n", queryPageParam)
var errObjectIdParamMsg = fmt.Sprint("ObjectIDs must be exactly 12 bytes long\n")
