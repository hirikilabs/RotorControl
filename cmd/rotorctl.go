package main

import (
	"log"
	"rotorctl/internal/rotors"
	"rotorctl/internal/rotserver"
	"rotorctl/internal/config"
	"rotorctl/pkg/rot2prog"
)

func main() {
	log.Println("▶️  Starting rotor control...")

	// load configuration
	config := config.Config{}
	err := config.Load("config.json")
	if err != nil {
		log.Fatal("Problem with config:", err)
	}

	// interface for rotor models
	var rotor rotors.Rotor
	
	// create rotor
	switch config.RotorModel {
	case "rot2prog":
		rotor = &rot2prog.Rot2Prog{Device: config.Device}
		err = rotor.Init()
		if err != nil {
			log.Fatal("Problem with the rotor interface: ", err)
		}
	default:
		log.Fatal("Unrecognized rotor model")
		
	}

	// start TCP server
	server := rotserver.NewRotServer(config.ServerAddr, rotor)
	log.Fatal(server.Start())
}


