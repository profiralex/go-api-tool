package migrations

import (
	"fmt"
	"github.com/profiralex/go-api-tool/pkg/files"
	"github.com/profiralex/go-api-tool/pkg/gen/utils"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const migrationsDirectory = "migrations"

type migrationsGenerator struct {
	fp      files.Provider
	dialect SqlDialect
}

func NewMigrationsGenerator(fp files.Provider, dialect SqlDialect) *migrationsGenerator {
	return &migrationsGenerator{
		fp:      fp,
		dialect: dialect,
	}
}

func (g *migrationsGenerator) Generate(projectPath string, models []Table) error {
	migrationsPath := filepath.Join(projectPath, migrationsDirectory)
	err := utils.EnsureDirectoryExists(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to create migrations directory: %w", err)
	}

	newModels, _ /*modelsUpdates*/, deletedModels, err := g.getModelsDiff(migrationsPath, models)
	if err != nil {
		return fmt.Errorf("failed to parse current models: %w", err)
	}

	for _, model := range newModels {
		upSql, err := g.dialect.CreateSql(model)
		if err != nil {
			return fmt.Errorf("failed to generate creation up sql for model %s: %w", model.Name, err)
		}
		downSql, err := g.dialect.DeleteSql(model)
		if err != nil {
			return fmt.Errorf("failed to generate creation down sql for model %s: %w", model.Name, err)
		}

		err = g.createModelMigrationFiles(model, "create", upSql, downSql)
		if err != nil {
			return err
		}
	}

	for _, model := range deletedModels {
		upSql, err := g.dialect.DeleteSql(model)
		if err != nil {
			return fmt.Errorf("failed to generate deletion up sql for model %s: %w", model.Name, err)
		}
		downSql, err := g.dialect.CreateSql(model)
		if err != nil {
			return fmt.Errorf("failed to generate deletion down sql for model %s: %w", model.Name, err)
		}

		err = g.createModelMigrationFiles(model, "delete", upSql, downSql)
		if err != nil {
			return err
		}
	}

	//for _, modelUpdate := range modelsUpdates {
	//	if len(modelUpdate.ColumnsToDelete) > 0 {
	//		deleteFieldsSqlUp, err := g.dialect.DeleteFieldsSql(modelUpdate.Table, modelUpdate.ColumnsToDelete)
	//		if err != nil {
	//			return fmt.Errorf("failed to generate remove fields up sql for model %s: %w", modelUpdate.Table.Name, err)
	//		}
	//
	//		deleteFieldsSqlDown, err := g.dialect.AddFieldsSql(modelUpdate.Table, modelUpdate.ColumnsToDelete)
	//		if err != nil {
	//			return fmt.Errorf("failed to generate remove fields down sql for model %s: %w", modelUpdate.Table.Name, err)
	//		}
	//
	//		err = g.createModelMigrationFiles(modelUpdate.Table, "update", deleteFieldsSqlUp, deleteFieldsSqlDown)
	//		if err != nil {
	//			return err
	//		}
	//	}
	//
	//	if len(modelUpdate.ColumnsToAdd) > 0 {
	//		addFieldsSqlUp, err := g.dialect.AddFieldsSql(modelUpdate.Table, modelUpdate.ColumnsToAdd)
	//		if err != nil {
	//			return fmt.Errorf("failed to generate add fields up sql for model %s: %w", modelUpdate.Table.Name, err)
	//		}
	//
	//		addFieldsSqlDown, err := g.dialect.DeleteFieldsSql(modelUpdate.Table, modelUpdate.ColumnsToAdd)
	//		if err != nil {
	//			return fmt.Errorf("failed to generate add fields down sql for model %s: %w", modelUpdate.Table.Name, err)
	//		}
	//
	//		err = g.createModelMigrationFiles(modelUpdate.Table, "update", addFieldsSqlUp, addFieldsSqlDown)
	//		if err != nil {
	//			return err
	//		}
	//	}
	//}

	return nil
}

func (g *migrationsGenerator) getModelsDiff(migrationsPath string, models []Table) ([]Table, []apiModelUpdate, []Table, error) {
	currentModels, err := g.parseCurrentModels(migrationsPath)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get current models: %w", err)
	}

	var modelsToCreate []Table
	var modelsToDelete []Table
	var modelsUpdates []apiModelUpdate

	//check for new or changed models
	for _, model := range models {
		var found bool
		var currentModel Table
		for _, m := range currentModels {
			if m.Name == model.Name {
				currentModel = m
				found = true
			}
		}
		if !found {
			modelsToCreate = append(modelsToCreate, model)
			continue
		}

		modelUpdates := apiModelUpdate{Table: model}

		//check for new or updated fields
		for _, field := range model.Columns {
			if !currentModel.HasField(field.Name) {
				modelUpdates.ColumnsToAdd = append(modelUpdates.ColumnsToAdd, field)
			}

			// check for changed constraints
		}

		//check for deleted fields
		for _, currentField := range currentModel.Columns {
			if !model.HasField(currentField.Name) {
				modelUpdates.ColumnsToDelete = append(modelUpdates.ColumnsToDelete, currentField)
			}
		}

		if len(modelUpdates.ColumnsToAdd) > 0 || len(modelUpdates.ColumnsToDelete) > 0 {
			modelsUpdates = append(modelsUpdates, modelUpdates)
		}
	}

	//check for deleted models
	for _, currentModel := range currentModels {
		var found bool
		for _, m := range models {
			if m.Name == currentModel.Name {
				found = true
			}
		}

		if !found {
			modelsToDelete = append(modelsToDelete, currentModel)
		}
	}

	return modelsToCreate, modelsUpdates, modelsToDelete, nil
}

func (g *migrationsGenerator) modelChanged(model1 Table, model2 Table) bool {
	return false
}

func (g *migrationsGenerator) parseCurrentModels(migrationsPath string) ([]Table, error) {
	queries := map[string]string{}
	var migrationFiles []string

	err := filepath.Walk(migrationsPath, func(path string, info os.FileInfo, err error) error {
		if !strings.HasSuffix(path, ".up.sql") {
			return nil
		}

		_, filename := filepath.Split(path)
		migrationFiles = append(migrationFiles, filename)
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read %s file: %w", path, err)
		}

		queries[filename] = string(content)

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	sort.Strings(migrationFiles)

	modelsMap := map[string]Table{}
	for _, file := range migrationFiles {
		parts := strings.SplitN(strings.TrimSuffix(file, "_table.up.sql"), "_", 3)
		op := parts[1]
		modelName := utils.GetModelNameFromSqlTable(parts[2])
		query := queries[file]

		switch op {
		case "create":
			_, ok := modelsMap[modelName]
			if ok {
				return nil, fmt.Errorf("model %s for create query %s already exists", modelName, file)
			}

			model := Table{Name: modelName}
			err = g.dialect.ParseCreateQuery(&model, query)
			if err != nil {
				return nil, fmt.Errorf("failed to parse update query %s: %w", file, err)
			}
			modelsMap[modelName] = model
		case "update":
			model, ok := modelsMap[modelName]
			if !ok {
				return nil, fmt.Errorf("model %s for update query %s not found", modelName, file)
			}

			err = g.dialect.ParseUpdateQuery(&model, query)
			if err != nil {
				return nil, fmt.Errorf("failed to parse update query %s: %w", file, err)
			}
			modelsMap[modelName] = model
		case "delete":
			_, ok := modelsMap[modelName]
			if !ok {
				return nil, fmt.Errorf("model %s for delete query %s not found", modelName, file)
			}
			delete(modelsMap, modelName)
		}
	}

	var models []Table
	for _, model := range modelsMap {
		models = append(models, model)
	}

	return models, nil
}

func (g *migrationsGenerator) createModelMigrationFiles(model Table, migrationOp string, upSql []byte, downSql []byte) error {
	ts := time.Now().UnixNano()
	upSqlFilename := fmt.Sprintf("%d_%s_%s_table.up.sql", ts, migrationOp, utils.GetModelSqlTable(model.Name))
	downSqlFilename := fmt.Sprintf("%d_%s_%s_table.down.sql", ts, migrationOp, utils.GetModelSqlTable(model.Name))

	err := utils.GenerateFile(upSqlFilename, upSql)
	if err != nil {
		return err
	}

	err = utils.GenerateFile(downSqlFilename, downSql)
	if err != nil {
		return err
	}

	return nil
}
