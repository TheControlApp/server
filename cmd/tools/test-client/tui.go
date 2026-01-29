package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
	"github.com/thecontrolapp/server/internal/models"
)

type model struct {
	client      *Client
	viewport    viewport.Model
	mockDisplay viewport.Model
	logs        []logMessage
	width       int
	height      int
	lastCommand *models.Command
	logFile     *os.File
	ready       bool
}

type logMessage struct {
	timestamp time.Time
	message   string
	isError   bool
}

type wsMsg struct {
	data string
	err  error
}

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00FF00")).
			MarginBottom(1)

	logStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 2)

	displayStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FF6B6B")).
			Padding(1, 2)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00"))

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00BFFF"))
)

func initialModel(client *Client) model {
	// Open log file
	logFile, err := os.OpenFile("test-client.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		logFile = nil
	}

	m := model{
		client:  client,
		logs:    []logMessage{},
		logFile: logFile,
	}

	// Request pending commands on connect
	m.requestPendingCommands()

	return m
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		waitForWebSocketMessage(m.client),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		if !m.ready {
			// Initialize viewports
			logWidth := (m.width - 8) / 2
			displayWidth := m.width - logWidth - 8

			m.viewport = viewport.New(logWidth, m.height-8)
			m.mockDisplay = viewport.New(displayWidth, m.height-8)
			m.ready = true
		} else {
			// Update viewport sizes
			logWidth := (m.width - 8) / 2
			displayWidth := m.width - logWidth - 8

			m.viewport.Width = logWidth
			m.viewport.Height = m.height - 8
			m.mockDisplay.Width = displayWidth
			m.mockDisplay.Height = m.height - 8
		}

		m.updateViewports()

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			if m.logFile != nil {
				m.logFile.Close()
			}
			return m, tea.Quit
		}

	case wsMsg:
		if msg.err != nil {
			m.addLog(fmt.Sprintf("WebSocket error: %v", msg.err), true)
			return m, tea.Quit
		}

		// Process the received message
		m.processCommand(msg.data)
		return m, waitForWebSocketMessage(m.client)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	// Build log section
	logSection := titleStyle.Render("📨 Received Commands Log") + "\n\n" + m.viewport.View()

	// Build mock display section
	displaySection := titleStyle.Render("🖥️  Mock Client Display") + "\n\n" + m.mockDisplay.View()

	// Combine log and display side by side
	logBox := logStyle.Width(m.viewport.Width + 4).Render(logSection)
	displayBox := displayStyle.Width(m.mockDisplay.Width + 4).Render(displaySection)

	content := lipgloss.JoinHorizontal(lipgloss.Top, logBox, displayBox)

	footer := "\n" + infoStyle.Render("Press q to quit")

	return content + footer
}

func (m *model) addLog(message string, isError bool) {
	logMsg := logMessage{
		timestamp: time.Now(),
		message:   message,
		isError:   isError,
	}
	m.logs = append(m.logs, logMsg)

	// Write to log file
	if m.logFile != nil {
		logType := "INFO"
		if isError {
			logType = "ERROR"
		}
		fmt.Fprintf(m.logFile, "[%s] %s: %s\n", logMsg.timestamp.Format("2006-01-02 15:04:05"), logType, message)
	}

	// Keep only last 100 messages
	if len(m.logs) > 100 {
		m.logs = m.logs[1:]
	}

	m.updateViewports()
}

func (m *model) updateViewports() {
	// Update log viewport
	var logContent strings.Builder
	logContent.WriteString("Command Log\n")
	logContent.WriteString(strings.Repeat("─", 60) + "\n\n")

	for _, log := range m.logs {
		timestamp := log.timestamp.Format("15:04:05")
		if log.isError {
			logContent.WriteString(errorStyle.Render(fmt.Sprintf("[%s] %s", timestamp, log.message)) + "\n")
		} else {
			logContent.WriteString(infoStyle.Render(fmt.Sprintf("[%s] ", timestamp)) + log.message + "\n")
		}
	}

	if len(m.logs) == 0 {
		logContent.WriteString(infoStyle.Render("Waiting for commands..."))
	}

	m.viewport.SetContent(logContent.String())
	m.viewport.GotoBottom()

	// Update mock display viewport
	m.updateMockDisplay()
}

func (m *model) updateMockDisplay() {
	var content strings.Builder

	if m.lastCommand == nil {
		content.WriteString(infoStyle.Render("No command received yet\n\n"))
		content.WriteString("Waiting for commands from the server...")
	} else {
		// Display the last command in a user-friendly format
		content.WriteString(successStyle.Render("✓ Command Received\n"))
		content.WriteString(strings.Repeat("─", 60) + "\n\n")

		for i, instruction := range m.lastCommand.Instructions {
			content.WriteString(fmt.Sprintf("Instruction %d: %s\n", i+1, instruction.Type))
			content.WriteString(strings.Repeat("─", 40) + "\n")

			// Mock display based on instruction type
			switch instruction.Type {
			case "std_ping":
				content.WriteString(successStyle.Render("🏓 PING!\n"))
				content.WriteString("Responding to ping request...\n")

			case "std_notification":
				content.WriteString("📢 NOTIFICATION\n")
				if contentMap, ok := instruction.Content.(map[string]interface{}); ok {
					if title, ok := contentMap["title"].(string); ok {
						content.WriteString(fmt.Sprintf("Title: %s\n", title))
					}
					if body, ok := contentMap["body"].(string); ok {
						content.WriteString(fmt.Sprintf("Body: %s\n", body))
					}
					if duration, ok := contentMap["duration"].(float64); ok {
						content.WriteString(fmt.Sprintf("Duration: %.0fs\n", duration))
					}
				}

			case "std_popup":
				content.WriteString("💬 POPUP\n")
				if contentMap, ok := instruction.Content.(map[string]interface{}); ok {
					if title, ok := contentMap["title"].(string); ok {
						content.WriteString(fmt.Sprintf("Title: %s\n", title))
					}
					if message, ok := contentMap["message"].(string); ok {
						content.WriteString(fmt.Sprintf("Message: %s\n", message))
					}
				}

			case "std_display_text":
				content.WriteString("📝 DISPLAY TEXT\n")
				if contentMap, ok := instruction.Content.(map[string]interface{}); ok {
					if title, ok := contentMap["title"].(string); ok {
						content.WriteString(fmt.Sprintf("Title: %s\n", title))
					}
					if text, ok := contentMap["text"].(string); ok {
						content.WriteString(fmt.Sprintf("Text: %s\n", text))
					}
				}

			case "std_timer":
				content.WriteString("⏱️  TIMER\n")
				if contentMap, ok := instruction.Content.(map[string]interface{}); ok {
					if duration, ok := contentMap["duration"].(float64); ok {
						content.WriteString(fmt.Sprintf("Duration: %.0fs\n", duration))
					}
					if title, ok := contentMap["title"].(string); ok {
						content.WriteString(fmt.Sprintf("Title: %s\n", title))
					}
				}

			case "std_choice":
				content.WriteString("❓ CHOICE PROMPT\n")
				if contentMap, ok := instruction.Content.(map[string]interface{}); ok {
					if question, ok := contentMap["question"].(string); ok {
						content.WriteString(fmt.Sprintf("Question: %s\n", question))
					}
					if options, ok := contentMap["options"].([]interface{}); ok {
						content.WriteString("Options:\n")
						for _, opt := range options {
							if optMap, ok := opt.(map[string]interface{}); ok {
								if text, ok := optMap["text"].(string); ok {
									content.WriteString(fmt.Sprintf("  - %s\n", text))
								}
							}
						}
					}
				}

			case "std_open_url":
				content.WriteString("🌐 OPEN URL\n")
				if contentMap, ok := instruction.Content.(map[string]interface{}); ok {
					if url, ok := contentMap["url"].(string); ok {
						content.WriteString(fmt.Sprintf("URL: %s\n", url))
					}
				}

			case "kink_message":
				content.WriteString("💜 KINK MESSAGE\n")
				if contentMap, ok := instruction.Content.(map[string]interface{}); ok {
					if title, ok := contentMap["title"].(string); ok {
						content.WriteString(fmt.Sprintf("Title: %s\n", title))
					}
					if message, ok := contentMap["message"].(string); ok {
						content.WriteString(fmt.Sprintf("Message: %s\n", message))
					}
				}

			case "kink_tts":
				content.WriteString("🔊 TEXT TO SPEECH\n")
				if contentMap, ok := instruction.Content.(map[string]interface{}); ok {
					if text, ok := contentMap["text"].(string); ok {
						content.WriteString(fmt.Sprintf("Text: %s\n", text))
					}
				}

			default:
				content.WriteString(fmt.Sprintf("⚠️  Unknown instruction type: %s\n", instruction.Type))
				if instruction.Content != nil {
					contentJSON, _ := json.MarshalIndent(instruction.Content, "", "  ")
					content.WriteString(fmt.Sprintf("Content: %s\n", string(contentJSON)))
				}
			}

			content.WriteString("\n")
		}
	}

	m.mockDisplay.SetContent(content.String())
}

func (m *model) processCommand(data string) {
	// Try to pretty print JSON first
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(data), "", "  "); err == nil {
		m.addLog(fmt.Sprintf("📥 Raw: %s", prettyJSON.String()), false)
	} else {
		m.addLog(fmt.Sprintf("📥 Raw: %s", data), false)
	}

	// Try to parse as Command
	var cmd models.Command
	if err := json.Unmarshal([]byte(data), &cmd); err != nil {
		m.addLog(fmt.Sprintf("⚠️  Failed to parse command: %v", err), true)
		return
	}

	m.lastCommand = &cmd
	m.addLog(fmt.Sprintf("✓ Parsed command with %d instruction(s)", len(cmd.Instructions)), false)
	m.updateViewports()

	// Send acknowledgment to server
	m.sendCommandAck(cmd.ID)
}

// requestPendingCommands sends a message to server requesting all pending commands
func (m *model) requestPendingCommands() {
	msg := map[string]string{
		"type": "fetch_pending_commands",
	}

	data, err := json.Marshal(msg)
	if err != nil {
		m.addLog(fmt.Sprintf("Failed to marshal fetch request: %v", err), true)
		return
	}

	if err := m.client.conn.WriteMessage(1, data); err != nil {
		m.addLog(fmt.Sprintf("Failed to send fetch request: %v", err), true)
	} else {
		m.addLog("📬 Requested pending commands from server", false)
	}
}

// sendCommandAck sends acknowledgment to server that command was received
func (m *model) sendCommandAck(commandID uuid.UUID) {
	ack := map[string]string{
		"type":       "command_ack",
		"command_id": commandID.String(),
	}

	data, err := json.Marshal(ack)
	if err != nil {
		m.addLog(fmt.Sprintf("Failed to marshal ack: %v", err), true)
		return
	}

	if err := m.client.conn.WriteMessage(1, data); err != nil {
		m.addLog(fmt.Sprintf("Failed to send ack: %v", err), true)
	} else {
		m.addLog(fmt.Sprintf("✅ Acknowledged command %s", commandID.String()[:8]), false)
	}
}

func waitForWebSocketMessage(c *Client) tea.Cmd {
	return func() tea.Msg {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			return wsMsg{err: err}
		}

		return wsMsg{data: string(message)}
	}
}

func RunTUI(client *Client) error {
	p := tea.NewProgram(
		initialModel(client),
		tea.WithAltScreen(),
	)

	_, err := p.Run()
	return err
}
