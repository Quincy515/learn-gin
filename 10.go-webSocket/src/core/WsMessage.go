package core

type WsMessage struct {
	MessageType int
	MessageData []byte
}

func NewWsMessage(messageType int, messageData []byte) *WsMessage {
	return &WsMessage{MessageType: messageType, MessageData: messageData}
}
