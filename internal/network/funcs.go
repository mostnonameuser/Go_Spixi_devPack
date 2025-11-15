package network

import (
	"encoding/json"
	"fmt"
	"github.com/mostnonameuser/Go_Spixi_devPack/internal/config"
	"sync"
)

func NewMQQTBroker (conf *config.Config) Broker {
	return &MQTTBroker{
			qUixiLink: fmt.Sprintf("http://%s", conf.QuixiApi),
			config:   conf,
			handlers: make(map[string]MessageHandlerFunc),
	}
}
func NewWsBroker (conf *config.Config) Broker {
	return &WsBroker{
		config:   conf,
		conns: &sync.Map{},
	}
}
func MessageProcessor (msg *AiMessage) []byte {
	switch msg.Action {
	case "ping":
		fmt.Println("Ping received")
		readyPayload := map[string]interface{}{
			"action":       "debug",
		}
		payloadJSON, err := json.Marshal(readyPayload)
		if err != nil {
			fmt.Println("marshal ready payload:", err)
		}
		return payloadJSON
	case "sendMessage":
		fmt.Println("SendMessageDetect")
		hstr := fmt.Sprintf("DemoMessage, u wrote : %s \n", msg.Text)
		readyPayload := map[string]interface{}{
			"action":         "response",
			"sender":   "ai",
			"message":  hstr,
			"messageId" : msg.MessageId,
		}
		payloadJSON, err := json.Marshal(readyPayload)
		if err != nil {
			fmt.Println("marshal answer ready payload:", err)
		}
		return payloadJSON
	default:
		fmt.Println("Some other activity: ", msg.Action)
		return nil
	}
}