package main

import (
	"gopkg.in/gcfg.v1"
)

type Config struct {
	Profile map[string]*ConfigSection
}

type ConfigSection struct {
	DBAddr     string
	DBID       int
	DBPassword string
	KeyPrefix  string
}

func initialiseConfig(file string) (*ConfigSection, error) {

	var config = Config{}
	var sectn = ConfigSection{}
	profilename := basename

	err := gcfg.ReadFileInto(&config, file)
	if err != nil {
		return nil, err
	}

	if config.Profile[profilename].DBAddr == "" {
		sectn.DBAddr = "localhost:6379"
	} else {
		sectn.DBAddr = config.Profile[profilename].DBAddr
	}
	if config.Profile[profilename].DBID != 0 {
		sectn.DBID = config.Profile[profilename].DBID
	}
	if config.Profile[profilename].DBPassword != "" {
		sectn.DBPassword = config.Profile[profilename].DBPassword
	}
	if config.Profile[profilename].KeyPrefix == "" {
		sectn.KeyPrefix = basename
	} else {
		sectn.KeyPrefix = config.Profile[profilename].KeyPrefix
	}

	return &sectn, nil

}
