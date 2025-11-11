package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mostnonameuser/Go_Spixi_devPack/internal/config"
	"github.com/mostnonameuser/Go_Spixi_devPack/internal/network"
)

type Service struct {
	conf         *config.Config
	Listener     *MQTTBroker
	DevService   *network.DevService
}

type Message struct {
	Topic   string
	Payload []byte
}
type MQTTBroker struct {
	client   mqtt.Client
	config   *config.Config
	handlers map[string]MessageHandlerFunc
	qUixiLink string
	devService *network.DevService
}

type MessageHandlerFunc func(client mqtt.Client, msg mqtt.Message)

type DirectChatMessage struct {
	Type       int         `json:"type"`
	RealSender interface{} `json:"realSender"`
	Sender     struct {
		Version             int         `json:"version"`
		AddressWithChecksum string      `json:"addressWithChecksum"`
		Base58Address       string      `json:"base58Address"`
		AddressNoChecksum   string      `json:"addressNoChecksum"`
		SectorPrefix        string      `json:"sectorPrefix"`
		Nonce               interface{} `json:"nonce"`
		PubKey              interface{} `json:"pubKey"`
	} `json:"sender"`
	Recipient struct {
		Version             int         `json:"version"`
		AddressWithChecksum string      `json:"addressWithChecksum"`
		Base58Address       string      `json:"base58Address"`
		AddressNoChecksum   string      `json:"addressNoChecksum"`
		SectorPrefix        string      `json:"sectorPrefix"`
		Nonce               interface{} `json:"nonce"`
		PubKey              interface{} `json:"pubKey"`
	} `json:"recipient"`
	Data struct {
		Type    int    `json:"type"`
		Channel int    `json:"channel"`
		Data    string `json:"data"`
	} `json:"data"`
	OriginalData           string      `json:"originalData"`
	OriginalChecksum       string      `json:"originalChecksum"`
	Signature              interface{} `json:"signature"`
	EncryptionType         int         `json:"encryptionType"`
	Encrypted              bool        `json:"encrypted"`
	Id                     string      `json:"id"`
	Timestamp              int         `json:"timestamp"`
	RequireRcvConfirmation bool        `json:"requireRcvConfirmation"`
	Version                int         `json:"version"`
}

type QuixiAnswer struct {
	Type       int         `json:"type"`
	RealSender interface{} `json:"realSender"`
	Sender     struct {
		Version             int         `json:"version"`
		AddressWithChecksum string      `json:"addressWithChecksum"`
		Base58Address       string      `json:"base58Address"`
		AddressNoChecksum   string      `json:"addressNoChecksum"`
		SectorPrefix        string      `json:"sectorPrefix"`
		Nonce               interface{} `json:"nonce"`
		PubKey              interface{} `json:"pubKey"`
	} `json:"sender"`
	Recipient struct {
		Version             int         `json:"version"`
		AddressWithChecksum string      `json:"addressWithChecksum"`
		Base58Address       string      `json:"base58Address"`
		AddressNoChecksum   string      `json:"addressNoChecksum"`
		SectorPrefix        string      `json:"sectorPrefix"`
		Nonce               interface{} `json:"nonce"`
		PubKey              interface{} `json:"pubKey"`
	} `json:"recipient"`
	Data struct {
		Type    int `json:"type"`
		Channel int `json:"channel"`
		Data    struct {
			SessionId string          `json:"sessionId"`
			Data      string          `json:"data"`
			AppId     string 		  `json:"appId"`
			MaxProtocolVersion int    `json:"maxProtocolVersion"`
			PubKey             string `json:"pubKey"`
		} `json:"data"`
	} `json:"data"`
	OriginalData           interface{} `json:"originalData"`
	OriginalChecksum       interface{} `json:"originalChecksum"`
	Signature              string      `json:"signature"`
	EncryptionType         int         `json:"encryptionType"`
	Encrypted              bool        `json:"encrypted"`
	Id                     string      `json:"id"`
	Timestamp              int         `json:"timestamp"`
	RequireRcvConfirmation bool        `json:"requireRcvConfirmation"`
	Version                int         `json:"version"`
}


