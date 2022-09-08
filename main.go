package main

import (
	"fmt"
	"os"

	"github.com/76creates/stickers"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"

	"github.com/nore-dev/fman/model"
	"github.com/nore-dev/fman/theme"
)

type App struct {
	listModel    model.ListModel
	entryModel   model.EntryModel
	toolbarModel model.ToolbarModel

	flexBox *stickers.FlexBox
}

func (app *App) Init() tea.Cmd {
	return nil
}

func (app *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return app, tea.Quit
		}

	case tea.WindowSizeMsg:
		toolbarHeight := 3

		app.flexBox.SetHeight(msg.Height - toolbarHeight)
		app.flexBox.SetWidth(msg.Width)

		app.flexBox.ForceRecalculate()

		app.listModel.Width = app.flexBox.Row(0).Cell(0).GetWidth()
		app.entryModel.Width = app.flexBox.Row(0).Cell(1).GetWidth()

	}

	var listCmd, toolbarCmd, entryCmd tea.Cmd

	app.listModel, listCmd = app.listModel.Update(msg)
	app.toolbarModel, toolbarCmd = app.toolbarModel.Update(msg)
	app.entryModel, entryCmd = app.entryModel.Update(msg)

	return app, tea.Batch(listCmd, toolbarCmd, entryCmd)
}

func (app *App) View() string {
	app.flexBox.ForceRecalculate()

	row := app.flexBox.Row(0)

	// Set content of list view
	row.Cell(0).SetContent(app.listModel.View())

	// Set content of entry view
	row.Cell(1).SetContent(app.entryModel.View())

	return zone.Scan(lipgloss.JoinVertical(
		lipgloss.Top,
		app.toolbarModel.View(),
		app.flexBox.Render()))
}

func main() {
	// Initialize Bubblezone
	zone.NewGlobal()
	defer zone.Close()

	app := App{
		listModel:    model.NewListModel(),
		entryModel:   model.NewEntryModel(),
		toolbarModel: model.NewToolbarModel(),
		flexBox:      stickers.NewFlexBox(0, 0),
	}

	rows := []*stickers.FlexBoxRow{
		app.flexBox.NewRow().AddCells(
			[]*stickers.FlexBoxCell{
				stickers.NewFlexBoxCell(7, 1).SetStyle(lipgloss.NewStyle().Padding(1)), // List
				stickers.NewFlexBoxCell(3, 1).SetStyle(theme.ContainerStyle),           // Entry Information
			},
		),
	}

	app.flexBox.AddRows(rows)

	p := tea.NewProgram(&app, tea.WithAltScreen(), tea.WithMouseAllMotion())

	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
