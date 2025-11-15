package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

func NewConfig() *Config {
	var conf = new(Config)
	conf.GetConfig()
	return conf
}
func (c *Config) GetConfig() {
	var err error
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("failed to get working directory: %w", err)
		return
	}
	configFile := filepath.Join(wd, "config.json")
	for {
		func() {
			file, cerr := os.Open(configFile)
			if cerr != nil {
				fmt.Println("check if config.json exist: ", cerr)
				return
			}
			decoder := json.NewDecoder(file)
			err = decoder.Decode(&c)
			if err != nil {
				fmt.Println("error decode config.json: ", err)
				return
			}
			err := file.Close()
			if err != nil {
				fmt.Println("error closing config.json: ", err)
				return
			}
		}()
		time.Sleep(1 * time.Second)
		if err == nil {
			break
		}
	}
}
func (c *Config) WriteConfig() {
	var err error
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("failed to get working directory: ", err)
	}
	cfile := filepath.Join(wd, "config.json")
	for {
		func() {
			file, err := json.MarshalIndent(c, "", " ")
			if err != nil {
				fmt.Println("ConfigParseWriteError: ", err)
				return
			}
			if runtime.GOOS == "windows" {
				err = os.WriteFile(cfile, file, 0644)
				if err != nil {
					fmt.Println("ConfigWriteError: ", err)
					return
				}
			} else {
				err = os.WriteFile(cfile, file, 0644)
				if err != nil {
					fmt.Println("ConfigWriteError: ", err)
					return
				}
			}
			time.Sleep(1 * time.Second)
		}()
		if err == nil {
			break
		}
	}
}
