package internal

import (
	"io/fs"
	"lproxy/internal/admin"
	"lproxy/internal/auth"
	"lproxy/internal/proxy"
	"lproxy/pkg/file"
	"lproxy/pkg/log"
	"path/filepath"
)

type Configuration struct {
	Logs struct {
		File    *log.FileLoggerConfig   `json:"file"`
		Console *log.StdoutLoggerConfig `json:"console"`
	} `json:"logs"`
	Proxy *proxy.Config `json:"proxy"`
	Auth  *auth.Config  `json:"auth"`
	Admin *admin.Config `json:"admin"`
}

func ParseConfiguration(main string) (*Configuration, error) {
	config := &Configuration{}
	err := file.LoadStructureFromJsonFile(main, config)
	if err != nil {
		return nil, err
	}
	var dir string
	if !filepath.IsAbs(config.Proxy.PointsDir) {
		mainFileDir := filepath.Dir(main)
		dir = mainFileDir + "/" + config.Proxy.PointsDir
	} else {
		dir = config.Proxy.PointsDir
	}

	err = filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		point := &proxy.PointConfig{}
		if !info.IsDir() {
			unmarshalErr := file.LoadStructureFromJsonFile(path, point)
			if unmarshalErr == nil {
				config.Proxy.Points = append(config.Proxy.Points, point)
			}
		}

		return nil
	})

	return config, err
}
