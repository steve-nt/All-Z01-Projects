package server

import (
	"fmt"
	"forum/src/models"
	"forum/src/utils"
	"log"
	"os"
	"path/filepath"
)

func usage(programName string, toStdErr bool) {
	fd := os.Stderr
	if !toStdErr {
		fd = os.Stdout
	}
	fmt.Fprintf(fd, "Usage:\n\n")
	fmt.Fprintf(fd, "\t%s [--help]\n", programName)
	fmt.Fprintf(fd, "\t%s [--db-path <path>] [ip] [port]\n", programName)
	fmt.Fprintf(fd, "\t%s [--db-path <path>] [port]\n", programName)
	fmt.Fprintf(fd, "\t%s [--db-path <path>]\n", programName)
}

func Main(args []string) {
	config := utils.DefaultConfiguration()

	var positionalArgs []string

	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--db-path":
			config.DbPath = args[i+1]
			i++
		case "--help", "-h":
			programName := filepath.Base(args[0])
			usage(programName, false)
			os.Exit(0)
		case "--version":
			programName := filepath.Base(args[0])
			version := utils.GetVersion()
			fmt.Println(programName+"-"+version)
			os.Exit(0)
		default:
			positionalArgs = append(positionalArgs, args[i])
		}
	}

	if len(positionalArgs) == 2 {
		config.Ip = positionalArgs[0]
		config.Port = positionalArgs[1]
	} else if len(positionalArgs) == 1 {
		config.Port = positionalArgs[0]
	} else if len(positionalArgs) == 0 {
		// keep defaults
	} else {
		// unexpected argument length
		programName := filepath.Base(args[0])
		usage(programName, true)
		os.Exit(1)
	}

	if err := models.InitDB(config.DbPath); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	if err := models.InitTemplates(); err != nil {
		fmt.Printf("Error initializing templates: %s\n", err.Error())
		os.Exit(1)
	}

	log.Printf("http://%s:%s/", config.Ip, config.Port)

	if err := startServer(config.Ip, config.Port); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
}
