package main

import (
	"github.com/brutella/hap"

	"context"
	"hap-ui/air_conditioner"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config, err := airconditioner.LoadConfig("./config.toml")
	if err != nil {
		panic(err)
	}

	handler := airconditioner.NewHandler(config)

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
