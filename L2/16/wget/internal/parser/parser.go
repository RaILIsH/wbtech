package parser

import (
	"bytes"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) ExtractLinks(content []byte, parentURL, baseURL string) ([]string, error) {
	var links []string

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(content))
	if err != nil {
		return nil, err
	}

	// Извлекаем все ссылки
	doc.Find("a[href], link[href], script[src], img[src], iframe[src]").Each(func(i int, s *goquery.Selection) {
		var link string
		if href, exists := s.Attr("href"); exists {
			link = href
		} else if src, exists := s.Attr("src"); exists {
			link = src
		}

		if link != "" {
			normalized, err := p.normalizeLink(link, parentURL, baseURL)
			if err == nil && normalized != "" {
				links = append(links, normalized)
			}
		}
	})

	return p.removeDuplicates(links), nil
}

func (p *Parser) normalizeLink(link, parentURL, baseURL string) (string, error) {
	// Игнорируем якорные ссылки и javascript
	if strings.HasPrefix(link, "#") || strings.HasPrefix(link, "javascript:") {
		return "", nil
	}

	// Если ссылка уже абсолютная
	if strings.HasPrefix(link, "http://") || strings.HasPrefix(link, "https://") {
		return link, nil
	}

	base, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	parent, err := url.Parse(parentURL)
	if err != nil {
		return "", err
	}

	// Обрабатываем относительные ссылки
	if strings.HasPrefix(link, "/") {
		// Абсолютная от корня
		return base.Scheme + "://" + base.Host + link, nil
	} else if strings.HasPrefix(link, "./") {
		// Относительная от текущей директории
		parentDir := parent.Path[:strings.LastIndex(parent.Path, "/")+1]
		return base.Scheme + "://" + base.Host + parentDir + link[2:], nil
	} else if !strings.Contains(link, "://") {
		// Относительный путь
		parentDir := parent.Path
		if !strings.HasSuffix(parentDir, "/") {
			parentDir = parentDir[:strings.LastIndex(parentDir, "/")+1]
		}
		return base.Scheme + "://" + base.Host + parentDir + link, nil
	}

	return link, nil
}

func (p *Parser) removeDuplicates(links []string) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, link := range links {
		if !seen[link] {
			seen[link] = true
			result = append(result, link)
		}
	}

	return result
}
