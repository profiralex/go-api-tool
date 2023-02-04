package main

import (
	"flag"
	"github.com/profiralex/go-api-tool/pkg/gen"
	"github.com/profiralex/go-api-tool/pkg/logs"
	log "github.com/sirupsen/logrus"
)

func main() {
	// parse args
	projectPath := flag.String("proj", "./", "project path")
	logLevel := flag.String("log", "warn", "log level: info, warning, error")
	flag.Parse()

	logs.Init(*logLevel)

	generator := gen.NewGenerator(*projectPath)
	err := generator.Init()
	if err != nil {
		log.Errorf("failed to read the current setup: %s", err)
	}

	err = generator.Generate()
	if err != nil {
		log.Errorf("api generation failed: %s", err)
	}
}
