package registry

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/theydontwantyoutovibecode/dade/internal/config"
)

type Project struct {
	Port     int    `json:"port"`
	Path     string `json:"path"`
	Template string `json:"template"`
	Created  string `json:"created"`
}

type Entry struct {
	Name    string
	Project Project
}

func Load(path string) (map[string]Project, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return map[string]Project{}, nil
		}
		return nil, err
	}
	if len(data) == 0 {
		return map[string]Project{}, nil
	}

	var projects map[string]Project
	if err := json.Unmarshal(data, &projects); err != nil {
		return map[string]Project{}, nil
	}
	if projects == nil {
		projects = map[string]Project{}
	}
	return projects, nil
}

func Save(path string, projects map[string]Project) error {
	if projects == nil {
		projects = map[string]Project{}
	}
	data, err := json.MarshalIndent(projects, "", "  ")
	if err != nil {
		return err
	}
	return writeAtomically(path, data, 0644)
}

func Register(path, name string, port int, projectPath, template string) (Project, error) {
	projects, err := Load(path)
	if err != nil {
		return Project{}, err
	}
	created := time.Now().UTC().Format(time.RFC3339)
	if existing, ok := projects[name]; ok && existing.Created != "" {
		created = existing.Created
	}
	project := Project{Port: port, Path: projectPath, Template: template, Created: created}
	projects[name] = project
	if err := Save(path, projects); err != nil {
		return Project{}, err
	}
	return project, nil
}

func Unregister(path, name string) (bool, error) {
	projects, err := Load(path)
	if err != nil {
		return false, err
	}
	if _, ok := projects[name]; !ok {
		return false, nil
	}
	delete(projects, name)
	if err := Save(path, projects); err != nil {
		return false, err
	}
	return true, nil
}

func Get(path, name string) (Project, bool, error) {
	projects, err := Load(path)
	if err != nil {
		return Project{}, false, err
	}
	project, ok := projects[name]
	return project, ok, nil
}

func GetByPath(path string, projectPath string) (Entry, bool, error) {
	projects, err := Load(path)
	if err != nil {
		return Entry{}, false, err
	}
	for name, project := range projects {
		if project.Path == projectPath {
			return Entry{Name: name, Project: project}, true, nil
		}
	}
	return Entry{}, false, nil
}

func List(path string) ([]Entry, error) {
	projects, err := Load(path)
	if err != nil {
		return nil, err
	}
	entries := make([]Entry, 0, len(projects))
	for name, project := range projects {
		entries = append(entries, Entry{Name: name, Project: project})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name < entries[j].Name
	})
	return entries, nil
}

func GetPort(path string, name string) (int, bool, error) {
	project, ok, err := Get(path, name)
	if err != nil || !ok {
		return 0, ok, err
	}
	return project.Port, true, nil
}

func GetPath(path string, name string) (string, bool, error) {
	project, ok, err := Get(path, name)
	if err != nil || !ok {
		return "", ok, err
	}
	return project.Path, true, nil
}

func GetTemplate(path string, name string) (string, bool, error) {
	project, ok, err := Get(path, name)
	if err != nil || !ok {
		return "", ok, err
	}
	return project.Template, true, nil
}

func Exists(path string, name string) (bool, error) {
	_, ok, err := Get(path, name)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func ListNames(path string) ([]string, error) {
	entries, err := List(path)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		names = append(names, entry.Name)
	}
	return names, nil
}

func NextPort(path string) (int, error) {
	projects, err := Load(path)
	if err != nil {
		return 0, err
	}
	maxPort := config.BasePort - 1
	for _, project := range projects {
		if project.Port > maxPort {
			maxPort = project.Port
		}
	}
	if maxPort < config.BasePort {
		return config.BasePort, nil
	}
	return maxPort + 1, nil
}

func IsPortAvailable(port int, projects map[string]Project) bool {
	for _, project := range projects {
		if project.Port == port {
			return false
		}
	}
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return false
	}
	_ = listener.Close()
	return true
}

// UpdatePort updates the port for an existing project in the registry.
func UpdatePort(path, name string, port int) (Project, error) {
	projects, err := Load(path)
	if err != nil {
		return Project{}, err
	}
	project, ok := projects[name]
	if !ok {
		return Project{}, fmt.Errorf("project '%s' not found", name)
	}
	project.Port = port
	projects[name] = project
	if err := Save(path, projects); err != nil {
		return Project{}, err
	}
	return project, nil
}

func writeAtomically(path string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	file, err := os.CreateTemp(dir, "projects-*.json")
	if err != nil {
		return err
	}
	defer func() {
		_ = os.Remove(file.Name())
	}()

	if _, err := file.Write(data); err != nil {
		_ = file.Close()
		return err
	}
	if err := file.Chmod(perm); err != nil {
		_ = file.Close()
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	return os.Rename(file.Name(), path)
}
