package network

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mostnonameuser/Go_Spixi_devPack/internal/config"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)


var wg sync.WaitGroup
type MQTTBroker struct {
	client   mqtt.Client
	config   *config.Config
	handlers map[string]MessageHandlerFunc
	qUixiLink string
	messageChan chan Message
}
func (b *MQTTBroker) Start (ctx context.Context) error {
	fmt.Println("MQTT MODE")
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				log.Println("GetMessages stopped due to cancellation")
				return
			case msg, ok := <- b.messageChan:
				if !ok {
					log.Println("Message queue closed")
					return
				}
				switch msg.Topic {
				case "Chat":
					log.Println("InChatSpixiMessage",msg.Topic)
					continue
				case "RequestAdd2":
					log.Println(msg.Topic)
					b.AcceptContact(msg)
					continue
				case "FriendStatusUpdate":
					log.Println(msg.Topic)
					continue
				case "AppRequest":
					log.Println(msg.Topic)
					continue
				case "AppData":
					log.Println(msg.Topic)
					continue
				case "AppProtocolData":
					log.Println(msg.Topic)
					b.ProceedAppMessage(msg)
					continue
				case "AppEndSession":
					log.Println(msg.Topic)
					continue
				case "SentFunds":
					log.Println(msg.Topic)
					continue
				case "TransactionStatusUpdate":
					log.Println(msg.Topic)
					continue
				default:
					log.Println("Unknown topic", msg.Topic)
					continue
				}
			}
		}
	}()
	return nil
}

func (b *MQTTBroker) ProceedAppMessage(msg Message){
	body := QuixiAnswer{}
	_ = json.Unmarshal(msg.Payload, &body)
	decodedBytes, errDbytes := base64.StdEncoding.DecodeString(body.Data.Data.Data)
	if errDbytes != nil {
		log.Println("Error decodeBytes ProceedAppMessage ", errDbytes)
	}
	mesg := &AiMessage{}
	erss := json.Unmarshal(decodedBytes, &mesg)
	if erss != nil {
		log.Println("Error Unmarshal Message ProceedAppMessage ", erss)
	}
	jsonBody := MessageProcessor(mesg)
	if jsonBody != nil {
		postIxiErr := b.Publish( "sendAppData", body.Sender.Base58Address, jsonBody)
		if postIxiErr != nil {
			log.Printf("Processing context messages MQTT answer error %v\n", postIxiErr)
		}
	}
}
func (b *MQTTBroker) AcceptContact (msg Message){
	body := QuixiAnswer{}
	_ = json.Unmarshal(msg.Payload, &body)
	baseURL := fmt.Sprintf("%s/acceptContact", b.qUixiLink)
	params := url.Values{}
	params.Add("address", body.Sender.Base58Address)
	fullURL := baseURL + "?" + params.Encode()
	if  sendGetToQuIXI(fullURL) {
		log.Println("Accepted Request2 contact ", body.Sender.Base58Address)
	}
}
func (b *MQTTBroker) Connect() error {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", b.config.Mqtt))
	b.client = mqtt.NewClient(opts)
	if token := b.client.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to connect: %v", token.Error())
	}
	b.messageChan = make(chan Message, 100)
	b.client.Subscribe("#", 0, func(client mqtt.Client, msg mqtt.Message) {
		b.messageChan <- Message{
			Topic:   msg.Topic(),
			Payload: msg.Payload(),
		}
	})
	return nil
}

func (b *MQTTBroker) Subscribe(topic string, handler MessageHandlerFunc) error {
	log.Printf("Subscribed to topic: %s", topic)
	b.handlers[topic] = handler
	return nil
}

func (b *MQTTBroker) GetMessageChannel() <-chan Message {
	return b.messageChan
}

func (b *MQTTBroker) Publish (method, addr string, payload []byte) error {
	reqBody := SendAppDataRequest{
		Method: method,
		Params: SendAppDataParams{
			Address:   addr,
			ProtocolId: b.config.ClientId,
			Data:      string(payload),
		},
	}
	jsonBody, err := json.Marshal(reqBody)
	baseURL := fmt.Sprintf("%s/%s", b.qUixiLink, method)
	req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("create HTTP request: %w", err)
	}
	req.Header.Set("User-Agent", "curl/8.5.0")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			DisableCompression: true,
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("QuIXI returned status %d", resp.StatusCode)
	}

	return nil
}

func (b *MQTTBroker) Disconnect() error {
	b.client.Disconnect(250)
	return nil
}

func sendGetToQuIXI(command string) bool  {
	resp, err := http.Get(command)
	if err != nil {
		fmt.Printf("Ошибка запроса: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	_ , err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Ошибка чтения ответа: %v\n", err)
		return false
	}
	return true
}