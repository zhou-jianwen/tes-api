package testutil

import (
	"fmt"
	"strings"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// NewSqlMock ...
func NewSqlMock() (sqlmock.Sqlmock, *gorm.DB) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		panic(err)
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}))
	if err != nil {
		panic(err)
	}
	return mock, gormDB
}

// GenInsertSql ...
func GenInsertSql(rows []string) string {
	names := make([]string, 0, len(rows))
	questions := make([]string, 0, len(rows))
	for _, row := range rows {
		names = append(names, fmt.Sprintf("`%s`", row))
		questions = append(questions, "?")
	}
	return fmt.Sprintf("(%s) VALUES (%s)", strings.Join(names, ","), strings.Join(questions, ","))
}

// GenSelectFieldsSql ...
func GenSelectFieldsSql(tableName string, rows []string) string {
	fields := make([]string, 0, len(rows))
	for _, row := range rows {
		fields = append(fields, fmt.Sprintf("`%s`.`%s`", tableName, row))
	}
	return strings.Join(fields, ",")
}

// GenUpdateSql ...
func GenUpdateSql(rows []string) string {
	items := make([]string, 0, len(rows))
	for _, row := range rows {
		items = append(items, fmt.Sprintf("`%s`=?", row))
	}
	return strings.Join(items, ",")
}

// GenDuplicateKeySql ...
func GenDuplicateKeySql(rows []string) string {
	items := make([]string, 0, len(rows))
	for _, row := range rows {
		items = append(items, fmt.Sprintf("`%s`=VALUES(`%s`)", row, row))
	}
	return fmt.Sprintf("ON DUPLICATE KEY UPDATE %s", strings.Join(items, ","))
}

// NewCountRows ...
func NewCountRows(count int) *sqlmock.Rows {
	return sqlmock.NewRows([]string{"count"}).AddRow(count)
}
