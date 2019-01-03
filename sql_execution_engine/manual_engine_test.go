package sql_execution_engine

import (
	"fmt"
	"github.com/moiot/gravity/schema_store"
	"testing"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/stretchr/testify/require"

	"github.com/moiot/gravity/pkg/core"
	"github.com/moiot/gravity/pkg/mysql_test"
)

func TestManualSQLEngine(t *testing.T) {
	r := require.New(t)

	db, mock, err := sqlmock.New()
	r.NoError(err)

	newName := "fake data"

	mock.ExpectExec(fmt.Sprintf("UPDATE `%s`.`%s` SET name = ?", t.Name(), mysql_test.TestTableName)).WithArgs(newName).WillReturnResult(sqlmock.NewResult(1, 1))

	SQLTemplate := "UPDATE `{{.TargetTable.Schema}}`.`{{.TargetTable.Name}}` SET name = ?"

	engineConfig := MySQLExecutionEngineConfig{
		SQLTemplate: SQLTemplate,
		SQLArgExpr:  []string{"name"},
	}
	engine := NewManualSQLEngine(db, engineConfig)

	msg := core.Msg{
		Type: core.MsgDML,
		DmlMsg: &core.DMLMsg{
			Data: map[string]interface{}{
				"name": newName,
			},
		},
	}

	tableDef := schema_store.Table{
		Schema: t.Name(),
		Name: mysql_test.TestTableName,
	}
	r.NoError(engine.Execute([]*core.Msg{&msg}, &tableDef))
}
