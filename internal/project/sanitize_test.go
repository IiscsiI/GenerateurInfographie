package project

import (
	"strings"
	"testing"
)

func TestSanitizeHTML(t *testing.T) {
	cases := []struct {
		name, in, want string
	}{
		{"balise autorisée conservée", "<strong>NE PAS ÉTEINDRE</strong>", "<strong>NE PAS ÉTEINDRE</strong>"},
		{"script supprimé", `avant<script>alert(1)</script>après`, "avantalert(1)après"},
		{"gestionnaire d'événement retiré", `<span onclick="evil()">x</span>`, "<span>x</span>"},
		{"style whitelisté conservé", `<span style="color:red">x</span>`, `<span style="color:red">x</span>`},
		{"style dangereux filtré", `<span style="background:url(javascript:1)">x</span>`, "<span>x</span>"},
		{"chevrons orphelins échappés", "a < b > c", "a &lt; b &gt; c"},
		{"iframe supprimée", `<iframe src="https://evil"></iframe>ok`, "ok"},
		{"input checkbox conservé", `<input type="checkbox" checked="checked">`, `<input type="checkbox" checked="checked" />`},
		{"input non-checkbox rejeté", `<input type="text" value="x">`, ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := SanitizeHTML(c.in)
			if got != c.want {
				t.Errorf("SanitizeHTML(%q)\n  obtenu : %q\n  attendu: %q", c.in, got, c.want)
			}
		})
	}
}

func TestSanitizeHTML_NeverLeavesScriptClosing(t *testing.T) {
	// Garantie pour l'injection dans les blocs <script> des exports :
	// aucune séquence </script ne doit survivre à la sanitisation.
	inputs := []string{
		`</script><script>alert(1)</script>`,
		`x</ScRiPt>y`,
		`<strong></script></strong>`,
	}
	for _, in := range inputs {
		got := SanitizeHTML(in)
		if strings.Contains(strings.ToLower(got), "</script") {
			t.Errorf("séquence </script survivante dans %q -> %q", in, got)
		}
	}
}

func TestSanitizeProject_CoversAllUserFields(t *testing.T) {
	p := DefaultProject()
	payload := `<script>alert(1)</script>reste`
	p.Content.Title = payload
	p.Content.Subtitle = payload
	p.Content.EmergencyMessage = payload
	p.Content.FooterMessage = payload
	p.Elements.Steps = []Step{{
		ID: 1, Title: "t", Actions: []Action{{ID: 1, Type: "info", Text: payload}},
	}}

	SanitizeProject(&p)

	for name, v := range map[string]string{
		"title":     p.Content.Title,
		"subtitle":  p.Content.Subtitle,
		"emergency": p.Content.EmergencyMessage,
		"footer":    p.Content.FooterMessage,
		"action":    p.Elements.Steps[0].Actions[0].Text,
	} {
		if strings.Contains(v, "<script") {
			t.Errorf("champ %s non sanitisé : %q", name, v)
		}
	}
}

func TestSanitizeHTML_Idempotent(t *testing.T) {
	// La sanitisation doit être stable : sanitize(sanitize(x)) == sanitize(x).
	// Sans cela, chaque cycle sauvegarde/import corrompt le texte.
	inputs := []string{
		"l'étape suivante",
		`dire "non"`,
		"Tom & Jerry",
		"déjà échappé : &#39; et &amp;",
		"<strong>gras</strong> et l'apostrophe",
		"a < b > c",
	}
	for _, in := range inputs {
		once := SanitizeHTML(in)
		twice := SanitizeHTML(once)
		if once != twice {
			t.Errorf("sanitisation non idempotente pour %q :\n  1x: %q\n  2x: %q", in, once, twice)
		}
	}
}

func TestSanitizeHTML_PreservesApostrophes(t *testing.T) {
	got := SanitizeHTML("l'étape d'urgence")
	if got != "l'étape d'urgence" {
		t.Errorf("apostrophes altérées : %q", got)
	}
}
