package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/shenyisyn/goft-gin/goft"
	"net/http"
)

type CrosMiddleware struct {

}

func NewCrosMiddleware() *CrosMiddleware {
	return &CrosMiddleware{}
}
func(*CrosMiddleware) OnRequest(c *gin.Context) error{
	method := c.Request.Method
	if method != "" {
		c.Header("Access-Control-Allow-Origin", "*")  // 可将将 * 替换为指定的域名
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization,X-Token")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
	}

	if method == "OPTIONS" {
		c.Set(goft.HTTP_STATUS,http.StatusNoContent) //设置响应 httpcode
		panic("")
	}
	return nil

}
func(*CrosMiddleware) OnResponse(result interface{}) (interface{}, error){
	return result,nil
}
