package main

import (
	"gopkg.in/gcfg.v1"
)

type Config struct {
	Server ServerSection
	Db     map[string]*DbSection
}

type ServerSection struct {
	BindAddr string
	BindPort int
	HtmlPath string
}
type DbSection struct {
	DbAddr     string
	DbID       int
	DbPassword string
	Prefix     []string
}

func (c *Config) LoadFromFile(file string) error {

	err := gcfg.ReadFileInto(c, file)
	if err != nil {
		return err
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
