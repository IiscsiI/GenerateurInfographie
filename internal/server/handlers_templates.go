package server

import (
	"encoding/json"
	"net/http"

	"infographic-generator/internal/project"
)

// ==========================================================================
// API DES MODÈLES PERSONNALISÉS
// ==========================================================================

// handleListTemplates : GET /api/templates
func (s *Server) handleListTemplates(w http.ResponseWriter, r *http.Request) {
	list, err := s.store.ListTemplates()
	if err != nil {
		s.jsonError(w, 500, "liste des modeles: "+err.Error())
		return
	}
	s.jsonResponse(w, 200, list)
}

// handleSaveTemplate : POST /api/templates {name, project}
func (s *Server) handleSaveTemplate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name    string          `json:"name"`
		Project json.RawMessage `json:"project"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.jsonError(w, 400, "JSON invalide: "+err.Error())
		return
	}
	p, err := project.ImportFromJSON(req.Project)
	if err != nil {
		s.jsonError(w, 400, "projet invalide: "+err.Error())
		return
	}
	id, err := s.store.SaveTemplate(req.Name, p)
	if err != nil {
		s.jsonError(w, 400, "sauvegarde du modele: "+err.Error())
		return
	}
	s.jsonResponse(w, 200, map[string]string{"id": id})
}

// handleGetTemplate : GET /api/templates/{id}
func (s *Server) handleGetTemplate(w http.ResponseWriter, r *http.Request) {
	p, err := s.store.LoadTemplate(r.PathValue("id"))
	if err != nil {
		s.jsonError(w, 404, err.Error())
		return
	}
	s.jsonResponse(w, 200, p)
}

// handleDeleteTemplate : DELETE /api/templates/{id}
func (s *Server) handleDeleteTemplate(w http.ResponseWriter, r *http.Request) {
	if err := s.store.DeleteTemplate(r.PathValue("id")); err != nil {
		s.jsonError(w, 404, err.Error())
		return
	}
	s.jsonResponse(w, 200, map[string]bool{"ok": true})
}
