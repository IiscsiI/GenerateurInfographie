package project

import (
	"html"
	"regexp"
	"strings"
)

// SanitizeProject normalizes all user-supplied HTML fields in a project.
// Applied on save (defense: clean storage) and on render (defense in depth).
func SanitizeProject(p *Project) {
	p.Content.Title = SanitizeHTML(p.Content.Title)
	p.Content.Subtitle = SanitizeHTML(p.Content.Subtitle)
	p.Content.EmergencyMessage = SanitizeHTML(p.Content.EmergencyMessage)
	p.Content.FooterMessage = SanitizeHTML(p.Content.FooterMessage)
	for i := range p.Elements.Steps {
		for j := range p.Elements.Steps[i].Actions {
			p.Elements.Steps[i].Actions[j].Text = SanitizeHTML(p.Elements.Steps[i].Actions[j].Text)
		}
	}
	for i := range p.Elements.Logos {
		sanitizeLogo(&p.Elements.Logos[i])
	}
}

// sanitizeLogo borne les valeurs numériques d'un logo.
// Les coordonnées libres sont des pourcentages du conteneur : [0;100].
func sanitizeLogo(l *Logo) {
	if l.Size < 20 {
		l.Size = 20
	}
	if l.Size > 600 {
		l.Size = 600
	}
	l.X = clampPct(l.X)
	l.Y = clampPct(l.Y)
}

func clampPct(v *float64) *float64 {
	if v == nil {
		return nil
	}
	c := *v
	if c < 0 {
		c = 0
	}
	if c > 100 {
		c = 100
	}
	return &c
}

// SanitizeHTML filters user-supplied HTML with a strict whitelist.
// Uses regex scanning (not full DOM parsing) to avoid external dependencies.
// Mirrors the JS sanitizer in the editor.
//
// Strategy:
//  1. Tokenize the input into text and tag segments.
//  2. For each tag, drop it entirely if not whitelisted.
//     If whitelisted, rewrite the tag with only whitelisted attributes.
//  3. Any text between tags is HTML-escaped if it contains raw < or >.
func SanitizeHTML(input string) string {
	if input == "" {
		return ""
	}

	var out strings.Builder
	pos := 0
	for {
		m := tagScanRE.FindStringIndex(input[pos:])
		if m == nil {
			// Remaining text
			out.WriteString(escapeText(input[pos:]))
			break
		}
		// Emit preceding text
		out.WriteString(escapeText(input[pos : pos+m[0]]))

		tagStr := input[pos+m[0] : pos+m[1]]
		clean := sanitizeTag(tagStr)
		out.WriteString(clean)
		pos += m[1]
	}
	return out.String()
}

// Matches <tag ...> or </tag> (non-greedy, no nested <>)
var tagScanRE = regexp.MustCompile(`<\/?[a-zA-Z][^<>]*>`)

// Matches tag name and determines if it's a closing tag
var tagNameRE = regexp.MustCompile(`<(\/?)([a-zA-Z][a-zA-Z0-9]*)`)

// Matches attribute: name="value" or name='value'
var attrRE = regexp.MustCompile(`([a-zA-Z_\-:]+)\s*=\s*(?:"([^"]*)"|'([^']*)')`)

var allowedTagsGo = map[string]bool{
	"strong": true, "b": true, "em": true, "i": true, "u": true,
	"s": true, "strike": true, "br": true, "p": true, "div": true, "span": true,
	"label": true, "input": true,
}

var allowedAttrsGo = map[string]map[string]bool{
	"p":     {"style": true},
	"div":   {"style": true},
	"span":  {"style": true},
	"input": {"type": true, "checked": true, "disabled": true},
}

var allowedStylesGo = map[string]bool{
	"color": true, "text-align": true, "font-weight": true,
	"font-style": true, "text-decoration": true,
}

var styleBlockREGo = regexp.MustCompile(`(?i)expression\s*\(|url\s*\(|@import|javascript:`)
var styleSafeREGo = regexp.MustCompile(`^[a-zA-Z0-9#.,()\s%\-_']+$`)

// sanitizeTag takes a raw tag string like `<span style="color:red">` or `</strong>`
// and returns either the sanitized tag or an empty string if disallowed.
func sanitizeTag(raw string) string {
	nm := tagNameRE.FindStringSubmatch(raw)
	if nm == nil {
		return ""
	}
	closing := nm[1] == "/"
	tag := strings.ToLower(nm[2])
	if !allowedTagsGo[tag] {
		return ""
	}
	if closing {
		return "</" + tag + ">"
	}

	// Special rule: <input> must have type="checkbox", nothing else.
	if tag == "input" {
		hasCheckbox := false
		for _, m := range attrRE.FindAllStringSubmatch(raw, -1) {
			if strings.ToLower(m[1]) == "type" {
				val := m[2]
				if val == "" {
					val = m[3]
				}
				if strings.ToLower(val) == "checkbox" {
					hasCheckbox = true
					break
				}
			}
		}
		if !hasCheckbox {
			return ""
		}
	}

	// Rebuild opening tag with only whitelisted attributes
	var b strings.Builder
	b.WriteString("<")
	b.WriteString(tag)

	allowed := allowedAttrsGo[tag]
	if allowed != nil {
		for _, m := range attrRE.FindAllStringSubmatch(raw, -1) {
			name := strings.ToLower(m[1])
			if !allowed[name] {
				continue
			}
			val := m[2]
			if val == "" {
				val = m[3]
			}
			if name == "style" {
				val = sanitizeStyle(val)
				if val == "" {
					continue
				}
			}
			// Force input type=checkbox
			if tag == "input" && name == "type" {
				val = "checkbox"
			}
			// Escape the value
			b.WriteString(` `)
			b.WriteString(name)
			b.WriteString(`="`)
			b.WriteString(html.EscapeString(val))
			b.WriteString(`"`)
		}
	}

	// Self-closing tags
	if tag == "br" || tag == "input" {
		b.WriteString(" />")
	} else {
		b.WriteString(">")
	}
	return b.String()
}

func sanitizeStyle(s string) string {
	if s == "" {
		return ""
	}
	var out []string
	for _, decl := range strings.Split(s, ";") {
		idx := strings.Index(decl, ":")
		if idx < 0 {
			continue
		}
		prop := strings.ToLower(strings.TrimSpace(decl[:idx]))
		val := strings.TrimSpace(decl[idx+1:])
		if !allowedStylesGo[prop] {
			continue
		}
		if styleBlockREGo.MatchString(val) {
			continue
		}
		if !styleSafeREGo.MatchString(val) {
			continue
		}
		out = append(out, prop+":"+val)
	}
	return strings.Join(out, ";")
}

// textEscaper neutralise uniquement les chevrons orphelins dans le texte.
// SURTOUT PAS html.EscapeString ici : celui-ci échappe aussi & ' " — la
// sanitisation ne serait plus idempotente et chaque sauvegarde ou import
// ajouterait un niveau d'échappement (l'étape -> l&#39;étape -> l&amp;#39;étape...).
var textEscaper = strings.NewReplacer("<", "&lt;", ">", "&gt;")

func escapeText(s string) string {
	return textEscaper.Replace(s)
}
