package project

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// Storage handles filesystem-based project persistence.
type Storage struct {
	BaseDir    string
	ProjectDir string
	ExportDir  string
}

// ProjectInfo is a summary for listing projects.
type ProjectInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Modified string `json:"modified"`
}

// NewStorage initializes the data directory structure.
func NewStorage(baseDir string) (*Storage, error) {
	s := &Storage{
		BaseDir:    baseDir,
		ProjectDir: filepath.Join(baseDir, "projects"),
		ExportDir:  filepath.Join(baseDir, "exports"),
	}

	for _, dir := range []string{s.ProjectDir, s.ExportDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("creation repertoire %s: %w", dir, err)
		}
	}

	return s, nil
}

// List returns all project summaries, sorted by modification date (newest first).
func (s *Storage) List() ([]ProjectInfo, error) {
	entries, err := os.ReadDir(s.ProjectDir)
	if err != nil {
		return nil, err
	}

	projects := []ProjectInfo{}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		jsonPath := filepath.Join(s.ProjectDir, entry.Name(), "project.json")
		info, err := os.Stat(jsonPath)
		if err != nil {
			continue
		}

		// Read project name from JSON
		name := entry.Name()
		data, err := os.ReadFile(jsonPath)
		if err == nil {
			var p Project
			if json.Unmarshal(data, &p) == nil && p.Metadata.Name != "" {
				name = p.Metadata.Name
			}
		}

		projects = append(projects, ProjectInfo{
			ID:       entry.Name(),
			Name:     name,
			Modified: info.ModTime().Format(time.RFC3339),
		})
	}

	sort.Slice(projects, func(i, j int) bool {
		return projects[i].Modified > projects[j].Modified
	})

	return projects, nil
}

// Load reads a project by ID.
func (s *Storage) Load(id string) (*Project, error) {
	id = sanitizeID(id)
	jsonPath := filepath.Join(s.ProjectDir, id, "project.json")

	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("projet introuvable: %s", id)
	}

	var p Project
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("JSON invalide pour %s: %w", id, err)
	}

	return &p, nil
}

// Save writes a project. Returns the project ID.
func (s *Storage) Save(id string, p *Project) (string, error) {
	if id == "" {
		id = slugify(p.Metadata.Name)
		if id == "" {
			id = fmt.Sprintf("projet-%d", time.Now().Unix())
		}
	}
	id = sanitizeID(id)

	projDir := filepath.Join(s.ProjectDir, id)
	assetsDir := filepath.Join(projDir, "assets")

	if err := os.MkdirAll(assetsDir, 0755); err != nil {
		return "", err
	}

	p.Version = "4.0"
	p.Generated = time.Now().Format(time.RFC3339)
	p.Generator = "Generateur Infographie v2.1 (Go)"

	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return "", err
	}

	jsonPath := filepath.Join(projDir, "project.json")
	if err := os.WriteFile(jsonPath, data, 0644); err != nil {
		return "", err
	}

	return id, nil
}

// Delete removes a project directory.
func (s *Storage) Delete(id string) error {
	id = sanitizeID(id)
	projDir := filepath.Join(s.ProjectDir, id)

	if _, err := os.Stat(projDir); os.IsNotExist(err) {
		return fmt.Errorf("projet introuvable: %s", id)
	}

	return os.RemoveAll(projDir)
}

// SaveAsset stores an uploaded file in the project's assets directory.
func (s *Storage) SaveAsset(projectID, filename string, r io.Reader) (string, error) {
	projectID = sanitizeID(projectID)
	filename = sanitizeFilename(filename)

	assetsDir := filepath.Join(s.ProjectDir, projectID, "assets")
	if err := os.MkdirAll(assetsDir, 0755); err != nil {
		return "", err
	}

	filePath := filepath.Join(assetsDir, filename)
	f, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := io.Copy(f, r); err != nil {
		return "", err
	}

	// Return the URL path for the frontend
	return fmt.Sprintf("/api/projects/%s/assets/%s", projectID, filename), nil
}

// GetAssetPath returns the filesystem path for a project asset.
func (s *Storage) GetAssetPath(projectID, filename string) string {
	return filepath.Join(s.ProjectDir, sanitizeID(projectID), "assets", sanitizeFilename(filename))
}

// SaveExport stores an export file and returns the filesystem path.
func (s *Storage) SaveExport(filename string, data []byte) (string, error) {
	filename = sanitizeFilename(filename)
	filePath := filepath.Join(s.ExportDir, filename)

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return "", err
	}

	return filePath, nil
}

var slugRe = regexp.MustCompile(`[^a-zA-Z0-9]+`)

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = slugRe.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if len(s) > 60 {
		s = s[:60]
	}
	return s
}

func sanitizeID(id string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9_\-]`)
	return re.ReplaceAllString(id, "")
}

func sanitizeFilename(name string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9._\-]`)
	return re.ReplaceAllString(name, "_")
}
