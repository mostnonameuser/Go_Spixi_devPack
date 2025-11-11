package network

import (
	"github.com/mostnonameuser/Go_Spixi_devPack/internal/config"
	"sync"
)

type Broker struct {
	user string
	password string
	address string

}
type DevService struct {
	Conf *config.Config
	Conns  *sync.Map
	Mu sync.Mutex
}
type AiMessage struct {
	Action    string `json:"action"`
	Text      string `json:"text"`
	Chat      string `json:"chat"`
	Ts 		  int `json:"ts"`
	MessageId int `json:"messageId"`
	ChatId    int `json:"chatId"`
}
type SendAppDataRequest struct {
	Method string `json:"method"`
	Params SendAppDataParams `json:"params"`
}
type SendAppDataParams struct {
	Address     string `json:"address"`
	ProtocolId  string `json:"protocolId"`
	Data        string `json:"data"`
}