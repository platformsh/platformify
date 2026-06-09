package models

import (
	_ "embed"
	"encoding/json"
	"log"
	"sort"
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
	sort.Slice(Runtimes, func(i, j int) bool {
		return Runtimes[i].Name < Runtimes[j].Name
	})

	allServices := map[string]*ServiceName{}
	if err := json.Unmarshal(registry, &allServices); err != nil {
		log.Fatal(err)
	}
	for _, s := range allServices {
		if !s.Runtime {
			ServiceNames = append(ServiceNames, s)
		}
	}
	sort.Slice(ServiceNames, func(i, j int) bool {
		return ServiceNames[i].Name < ServiceNames[j].Name
	})
}
