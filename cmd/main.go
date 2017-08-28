package main

import (
	"log"
	"fmt"
	"github.com/zpatrick/go-config"
	"github.com/qframe/collector-docker-events"
	"github.com/qframe/types/qchannel"
)

const (
	dockerHost = "unix:///var/run/docker.sock"
	dockerAPI = "v1.31"
)


func main() {
	qChan := qtypes_qchannel.NewQChan()
	qChan.Broadcast()
	cfgMap := map[string]string{
		"log.level": "info",
		"collector.tcp.port": "10001",
		"collector.tcp.docker-host": "unix:///var/run/docker.sock",
		"filter.inventory.inputs": "docker-events",
		"filter.inventory.ticker-ms": "2500",	}

	cfg := config.NewConfig(
		[]config.Provider{
			config.NewStatic(cfgMap),
		},
	)
	// Start docker-events
	inv := qcollector_docker_events.NewInventory()
	pde, err := qcollector_docker_events.New(qChan, inv, *cfg, "docker-events")
	if err != nil {
		log.Printf("[EE] Failed to create collector: %v", err)
		return
	}
	go pde.Run()
	bg := qChan.Data.Join()
	done := false
	for {
		select {
		case val := <- bg.Read:
			switch val.(type) {
			default:
				fmt.Printf("%v\n", val)
			}
		}
		if done {
			break
		}
	}
}

