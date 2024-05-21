package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	asciiArt = `
	██████╗░░█████╗░███████╗████████╗██████╗░░█████╗░██████╗░███████╗██████╗░
	██╔══██╗██╔══██╗██╔════╝╚══██╔══╝██╔══██╗██╔══██╗██╔══██╗██╔════╝██╔══██╗
	██║░░██║███████║█████╗░░░░░██║░░░██████╔╝███████║██║░░██║█████╗░░██████╔╝
	██║░░██║██╔══██║██╔══╝░░░░░██║░░░██╔══██╗██╔══██║██║░░██║██╔══╝░░██╔══██╗
	██████╔╝██║░░██║███████╗░░░██║░░░██║░░██║██║░░██║██████╔╝███████╗██║░░██║
	╚═════╝░╚═╝░░╚═╝╚══════╝░░░╚═╝░░░╚═╝░░╚═╝╚═╝░░╚═╝╚═════╝░╚══════╝╚═╝░░╚═╝
	`

	headerHeight = 11

	headerStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Width(90).
			Height(headerHeight)

	statusStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Width(28).
			Align(lipgloss.Right).
			Height(headerHeight)

	bodyStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Width(120)
)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	mainMenu    list.Model
	configMenu  list.Model
	currentMenu string
	filepath    string
}

func initialModel() model {
	mainItems := []list.Item{
		item{title: "Trade", desc: "Go to trading menu"},
		item{title: "Configs", desc: "Create and edit configurations"},
		item{title: "Version", desc: "Show version information"},
		item{title: "Exit", desc: "Exit the application"},
	}

	configItems := []list.Item{
		item{title: "Config 1", desc: "Config option 1"},
		item{title: "Config 2", desc: "Config option 2"},
		item{title: "Back", desc: "Return to main menu"},
	}

	const defaultWidth = 40
	const defaultHeight = 14

	mainMenu := list.New(mainItems, list.NewDefaultDelegate(), defaultWidth, defaultHeight*2) // Double the height
	mainMenu.Title = "Main Menu"

	configMenu := list.New(configItems, list.NewDefaultDelegate(), defaultWidth, defaultHeight*2) // Double the height
	configMenu.Title = "Configs Menu"

	return model{
		mainMenu:    mainMenu,
		configMenu:  configMenu,
		currentMenu: "main",
		filepath:    "home/menu",
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.currentMenu == "main" {
				selectedItem := m.mainMenu.SelectedItem().(item)
				switch selectedItem.title {
				case "Exit":
					return m, tea.Quit
				case "Configs":
					m.currentMenu = "configs"
					m.filepath = "home/menu/configs"
				case "Version":
					fmt.Println("Version 1.0.0")
				default:
					fmt.Printf("Selected option: %s\n", selectedItem.title)
				}
			} else if m.currentMenu == "configs" {
				selectedItem := m.configMenu.SelectedItem().(item)
				switch selectedItem.title {
				case "Back":
					m.currentMenu = "main"
					m.filepath = "home/menu"
				default:
					fmt.Printf("Selected config: %s\n", selectedItem.title)
				}
			}
		}
	}

	var newListModel list.Model
	var cmd tea.Cmd

	if m.currentMenu == "main" {
		newListModel, cmd = m.mainMenu.Update(msg)
		m.mainMenu = newListModel
	} else if m.currentMenu == "configs" {
		newListModel, cmd = m.configMenu.Update(msg)
		m.configMenu = newListModel
	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	header := headerStyle.Render(fmt.Sprintf("%s\nFilepath: %s", asciiArt, m.filepath))
	status := statusStyle.Render("Status: Placeholder")

	var body string
	if m.currentMenu == "main" {
		body = bodyStyle.Render(m.mainMenu.View())
	} else if m.currentMenu == "configs" {
		body = bodyStyle.Render(m.configMenu.View())
	}

	ui := lipgloss.JoinHorizontal(lipgloss.Top, header, status) + "\n" + body
	return ui
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
