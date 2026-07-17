package project

import "encoding/json"

// Project is the root data structure, matching JSON v4.0.
type Project struct {
	Version   string   `json:"_version,omitempty"`
	Generated string   `json:"_generated,omitempty"`
	Generator string   `json:"_generator,omitempty"`
	Metadata  Metadata `json:"metadata"`
	Content   Content  `json:"content"`
	Theme     Theme    `json:"theme"`
	Elements  Elements `json:"elements"`
}

// Metadata holds project metadata.
type Metadata struct {
	Name         string `json:"name"`
	Author       string `json:"author"`
	Organization string `json:"organization"`
	License      string `json:"license"`
}

// Content holds text content fields.
type Content struct {
	Title            string `json:"title"`
	Subtitle         string `json:"subtitle"`
	ShowHeaderIcon   bool   `json:"showHeaderIcon"`
	EmergencyMessage string `json:"emergencyMessage"`
	FooterMessage    string `json:"footerMessage"`
}

// Theme holds color settings and typography.
type Theme struct {
	Colors     Colors `json:"colors"`
	FontSize   int    `json:"fontSize,omitempty"`
	FontFamily string `json:"fontFamily,omitempty"`
}

// Colors holds the color palette.
type Colors struct {
	Header    string `json:"header"`
	Timeline  string `json:"timeline"`
	StepNum   string `json:"stepNumber"`
	Footer    string `json:"footer"`
	Emergency string `json:"emergency"`
	BgColor   string `json:"background"`
	ContentBg string `json:"contentBg"`
}

// Elements holds logos, timeline, and steps.
type Elements struct {
	Logos    []Logo         `json:"logos"`
	Timeline []TimelineItem `json:"timeline"`
	Steps    []Step         `json:"steps"`
}

// Logo represents a logo element.
type Logo struct {
	ID       int    `json:"id"`
	URL      string `json:"url,omitempty"`
	File     string `json:"file,omitempty"`
	Position string `json:"position"`
	Size     int    `json:"size"`
}

// TimelineItem represents a timeline element.
type TimelineItem struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Type string `json:"type"` // "item" or "separator"
}

// Step represents a process step.
type Step struct {
	ID           int              `json:"id"`
	Number       int              `json:"number"`
	Icon         string           `json:"icon"`
	Title        string           `json:"title"`
	Category     string           `json:"category"`
	CustomColors *StepColors      `json:"customColors"`
	Actions      []Action         `json:"actions"`
}

// StepColors holds optional custom colors for a step.
type StepColors struct {
	Border   *string `json:"border"`
	Bg       *string `json:"background"`
	Text     *string `json:"text"`
	ActionBg *string `json:"actionBg"`
}

// Action represents an action within a step.
type Action struct {
	ID           int              `json:"id"`
	Type         string           `json:"type"` // "critical", "important", "info", "success"
	Text         string           `json:"text"`
	CustomColors json.RawMessage  `json:"customColors"` // null or object
}

// DefaultProject returns a new project with default content.
func DefaultProject() Project {
	return Project{
		Version:   "4.0",
		Generator: "Generateur Infographie v2.0 (Go)",
		Metadata: Metadata{
			Name:    "Nouveau Projet",
			Author:  "",
			License: "AGPL-3.0 / Commercial",
		},
		Content: Content{
			Title:            "CYBERATTAQUE DETECTEE",
			Subtitle:         "Procedure d'urgence - Actions immediates",
			ShowHeaderIcon:   true,
			EmergencyMessage: "EN CAS D'URGENCE : Contactez immediatement votre RSSI",
			FooterMessage:    "RAPPEL IMPORTANT : Chaque minute compte en cas de cyberattaque",
		},
		Theme: Theme{
			Colors: Colors{
				Header:    "#dc3545",
				Timeline:  "#3498db",
				StepNum:   "#0056b3",
				Footer:    "#2c3e50",
				Emergency: "#ff6b6b",
				BgColor:   "#f8f9fa",
				ContentBg: "#ffffff",
			},
			FontSize:   100,
			FontFamily: "system",
		},
		Elements: Elements{
			Logos:    []Logo{},
			Timeline: []TimelineItem{},
			Steps:    []Step{},
		},
	}
}
