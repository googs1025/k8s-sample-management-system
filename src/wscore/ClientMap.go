package wscore

import (
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
)
//外部公共使用,保存所有客户端对象的Map
var ClientMap *ClientMapStruct

func init() {
	ClientMap = &ClientMapStruct{}
}

type ClientMapStruct struct {
	data sync.Map  //  key 是客户端IP  value 就是 WsClient连接对象
}

func(c *ClientMapStruct) Store(conn *websocket.Conn){
	 wsClient:=NewWsClient(conn)
	 c.data.Store(conn.RemoteAddr().String(),wsClient)
	 go wsClient.Ping(time.Second*1)
	 go wsClient.ReadLoop() //处理读 循环
	// go wsClient.HandlerLoop() //处理 总控制循环
}


//向所有客户端 发送消息--发送deployment列表
func(c *ClientMapStruct) SendAllDepList(v interface{}){
	c.data.Range(func(key, value interface{}) bool {
		con:=value.(*WsClient).conn
 			err:=con.WriteJSON(v)
			if err!=nil{
				c.Remove(con)
				log.Println(err)
			}
			return true
	})
}

func(c *ClientMapStruct) Remove(conn *websocket.Conn){
	c.data.Delete(conn.RemoteAddr().String())
}


