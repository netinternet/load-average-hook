package main

import (
	"net/url"
	"os"
	"strings"

	"gopkg.in/ini.v1"
)

func init() {
	config := "/etc/loadhook.conf"
	if _, err := os.Stat(config); os.IsNotExist(err) {
		config = "./loadhook.conf"
	}
	if _, err := os.Stat(config); os.IsNotExist(err) {
		panic("configuration file missing!")
	}
	c, err := ini.Load(config)
	if err != nil {
		panic("Fail to read file: " + err.Error())
	}
	interval, err = c.Section("").Key("interval").Int64()
	if interval <= 0 {
		interval = int64(defaultInterval)
	}
	standby, err = c.Section("").Key("standby").Int64()
	if standby <= 0 {
		standby = int64(defaultStandby)
	}
	loadLimit, err = c.Section("").Key("load_limit").Float64()
	if loadLimit <= 0 {
		loadLimit = defaultLoadLimit
	}
	webhook = c.Section("").Key("webhook").String()
	if _, err := url.ParseRequestURI(webhook); err != nil {
		panic("Configuration Error:" + err.Error())
	}
	method = c.Section("").Key("method").String()
	method = strings.ToUpper(method)
	var validMethod bool
	validMethods := []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
	for _, m := range validMethods {
		if method == m {
			validMethod = true
			break
		}
	}
	if !validMethod {
		panic("Configuration Error: http method is not correct.")
	}
	insecure, _ = c.Section("").Key("insecure").Bool()
}
