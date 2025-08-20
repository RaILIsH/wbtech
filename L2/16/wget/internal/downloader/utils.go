package downloader

import (
	"net/url"
	"path"
	"strings"
)

func normalizeURL(baseURL, relativeURL string) (string, error) {
	if strings.HasPrefix(relativeURL, "http://") || strings.HasPrefix(relativeURL, "https://") {
		return relativeURL, nil
	}

	base, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	// Если ссылка абсолютная от корня
	if strings.HasPrefix(relativeURL, "/") {
		return base.Scheme + "://" + base.Host + relativeURL, nil
	}

	// Если ссылка относительная
	if !strings.Contains(relativeURL, "://") {
		// Получаем директорию из baseURL
		basePath := path.Dir(base.Path)
		if basePath == "/" {
			basePath = ""
		}
		return base.Scheme + "://" + base.Host + basePath + "/" + relativeURL, nil
	}

	return relativeURL, nil
}

func getFileExtension(contentType string) string {
	switch {
	case strings.Contains(contentType, "text/html"):
		return ".html"
	case strings.Contains(contentType, "text/css"):
		return ".css"
	case strings.Contains(contentType, "application/javascript"):
		return ".js"
	case strings.Contains(contentType, "image/jpeg"):
		return ".jpg"
	case strings.Contains(contentType, "image/png"):
		return ".png"
	case strings.Contains(contentType, "image/gif"):
		return ".gif"
	case strings.Contains(contentType, "image/svg+xml"):
		return ".svg"
	default:
		return ".bin"
	}
}
