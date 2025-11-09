package network

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)
const address = "MyTestAddress"
func (ds *DevService) AddConn(id string, conn *websocket.Conn) {
	ds.Conns.Store(id, conn)
}

func (ds *DevService) RemoveConn(id string) {
	ds.Conns.Delete(id)
}

func (ds *DevService) SendToConn(id, msg string) error {
	if raw, ok := ds.Conns.Load(id); ok {
		if conn, ok := raw.(*websocket.Conn); ok {
			return conn.WriteMessage(websocket.TextMessage, []byte(msg))
		}
	}
	return fmt.Errorf("connection %s not found", id)
}
func (ds *DevService) Start(ctx context.Context) error {
	addr := ds.Conf.DevServerWs
	server := &http.Server{
		Addr: addr,
	}
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("WebSocket upgrade error: %v", err)
			return
		}
		defer func() {
			ds.Conns.Delete(address)
			conn.Close()
		}()
		ds.Conns.Store(address, conn)
		log.Println("WebSocket client connected")

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("WebSocket read error: %v", err)
				}
				break
			}

			log.Printf("Received from browser: %s", msg)
			rcvd := SendAppDataParams{}
			erss := json.Unmarshal(msg, &rcvd)
			if erss != nil {
				log.Println("Error Unmarch Initial Message from chat")
			}
			mesg := AiMessage{}
			ersA := json.Unmarshal([]byte(rcvd.Data), &mesg)
			if ersA != nil {
				log.Println("Error Unmarch Message from WsChat")
			}
			switch mesg.Action {
			case "ping":
				fmt.Println("Ping here")
				continue
			case "unlock":
				fmt.Println("Unlock")
				continue
			case "sendMessage":
				fmt.Println("SendMessageDetect")
				hstr := fmt.Sprintf("DemoMessage, u wrote : %s \n", mesg.Text)
				readyPayload := map[string]interface{}{
					"action":         "response",
					"sender":   "ai",
					"message":  hstr,
					"messageId" : mesg.MessageId,
				}
				payloadJSON, err := json.Marshal(readyPayload)
				if err != nil {
					fmt.Println("marshal answer ready payload:", err)
				}
				ds.SendWsMessage(payloadJSON)
				continue
			case "help":
				ds.HelpMessage() // example if app can send "help" to server
				continue
			default:
				fmt.Println(mesg.Action)
			}
			if err := conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"ack","message":"ok"}`)); err != nil {
				log.Printf("WebSocket write error: %v", err)
				break
			}
		}

		log.Println("WebSocket client disconnected")
	})
	go func() {
		<-ctx.Done()
		log.Println("Shutting down dev WebSocket server...")
		server.Shutdown(context.Background())
	}()

	log.Printf("Dev WebSocket server starting on %s", addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("dev server failed: %w", err)
	}

	return nil
}
func (ds *DevService) HelpMessage (){
	hstr := fmt.Sprint("It's a help message")
	readyPayload := map[string]interface{}{
		"action":         "system",
		"text":  hstr,
		"messageId" : fmt.Sprint(time.Now().UnixMilli()),
	}
	payloadJSON, err := json.Marshal(readyPayload)
	if err != nil {
		fmt.Println("marshal ready payload:", err)
	}
	ds.SendWsMessage(payloadJSON)
}

func (ds *DevService) SendWsMessage (src []byte){
	reqBody := SendAppDataParams{
		Address:   address,
		ProtocolId: ds.Conf.ClientId,
		Data:      string(src),
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		log.Printf("marshal devService request error: %v \n\n", err)
	}
	if raw, ok := ds.Conns.Load(address); ok {
		if conn, ok := raw.(*websocket.Conn); ok {
			ds.Mu.Lock()
			err = conn.WriteMessage(websocket.TextMessage, jsonBody)
			if err != nil {
				fmt.Println("Error WS Send ", err)
			}
			ds.Mu.Unlock()
		}
	}
}