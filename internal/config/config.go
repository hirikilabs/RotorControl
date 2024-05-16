package config

import (
	"encoding/json"
    "io/ioutil"
)

type Config struct {
	RotorModel string `json:"rotor_model"`
	Device     string `json:"device"`
	ServerAddr string `json:"server_addr"`
}

func (c *Config) Load(filename string) error {
	content, err := ioutil.ReadFile(filename)
    if err != nil {
		return err
    }

	err = json.Unmarshal(content, c)
    if err != nil {
        return err
    }

	return nil
}
