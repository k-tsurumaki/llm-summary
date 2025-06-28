package main

import (
	"log"

	"gopkg.in/go-ini/ini.v1"
)

type ConfigList struct {
	PORT         string
	LLM_API_KEY  string
	LLM_BASE_URL string
}

var Config ConfigList

func init() {
	LoadConfig()
}

func LoadConfig() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatalf("Failed to load config file: %v", err)
	}
	Config = ConfigList{
		PORT:         cfg.Section("web").Key("port").String(),
		LLM_API_KEY:  cfg.Section("llm").Key("api_key").String(),
		LLM_BASE_URL: cfg.Section("llm").Key("base_url").String(),
	}
}
