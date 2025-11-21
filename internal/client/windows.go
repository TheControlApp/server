package client

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Windows-specific kink command implementations
// These provide the actual platform functionality for consensual kink commands

// RegisterWindowsHandlers registers Windows-specific command handlers
func (c *Client) RegisterWindowsHandlers() {
	c.RegisterCommandHandler("kink_message", c.handleWindowsMessage)
	c.RegisterCommandHandler("kink_open_link", c.handleWindowsOpenLink)
	c.RegisterCommandHandler("kink_download_file", c.handleWindowsDownloadFile)
	c.RegisterCommandHandler("kink_change_wallpaper", c.handleWindowsChangeWallpaper)
	c.RegisterCommandHandler("kink_run_file", c.handleWindowsRunFile)
	c.RegisterCommandHandler("kink_play_audio", c.handleWindowsPlayAudio)
	c.RegisterCommandHandler("kink_tts", c.handleWindowsTTS)
	c.RegisterCommandHandler("kink_popup_image", c.handleWindowsPopupImage)
	c.RegisterCommandHandler("kink_lock_screen", c.handleWindowsLockScreen)
}

// handleWindowsMessage shows a Windows message box
func (c *Client) handleWindowsMessage(cmd Command) CommandResult {
	var payload struct {
		Message string `json:"message"`
		Title   string `json:"title"`
		Style   string `json:"style"` // "info", "warning", "error", "question"
	}

	if err := ParsePayload(cmd.Content, &payload); err != nil {
		return CommandResult{
			Status: "failed",
			Error:  fmt.Sprintf("Invalid payload: %v", err),
		}
	}

	// Default values
	if payload.Title == "" {
		payload.Title = "ControlApp Message"
	}
	if payload.Style == "" {
		payload.Style = "info"
	}

	// Convert style to Windows MessageBox flags
	var flags uint32 = 0x00000000 // MB_OK
	switch payload.Style {
	case "warning":
		flags |= 0x00000030 // MB_ICONWARNING
	case "error":
		flags |= 0x00000010 // MB_ICONERROR
	case "question":
		flags |= 0x00000020 // MB_ICONQUESTION
	default: // info
		flags |= 0x00000040 // MB_ICONINFORMATION
	}

	// Call Windows MessageBox API
	user32 := windows.NewLazyDLL("user32.dll")
	messageBox := user32.NewProc("MessageBoxW")

	titlePtr, _ := windows.UTF16PtrFromString(payload.Title)
	messagePtr, _ := windows.UTF16PtrFromString(payload.Message)

	ret, _, _ := messageBox.Call(
		0, // hwnd
		uintptr(unsafe.Pointer(messagePtr)),
		uintptr(unsafe.Pointer(titlePtr)),
		uintptr(flags),
	)

	return CommandResult{
		Status: "completed",
		Result: map[string]interface{}{
			"button_clicked": int(ret),
			"message_length": len(payload.Message),
		},
	}
}

// handleWindowsOpenLink opens a URL in the default browser
func (c *Client) handleWindowsOpenLink(cmd Command) CommandResult {
	var payload struct {
		URL string `json:"url"`
	}

	if err := ParsePayload(cmd.Content, &payload); err != nil {
		return CommandResult{
			Status: "failed",
			Error:  fmt.Sprintf("Invalid payload: %v", err),
		}
	}

	if payload.URL == "" {
		return CommandResult{
			Status: "failed",
			Error:  "URL is required",
		}
	}

	// Use rundll32 to open URL in default browser
	err := exec.Command("rundll32", "url.dll,FileProtocolHandler", payload.URL).Start()
	if err != nil {
		return CommandResult{
			Status: "failed",
			Error:  fmt.Sprintf("Failed to open URL: %v", err),
		}
	}

	return CommandResult{
		Status: "completed",
		Result: map[string]interface{}{
			"url_opened": payload.URL,
		},
	}
}

// handleWindowsDownloadFile downloads a file to the specified location
func (c *Client) handleWindowsDownloadFile(cmd Command) CommandResult {
	var payload struct {
		URL            string `json:"url"`
		Filename       string `json:"filename"`
		DownloadFolder string `json:"download_folder"`
	}

	if err := ParsePayload(cmd.Content, &payload); err != nil {
		return CommandResult{
			Status: "failed",
			Error:  fmt.Sprintf("Invalid payload: %v", err),
		}
	}

	if payload.URL == "" {
		return CommandResult{
			Status: "failed",
			Error:  "URL is required",
		}
	}

	// Use config download folder if not specified
	if payload.DownloadFolder == "" {
		payload.DownloadFolder = c.config.DownloadFolder
	}

	// Create download folder if it doesn't exist
	if err := os.MkdirAll(payload.DownloadFolder, 0755); err != nil {
		return CommandResult{
			Status: "failed",
			Error:  fmt.Sprintf("Failed to create download folder: %v", err),
		}
	}

	// Generate filename if not provided
	if payload.Filename == "" {
		parts := strings.Split(payload.URL, "/")
		payload.Filename = parts[len(parts)-1]
		if payload.Filename == "" {
			payload.Filename = "download_" + strconv.FormatInt(time.Now().Unix(), 10)
		}
	}

	// Full file path
	filePath := filepath.Join(payload.DownloadFolder, payload.Filename)

	// Use PowerShell to download file (more reliable than curl on Windows)
	psCmd := fmt.Sprintf(`Invoke-WebRequest -Uri "%s" -OutFile "%s"`, payload.URL, filePath)
	execCmd := exec.Command("powershell", "-Command", psCmd)

	if err := execCmd.Run(); err != nil {
		return CommandResult{
			Status: "failed",
			Error:  fmt.Sprintf("Download failed: %v", err),
		}
	}

	// Check if file was created
	if _, err := os.Stat(filePath); err != nil {
		return CommandResult{
			Status: "failed",
			Error:  fmt.Sprintf("File not found after download: %v", err),
		}
	}

	return CommandResult{
		Status: "completed",
		Result: map[string]interface{}{
			"file_path":       filePath,
			"url":             payload.URL,
			"filename":        payload.Filename,
			"download_folder": payload.DownloadFolder,
		},
	}
}

// handleWindowsChangeWallpaper changes the desktop wallpaper
func (c *Client) handleWindowsChangeWallpaper(cmd Command) CommandResult {
	var payload struct {
		ImagePath string `json:"image_path"`
		Style     string `json:"style"` // "fill", "fit", "stretch", "tile", "center"
	}

	if err := ParsePayload(cmd.Content, &payload); err != nil {
		return CommandResult{
			Status: "failed",
			Error:  fmt.Sprintf("Invalid payload: %v", err),
		}
	}

	if payload.ImagePath == "" {
		return CommandResult{
			Status: "failed",
			Error:  "Image path is required",
		}
	}

	// Check if file exists
	if _, err := os.Stat(payload.ImagePath); err != nil {
		return CommandResult{
			Status: "failed",
			Error:  fmt.Sprintf("Image file not found: %v", err),
		}
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(payload.ImagePath)
	if err != nil {
		return CommandResult{
			Status: "failed",
			Error:  fmt.Sprintf("Failed to get absolute path: %v", err),
		}
	}

	// Use Windows API to set wallpaper
	user32 := windows.NewLazyDLL("user32.dll")
	systemParametersInfo := user32.NewProc("SystemParametersInfoW")

	pathPtr, _ := windows.UTF16PtrFromString(absPath)

	ret, _, _ := systemParametersInfo.Call(
		0x0014, // SPI_SETDESKWALLPAPER
		0,
		uintptr(unsafe.Pointer(pathPtr)),
		0x0002, // SPIF_SENDCHANGE
	)

	if ret == 0 {
		return CommandResult{
			Status: "failed",
			Error:  "Failed to set wallpaper via Windows API",
		}
	}

	return CommandResult{
		Status: "completed",
		Result: map[string]interface{}{
			"wallpaper_path": absPath,
			"style":          payload.Style,
		},
	}
}

// handleWindowsRunFile executes a file (with safety checks)
func (c *Client) handleWindowsRunFile(cmd Command) CommandResult {
	var payload struct {
		FilePath  string            `json:"file_path"`
		Arguments []string          `json:"arguments"`
		WorkDir   string            `json:"work_dir"`
		Wait      bool              `json:"wait"`
		Env       map[string]string `json:"env"`
	}

	if err := ParsePayload(cmd.Content, &payload); err != nil {
		return CommandResult{
			Status: "failed",
			Error:  fmt.Sprintf("Invalid payload: %v", err),
		}
	}

	if payload.FilePath == "" {
		return CommandResult{
			Status: "failed",
			Error:  "File path is required",
		}
	}

	// Safety check: require consent for this dangerous command
	if !c.config.IsCommandAllowed("kink_run_file") {
		return CommandResult{
			Status: "failed",
			Error:  "kink_run_file command is not allowed by user consent",
		}
	}

	// Check if file exists
	if _, err := os.Stat(payload.FilePath); err != nil {
		return CommandResult{
			Status: "failed",
			Error:  fmt.Sprintf("File not found: %v", err),
		}
	}

	// Create command
	execCmd := exec.Command(payload.FilePath, payload.Arguments...)

	if payload.WorkDir != "" {
		execCmd.Dir = payload.WorkDir
	}

	// Set environment variables
	if len(payload.Env) > 0 {
		execCmd.Env = os.Environ()
		for k, v := range payload.Env {
			execCmd.Env = append(execCmd.Env, fmt.Sprintf("%s=%s", k, v))
		}
	}

	// Execute
	var err error
	if payload.Wait {
		err = execCmd.Run()
	} else {
		err = execCmd.Start()
	}

	if err != nil {
		return CommandResult{
			Status: "failed",
			Error:  fmt.Sprintf("Failed to execute file: %v", err),
		}
	}

	result := map[string]interface{}{
		"file_path": payload.FilePath,
		"waited":    payload.Wait,
	}

	if execCmd.Process != nil {
		result["process_id"] = execCmd.Process.Pid
	}

	return CommandResult{
		Status: "completed",
		Result: result,
	}
}

// handleWindowsPlayAudio plays an audio file
func (c *Client) handleWindowsPlayAudio(cmd Command) CommandResult {
	var payload struct {
		AudioPath string  `json:"audio_path"`
		Volume    float64 `json:"volume"` // 0.0 to 1.0
		Loop      bool    `json:"loop"`
	}

	if err := ParsePayload(cmd.Content, &payload); err != nil {
		return CommandResult{
			Status: "failed",
			Error:  fmt.Sprintf("Invalid payload: %v", err),
		}
	}

	if payload.AudioPath == "" {
		return CommandResult{
			Status: "failed",
			Error:  "Audio path is required",
		}
	}

	// Check if file exists
	if _, err := os.Stat(payload.AudioPath); err != nil {
		return CommandResult{
			Status: "failed",
			Error:  fmt.Sprintf("Audio file not found: %v", err),
		}
	}

	// Use Windows Media Player or PowerShell to play audio
	// For simplicity, use PowerShell with Windows Media Player COM object
	psCmd := fmt.Sprintf(`
		$player = New-Object -ComObject WMPlayer.OCX
		$playlist = $player.newPlaylist("playlist", @())
		$media = $player.newMedia("%s")
		$playlist.appendItem($media)
		$player.currentPlaylist = $playlist
		$player.controls.play()
	`, payload.AudioPath)

	execCmd := exec.Command("powershell", "-Command", psCmd)

	if err := execCmd.Start(); err != nil {
		return CommandResult{
			Status: "failed",
			Error:  fmt.Sprintf("Failed to play audio: %v", err),
		}
	}

	return CommandResult{
		Status: "completed",
		Result: map[string]interface{}{
			"audio_path": payload.AudioPath,
			"volume":     payload.Volume,
			"loop":       payload.Loop,
		},
	}
}

// handleWindowsTTS performs text-to-speech
func (c *Client) handleWindowsTTS(cmd Command) CommandResult {
	var payload struct {
		Text   string `json:"text"`
		Voice  string `json:"voice"`
		Rate   int    `json:"rate"`   // -10 to 10
		Volume int    `json:"volume"` // 0 to 100
	}

	if err := ParsePayload(cmd.Content, &payload); err != nil {
		return CommandResult{
			Status: "failed",
			Error:  fmt.Sprintf("Invalid payload: %v", err),
		}
	}

	if payload.Text == "" {
		return CommandResult{
			Status: "failed",
			Error:  "Text is required",
		}
	}

	// Default values
	if payload.Volume == 0 {
		payload.Volume = 80
	}

	// Use PowerShell with SAPI for TTS
	psCmd := fmt.Sprintf(`
		$voice = New-Object -ComObject SAPI.SpVoice
		$voice.Volume = %d
		$voice.Rate = %d
		$voice.Speak("%s")
	`, payload.Volume, payload.Rate, strings.ReplaceAll(payload.Text, "\"", "`\""))

	execCmd := exec.Command("powershell", "-Command", psCmd)

	if err := execCmd.Run(); err != nil {
		return CommandResult{
			Status: "failed",
			Error:  fmt.Sprintf("TTS failed: %v", err),
		}
	}

	return CommandResult{
		Status: "completed",
		Result: map[string]interface{}{
			"text":   payload.Text,
			"voice":  payload.Voice,
			"rate":   payload.Rate,
			"volume": payload.Volume,
		},
	}
}

// handleWindowsPopupImage displays an image in a popup window
func (c *Client) handleWindowsPopupImage(cmd Command) CommandResult {
	var payload struct {
		ImagePath string `json:"image_path"`
		Title     string `json:"title"`
		Duration  int    `json:"duration"` // seconds, 0 for indefinite
	}

	if err := ParsePayload(cmd.Content, &payload); err != nil {
		return CommandResult{
			Status: "failed",
			Error:  fmt.Sprintf("Invalid payload: %v", err),
		}
	}

	if payload.ImagePath == "" {
		return CommandResult{
			Status: "failed",
			Error:  "Image path is required",
		}
	}

	// Check if file exists
	if _, err := os.Stat(payload.ImagePath); err != nil {
		return CommandResult{
			Status: "failed",
			Error:  fmt.Sprintf("Image file not found: %v", err),
		}
	}

	// Use Windows Photo Viewer or default image viewer
	err := exec.Command("rundll32.exe", "shimgvw.dll,ImageView_Fullscreen", payload.ImagePath).Start()
	if err != nil {
		return CommandResult{
			Status: "failed",
			Error:  fmt.Sprintf("Failed to display image: %v", err),
		}
	}

	return CommandResult{
		Status: "completed",
		Result: map[string]interface{}{
			"image_path": payload.ImagePath,
			"title":      payload.Title,
		},
	}
}

// handleWindowsLockScreen locks the Windows screen
func (c *Client) handleWindowsLockScreen(cmd Command) CommandResult {
	// Lock the workstation
	user32 := windows.NewLazyDLL("user32.dll")
	lockWorkStation := user32.NewProc("LockWorkStation")

	ret, _, _ := lockWorkStation.Call()
	if ret == 0 {
		return CommandResult{
			Status: "failed",
			Error:  "Failed to lock workstation",
		}
	}

	return CommandResult{
		Status: "completed",
		Result: map[string]interface{}{
			"locked_at": time.Now().Format(time.RFC3339),
		},
	}
}
