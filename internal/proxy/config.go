package proxy

type Config struct {
	Host      string         `json:"host"`
	Port      string         `json:"port"`
	PointsDir string         `json:"points_dir"`
	Points    []*PointConfig `json:"points"`
	Tls       *struct {
		CertFile string `json:"cert_file"`
		KeyFile  string `json:"key_file"`
	} `json:"tls"`
}

type PointConfig struct {
	Code  string `json:"code"` // Point code for identification
	Match struct {
		Host string `json:"host"` // Host - regex
		Path string `json:"path"` // Path - regex
	} `json:"match"` // Pattern for identify proxy to point
	TokenHeader                 string   `json:"token_header"`                    // Header that contains token for authentication
	Destination                 string   `json:"destination"`                     // Proxy destination - URL
	Providers                   []string `json:"providers"`                       // Authorization providers see internal/auth/providers
	ProxyHeaders                []string `json:"proxy_headers"`                   // Headers that we proxy to destination
	ProxyAppendedDataHeaderName string   `json:"proxy_appended_data_header_name"` // In this header we add additional info like username, roles and etc.
	IsHeaderCrypt               bool     `json:"is_header_crypt"`                 // Crypt headers in AES
	HeaderCryptMasterPass       string   `json:"header_crypt_master_pass"`        // Master password for crypt headers if IsHeaderFlag is enabled
	SignHeaderToken             string   `json:"sign_header_token"`               // Header name for sign (hash secret token and request data)
	Secret                      string   `json:"secret"`                          // Secret token. With this token we sign data for proxy
}
