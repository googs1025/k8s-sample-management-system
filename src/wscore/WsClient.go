package wscore

import (
	"github.com/gorilla/websocket"
	"time"
)

type WsClient struct {
	conn *websocket.Conn
	readChan chan *WsMessage  // 读队列 (chan)
	closeChan chan byte  	  // 失败chan
}

func NewWsClient(conn *websocket.Conn) *WsClient {
	return &WsClient{conn: conn,readChan:make(chan *WsMessage),closeChan:make(chan byte)}
}

func(c *WsClient) Ping(waittime time.Duration){
	for {
		time.Sleep(waittime)
		err := c.conn.WriteMessage(websocket.TextMessage,[]byte("ping"))
		if err != nil {
			ClientMap.Remove(c.conn)
			return
		}
	}
}


func(c *WsClient) ReadLoop(){
	for {
		t, data, err := c.conn.ReadMessage()
		if err != nil {
			c.conn.Close()
			ClientMap.Remove(c.conn)
			c.closeChan <- 1
			break
		}

		c.readChan <- NewWsMessage(t,data)
	}
}



