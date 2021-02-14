package commands

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

type BasicResponse struct {
	Status int `json:"status"`
}

type GroupPayload struct {
	GroupName string `json:"group_name"`
}

func PrintTable(columNames table.Row, records []table.Row) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(columNames)
	t.AppendRows(records)
	t.Render()
}
