package main

import (
	"encoding/json"
	"os"
)

type ConfigFile struct {
	Repos   map[string][]*columnDef `json:"repos"`
	Default []*columnDef            `json:"default"`
}

func ReadConfig(path string) *ConfigFile {
	f, err := os.Open(os.ExpandEnv(path))
	if err != nil {
		return nil
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	var conf ConfigFile
	err = dec.Decode(&conf)
	if err != nil {
		log.Fatalf("Invalid config file: %#v", err)
		return nil
	}

	if conf.Default == nil {
		conf.Default = []*columnDef{}
	}
	return &conf
}
