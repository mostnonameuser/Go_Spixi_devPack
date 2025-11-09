package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
)
func NewConfig () *Config  {
	var conf = new(Config)
	conf.GetConfig()
	return conf
}
func (c *Config) GetConfig()  {
	var err error
	wd, err := os.Getwd()
	if err != nil {
		fmt.Errorf("failed to get working directory: %w", err)
	}
	configFile := filepath.Join(wd, "config.json")
	for {
		func() {
			file, cerr := os.Open(configFile)
			if cerr != nil {
				fmt.Errorf("check if config.json exist: %w", cerr)
			}
			decoder := json.NewDecoder(file)
			err = decoder.Decode(&c)
			if err != nil {
				fmt.Errorf("error decode config.json: %w", err)
			}
			file.Close()
		}()
		time.Sleep(1*time.Second)
		if err == nil {
			break
		}
	}
}
func (c *Config) WriteConfig() {
	var err error
	wd, err := os.Getwd()
	if err != nil {
		fmt.Errorf("failed to get working directory: %w", err)
	}
	cfile := filepath.Join(wd, "config.json")
	for {
		func() {
			//defer Mu.Unlock()
			file, err := json.MarshalIndent(c, "", " ")
			if err != nil {
				fmt.Errorf("ConfigParseWriteError: %w", err)
			}
			if runtime.GOOS == "windows" {
				err = os.WriteFile(cfile, file, 0644)
				if err != nil {
					fmt.Errorf("ConfigWriteError: %w", err)
				}
			} else {
				err = os.WriteFile(cfile, file, 0644)
				if err != nil {
					fmt.Errorf("ConfigWriteError: %w", err)
				}
			}
			time.Sleep(1 * time.Second)
		}()
		if err == nil {
			break
		}
	}
}
