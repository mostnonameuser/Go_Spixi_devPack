package network

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/mostnonameuser/Go_Spixi_devPack/internal/config"
	"log"
	"net/http"
	"sync"
	"time"
)

const address = "MyTestAddress"
type WsBroker struct {
	config  *config.Config
	conns  *sync.Map
	mu    sync.Mutex
	messageChan chan Message
}

func (w *WsBroker) Connect() error {
	w.messageChan = make(chan Message, 100)
	http.HandleFunc("/ws", w.handleWS)
	go http.ListenAndServe( w.config.DevServerWs, nil)
	return nil
}

func (w *WsBroker) handleWS(wr http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(wr, r, nil, 1024, 1024)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()
	log.Println("WebSocket client connected")
	w.conns.Store(address, conn)
	defer w.conns.Delete(conn)
	for {
		_, payload, err := conn.ReadMessage()
		if err != nil {
			break
		}
		var rcvd SendAppDataParams
		if err := json.Unmarshal(payload, &rcvd); err != nil {
			log.Printf("Invalid JSON from WS: %v", err)
			continue
		}
		w.messageChan <- Message{
			Topic:   "AppMessage",
			Payload: payload,
		}
	}
}

func (w *WsBroker) Subscribe(topic string, handler MessageHandlerFunc) error {
	return nil
}

func (w *WsBroker) GetMessageChannel() <-chan Message {
	return w.messageChan
}

func (w *WsBroker) Publish(topic, addr string, payload []byte) error {
	reqBody := SendAppDataParams{
		Address:   addr,
		ProtocolId: w.config.ClientId,
		Data:      string(payload),
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		log.Printf("marshal devService request error: %v \n\n", err)
	}
	if raw, ok := w.conns.Load(address); ok {
		if conn, ok := raw.(*websocket.Conn); ok {
			w.mu.Lock()
			err = conn.WriteMessage(websocket.TextMessage, jsonBody)
			if err != nil {
				fmt.Println("Error WS Send ", err)
			}
			w.mu.Unlock()
		}
	}
	return nil
}

func (w *WsBroker) Disconnect() error {
	return nil
}

func (w *WsBroker) Start (ctx context.Context) error {
	fmt.Println("WebSocket MODE")
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				log.Println("GetMessages stopped due to cancellation")
				return
			case msg, ok := <- w.messageChan:
				if !ok {
					log.Println("Message queue closed")
					return
				}
				rcvd := SendAppDataParams{}
				erss := json.Unmarshal(msg.Payload, &rcvd)
				if erss != nil {
					log.Println("Error Unmarch Initial Message from chat")
				}
				fmt.Println("WS MESSAGE", rcvd.ProtocolId, rcvd.Data)
				mesg := &AiMessage{}
				ersA := json.Unmarshal([]byte(rcvd.Data), &mesg)
				if ersA != nil {
					log.Println("Error Unmarch Message from WsChat")
				}
				jsb, _ := json.Marshal(mesg)
				fmt.Println("WS MESSAGE", string(jsb))
				jsonBody := MessageProcessor(mesg)
				if jsonBody != nil {
					postIxiErr := w.Publish( "sendAppData", address, jsonBody)
					if postIxiErr != nil {
						log.Printf("Processing context messages WS answer error %v\n", postIxiErr)
					}
				}
			}
		}
	}()
	go func() {
		<-ctx.Done()
	}()
	return nil
}

func (w *WsBroker) HelpMessage (){
	hstr := fmt.Sprint("It's a help message")
	readyPayload := map[string]interface{}{
		"action":         "system",
		"text":  			hstr,
		"messageId" : fmt.Sprint(time.Now().UnixMilli()),
	}
	payloadJSON, err := json.Marshal(readyPayload)
	if err != nil {
		fmt.Println("marshal ready payload:", err)
	}
	w.SendWsMessage(payloadJSON)
}

func (w *WsBroker) SendWsMessage (src []byte){
	reqBody := SendAppDataParams{
		Address:   address,
		ProtocolId: w.config.ClientId,
		Data:      string(src),
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		log.Printf("marshal devService request error: %v \n\n", err)
	}
	if raw, ok := w.conns.Load(address); ok {
		if conn, ok := raw.(*websocket.Conn); ok {
			w.mu.Lock()
			err = conn.WriteMessage(websocket.TextMessage, jsonBody)
			if err != nil {
				fmt.Println("Error WS Send ", err)
			}
			w.mu.Unlock()
		}
	}
}