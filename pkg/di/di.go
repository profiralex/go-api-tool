package di

import (
	"github.com/golobby/container"
	"github.com/profiralex/go-api-tool/pkg/files"
	"github.com/profiralex/go-api-tool/pkg/gen/logs"
)

func RegisterDependencies() {
	container.Singleton(func() files.Provider {
		return files.NewBoxProvider()
	})

	container.Transient(func(fp files.Provider) *logs.Generator {
		return logs.NewGenerator(fp)
	})
}

func UnregisterDependencies() {
	container.Reset()
}

func Make(receiver interface{}) {
	container.Make(receiver)
}
