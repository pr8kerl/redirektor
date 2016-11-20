package redirektor

import (
	"gopkg.in/gcfg.v1"
)

type Config struct {
	Server ServerSection
	Db     map[string]*DBSection
}

type ServerSection struct {
	BindAddr string
	HtmlPath string
}
type DbSection struct {
	DbAddr     string
	DbID       int
	DbPassword string
	KeyPrefix  string
}

func (c *Config) LoadFromFile(file string) error {

	err := gcfg.ReadFileInto(c, file)
	if err != nil {
		return nil, err
	}

	if c.Server.BindAddr == "" {
		c.Server.BindAddr = "localhost"
	}
	if c.Server.BindPort == 0 {
		c.Server.BindPort = 5000
	}
	if c.Server.HtmlPath == "" {
		c.Server.HtmlPath = "/public"
	}

	return nil

}
