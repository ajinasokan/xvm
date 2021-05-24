package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func fileExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// loads config file with key=value pairs into a map
func loadConfig(path string) (map[string]string, error) {
	m := map[string]string{}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return m, err
	}
	lines := strings.Split(string(b), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		i := strings.Index(line, "=")
		if i == -1 {
			return m, fmt.Errorf("invalid format %s in %s", line, path)
		}
		m[line[:i]] = line[i+1:]
	}
	return m, nil
}

// saves map to config file with key=value pairs
func saveConfig(path string, m map[string]string) error {
	s := ""
	for k, v := range m {
		s += fmt.Sprintf("%s=%s", k, v)
	}
	return ioutil.WriteFile(path, []byte(s), 0644)
}

func readGlobalConf() map[string]string {
	currentUser, err := user.Current()
	if err != nil {
		exitWithError("unable to access home dir")
	}
	linksDir := filepath.Join(currentUser.HomeDir, linkDirName)
	globalConfPath := filepath.Join(linksDir, confName)
	globalConf, _ := loadConfig(globalConfPath)
	return globalConf
}
