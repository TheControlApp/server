package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/thecontrolapp/server/internal/models"
)

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	menuStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(1, 2)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true).
			PaddingLeft(2)

	normalStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	logStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#00D787")).
			Padding(1, 2)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00D787")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5555")).
			Bold(true)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#50FA7B"))
)

type commandItem struct {
	name        string
	description string
	action      func(*model) tea.Cmd
}

type wsMessage struct {
	timestamp time.Time
	data      string
	isError   bool
}

type model struct {
	commander     *Commander
	commands      []commandItem
	cursor        int
	viewport      viewport.Model
	logs          []wsMessage
	ready         bool
	width         int
	height        int
	statusMessage string
	statusIsError bool
}

// Messages
type wsMsg struct {
	data string
	err  error
}

type statusMsg struct {
	message string
	isError bool
}

func initialModel(commander *Commander) model {
	commands := []commandItem{
		{
			name:        "Ping",
			description: "Test connectivity",
			action:      sendPing,
		},
		{
			name:        "Notification",
			description: "Send notification",
			action:      sendNotification,
		},
		{
			name:        "Popup",
			description: "Display popup message",
			action:      sendPopup,
		},
		{
			name:        "Display Text",
			description: "Show text content",
			action:      sendDisplayText,
		},
		{
			name:        "Timer",
			description: "Start countdown timer",
			action:      sendTimer,
		},
		{
			name:        "Multiple Choice",
			description: "Present options",
			action:      sendChoice,
		},
		{
			name:        "Open URL",
			description: "Open web link",
			action:      sendOpenURL,
		},
		{
			name:        "Kink Message",
			description: "Send kink message",
			action:      sendKinkMessage,
		},
		{
			name:        "Kink TTS",
			description: "Text-to-speech",
			action:      sendKinkTTS,
		},
		{
			name:        "Multi-Instruction",
			description: "Multiple commands",
			action:      sendMultiInstruction,
		},
	}

	vp := viewport.New(80, 20)
	vp.SetContent("WebSocket Message Log\n\nWaiting for messages...")

	return model{
		commander: commander,
		commands:  commands,
		cursor:    0,
		viewport:  vp,
		logs:      []wsMessage{},
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		waitForWebSocketMessage(m.commander),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true

		// Update viewport size (right side takes ~60% of width)
		vpWidth := int(float64(msg.Width) * 0.6)
		vpHeight := msg.Height - 6
		m.viewport.Width = vpWidth
		m.viewport.Height = vpHeight
		m.updateViewport()

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.commands)-1 {
				m.cursor++
			}

		case "enter", " ":
			// Execute selected command
			if m.cursor < len(m.commands) {
				return m, m.commands[m.cursor].action(&m)
			}
		}

	case wsMsg:
		if msg.err != nil {
			m.addLog(fmt.Sprintf("❌ Error: %v", msg.err), true)
		} else {
			m.addLog(msg.data, false)
		}
		return m, waitForWebSocketMessage(m.commander)

	case statusMsg:
		m.statusMessage = msg.message
		m.statusIsError = msg.isError
		return m, nil
	}

	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	// Calculate dimensions
	menuWidth := int(float64(m.width) * 0.35)

	// Build menu
	var menu strings.Builder
	menu.WriteString(titleStyle.Render("📡 Test Commander") + "\n\n")

	for i, cmd := range m.commands {
		cursor := "  "
		style := normalStyle

		if m.cursor == i {
			cursor = "▶ "
			style = selectedStyle
		}

		menu.WriteString(style.Render(fmt.Sprintf("%s%d. %s", cursor, i+1, cmd.name)) + "\n")
		menu.WriteString(normalStyle.Foreground(lipgloss.Color("#666")).Render(fmt.Sprintf("     %s", cmd.description)) + "\n")
	}

	menu.WriteString("\n" + normalStyle.Render("↑/↓: Navigate • Enter: Execute • q: Quit"))

	// Build log section
	logSection := titleStyle.Render("📨 WebSocket Log") + "\n\n" + m.viewport.View()

	// Status bar
	statusBar := ""
	if m.statusMessage != "" {
		if m.statusIsError {
			statusBar = errorStyle.Render("✗ " + m.statusMessage)
		} else {
			statusBar = successStyle.Render("✓ " + m.statusMessage)
		}
	}

	// Combine menu and log side by side
	menuBox := menuStyle.Width(menuWidth).Render(menu.String())
	logBox := logStyle.Width(m.width - menuWidth - 6).Render(logSection)

	content := lipgloss.JoinHorizontal(lipgloss.Top, menuBox, logBox)

	if statusBar != "" {
		content = lipgloss.JoinVertical(lipgloss.Left, content, "\n"+statusBar)
	}

	return content
}

func (m *model) addLog(message string, isError bool) {
	m.logs = append(m.logs, wsMessage{
		timestamp: time.Now(),
		data:      message,
		isError:   isError,
	})

	// Keep only last 100 messages
	if len(m.logs) > 100 {
		m.logs = m.logs[1:]
	}

	m.updateViewport()
}

func (m *model) updateViewport() {
	var content strings.Builder
	content.WriteString("WebSocket Message Log\n")
	content.WriteString(strings.Repeat("─", 80) + "\n\n")

	for _, log := range m.logs {
		timestamp := log.timestamp.Format("15:04:05")
		if log.isError {
			content.WriteString(errorStyle.Render(fmt.Sprintf("[%s] %s", timestamp, log.data)) + "\n")
		} else {
			content.WriteString(infoStyle.Render(fmt.Sprintf("[%s] ", timestamp)) + log.data + "\n")
		}
	}

	if len(m.logs) == 0 {
		content.WriteString(normalStyle.Foreground(lipgloss.Color("#666")).Render("Waiting for messages..."))
	}

	m.viewport.SetContent(content.String())
	m.viewport.GotoBottom()
}

func waitForWebSocketMessage(c *Commander) tea.Cmd {
	return func() tea.Msg {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			return wsMsg{err: err}
		}

		// Try to pretty print JSON
		var prettyJSON bytes.Buffer
		if err := json.Indent(&prettyJSON, message, "", "  "); err == nil {
			return wsMsg{data: prettyJSON.String()}
		}
		return wsMsg{data: string(message)}
	}
}

// Command actions
func sendPing(m *model) tea.Cmd {
	return func() tea.Msg {
		cmd := models.Command{
			Instructions: []models.Instruction{
				{
					Type: "std_ping",
					Content: map[string]interface{}{
						"timestamp": time.Now().Format(time.RFC3339),
					},
				},
			},
			Tags: "test",
		}

		if err := m.commander.SendCommand(cmd); err != nil {
			return statusMsg{message: fmt.Sprintf("Failed to send ping: %v", err), isError: true}
		}
		return statusMsg{message: "Ping command sent", isError: false}
	}
}

func sendNotification(m *model) tea.Cmd {
	return func() tea.Msg {
		cmd := models.Command{
			Instructions: []models.Instruction{
				{
					Type: "std_notification",
					Content: map[string]interface{}{
						"title":    "Test Notification",
						"body":     "This is a test notification from the commander",
						"duration": 5,
					},
				},
			},
			Tags: "test",
		}

		if err := m.commander.SendCommand(cmd); err != nil {
			return statusMsg{message: fmt.Sprintf("Failed to send notification: %v", err), isError: true}
		}
		return statusMsg{message: "Notification sent", isError: false}
	}
}

func sendPopup(m *model) tea.Cmd {
	return func() tea.Msg {
		cmd := models.Command{
			Instructions: []models.Instruction{
				{
					Type: "std_popup",
					Content: map[string]interface{}{
						"body":    "This is a test popup message. Please acknowledge.",
						"title":   "Test Popup",
						"timeout": 30,
					},
				},
			},
			Tags: "test",
		}

		if err := m.commander.SendCommand(cmd); err != nil {
			return statusMsg{message: fmt.Sprintf("Failed: %v", err), isError: true}
		}
		return statusMsg{message: "Popup sent", isError: false}
	}
}

func sendDisplayText(m *model) tea.Cmd {
	return func() tea.Msg {
		cmd := models.Command{
			Instructions: []models.Instruction{
				{
					Type: "std_display_text",
					Content: map[string]interface{}{
						"text":     "Welcome to ControlApp!\n\nThis is a test message.",
						"title":    "Test Display",
						"closable": true,
					},
				},
			},
			Tags: "test",
		}

		if err := m.commander.SendCommand(cmd); err != nil {
			return statusMsg{message: fmt.Sprintf("Failed: %v", err), isError: true}
		}
		return statusMsg{message: "Display text sent", isError: false}
	}
}

func sendTimer(m *model) tea.Cmd {
	return func() tea.Msg {
		cmd := models.Command{
			Instructions: []models.Instruction{
				{
					Type: "std_timer",
					Content: map[string]interface{}{
						"duration":      60,
						"title":         "Test Timer",
						"show_progress": true,
					},
				},
			},
			Tags: "test",
		}

		if err := m.commander.SendCommand(cmd); err != nil {
			return statusMsg{message: fmt.Sprintf("Failed: %v", err), isError: true}
		}
		return statusMsg{message: "Timer started", isError: false}
	}
}

func sendChoice(m *model) tea.Cmd {
	return func() tea.Msg {
		cmd := models.Command{
			Instructions: []models.Instruction{
				{
					Type: "std_choice",
					Content: map[string]interface{}{
						"question": "How is the client working?",
						"options": []map[string]interface{}{
							{"id": "great", "text": "Great! Everything works"},
							{"id": "good", "text": "Good, minor issues"},
							{"id": "bad", "text": "Bad, not working"},
						},
						"timeout": 60,
					},
				},
			},
			Tags: "test",
		}

		if err := m.commander.SendCommand(cmd); err != nil {
			return statusMsg{message: fmt.Sprintf("Failed: %v", err), isError: true}
		}
		return statusMsg{message: "Choice sent", isError: false}
	}
}

func sendOpenURL(m *model) tea.Cmd {
	return func() tea.Msg {
		cmd := models.Command{
			Instructions: []models.Instruction{
				{
					Type: "std_open_url",
					Content: map[string]interface{}{
						"url":     "https://github.com/TheControlApp/server",
						"confirm": true,
					},
				},
			},
			Tags: "test",
		}

		if err := m.commander.SendCommand(cmd); err != nil {
			return statusMsg{message: fmt.Sprintf("Failed: %v", err), isError: true}
		}
		return statusMsg{message: "Open URL sent", isError: false}
	}
}

func sendKinkMessage(m *model) tea.Cmd {
	return func() tea.Msg {
		cmd := models.Command{
			Instructions: []models.Instruction{
				{
					Type: "kink_message",
					Content: map[string]interface{}{
						"message":  "Your commander wants your attention",
						"title":    "Control Message",
						"style":    "info",
						"duration": 5000,
					},
				},
			},
			Tags: "kink",
		}

		if err := m.commander.SendCommand(cmd); err != nil {
			return statusMsg{message: fmt.Sprintf("Failed: %v", err), isError: true}
		}
		return statusMsg{message: "Kink message sent", isError: false}
	}
}

func sendKinkTTS(m *model) tea.Cmd {
	return func() tea.Msg {
		cmd := models.Command{
			Instructions: []models.Instruction{
				{
					Type: "kink_tts",
					Content: map[string]interface{}{
						"text":   "Hello, this is a test of the text to speech system",
						"voice":  "default",
						"volume": 0.8,
					},
				},
			},
			Tags: "kink",
		}

		if err := m.commander.SendCommand(cmd); err != nil {
			return statusMsg{message: fmt.Sprintf("Failed: %v", err), isError: true}
		}
		return statusMsg{message: "TTS sent", isError: false}
	}
}

func sendMultiInstruction(m *model) tea.Cmd {
	return func() tea.Msg {
		cmd := models.Command{
			Instructions: []models.Instruction{
				{
					Type: "std_notification",
					Content: map[string]interface{}{
						"title":    "Multi-Step Command",
						"body":     "Starting sequence...",
						"duration": 3,
					},
				},
				{
					Type: "std_display_text",
					Content: map[string]interface{}{
						"text":     "Step 1: First instruction",
						"title":    "Multi-Command Test",
						"closable": true,
					},
				},
			},
			Tags: "test,multi",
		}

		if err := m.commander.SendCommand(cmd); err != nil {
			return statusMsg{message: fmt.Sprintf("Failed: %v", err), isError: true}
		}
		return statusMsg{message: "Multi-instruction sent", isError: false}
	}
}

// RunTUI starts the Bubble Tea TUI
func RunTUI(commander *Commander) error {
	p := tea.NewProgram(
		initialModel(commander),
		tea.WithAltScreen(),
	)

	_, err := p.Run()
	return err
}
