package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"infographic-generator/internal/project"
)

func (s *Server) jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (s *Server) jsonError(w http.ResponseWriter, status int, msg string) {
	s.jsonResponse(w, status, map[string]string{"error": msg})
}

// --- Projects CRUD ---

func (s *Server) handleListProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := s.store.List()
	if err != nil {
		s.jsonError(w, 500, "Erreur de lecture des projets")
		return
	}
	s.jsonResponse(w, 200, projects)
}

func (s *Server) handleGetProject(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	p, err := s.store.Load(id)
	if err != nil {
		s.jsonError(w, 404, err.Error())
		return
	}
	s.jsonResponse(w, 200, p)
}

func (s *Server) handleCreateProject(w http.ResponseWriter, r *http.Request) {
	var p project.Project
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		s.jsonError(w, 400, "JSON invalide")
		return
	}

	// Sanitize user-supplied HTML before persisting
	project.SanitizeProject(&p)

	id, err := s.store.Save("", &p)
	if err != nil {
		s.jsonError(w, 500, "Erreur de sauvegarde")
		return
	}

	s.jsonResponse(w, 201, map[string]string{"id": id})
}

func (s *Server) handleUpdateProject(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var p project.Project
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		s.jsonError(w, 400, "JSON invalide")
		return
	}

	// Sanitize user-supplied HTML before persisting
	project.SanitizeProject(&p)

	savedID, err := s.store.Save(id, &p)
	if err != nil {
		s.jsonError(w, 500, "Erreur de sauvegarde")
		return
	}

	s.jsonResponse(w, 200, map[string]string{"id": savedID})
}

func (s *Server) handleDeleteProject(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.store.Delete(id); err != nil {
		s.jsonError(w, 404, err.Error())
		return
	}
	s.jsonResponse(w, 200, map[string]string{"status": "deleted"})
}

func (s *Server) handleDuplicateProject(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	p, err := s.store.Load(id)
	if err != nil {
		s.jsonError(w, 404, err.Error())
		return
	}

	p.Metadata.Name = p.Metadata.Name + " (copie)"
	newID, err := s.store.Save("", p)
	if err != nil {
		s.jsonError(w, 500, "Erreur de duplication")
		return
	}

	s.jsonResponse(w, 201, map[string]string{"id": newID})
}

// --- Assets ---

// allowedAssetExts : seules les images matricielles sûres sont acceptées.
// SVG volontairement exclu : c'est un format scriptable — un SVG malveillant
// resservi depuis l'origine de l'éditeur serait un XSS stocké contournant
// toute la sanitisation HTML.
var allowedAssetExts = map[string]bool{
	".png": true, ".jpg": true, ".jpeg": true, ".gif": true, ".webp": true,
}

func (s *Server) handleUploadAsset(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	// Limite dure sur le corps de la requête (l'ancien ParseMultipartForm
	// ne plafonnait que la partie gardée en mémoire, pas la taille reçue).
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20) // 10 Mo
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		s.jsonError(w, 413, "Fichier trop volumineux (10 Mo max)")
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		s.jsonError(w, 400, "Fichier manquant")
		return
	}
	defer file.Close()

	// 1. Contrôle de l'extension (liste blanche)
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedAssetExts[ext] {
		s.jsonError(w, 400, "Format non autorise (png, jpg, jpeg, gif, webp)")
		return
	}

	// 2. Contrôle du contenu réel (magic bytes) : l'extension ne suffit pas,
	// un fichier HTML renommé en .png doit être rejeté.
	head := make([]byte, 512)
	n, _ := io.ReadFull(file, head)
	head = head[:n]
	contentType := http.DetectContentType(head)
	if !strings.HasPrefix(contentType, "image/") {
		s.jsonError(w, 400, "Le contenu du fichier n'est pas une image valide")
		return
	}

	// Réassembler le flux complet (en-tête déjà lu + reste)
	reader := io.MultiReader(bytes.NewReader(head), file)

	url, err := s.store.SaveAsset(id, header.Filename, reader)
	if err != nil {
		s.jsonError(w, 500, "Erreur upload")
		return
	}

	s.jsonResponse(w, 201, map[string]string{"url": url})
}

func (s *Server) handleGetAsset(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	filename := r.PathValue("filename")

	// Même liste blanche qu'à l'upload : un fichier déposé sur disque par un
	// autre canal ne doit pas devenir servable pour autant.
	if !allowedAssetExts[strings.ToLower(filepath.Ext(filename))] {
		http.Error(w, "Format non autorise", http.StatusForbidden)
		return
	}

	path := s.store.GetAssetPath(id, filename)
	http.ServeFile(w, r, path)
}

// --- Import ---

func (s *Server) handleImport(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	file, header, err := r.FormFile("file")
	if err != nil {
		s.jsonError(w, 400, "Fichier manquant")
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		s.jsonError(w, 400, "Erreur de lecture")
		return
	}

	var p *project.Project

	if strings.HasSuffix(header.Filename, ".json") {
		p, err = project.ImportFromJSON(data)
	} else if strings.HasSuffix(header.Filename, ".html") {
		p, err = project.ImportFromHTML(string(data))
	} else {
		s.jsonError(w, 400, "Format non supporte (JSON ou HTML attendu)")
		return
	}

	if err != nil {
		s.jsonError(w, 400, fmt.Sprintf("Erreur d'import: %v", err))
		return
	}

	// Sanitize imported content before persisting
	project.SanitizeProject(p)

	id, err := s.store.Save("", p)
	if err != nil {
		s.jsonError(w, 500, "Erreur de sauvegarde apres import")
		return
	}

	s.jsonResponse(w, 201, map[string]interface{}{
		"id":      id,
		"project": p,
	})
}

// --- Browser Status ---

func (s *Server) handleBrowserStatus(w http.ResponseWriter, r *http.Request) {
	available := s.browser != nil && s.browser.Path != ""
	resp := map[string]interface{}{
		"available": available,
	}
	if available {
		resp["name"] = s.browser.Name
		resp["path"] = s.browser.Path
	}
	s.jsonResponse(w, 200, resp)
}

// --- Shutdown ---

func (s *Server) handleShutdown(w http.ResponseWriter, r *http.Request) {
	s.jsonResponse(w, 200, map[string]string{"status": "shutting_down"})

	// Trigger shutdown after response is sent
	go func() {
		time.Sleep(500 * time.Millisecond)
		if s.OnShutdown != nil {
			s.OnShutdown()
		}
	}()
}
