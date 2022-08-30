package configuration

type ServerConfig struct {
	Ip   string `json:"ip"`
	Port string `json:"port"`
	Tick int    `json:"tick"`
}
