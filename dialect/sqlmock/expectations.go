package sqlmock

import (
	"fmt"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func ExpectRun(mock sqlmock.Sqlmock) {
	mock.ExpectBegin()
	tableExistsRows := sqlmock.NewRows([]string{"table_exists"}).AddRow(false)
	mock.ExpectQuery(tableExists).WillReturnRows(tableExistsRows)
	mock.ExpectRollback()

	expectRevision := func(id int, sql string) {
		mock.ExpectBegin()

		mock.ExpectExec(sql).
			WillReturnResult(sqlmock.NewResult(0, 0))

		mock.ExpectExec(setID).
			WithArgs("revision", fmt.Sprintf(`{"ID":%d}`, id)).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectCommit()
	}

	expectRevision(1, createTable)

	mock.ExpectBegin()

	mock.ExpectExec(setID).
		WithArgs("version", fmt.Sprintf(`{"ID":%d}`, version)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectCommit()
}
