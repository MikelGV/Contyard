package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss/table"
)

func CreateTable() *table.Table {
    t := table.New().
        Headers("containers", "pods").
        Row("cpu_usage").
        Row("memory_usage").
        Row("memory_limit")

    fmt.Println(t)
    return  t
}
