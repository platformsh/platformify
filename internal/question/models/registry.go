package models

import (
	_ "embed"
	"encoding/json"
	"log"
)

//go:embed registry.json
var registry []byte

var Runtimes RuntimeList

var ServiceNames ServiceNameList

func init() {
	allRuntimes := map[string]*Runtime{}
	if err := json.Unmarshal(registry, &allRuntimes); err != nil {
		log.Fatal(err)
	}
	for _, r := range allRuntimes {
		if r.Runtime {
			Runtimes = append(Runtimes, r)
		}
	}

	allServices := map[string]*ServiceName{}
	if err := json.Unmarshal(registry, &allServices); err != nil {
		log.Fatal(err)
	}
	for _, s := range allServices {
		if !s.Runtime {
			ServiceNames = append(ServiceNames, s)
		}
	}
}
