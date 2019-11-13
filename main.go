package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/GameBench/gba-client-go"
)

func main() {
	config := &gba.Config{BaseUrl: "http://localhost:8000"}
	client := gba.New(config)

	listDevicesCommand := flag.NewFlagSet("list-devices", flag.ExitOnError)

	if len(os.Args) < 2 {
		os.Exit(1)
	}

	switch os.Args[1] {
	case "list-devices":
		listDevicesCommand.Parse(os.Args[2:])
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	if listDevicesCommand.Parsed() {
		devices, err := client.ListDevices()
		if err != nil {
			panic(err)
		}

		for _, device := range devices {
			fmt.Printf("%s\t%s\n", device.Id, device.Name)
		}
	}
}
