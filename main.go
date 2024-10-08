package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
	"sort"
	"time"

	log "github.com/sirupsen/logrus"

	"gopkg.in/yaml.v3"
)

// Watcher holds the configuration for each DNS record
type Watcher struct {
	Records []string `yaml:"records"`
	Command string   `yaml:"command"`
	Args    []string `yaml:"args"`
}

// Config holds the entire configuration
type Config struct {
	Sleep    int64     `yaml:"sleep"`
	LogLevel string    `yaml:"log_level"`
	Watchers []Watcher `yaml:"watchers"`
}

func main() {
	// Load flags
	configFile := flag.String("config", "config.yaml", "Path to the config file")
	showExampleConfig := flag.Bool("example", false, "Print example config and exit")
	flag.Parse()

	if *showExampleConfig {
		data, err := renderExampleConfig()
		if err != nil {
			log.Fatalf("Error rendering example config: %v", err)
		}
		fmt.Println(string(data))
		return
	}
	config, err := LoadConfig(*configFile)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	logsLevel, err := log.ParseLevel(config.LogLevel)
	if err != nil {
		log.Fatalf("Error setting log level: %v", err)
	}
	log.SetLevel(logsLevel)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Store the last known state of each watcher
	lastState := map[int]map[string][]string{}
	doChecks(config, lastState, false)
	time.Sleep(time.Second * time.Duration(config.Sleep))

	for {
		doChecks(config, lastState, true)
		// Sleep before next check (adjust the duration as needed)
		log.Debugln("Last state:", lastState)
		time.Sleep(time.Second * time.Duration(config.Sleep))
	}
}

func doChecks(config *Config, lastState map[int]map[string][]string, noAction bool) {
	for idx, watcher := range config.Watchers {
		triggerCommand := false
		for _, record := range watcher.Records {
			resolved, err := checkDNS(record)
			sort.Strings(resolved)
			log.Debugf("Resolved DNS record %s to %v", record, resolved)

			if err != nil {
				log.Errorf("Error resolving DNS record %s: %v", record, err)
				continue
			}

			last, exists := lastState[idx][record]
			if !exists || !isEqual(last, resolved) {
				if noAction {
					log.Infof("DNS record %s changed. Executing command.", record)
					triggerCommand = true
				}

				if _, ok := lastState[idx]; !ok {
					lastState[idx] = map[string][]string{}
				}

				lastState[idx][record] = resolved
			} else {
				log.Infof("DNS record %s unchanged", record)
			}
		}
		if triggerCommand {

			if err := executeCommand(watcher.Command, watcher.Args...); err != nil {
				log.Errorf("Error executing command: %v", err)
			}
		}
	}
}

// isEqual checks if two slices of strings are equal
func isEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func checkDNS(record string) ([]string, error) {
	var resolved []string

	ips, err := net.LookupIP(record)
	if err != nil {
		return nil, err
	}
	for _, ip := range ips {
		if ipv4 := ip.To4(); ipv4 != nil {
			resolved = append(resolved, ipv4.String())
		}
	}

	return resolved, nil
}

// executeCommand runs the given command
func executeCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// LoadConfig loads YAML config file into Config struct
func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
