package main

import (
	"context"
	"github.com/brutella/hap"
	"github.com/eclipse/paho.mqtt.golang"
	"hap-ui/ac"
	"hap-ui/config"
	"hap-ui/heater"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const (
	HAPUICONFIG_ENV_KEY       = "HAPUICONFIG"
	HAPUICONFIG_DEFAULT_FPATH = "./config.toml"
)

func main() {
	fpath, exist := os.LookupEnv(HAPUICONFIG_ENV_KEY)
	if !exist {
		fpath = HAPUICONFIG_DEFAULT_FPATH
	}

	config, err := config.LoadConfig(fpath)
	if err != nil {
		panic(err)
	}

	mqttClient := mqtt.NewClient(mqtt.NewClientOptions().AddBroker(config.Broker))
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	ac := ac.NewHandler(config.AC, mqttClient)
	heater := heater.NewHandler(config.Heater, mqttClient)

	server, err := hap.NewServer(hap.NewFsStore("./db"), ac.HAPAccessory.A, heater.HAPAccessory.A)
	if err != nil {
		log.Panic(err)
	}
	server.Pin = config.PIN

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-c
		signal.Stop(c)
		cancel()
	}()

	server.ListenAndServe(ctx)
}
