package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/kardianos/service"
	"github.com/shirou/gopsutil/load"
)

const (
	serviceName        = "Load Avarage Hook service"
	serviceDescription = "Load avarage monitor web hook service"
	defaultInterval    = 500 * time.Millisecond
	defaultStandby     = 60 * time.Second
	defaultLoadLimit   = float64(10)
)

var (
	interval  int64
	standby   int64
	webhook   string
	method    string
	insecure  bool
	loadLimit float64
)

type program struct{}

func (p program) Start(s service.Service) error {
	fmt.Println(s.String() + " started")
	go p.run()
	return nil
}

func (p program) Stop(s service.Service) error {
	fmt.Println(s.String() + " stopped")
	return nil
}

func (p program) run() {
	for {
		time.Sleep(time.Duration(interval) * time.Millisecond)
		load, err := load.Avg()
		if err != nil {
			panic(err)
		}
		if load.Load1 >= loadLimit {
			log.Println("Server load", load.Load1, " and starting web hook...")
			if err := sendWebHook(); err != nil {
				log.Println("Web hook error!")
			} else {
				log.Println("Web hook done!")
			}
			time.Sleep(time.Duration(standby) * time.Millisecond)
		}
	}
}

func sendWebHook() error {
	req, err := http.NewRequest(method, webhook, nil)
	if insecure {
		client := &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}
		_, err = client.Do(req)
		return err
	}
	_, err = http.DefaultClient.Do(req)
	return err
}
