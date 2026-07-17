package web

import (
	"embed"
	"fmt"
	"html"
	"io/fs"
	"regexp"
	"strings"

	"infographic-generator/internal/project"
)

// safeFontFamilyRE validates that a font-family CSS value contains only
// expected characters. Blocks injection attempts like `monospace;background:url(...)`.
var safeFontFamilyRE = regexp.MustCompile(`^[a-zA-Z0-9\s,'\-]+$`)

//go:embed editor/*
var editorFiles embed.FS

//go:embed render/*
var renderFiles embed.FS

// EditorFS returns the filesystem for serving the editor frontend.
func EditorFS() fs.FS {
	sub, _ := fs.Sub(editorFiles, "editor")
	return sub
}

// RenderTemplate returns a complete HTML page that renders the infographic
// from injected JSON data. Used by chromedp for export capture.
func RenderTemplate(projectJSON string) string {
	renderHTML, _ := fs.ReadFile(renderFiles, "render/index.html")
	// Replace the full placeholder including the fallback empty object
	return strings.Replace(string(renderHTML), "\"%%PROJECT_DATA%%\"", projectJSON, 1)
}

// RenderCSS returns the CSS used for infographic rendering.
func RenderCSS() string {
	return infographicCSS
}

// BuildStaticInfographic generates a static HTML representation of the infographic.
func BuildStaticInfographic(p *project.Project) string {
	var sb strings.Builder

	headerClass := "preview-header"
	if !p.Content.ShowHeaderIcon {
		headerClass += " hide-icon"
	}

	fontSize := p.Theme.FontSize
	if fontSize == 0 {
		fontSize = 100
	}
	fontFamilyStyle := ""
	if p.Theme.FontFamily != "" && p.Theme.FontFamily != "system" {
		// Sanitize font-family value: allow letters, digits, space, comma, apostrophe, hyphen
		if safeFontFamilyRE.MatchString(p.Theme.FontFamily) {
			fontFamilyStyle = fmt.Sprintf("; font-family: %s", p.Theme.FontFamily)
		}
	}

	// Container
	sb.WriteString(fmt.Sprintf(`<div class="preview-container" style="background: %s; font-size: %d%%%s">`, p.Theme.Colors.ContentBg, fontSize, fontFamilyStyle))

	// Header
	sb.WriteString(fmt.Sprintf(`<div class="%s" style="background: linear-gradient(135deg, %s, %s)">`,
		headerClass, p.Theme.Colors.Header, adjustBrightness(p.Theme.Colors.Header, -20)))

	// Header logos
	for _, logo := range p.Elements.Logos {
		src := logo.File
		if src == "" {
			src = logo.URL
		}
		if src != "" && strings.HasPrefix(logo.Position, "top") {
			sb.WriteString(fmt.Sprintf(`<img src="%s" class="logo %s" style="max-width: %dpx" alt="Logo">`,
				html.EscapeString(src), logo.Position, logo.Size))
		}
	}

	sb.WriteString(fmt.Sprintf(`<h1>%s</h1>`, project.SanitizeHTML(p.Content.Title)))
	sb.WriteString(fmt.Sprintf(`<p>%s</p>`, project.SanitizeHTML(p.Content.Subtitle)))
	sb.WriteString("</div>")

	// Timeline
	if len(p.Elements.Timeline) > 0 {
		sb.WriteString(`<div class="preview-timeline">`)
		for _, item := range p.Elements.Timeline {
			if item.Type == "item" {
				sb.WriteString(fmt.Sprintf(`<div class="preview-timeline-item" style="background: linear-gradient(135deg, %s, %s)">%s</div>`,
					p.Theme.Colors.Timeline, adjustBrightness(p.Theme.Colors.Timeline, -15),
					html.EscapeString(item.Text)))
			} else {
				sb.WriteString(fmt.Sprintf(`<div class="preview-timeline-item" style="background: transparent; color: %s; font-size: 1.2em;">%s</div>`,
					p.Theme.Colors.Timeline, html.EscapeString(item.Text)))
			}
		}
		sb.WriteString("</div>")
	}

	// Steps grid
	sb.WriteString(fmt.Sprintf(`<div class="preview-steps-grid" style="background: %s">`, p.Theme.Colors.ContentBg))
	for i, step := range p.Elements.Steps {
		borderColor := getStepColor(step, "border")
		bgColor := getStepColor(step, "background")
		textColor := getStepColor(step, "text")

		sb.WriteString(fmt.Sprintf(`<div class="preview-step" style="border-left-color: %s; background: %s; color: %s">`,
			borderColor, bgColor, textColor))
		sb.WriteString(fmt.Sprintf(`<div class="preview-step-number" style="background: linear-gradient(135deg, %s, %s)">%d</div>`,
			p.Theme.Colors.StepNum, adjustBrightness(p.Theme.Colors.StepNum, -20), i+1))
		sb.WriteString(fmt.Sprintf(`<div class="preview-step-icon">%s</div>`, step.Icon))
		sb.WriteString(fmt.Sprintf(`<h3 style="color: %s">%s</h3>`, textColor, html.EscapeString(step.Title)))

		sb.WriteString("<div>")
		for _, action := range step.Actions {
			aBorder := getActionBorderColor(action)
			aBg := getActionBgColor(action)
			aText := getActionTextColor(action)

			sb.WriteString(fmt.Sprintf(`<div class="preview-action-item" style="border-left-color: %s; background: %s; color: %s">%s</div>`,
				aBorder, aBg, aText, project.SanitizeHTML(action.Text)))
		}
		sb.WriteString("</div></div>")
	}
	sb.WriteString("</div>")

	// Emergency
	if p.Content.EmergencyMessage != "" {
		sb.WriteString(fmt.Sprintf(`<div class="preview-emergency-contact" style="background: linear-gradient(135deg, %s, %s)">%s</div>`,
			p.Theme.Colors.Emergency, adjustBrightness(p.Theme.Colors.Emergency, -20),
			project.SanitizeHTML(p.Content.EmergencyMessage)))
	}

	// Footer
	sb.WriteString(fmt.Sprintf(`<div class="preview-footer" style="background: linear-gradient(135deg, %s, %s)">`,
		p.Theme.Colors.Footer, adjustBrightness(p.Theme.Colors.Footer, -15)))

	for _, logo := range p.Elements.Logos {
		src := logo.File
		if src == "" {
			src = logo.URL
		}
		if src != "" && strings.HasPrefix(logo.Position, "bottom") {
			sb.WriteString(fmt.Sprintf(`<img src="%s" class="logo %s" style="max-width: %dpx" alt="Logo">`,
				html.EscapeString(src), logo.Position, logo.Size))
		}
	}

	// Footer message: sanitize once, newlines become <br>
	if p.Content.FooterMessage != "" {
		msg := strings.ReplaceAll(p.Content.FooterMessage, "\n", "<br>")
		sb.WriteString(fmt.Sprintf("<div>%s</div>", project.SanitizeHTML(msg)))
	}
	sb.WriteString("</div>")

	sb.WriteString("</div>") // close preview-container

	return sb.String()
}

// --- Color resolution helpers (mirroring JS ColorManager logic) ---

var categoryPresets = map[string]map[string]string{
	"immediate":     {"border": "#e74c3c", "background": "#ffffff", "text": "#2c3e50"},
	"management":    {"border": "#0056b3", "background": "#ffffff", "text": "#2c3e50"},
	"communication": {"border": "#f39c12", "background": "#ffffff", "text": "#2c3e50"},
	"continuity":    {"border": "#27ae60", "background": "#ffffff", "text": "#2c3e50"},
}

var actionTypePresets = map[string]map[string]string{
	"critical":  {"border": "#e74c3c", "background": "rgba(231, 76, 60, 0.08)", "text": "#2c3e50"},
	"important": {"border": "#f39c12", "background": "rgba(243, 156, 18, 0.08)", "text": "#2c3e50"},
	"info":      {"border": "#3498db", "background": "rgba(52, 152, 219, 0.08)", "text": "#2c3e50"},
	"success":   {"border": "#27ae60", "background": "rgba(39, 174, 96, 0.08)", "text": "#2c3e50"},
}

func getStepColor(step project.Step, colorType string) string {
	if step.CustomColors != nil {
		var val *string
		switch colorType {
		case "border":
			val = step.CustomColors.Border
		case "background":
			val = step.CustomColors.Bg
		case "text":
			val = step.CustomColors.Text
		}
		if val != nil && *val != "" {
			return *val
		}
	}
	if preset, ok := categoryPresets[step.Category]; ok {
		if c, ok := preset[colorType]; ok {
			return c
		}
	}
	defaults := map[string]string{"border": "#3498db", "background": "#ffffff", "text": "#2c3e50"}
	return defaults[colorType]
}

func getActionBorderColor(a project.Action) string {
	if preset, ok := actionTypePresets[a.Type]; ok {
		return preset["border"]
	}
	return "#3498db"
}

func getActionBgColor(a project.Action) string {
	if preset, ok := actionTypePresets[a.Type]; ok {
		return preset["background"]
	}
	return "rgba(52, 152, 219, 0.08)"
}

func getActionTextColor(a project.Action) string {
	return "#2c3e50"
}

func adjustBrightness(hexColor string, amount int) string {
	if len(hexColor) != 7 || hexColor[0] != '#' {
		return hexColor
	}
	var r, g, b int
	fmt.Sscanf(hexColor, "#%02x%02x%02x", &r, &g, &b)
	r = clamp(r+amount, 0, 255)
	g = clamp(g+amount, 0, 255)
	b = clamp(b+amount, 0, 255)
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

func clamp(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// infographicCSS is the complete CSS for rendering the infographic.
// Shared between the render page and the standalone HTML export.
const infographicCSS = `
* { margin: 0; padding: 0; box-sizing: border-box; }
body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; line-height: 1.6; color: #2c3e50; background: #f8f9fa; }

.preview-container { max-width: 1200px; margin: 0 auto; border-radius: 20px; box-shadow: 0 20px 40px rgba(0,0,0,0.15); overflow: hidden; }
.preview-header { color: white; text-align: center; padding: 50px 30px; position: relative; }
.preview-header::before { content: '\26A0\FE0F'; font-size: 80px; display: block; margin-bottom: 25px; }
.preview-header.hide-icon::before { display: none; }
.preview-header h1 { font-size: 3em; font-weight: 900; margin-bottom: 15px; text-shadow: 3px 3px 6px rgba(0,0,0,0.4); letter-spacing: 2px; }
.preview-header p { font-size: 1.3em; opacity: 0.95; font-weight: 500; }

.preview-header .logo, .preview-footer .logo { position: absolute; max-height: 100px; z-index: 10; filter: drop-shadow(0 2px 4px rgba(0,0,0,0.3)); }
.preview-header .logo.top-left { left: 25px; top: 25px; }
.preview-header .logo.top-right { right: 25px; top: 25px; }
.preview-header .logo.top-center { left: 50%; top: 25px; transform: translateX(-50%); }
.preview-footer .logo.bottom-left { left: 25px; bottom: 25px; }
.preview-footer .logo.bottom-right { right: 25px; bottom: 25px; }
.preview-footer .logo.bottom-center { left: 50%; bottom: 25px; transform: translateX(-50%); }

.preview-timeline { display: flex; align-items: center; justify-content: center; margin: 30px 0; flex-wrap: wrap; padding: 0 20px; }
.preview-timeline-item { color: white; padding: 12px 25px; border-radius: 30px; margin: 8px; font-weight: 700; font-size: 1em; box-shadow: 0 4px 15px rgba(0,0,0,0.2); }

.preview-steps-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(380px, 1fr)); gap: 35px; padding: 50px; }
.preview-step { background: white; border-radius: 20px; padding: 35px; box-shadow: 0 15px 35px rgba(0,0,0,0.1); border-left: 8px solid; position: relative; }
.preview-step-number { position: absolute; top: -20px; left: 25px; color: white; width: 50px; height: 50px; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-weight: 900; font-size: 1.4em; box-shadow: 0 6px 20px rgba(0,0,0,0.3); }
.preview-step-icon { font-size: 4em; margin-bottom: 25px; display: block; }
.preview-step h3 { font-size: 1.6em; margin-bottom: 20px; color: #2c3e50; font-weight: 800; line-height: 1.2; }

.preview-action-item { padding: 18px; border-radius: 12px; margin: 15px 0; border-left: 5px solid; font-weight: 600; box-shadow: 0 3px 10px rgba(0,0,0,0.1); }

.preview-emergency-contact { color: white; padding: 30px; border-radius: 15px; margin: 30px 50px; text-align: center; font-size: 1.2em; font-weight: 700; box-shadow: 0 10px 30px rgba(0,0,0,0.2); }
.preview-footer { color: white; text-align: center; padding: 40px 30px; position: relative; line-height: 1.6; }
`
