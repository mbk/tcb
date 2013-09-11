package config

import (
	"bufio"
	"fmt"
	"os"
	str "strings"
)

type config struct {
	inner map[string]string
}

type Configurable interface {
	GetOrElse(section, key, otherwise string) string
}

var configImpl = new(config)

func GetOrElse(key, otherwise string) string {
	value := configImpl.inner[key]
	if value == "" {
		return otherwise
	} else {
		return value
	}
}

func EnsureConfiguration(path string) {
	if configImpl.inner == nil {
		infile, err := os.Open(path)
		if err != nil {
			fmt.Println("Error reading config file, exiting: " + err.Error())
			os.Exit(-1)
		} else {
			configImpl.inner = make(map[string](string))
			scanner := bufio.NewScanner(infile)
			for scanner.Scan() {
				line := scanner.Text()
				if !(str.TrimSpace(line) == "") {
					if !str.HasPrefix(line, "#") {
						vals := str.Split(line, "=")
						if len(vals) >= 2 {
							configImpl.inner[vals[0]] = vals[1]
						}
					}
				}
			}
		}
	}
}
