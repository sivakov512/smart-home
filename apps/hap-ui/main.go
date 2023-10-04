package main

import (
	"github.com/brutella/hap"

	"context"
	"hap-ui/air_conditioner_v2"
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

	config, err := airconditionerv2.LoadConfig(fpath)
	if err != nil {
		panic(err)
	}

	handler := airconditionerv2.NewHandler(config)

	fs := hap.NewFsStore("./db")

	server, err := hap.NewServer(fs, handler.HAPAccessory.A)
	if err != nil {
		log.Panic(err)
	}
	server.Pin = "11122333"

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
