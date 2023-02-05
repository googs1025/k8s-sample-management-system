package wscore

// 前端消息
type WsMessage struct {
	MessageType int
	MessageData []byte  // JSON格式
}

func NewWsMessage(messageType int, messageData []byte) *WsMessage {
	return &WsMessage{MessageType: messageType, MessageData: messageData}
}


