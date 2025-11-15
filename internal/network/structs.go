package network

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