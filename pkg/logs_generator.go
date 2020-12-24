package pkg

import "path"

type logsGenerator struct {
	projectPath string
	apiSpec     apiSpec
}

var logsStatics = map[string]string{
	"logs/logs.go.tpl": path.Join(genDirectory, "logs", "logs.go"),
}

func NewLogsGenerator(projectPath string, spec apiSpec) *logsGenerator {
	return &logsGenerator{
		projectPath: projectPath,
		apiSpec:     spec,
	}
}

func (g *logsGenerator) Generate() error {
	return copyBoxStatics(g.projectPath, logsStatics)
}
