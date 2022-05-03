package admin

type Config struct {
	Host string `json:"host"`
	Port string `json:"port"`
	Tls  *struct {
		CertFile string `json:"cert_file"`
		KeyFile  string `json:"key_file"`
	} `json:"tls"`
}
