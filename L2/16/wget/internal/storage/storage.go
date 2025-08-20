package storage

import (
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type Storage struct {
	baseDir string
}

func NewStorage(baseDir string) *Storage {
	return &Storage{baseDir: baseDir}
}

func (s *Storage) Save(urlStr string, content []byte, contentType string) (string, error) {
	localPath, err := s.getLocalPath(urlStr, contentType)
	if err != nil {
		return "", err
	}

	// Создаем директории
	dir := filepath.Dir(localPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	// Сохраняем файл
	err = os.WriteFile(localPath, content, 0644)
	if err != nil {
		return "", err
	}

	return localPath, nil
}

func (s *Storage) getLocalPath(urlStr string, contentType string) (string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}

	// Создаем путь на основе URL
	path := u.Host + u.Path

	// Если путь заканчивается на / - добавляем index.html
	if strings.HasSuffix(path, "/") {
		path += "index.html"
	} else if filepath.Ext(path) == "" {
		// Если нет расширения - добавляем на основе content-type
		ext := s.getExtension(contentType)
		path += ext
	}

	return filepath.Join(s.baseDir, path), nil
}

func (s *Storage) getExtension(contentType string) string {
	parts := strings.Split(contentType, ";")
	mainType := strings.TrimSpace(parts[0])

	switch mainType {
	case "text/html":
		return ".html"
	case "text/css":
		return ".css"
	case "application/javascript", "text/javascript":
		return ".js"
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/svg+xml":
		return ".svg"
	case "application/json":
		return ".json"
	case "text/plain":
		return ".txt"
	default:
		return ".bin"
	}
}

func (s *Storage) FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func (s *Storage) CreateFile(path string) (io.WriteCloser, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	return os.Create(path)
}
