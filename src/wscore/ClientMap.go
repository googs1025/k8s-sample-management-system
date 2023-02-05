package wscore

import (
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
)

// 所有client端的集合
var ClientMap *ClientMapStruct

func init() {
	ClientMap = &ClientMapStruct{}
}

type ClientMapStruct struct {
	data sync.Map  //  key:客户端IP  value:WsClient连接对象
	lock sync.Mutex
}

func(c *ClientMapStruct) Store(conn *websocket.Conn){
	 wsClient := NewWsClient(conn)
	 c.data.Store(conn.RemoteAddr().String(),wsClient)

	 go wsClient.Ping(time.Second * 1)
	 go wsClient.ReadLoop() //处理读 循环
	// go wsClient.HandlerLoop() //处理 总控制循环
}


// 向所有客户端 发送消息--发送deployment列表
func(c *ClientMapStruct) SendAll(v interface{}){
	c.data.Range(func(key, value interface{}) bool {
		c.lock.Lock() //这里加个锁，因为目前不支持并发写wsClient
		defer c.lock.Unlock()
		func(){
			cc := value.(*WsClient).conn
			err := cc.WriteJSON(v)
			if err != nil {
				c.Remove(cc)
				log.Println(err)
			}
		}()

		return true
	})
}

func(c *ClientMapStruct) Remove(conn *websocket.Conn){
	c.data.Delete(conn.RemoteAddr().String())
}


