package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

func calculateWidth(min, width int) int {
	p := width / 10

	switch min {
	case XS:
		if p < XS {
			return XS
		}
		return p / 2
	case SM:
		if p < SM {
			return SM
		}
		return p / 2
	case MD:
		if p < MD {
			return MD
		}
		return p * 2
	case LG:
		if p < LG {
			return LG
		}
		return p * 2
	default:
		return p
	}
}

const (
	XS int = 1
	SM int = 3
	MD int = 5
	LG int = 10
)

func setupTable(tasks []task) table.Model {
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		log.Println("unable to calc height and width of term")
	}

	columns := []table.Column{
		{Title: "ID", Width: calculateWidth(XS, w)},
		{Title: "Name", Width: calculateWidth(LG, w)},
		{Title: "Project", Width: calculateWidth(MD, w)},
		{Title: "Status", Width: calculateWidth(SM, w)},
		{Title: "Created At", Width: calculateWidth(MD, w)},
	}

	var rows []table.Row
	for _, task := range tasks {
		rows = append(rows, table.Row{
			fmt.Sprintf("%d", task.ID),
			task.Name,
			task.Project,
			task.Status,
			task.Created.Format("2006-01-02"),
		})
	}

	tbl := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(len(tasks)),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240")).BorderBottom(true).Bold(false)
	tbl.SetStyles(s)
	return tbl
}
