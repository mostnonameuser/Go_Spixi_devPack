package config

type Config struct {
	EnableDevServer bool   `json:"enable_dev_server"`
	Debug 			bool `json:"debug"`
	DevServerWs     string `json:"dev_server_ws"`
	DevServerWeb     string `json:"dev_server_web"`
	QuixiApi        string `json:"quixi_api"`
	Mqtt         string `json:"mqtt"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	ClientId      string `json:"client_id"`
}

