package pkg

import (
	"path"
)

var logsContent = `/*Generated code do not modify it*/
package logs

import (
	log "github.com/sirupsen/logrus"
)

func Init(logLevel string) {
	level, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Warnf("Unknown log level %s defaulting to warning level", logLevel)
		level = log.WarnLevel
	}
	log.SetLevel(level)
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)
}
`

type logsGenerator struct {
	projectPath string
	apiSpec     apiSpec
}

func NewLogsGenerator(projectPath string, spec apiSpec) *logsGenerator {
	return &logsGenerator{
		projectPath: projectPath,
		apiSpec:     spec,
	}
}

func (g *logsGenerator) Generate() error {
	logsFilePath := path.Join(g.projectPath, genDirectory, "logs", "logs.go")
	return generateFile(logsFilePath, []byte(logsContent))
}
