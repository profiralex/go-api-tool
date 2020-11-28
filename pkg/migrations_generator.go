package pkg

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"
)

const migrationsDirectory = "migrations"

type migrationsGenerator struct {
	projectPath string
	apiSpec     apiSpec
	dialect     sqlDialect
}

type apiModelUpdate struct {
	model          apiModel
	fieldsToAdd    []apiModelField
	fieldsToUpdate []apiModelField
	fieldsToDelete []apiModelField
}

func NewMigrationsGenerator(projectPath string, spec apiSpec) *migrationsGenerator {
	return &migrationsGenerator{
		projectPath: projectPath,
		apiSpec:     spec,
		dialect:     newMysqlDialect(),
	}
}

func (g *migrationsGenerator) Generate() error {
	err := g.ensureMigrationsDirectoryExists()
	if err != nil {
		return fmt.Errorf("failed to create migrations directory: %w", err)
	}

	newModels := g.apiSpec.Models
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

	deletedModels := g.apiSpec.Models
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

	var modelUpdates []apiModelUpdate
	for _, model := range g.apiSpec.Models {
		modelUpdates = append(modelUpdates, apiModelUpdate{
			model:          model,
			fieldsToAdd:    model.Fields,
			fieldsToUpdate: model.Fields,
			fieldsToDelete: model.Fields,
		})
	}

	for _, modelUpdate := range modelUpdates {
		deleteFieldsSqlUp := ""
		deleteFieldsSqlDown := ""
		if len(modelUpdate.fieldsToDelete) > 0 {
			deleteFieldsSqlUp, err = g.dialect.DeleteFieldsSql(modelUpdate.model, modelUpdate.fieldsToDelete)
			if err != nil {
				return fmt.Errorf("failed to generate remove fields up sql for model %s: %w", modelUpdate.model.Name, err)
			}

			deleteFieldsSqlDown, err = g.dialect.AddFieldsSql(modelUpdate.model, modelUpdate.fieldsToDelete)
			if err != nil {
				return fmt.Errorf("failed to generate remove fields down sql for model %s: %w", modelUpdate.model.Name, err)
			}
		}

		addFieldsSqlUp := ""
		addFieldsSqlDown := ""
		if len(modelUpdate.fieldsToDelete) > 0 {
			addFieldsSqlUp, err = g.dialect.AddFieldsSql(modelUpdate.model, modelUpdate.fieldsToAdd)
			if err != nil {
				return fmt.Errorf("failed to generate add fields up sql for model %s: %w", modelUpdate.model.Name, err)
			}

			addFieldsSqlDown, err = g.dialect.DeleteFieldsSql(modelUpdate.model, modelUpdate.fieldsToAdd)
			if err != nil {
				return fmt.Errorf("failed to generate add fields down sql for model %s: %w", modelUpdate.model.Name, err)
			}
		}

		modifyFieldsUp := ""
		modifyFieldsDown := ""
		if len(modelUpdate.fieldsToDelete) > 0 {
			modifyFieldsUp, err = g.dialect.UpdateFieldsSql(modelUpdate.model, modelUpdate.fieldsToAdd)
			if err != nil {
				return fmt.Errorf("failed to generate update fields up sql for model %s: %w", modelUpdate.model.Name, err)
			}

			modifyFieldsDown, err = g.dialect.UpdateFieldsSql(modelUpdate.model, modelUpdate.fieldsToAdd)
			if err != nil {
				return fmt.Errorf("failed to generate update fields down sql for model %s: %w", modelUpdate.model.Name, err)
			}
		}

		finalUpdateFieldsUp := fmt.Sprintf("%s\n\n%s\n\n%s", deleteFieldsSqlUp, addFieldsSqlUp, modifyFieldsUp)
		finalUpdateFieldsDown := fmt.Sprintf("%s\n\n%s\n\n%s", deleteFieldsSqlDown, addFieldsSqlDown, modifyFieldsDown)

		err = g.createModelMigrationFiles(modelUpdate.model, "update", finalUpdateFieldsUp, finalUpdateFieldsDown)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *migrationsGenerator) ensureMigrationsDirectoryExists() error {
	migrationsDirectory := path.Join(g.projectPath, migrationsDirectory)
	fileStat, err := os.Stat(migrationsDirectory)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to check if migrations directory exists: %w", err)
		}
	} else {
		if !fileStat.IsDir() {
			// oops looks like there is a migrations file
			return fmt.Errorf("migrations file found in project root")
		}
		// directory already exists all good
		return nil
	}

	err = os.Mkdir(migrationsDirectory, 0755)
	if err != nil {
		return fmt.Errorf("failed to create migrations directory: %w", err)
	}

	return nil
}

func (g *migrationsGenerator) createModelMigrationFiles(model apiModel, migrationOp string, upSql string, downSql string) error {
	ts := time.Now().UnixNano()
	upSqlFilename := fmt.Sprintf("%d_%s_%s_table.up.sql", ts, migrationOp, getModelSqlTable(model.Name))
	downSqlFilename := fmt.Sprintf("%d_%s_%s_table.down.sql", ts, migrationOp, getModelSqlTable(model.Name))

	err := g.writeMigrationFile(upSqlFilename, upSql)
	if err != nil {
		return err
	}

	err = g.writeMigrationFile(downSqlFilename, downSql)
	if err != nil {
		return err
	}

	return nil
}

func (g *migrationsGenerator) writeMigrationFile(sqlFilename string, content string) error {
	upSqlFile := path.Join(g.projectPath, migrationsDirectory, sqlFilename)
	err := ioutil.WriteFile(upSqlFile, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write migration file %s: %w", sqlFilename, err)
	}
	return nil
}
