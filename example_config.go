package main

import "gopkg.in/yaml.v3"

var exampleConfig Config = Config{
	Sleep:    60,
	LogLevel: "info",
	Watchers: []Watcher{
		{
			Records: []string{
				"a.example.com",
				"b.example.com",
			},
			Command: "/bin/bash",
			Args:    []string{"echo", "example.com has changed!"},
		},
		{
			Records: []string{"c.example.com"},
			Command: "/bin/bash",
			Args:    []string{"systemctl", "restart", "nginx"},
		},
	},
}

func renderExampleConfig() ([]byte, error) {
	return yaml.Marshal(exampleConfig)
}
