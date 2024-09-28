package main

import "gopkg.in/yaml.v3"

var exampleConfig Config = Config{
	Sleep:    60,
	LogLevel: "info",
	Watchers: []Watcher{
		{
			Record:  "a.example.com",
			Command: "/bin/bash",
			Args:    []string{"echo", "example.com has changed!"},
		},
		{
			Record:  "b.example.com",
			Command: "/bin/bash",
			Args:    []string{"systemctl", "restart", "nginx"},
		},
	},
}

func renderExampleConfig() ([]byte, error) {
	return yaml.Marshal(exampleConfig)
}
