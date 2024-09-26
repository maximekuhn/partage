package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/maximekuhn/partage/internal/app/web"
	"gopkg.in/yaml.v3"
)

func main() {
	cfgPath := flag.String("config", "", "Configuration file path")
	flag.Parse()

	var config web.ServerConfig
	cfg := *cfgPath
	if cfg == "" {
		config = web.DefaultServerConfig()
	} else {
		cf, err := parseConfigFile(cfg)
		if err != nil {
			panic(err)
		}
		config = *cf
	}

	server, err := web.NewServer(config)
	if err != nil {
		panic(err)
	}

	log.Fatal(server.Run())
}

type configFile struct {
	DBFilepath string `yaml:"DBFilepath"`
}

func parseConfigFile(path string) (*web.ServerConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	fmt.Printf("data: %s\n", string(data))

	var cf configFile
	if err := yaml.Unmarshal(data, &cf); err != nil {
		return nil, err
	}

	return &web.ServerConfig{
		DBFilepath: cf.DBFilepath,
	}, nil
}
