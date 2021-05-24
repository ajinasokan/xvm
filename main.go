package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
)

const binaryName = "xvm"
const linkDirName = ".xvm"
const confName = ".xvm.conf"

//go:embed help.txt
var helpText string

func main() {
	// when symlinked binaries will have the symlink name as their arg[0]
	// this is used to identify whether xvm is executed as it is or as one
	// of the symlinks from ~/.xvm/
	if os.Args[0] != binaryName {
		pipeProcess()
	} else if len(os.Args) == 1 {
		printHelp()
	} else {
		manageCommands()
	}
}

func printHelp() {
	fmt.Println(helpText)
}

// find the correct path for the command by looping through directory
// tree from current directory to the root directory. if there is a
// .xvm.conf file with an override of the command then use that for exec.
// config in the current working dir gets most priority. if no override found
// for the command then default path specified while enabling xvm for the command
// is used.
func findCommandPath(command string, printScans bool) string {
	commandPath := ""

	dir, err := os.Getwd()
	if err != nil {
		exitWithError("unable to access current dir", err)
	}
	for {
		rcPath := path.Join(dir, confName)
		if printScans {
			fmt.Println("checking", rcPath)
		}
		if fileExist(rcPath) {
			m, err := loadConfig(rcPath)
			if err != nil {
				exitWithError("unable to read", rcPath, err)
			}
			if val, ok := m[command]; ok {
				commandPath = val
				if printScans {
					fmt.Println("found", commandPath)
				}
				return commandPath
			}
		}
		parent := path.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	if commandPath == "" {
		m := readGlobalConf()
		if val, ok := m[command]; ok {
			commandPath = val
			if printScans {
				fmt.Println("using default", commandPath)
			}
		} else {
			exitWithError("could not find", command)
		}
	}
	return commandPath
}

// find the path of the command and then start process with incoming args.
// proxy stdio of child with the host process and wait until the child process
// exits
func pipeProcess() {
	cmd := findCommandPath(os.Args[0], false)
	args := []string{}
	if len(os.Args) > 1 {
		args = os.Args[1:]
	}

	p := exec.Command(cmd, args...)
	p.Stdout = os.Stdout
	p.Stdin = os.Stdin
	p.Stderr = os.Stderr

	err := p.Start()
	if err != nil {
		exitWithError("unable to launch", cmd, args, err)
	}

	p.Wait()
}

func manageCommands() {
	command := os.Args[2]

	currentUser, err := user.Current()
	if err != nil {
		exitWithError("unable to access home dir")
	}

	linksDir := filepath.Join(currentUser.HomeDir, linkDirName)
	globalConfPath := filepath.Join(linksDir, confName)

	selfPath, err := os.Executable()
	if err != nil {
		exitWithError("unable to check xvm path", err)
	}

	if os.Args[1] == "enable" {
		createSymlink(linksDir, selfPath, os.Args[2])
		globalConf, _ := loadConfig(globalConfPath)
		globalConf[os.Args[2]] = os.Args[3]
		err = saveConfig(globalConfPath, globalConf)
		if err != nil {
			exitWithError("unable to save config", err)
		}
	}

	if os.Args[1] == "disable" {
		removeSymlink(linksDir, os.Args[2])
		globalConf, _ := loadConfig(globalConfPath)
		delete(globalConf, os.Args[2])
		err = saveConfig(globalConfPath, globalConf)
		if err != nil {
			exitWithError("unable to save config", err)
		}
	}

	if os.Args[1] == "set" {
		if len(os.Args) != 4 {
			exitWithError("invalid args")
		}
		localConf, _ := loadConfig(confName)
		localConf[os.Args[2]] = os.Args[3]
		err = saveConfig(confName, localConf)
		if err != nil {
			exitWithError("unable to save config", err)
		}
		fmt.Println("saved config", confName)
	}

	if os.Args[1] == "unset" {
		localConf, err := loadConfig(confName)
		if os.IsNotExist(err) {
			exitWithError("config file", confName, "not found")
			return
		}
		delete(localConf, os.Args[2])
		err = saveConfig(confName, localConf)
		if err != nil {
			exitWithError("unable to save config", err)
		}
		fmt.Println("saved config", confName)
	}

	if os.Args[1] == "get" {
		findCommandPath(command, true)
	}
}

func exitWithError(msg ...interface{}) {
	fmt.Fprintln(os.Stderr, msg...)
	os.Exit(1)
}
