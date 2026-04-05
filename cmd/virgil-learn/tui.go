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
	virgil  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7D56F4"))
	header  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFFFFF"))
	help    = lipgloss.NewStyle().Faint(true)
	errStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF5F5F"))
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

func (m tuiModel) Init() (tuiModel, tea.Cmd) {
	return m, tea.Batch(
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

func (m tuiModel) Update(msg tea.Msg) (tuiModel, tea.Cmd) {
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
		return tea.NewView(
			errStyle.Render(fmt.Sprintf("\n  Error: %v\n\n", m.err)) +
				help.Render("  Press q to quit\n"),
		)

	case stateAnalyzing:
		return tea.NewView(
			"\n  " + virgil.Render("virgil-learn") + "  " +
				m.spinner.View() +
				header.Render(" Analyzing...") +
				"\n\n" +
				help.Render("  Press q to quit\n"),
		)

	case stateReady:
		titleBar := virgil.Render(" virgil-learn ") +
			header.Render("— Analysis Results")
		helpBar := help.Render("  j/k / arrows scroll  · q quit")

		return tea.NewView(
			titleBar + "\n" +
				m.viewport.View() + "\n" +
				helpBar,
		)
	}

	return tea.NewView("")
}

// ----------------------------------------------------------------------------
// Helpers
// ----------------------------------------------------------------------------

// viewportHeight reserves 3 lines for the title bar and help bar.
func viewportHeight(totalHeight int) int {
	reserved := 3
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
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
