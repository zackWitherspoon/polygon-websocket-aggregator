package service

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"polygon-websocket-aggregator/model/web_socket"
)

const (
	// APIKey TODO: Before going into production, we want this to be retrieved form environment or volume
	APIKey = ""
)

type WebSocketClient interface {
	InitializeWSConnection(url string, tickerName string)
	ReadMessage() (messageType int, p []byte, err error)
	Close()
}

type TradeWebSocket struct {
	wsConn *websocket.Conn
}

var outgoingAuthenticationMessage = []byte(fmt.Sprintf("{\"action\":\"auth\",\"params\":\"%s\"}", APIKey))
var outgoingSubscribeMessage = "{\"action\":\"subscribe\",\"params\":\"T.%s\"}"

func (tradeWS *TradeWebSocket) ReadMessage() (messageType int, p []byte, err error) {
	return tradeWS.wsConn.ReadMessage()
}

func (tradeWS *TradeWebSocket) InitializeWSConnection(url string, tickerName string) {
	var responseMessage web_socket.WebSocketResponse

	logrus.Debug("Attempting to connect to websocket at url: " + url)
	wsConn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		logrus.Fatalf("Dial to the websocket at URL: %s has failed. The following error was returned: %s\n", url, err.Error())
	}

	err = wsConn.ReadJSON(&responseMessage)
	if err != nil {
		logrus.Fatalf("Unable to read the JSON from the connection to WebSocket. Please view the error, fix the code, and try agaian.\n"+
			"error: %s \n", err)
	}
	responseMessage.DebugResponse()

	//authenticate websocket
	authError := wsConn.WriteMessage(websocket.TextMessage, outgoingAuthenticationMessage)
	if authError != nil {
		logrus.Fatal(authError)
	}
	err = wsConn.ReadJSON(&responseMessage)
	if err != nil {
		logrus.Fatalf("Unable to read the JSON from the connection to WebSocket. Please view the error, fix the code, and try agaian.\n"+
			"error: %s \n", err)
	}
	responseMessage.DebugResponse()
	//subscribe to websocket
	var a = []byte(fmt.Sprintf(outgoingSubscribeMessage, tickerName))
	subscribeError := wsConn.WriteMessage(websocket.TextMessage, a)
	err = wsConn.ReadJSON(&responseMessage)
	if err != nil {
		logrus.Fatal(subscribeError)
	}
	responseMessage.DebugResponse()
	tradeWS.wsConn = wsConn
}

func (tradeWS *TradeWebSocket) Close() {
	tradeWS.wsConn.Close()
}
