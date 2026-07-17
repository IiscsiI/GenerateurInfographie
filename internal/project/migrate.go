package project

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ImportFromJSON parses a JSON file and migrates if needed.
func ImportFromJSON(data []byte) (*Project, error) {
	// First, check if it's old format (has "version" but not "_version")
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("JSON invalide: %w", err)
	}

	_, hasVersion := raw["version"]
	_, hasNewVersion := raw["_version"]

	if hasVersion && !hasNewVersion {
		return migrateV15(data)
	}

	// Modern v4.0 format
	var p Project
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("structure invalide: %w", err)
	}

	if err := validate(&p); err != nil {
		return nil, err
	}

	return &p, nil
}

// ImportFromHTML extracts project data from an exported HTML file.
func ImportFromHTML(html string) (*Project, error) {
	// Look for <script id="project-data" type="application/json">...</script>
	marker := `id="project-data"`
	idx := strings.Index(html, marker)
	if idx == -1 {
		return nil, fmt.Errorf("donnees projet introuvables dans le HTML")
	}

	// Find the opening > after the marker
	start := strings.Index(html[idx:], ">")
	if start == -1 {
		return nil, fmt.Errorf("balise script malformee")
	}
	start += idx + 1

	// Find closing </script>
	end := strings.Index(html[start:], "</script>")
	if end == -1 {
		return nil, fmt.Errorf("balise script non fermee")
	}
	end += start

	jsonData := strings.TrimSpace(html[start:end])
	return ImportFromJSON([]byte(jsonData))
}

func validate(p *Project) error {
	if p.Metadata.Name == "" {
		p.Metadata.Name = "Projet sans nom"
	}
	if p.Content.Title == "" {
		p.Content.Title = "Sans titre"
	}
	if p.Theme.Colors.Header == "" {
		p.Theme.Colors.Header = "#dc3545"
	}
	if p.Elements.Logos == nil {
		p.Elements.Logos = []Logo{}
	}
	if p.Elements.Timeline == nil {
		p.Elements.Timeline = []TimelineItem{}
	}
	if p.Elements.Steps == nil {
		p.Elements.Steps = []Step{}
	}
	if p.Theme.FontSize == 0 {
		p.Theme.FontSize = 100
	}
	if p.Theme.FontFamily == "" {
		p.Theme.FontFamily = "system"
	}
	return nil
}

// migrateV15 converts legacy v1.5 format to v4.0.
func migrateV15(data []byte) (*Project, error) {
	var old map[string]json.RawMessage
	if err := json.Unmarshal(data, &old); err != nil {
		return nil, err
	}

	p := DefaultProject()
	p.Metadata.Name = "Projet migre v1.5"

	// Extract project name
	if raw, ok := old["projectName"]; ok {
		var name string
		json.Unmarshal(raw, &name)
		if name != "" {
			p.Metadata.Name = name
		}
	}

	// Extract settings
	if raw, ok := old["settings"]; ok {
		var settings map[string]json.RawMessage
		if json.Unmarshal(raw, &settings) == nil {
			extractString(settings, "mainTitle", &p.Content.Title)
			extractString(settings, "mainSubtitle", &p.Content.Subtitle)
			extractString(settings, "emergencyMessage", &p.Content.EmergencyMessage)
			extractString(settings, "footerMessage", &p.Content.FooterMessage)

			// Colors
			if colorsRaw, ok := settings["colors"]; ok {
				var colors map[string]string
				if json.Unmarshal(colorsRaw, &colors) == nil {
					if v, ok := colors["header"]; ok { p.Theme.Colors.Header = v }
					if v, ok := colors["timeline"]; ok { p.Theme.Colors.Timeline = v }
					if v, ok := colors["stepNumber"]; ok { p.Theme.Colors.StepNum = v }
					if v, ok := colors["footer"]; ok { p.Theme.Colors.Footer = v }
					if v, ok := colors["emergency"]; ok { p.Theme.Colors.Emergency = v }
					if v, ok := colors["bodyBg"]; ok { p.Theme.Colors.BgColor = v }
					if v, ok := colors["previewBg"]; ok && p.Theme.Colors.BgColor == "#f8f9fa" { p.Theme.Colors.BgColor = v }
					if v, ok := colors["contentBg"]; ok { p.Theme.Colors.ContentBg = v }
				}
			}
		}
	}

	// Extract timeline
	if raw, ok := old["timeline"]; ok {
		var items []map[string]interface{}
		if json.Unmarshal(raw, &items) == nil {
			for i, item := range items {
				ti := TimelineItem{ID: i + 1}
				if v, ok := item["text"].(string); ok { ti.Text = v }
				if v, ok := item["type"].(string); ok { ti.Type = v }
				if ti.Type == "" { ti.Type = "item" }
				p.Elements.Timeline = append(p.Elements.Timeline, ti)
			}
		}
	}

	// Extract steps
	if raw, ok := old["steps"]; ok {
		var steps []json.RawMessage
		if json.Unmarshal(raw, &steps) == nil {
			idCounter := 100
			for i, stepRaw := range steps {
				var stepMap map[string]json.RawMessage
				if json.Unmarshal(stepRaw, &stepMap) != nil {
					continue
				}

				step := Step{
					ID:       idCounter,
					Number:   i + 1,
					Icon:     "📋",
					Title:    "Etape",
					Category: "management",
					CustomColors: &StepColors{},
					Actions: []Action{},
				}
				idCounter++

				extractInt(stepMap, "number", &step.Number)
				extractString2(stepMap, "icon", &step.Icon)
				extractString2(stepMap, "title", &step.Title)
				extractString2(stepMap, "category", &step.Category)

				// Actions
				if actionsRaw, ok := stepMap["actions"]; ok {
					var actions []map[string]json.RawMessage
					if json.Unmarshal(actionsRaw, &actions) == nil {
						for _, actMap := range actions {
							action := Action{ID: idCounter, Type: "info"}
							idCounter++
							extractString2(actMap, "type", &action.Type)
							extractString2(actMap, "text", &action.Text)
							if ccRaw, ok := actMap["customColors"]; ok {
								action.CustomColors = ccRaw
							} else {
								action.CustomColors = json.RawMessage("null")
							}
							step.Actions = append(step.Actions, action)
						}
					}
				}

				p.Elements.Steps = append(p.Elements.Steps, step)
			}
		}
	}

	return &p, nil
}

func extractString(m map[string]json.RawMessage, key string, target *string) {
	if raw, ok := m[key]; ok {
		var v string
		if json.Unmarshal(raw, &v) == nil && v != "" {
			*target = v
		}
	}
}

func extractString2(m map[string]json.RawMessage, key string, target *string) {
	extractString(m, key, target)
}

func extractInt(m map[string]json.RawMessage, key string, target *int) {
	if raw, ok := m[key]; ok {
		var v float64
		if json.Unmarshal(raw, &v) == nil {
			*target = int(v)
		}
	}
}
