package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mostnonameuser/Go_Spixi_devPack/internal/config"
	"github.com/mostnonameuser/Go_Spixi_devPack/internal/network"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)


var messageQueue = make(chan Message, 100) // Буферизованная очередь из 100 сообщений
var wg sync.WaitGroup

func (b *MQTTBroker) Connect() error {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", b.config.Mqtt))
	//opts.SetClientID(b.config.ClientID)

	if b.config.Username != "" && b.config.Password != "" {
		opts.SetUsername(b.config.Username)
		opts.SetPassword(b.config.Password)
	}

	b.client = mqtt.NewClient(opts)

	if token := b.client.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to connect: %v", token.Error())
	}

	log.Println("Connected to MQTT broker")
	return nil
}

func (b *MQTTBroker) Subscribe(topic string, handler mqtt.MessageHandler) error {
	if token := b.client.Subscribe(topic, 0, handler); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to subscribe: %v", token.Error())
	}

	log.Printf("Subscribed to topic: %s", topic)
	return nil
}

func (b *MQTTBroker) Publish(topic string, payload interface{}) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	token := b.client.Publish(topic, 0, false, jsonPayload)
	token.Wait()

	if token.Error() != nil {
		return fmt.Errorf("failed to publish: %v", token.Error())
	}

	return nil
}

func (b *MQTTBroker) Disconnect() {
	if b.client.IsConnected() {
		b.client.Disconnect(250)
	}
}

func handleMessage(client mqtt.Client, msg mqtt.Message) {
	messageQueue <- Message{
		Topic:   msg.Topic(),
		Payload: msg.Payload(),
	}
}
func (b *MQTTBroker) GetMessages (ctx context.Context){
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				log.Println("GetMessages stopped due to cancellation")
				return
			case msg, ok := <-messageQueue:
				if !ok {
					log.Println("Message queue closed")
					return
				}
				switch msg.Topic {
				case "Chat":
					b.ProceedDirectWalletMessage(msg)
					continue
				case "RequestAdd2":
					b.AcceptContact(msg)
					continue
				case "FriendStatusUpdate":
					continue
				case "AppRequest":
					continue
				case "AppData":
					continue
				case "AppProtocolData":
					continue
				case "AppEndSession":
					continue
				case "SentFunds":
					continue
				case "TransactionStatusUpdate":
					continue
				default:
					log.Println("Unknown topic", msg.Topic)
					continue
				}
				continue
			}
		}
	}()

}

func (b *MQTTBroker) ProceedDirectWalletMessage (msg Message){
	body := DirectChatMessage{}
	_ = json.Unmarshal(msg.Payload, &body)
	answer := ""
	switch body.Data.Data {
	case "/help" :
		answer = fmt.Sprint("Help Message\n")
	default:
		answer = fmt.Sprint("Type /help to get info.")
	}
	baseURL := fmt.Sprintf("%s/sendChatMessage", b.qUixiLink)
	params := url.Values{}
	params.Add("channel", fmt.Sprint(body.Data.Channel))
	params.Add("address", body.Sender.Base58Address)
	params.Add("message", answer)
	fullURL := baseURL + "?" + params.Encode()
	if  sendToQuIXI(fullURL) {
		fmt.Println("Proceed direct wallet message")
	}
}
func (b *MQTTBroker) ProceedAppMessage(msg Message){
	body := QuixiAnswer{}
	_ = json.Unmarshal(msg.Payload, &body)
	decodedBytes, errDbytes := base64.StdEncoding.DecodeString(body.Data.Data.Data)
	if errDbytes != nil {
		log.Println("Error decodeBytes ProceedAppMessage ", errDbytes)
	}
	mesg := network.AiMessage{}
	erss := json.Unmarshal(decodedBytes, &mesg)
	if erss != nil {
		log.Println("Error Unmarshal Message ProceedAppMessage ", erss)
	}
	switch mesg.Action {
	case "ping":
		fmt.Println("Ping received")
	default:
		fmt.Println("Some other activity: ", string(decodedBytes))
	}

}
func (b *MQTTBroker) AcceptContact (msg Message){
	body := QuixiAnswer{}
	_ = json.Unmarshal(msg.Payload, &body)
	baseURL := fmt.Sprintf("%s/acceptContact", b.qUixiLink)
	params := url.Values{}
	params.Add("address", body.Sender.Base58Address)
	fullURL := baseURL + "?" + params.Encode()
	if  sendToQuIXI(fullURL) {
		log.Println("Accepted Request2 contact ", body.Sender.Base58Address)
	}
}
func sendToQuIXI (command string) bool  {
	resp, err := http.Get(command)
	if err != nil {
		fmt.Printf("Error request sendToQuIX: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	_ , err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error read sendToQuIXI answer: %v\n", err)
		return false
	}
	return true
}
func postToQuIXI(url string, jsonBody []byte) error {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
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
func gracefulShutdown(cancel context.CancelFunc) {
	cancel()
	close(messageQueue)
	wg.Wait()
	log.Println("All workers stopped gracefully")
}

func NewMQTTBroker(config *config.Config, dserv *network.DevService) *MQTTBroker {
	return &MQTTBroker{
		config:   config,
		handlers: make(map[string]MessageHandlerFunc),
		qUixiLink: fmt.Sprintf("http://%s", config.QuixiApi),
		devService: dserv,
	}
}
func NewDevService(conf *config.Config) *network.DevService  {
	return &network.DevService{
		Conf: conf,
		Conns: &sync.Map{},
	}
}
func NewService(conf *config.Config) *Service {
	dserv := NewDevService(conf)
	return &Service{
		conf:   conf,
		Listener: NewMQTTBroker( conf, dserv),
		DevService:  dserv ,
	}
}


