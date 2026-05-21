package main

import (
	"fmt"
	"os"
	

	"brim/provider"
	"brim/utils"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	idx                      int
	title, url, desc, engine string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

// track state 
type viewState int

const (
	stateInput viewState = 0
	stateResults viewState = 1
)

type model struct {
	state       viewState
	textInput   textinput.Model
	result_list list.Model
	selectedItem item
	err         error
	width, height int
}

// Update returns a msg, 
// this is a custom msg 
type searchResultsMsg struct {
	items []list.Item
	err   error
}

// a Cmd is a function that returns a Msg.
// tea runs it on a goroutine; the Msg goes back to Update.
func doSearch(query string) tea.Cmd {
	return func() tea.Msg {
		s := provider.NewService()
		res, err := s.FetchWebResults(query, 25)
		if err != nil {
			return searchResultsMsg{err: err}
		}
		items := []list.Item{}
		for i, r := range res.Results {
			items = append(items, item{
				idx: i, title: r.Title, url: r.URL, engine: r.Engine, desc: r.Content,
			})
		}
		return searchResultsMsg{items: items}
	}
}

// func selectSearch(*items []list.Item, idx int )

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Type something to get started!"
	ti.Focus()
	ti.CharLimit = 256
	ti.SetWidth(40)

	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Search Results"

	return model{
		state:       stateInput,
		textInput:   ti,
		result_list: l,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			// check if in input state and have text 
			if m.state == stateInput && m.textInput.Value() != "" {
				return m, doSearch(m.textInput.Value())
			}
			if m.state == stateResults {
				if i, ok := m.result_list.SelectedItem().(item); ok {
					m.selectedItem = i
					return m, utils.OpenURL(i.url)
				}
			}		
	case "esc":
			// go back to input from results
			if m.state == stateResults {
				m.state = stateInput
				// m.selectedItem=nil
				return m, nil
			}
		}

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		h, v := docStyle.GetFrameSize()
		m.result_list.SetSize(msg.Width-h, msg.Height-v)

	// search cmd's msg  
	case searchResultsMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.result_list.SetItems(msg.items)
		m.state = stateResults
		return m, nil
	}

	if m.state == stateInput {
		m.textInput, cmd = m.textInput.Update(msg)
	} else {
		m.result_list, cmd = m.result_list.Update(msg)
	}
	return m, cmd
}

func (m model) View() tea.View {
	var content string
	if m.state == stateInput {
		content = fmt.Sprintf("Search:\n\n%s\n\n(enter to search, ctrl+c to quit)", m.textInput.View())
	} else {
		content = m.result_list.View() + "\n(esc to search again)"
	}
	v := tea.NewView(docStyle.Render(content))
	v.AltScreen = true
	return v
}

func main() {
	p := tea.NewProgram(initialModel())  
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
