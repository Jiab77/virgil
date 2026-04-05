package main

import (
	"fmt"

	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/jiab77/virgil/pkg/virgil/learning"
)

// ----------------------------------------------------------------------------
// Lipgloss styles
// ----------------------------------------------------------------------------

var (
	purple = lipgloss.Color("#7D56F4")
	white  = lipgloss.Color("#FFFFFF")
	red    = lipgloss.Color("#FF5F5F")
	dim    = lipgloss.Color("#555555")

	brandStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(white).
			Background(purple).
			Padding(0, 1)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(purple)

	pathStyle = lipgloss.NewStyle().
			Foreground(white).
			Faint(true)

	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(purple)

	helpStyle = lipgloss.NewStyle().
			Foreground(dim)

	helpKeyStyle = lipgloss.NewStyle().
			Foreground(purple).
			Bold(true)

	spinnerStyle = lipgloss.NewStyle().
			Foreground(purple)

	errStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(red)
)

// ----------------------------------------------------------------------------
// TUI states
// ----------------------------------------------------------------------------

type tuiState int

const (
	stateAnalyzing tuiState = iota
	stateReady
	stateErr
)

// ----------------------------------------------------------------------------
// Messages
// ----------------------------------------------------------------------------

type analysisDoneMsg struct {
	patterns []learning.CodePattern
}

type analysisErrMsg struct {
	err error
}

// ----------------------------------------------------------------------------
// Model
// ----------------------------------------------------------------------------

type tuiModel struct {
	spinner     spinner.Model
	viewport    viewport.Model
	state       tuiState
	codebasePath string
	useMarkdown bool
	err         error
	width       int
	height      int
	ready       bool
}

func initialModel(codebasePath string, useMarkdown bool) tuiModel {
	sp := spinner.New(spinner.WithSpinner(spinner.Dot))
	sp.Style = spinnerStyle
	return tuiModel{
		spinner:      sp,
		codebasePath: codebasePath,
		useMarkdown:  useMarkdown,
		state:        stateAnalyzing,
	}
}

// ----------------------------------------------------------------------------
// Init
// ----------------------------------------------------------------------------

func (m tuiModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		runAnalysis(m.codebasePath),
	)
}

// runAnalysis runs AnalyzeCodebase in a goroutine and returns the result as a Cmd.
func runAnalysis(codebasePath string) tea.Cmd {
	return func() tea.Msg {
		analyzer := &learning.BashAnalyzer{}
		patterns, err := analyzer.AnalyzeCodebase(codebasePath)
		if err != nil {
			return analysisErrMsg{err: err}
		}
		return analysisDoneMsg{patterns: patterns}
	}
}

// ----------------------------------------------------------------------------
// Update
// ----------------------------------------------------------------------------

func (m tuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		if m.state == stateReady && !m.ready {
			m.viewport.SetWidth(msg.Width)
			m.viewport.SetHeight(viewportHeight(msg.Height))
		}

	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}

	case analysisDoneMsg:
		patterns := msg.patterns
		patternsByFile, sortedFiles := groupPatterns(patterns)

		var content string
		if m.useMarkdown {
			content = renderMarkdown(m.codebasePath, patterns, patternsByFile, sortedFiles)
		} else {
			content = renderPlainText(m.codebasePath, patterns, patternsByFile, sortedFiles)
		}

		vp := viewport.New(viewport.WithWidth(m.width), viewport.WithHeight(viewportHeight(m.height)))
		vp.SetContent(content)
		m.viewport = vp
		m.state = stateReady
		m.ready = true
		return m, nil

	case analysisErrMsg:
		m.err = msg.err
		m.state = stateErr
		return m, nil

	case spinner.TickMsg:
		if m.state == stateAnalyzing {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	if m.state == stateReady {
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// ----------------------------------------------------------------------------
// View
// ----------------------------------------------------------------------------

func (m tuiModel) View() tea.View {
	switch m.state {
	case stateErr:
		errBox := borderStyle.
			Width(60).
			Render(
				errStyle.Render("Error\n\n") +
					fmt.Sprintf("%v", m.err),
			)
		quitHint := helpStyle.Render("  press ") +
			helpKeyStyle.Render("q") +
			helpStyle.Render(" to quit")

		v := tea.NewView("\n" + errBox + "\n\n" + quitHint + "\n")
		v.AltScreen = true
		return v

	case stateAnalyzing:
		brand := brandStyle.Render("virgil-learn")
		scanning := titleStyle.Render(" scanning")
		path := "  " + pathStyle.Render(m.codebasePath)
		spin := "  " + m.spinner.View()

		v := tea.NewView(
			"\n" +
				"  " + brand + scanning + "\n" +
				path + "\n\n" +
				spin + "\n",
		)
		v.AltScreen = true
		return v

	case stateReady:
		brand := brandStyle.Render("virgil-learn")
		path := pathStyle.Render("  " + m.codebasePath)

		titleBar := lipgloss.JoinHorizontal(
			lipgloss.Left,
			brand,
			path,
		)

		vpBordered := borderStyle.
			Width(m.width - 2).
			Height(viewportHeight(m.height)).
			Render(m.viewport.View())

		helpBar := helpStyle.Render("  scroll ") +
			helpKeyStyle.Render("j/k") +
			helpStyle.Render("  quit ") +
			helpKeyStyle.Render("q")

		v := tea.NewView(
			titleBar + "\n" +
				vpBordered + "\n" +
				helpBar,
		)
		v.AltScreen = true
		return v
	}

	v := tea.NewView("")
	v.AltScreen = true
	return v
}

// ----------------------------------------------------------------------------
// Helpers
// ----------------------------------------------------------------------------

// viewportHeight reserves lines for the title bar, help bar, and viewport border.
func viewportHeight(totalHeight int) int {
	// 1 title bar + 1 help bar + 2 border lines = 4 reserved
	reserved := 4
	h := totalHeight - reserved
	if h < 1 {
		return 1
	}
	return h
}

// groupPatterns deduplicates the grouping logic shared between renderers and TUI.
func groupPatterns(patterns []learning.CodePattern) (map[string][]learning.CodePattern, []string) {
	patternsByFile := make(map[string][]learning.CodePattern)
	for _, pattern := range patterns {
		patternsByFile[pattern.FilePath] = append(patternsByFile[pattern.FilePath], pattern)
	}
	sortedFiles := sortedKeys(patternsByFile)
	return patternsByFile, sortedFiles
}

// runTUI launches the bubbletea program for --tui mode.
func runTUI(codebasePath string, useMarkdown bool) error {
	m := initialModel(codebasePath, useMarkdown)
	p := tea.NewProgram(m)
	_, err := p.Run()
	return err
}
