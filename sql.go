package utils

import (
	"fmt"
	"strings"
)

type Column struct {
	Name          string
	Type          string
	AutoIncrement bool
	PrimaryKey    bool
	Nullable      bool
	Unique        bool
	Distributed   bool
	Default       *string
}

type CreateTable struct {
	TableName  string
	Columns    []Column
	Engine     string
	Properties []string
}

type DropTable struct {
	TableName string
}

type Field struct {
	Name  string
	Value any
}

type InsertTable struct {
	TableName string
	Fields    []Field
}

type InsertManyTable struct {
	TableName    string
	ColumnName   []string
	ColumnValues [][]string
}

func buildColumnDefinition(col Column) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("%s %s", col.Name, col.Type))

	builder.WriteString(buildNullConstraint(col))
	builder.WriteString(buildAutoIncrement(col))
	builder.WriteString(buildDefault(col))

	return builder.String()
}

func buildNullConstraint(col Column) string {
	if col.AutoIncrement || col.PrimaryKey {
		return " NOT NULL"
	}
	if col.Nullable {
		return " NULL"
	}
	return " NOT NULL"
}

func buildAutoIncrement(col Column) string {
	if col.AutoIncrement {
		return " AUTO_INCREMENT"
	}
	return ""
}

func buildDefault(col Column) string {
	if col.Default != nil {
		return " DEFAULT " + *col.Default
	}
	return ""
}

func collectColumnKeys(columns []Column) (distributed, primary, unique []string) {
	for _, col := range columns {
		if col.Distributed {
			distributed = append(distributed, col.Name)
		}
		if col.PrimaryKey {
			primary = append(primary, col.Name)
		}
		if col.Unique {
			unique = append(unique, col.Name)
		}
	}
	return distributed, primary, unique
}

func addTableConstraints(builder *strings.Builder, table CreateTable, primaryKeys, uniqueKeys, distributedColumns []string) {
	if table.Engine != "" {
		builder.WriteString(fmt.Sprintf("\nENGINE=%s", table.Engine))
	}

	if len(primaryKeys) > 0 {
		builder.WriteString(fmt.Sprintf("\nPRIMARY KEY (%s)", strings.Join(primaryKeys, ", ")))
	}

	if len(uniqueKeys) > 0 {
		builder.WriteString(fmt.Sprintf("\nUNIQUE KEY (%s)", strings.Join(uniqueKeys, ", ")))
	}

	if len(distributedColumns) > 0 {
		builder.WriteString(fmt.Sprintf("\nDISTRIBUTED BY HASH (%s)", strings.Join(distributedColumns, ", ")))
	}

	if len(primaryKeys) > 0 {
		builder.WriteString(fmt.Sprintf("\nORDER BY (%s)", strings.Join(primaryKeys, ", ")))
	}

	if len(table.Properties) > 0 {
		builder.WriteString(fmt.Sprintf("\nPROPERTIES (\n\t%s\n)", strings.Join(table.Properties, ", ")))
	}
}

func CreateTableQuery(table CreateTable) (string, error) {
	var builder strings.Builder

	distributedColumns, primaryKeys, uniqueKeys := collectColumnKeys(table.Columns)

	builder.WriteString(fmt.Sprintf("CREATE TABLE %s (\n", table.TableName))

	for i, col := range table.Columns {
		builder.WriteString("\t" + buildColumnDefinition(col))
		if i < len(table.Columns)-1 {
			builder.WriteString(",\n")
		}
	}

	builder.WriteString("\n)")
	addTableConstraints(&builder, table, primaryKeys, uniqueKeys, distributedColumns)
	builder.WriteString(";")

	return builder.String(), nil
}

func DropTableQuery(table DropTable) (string, error) {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("DROP TABLE IF EXISTS %s", table.TableName))

	return builder.String(), nil
}

func InsertTableQuery(data InsertTable) (string, error) {
	if len(data.Fields) == 0 {
		return "", fmt.Errorf("cannot create query with no field")
	}

	var builder strings.Builder
	var columns, values []string

	for _, field := range data.Fields {
		columns = append(columns, field.Name)
		values = append(values, fmt.Sprintf("'%v'", field.Value))
	}

	builder.WriteString(fmt.Sprintf("INSERT INTO %s (%s)\nVALUES (%s);", data.TableName, strings.Join(columns, ", "), strings.Join(values, ", ")))

	return builder.String(), nil
}

func InsertManyTableQuery(data InsertManyTable) (string, error) {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("INSERT INTO %s (%s)\nVALUES\n", data.TableName, strings.Join(data.ColumnName, ", ")))

	for i, values := range data.ColumnValues {
		builder.WriteString(fmt.Sprintf("\t(%s)", strings.Join(values, ", ")))

		if i < len(data.ColumnValues)-1 {
			builder.WriteString(",\n")
		}
	}

	builder.WriteString(";")

	return builder.String(), nil
}
