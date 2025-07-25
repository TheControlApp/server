package models

// Instruction represents a single instruction within a command
type Instruction struct {
	Type    string      `json:"type"`    // The instruction type (popup-msg, download-file, etc.)
	Content interface{} `json:"content"` // Arbitrary struct containing instruction-specific data
}

// InstructionContent represents different types of instruction content,
// this is just an example structure
// type InstructionContent struct {
// 	// For popup-msg
// 	Body   string `json:"body,omitempty"`
// 	Button string `json:"button,omitempty"`

// 	// For download-file
// 	FileHash string `json:"file_hash,omitempty"`
// 	FileName string `json:"file_name,omitempty"`

// 	// For display-text
// 	Text   string `json:"text,omitempty"`
// 	Format string `json:"format,omitempty"`

// 	// For timer
// 	Duration int    `json:"duration,omitempty"`
// 	Title    string `json:"title,omitempty"`

// 	// For notification
// 	// Title and Body already defined above

// 	// For open-url
// 	URL string `json:"url,omitempty"`

// 	// For form-input
// 	Fields   []FormField `json:"fields,omitempty"`
// 	SubmitTo string      `json:"submit_to,omitempty"`

// 	// For announcement
// 	Priority string `json:"priority,omitempty"`

// 	// For custom instructions - allows any additional fields
// 	// The interface{} content field in Instruction handles arbitrary data
// }

// FormField represents a field in a form-input instruction
type FormField struct {
	Name     string   `json:"name"`
	Label    string   `json:"label"`
	Type     string   `json:"type"`
	Options  []string `json:"options,omitempty"`
	Required bool     `json:"required,omitempty"`
}
