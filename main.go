package main

import (
	"fmt"
	"os"

	"brim/provider"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

// search result item
type item struct {
	idx                      int
	title, url, desc, engine string
}

// getters 
func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

// list of search result items
type model struct {
	result_list list.Model
}


func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}
	var cmd tea.Cmd
	m.result_list, cmd = m.result_list.Update(msg)
	return m, cmd
}

func (m model) View() tea.View {
	v := tea.NewView(docStyle.Render(m.result_list.View()))
	v.AltScreen = true
	return v
}

func main() {
	s := provider.NewService()
	res, err := s.FetchWebResults("rust borrow checker", 5)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	items := []list.Item{}
	for i, r := range res.Results {
		items = append(items, item{
			idx: i, title: r.Title, url: r.URL, engine: r.Engine, desc: "Not available",
		})
	}

	m := model{result_list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.result_list.Title = "Search Results"

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
