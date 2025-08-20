package downloader

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"wget/internal/parser"
	"wget/internal/storage"
)

type Downloader struct {
	workers int
	storage *storage.Storage
	parser  *parser.Parser
	visited sync.Map
	queue   chan downloadTask
	wg      sync.WaitGroup
	client  *http.Client
	baseURL string
}

type downloadTask struct {
	url   string
	depth int
}

func NewDownloader(workers int, storage *storage.Storage, parser *parser.Parser) *Downloader {
	return &Downloader{
		workers: workers,
		storage: storage,
		parser:  parser,
		queue:   make(chan downloadTask, 1000),
		client: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        workers,
				MaxIdleConnsPerHost: workers,
				IdleConnTimeout:     30 * time.Second,
			},
		},
	}
}

func (d *Downloader) DownloadRecursive(startURL string, maxDepth int, baseURL string) error {
	d.baseURL = baseURL

	// Запускаем воркеры
	for i := 0; i < d.workers; i++ {
		go d.worker()
	}

	// Добавляем начальную задачу
	d.queue <- downloadTask{url: startURL, depth: maxDepth}
	d.wg.Add(1)

	// Ждем завершения всех задач
	d.wg.Wait()
	close(d.queue)

	return nil
}

func (d *Downloader) worker() {
	for task := range d.queue {
		d.processTask(task)
		d.wg.Done()
	}
}

func (d *Downloader) processTask(task downloadTask) {
	// Проверяем, не посещали ли уже этот URL
	if _, visited := d.visited.LoadOrStore(task.url, true); visited {
		return
	}

	fmt.Printf("Загрузка: %s (глубина: %d)\n", task.url, task.depth)

	// Скачиваем контент
	content, contentType, err := d.downloadContent(task.url)
	if err != nil {
		log.Printf("Ошибка загрузки %s: %v", task.url, err)
		return
	}

	// Сохраняем файл
	_, err = d.storage.Save(task.url, content, contentType)
	if err != nil {
		log.Printf("Ошибка сохранения %s: %v", task.url, err)
		return
	}

	// Если это HTML и есть глубина для рекурсии - парсим ссылки
	if strings.Contains(contentType, "text/html") && task.depth > 0 {
		d.parseAndQueueLinks(content, task.url, task.depth-1)
	}
}

func (d *Downloader) downloadContent(urlStr string) ([]byte, string, error) {
	req, err := http.NewRequestWithContext(context.Background(), "GET", urlStr, nil)
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Wget/1.0)")

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("HTTP статус: %d", resp.StatusCode)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	return content, resp.Header.Get("Content-Type"), nil
}

func (d *Downloader) parseAndQueueLinks(content []byte, parentURL string, depth int) {
	links, err := d.parser.ExtractLinks(content, parentURL, d.baseURL)
	if err != nil {
		log.Printf("Ошибка парсинга ссылок: %v", err)
		return
	}

	for _, link := range links {
		// Проверяем, что ссылка в том же домене
		if d.isSameDomain(link) {
			d.wg.Add(1)
			d.queue <- downloadTask{url: link, depth: depth}
		}
	}
}

func (d *Downloader) isSameDomain(link string) bool {
	u1, err := url.Parse(d.baseURL)
	if err != nil {
		return false
	}

	u2, err := url.Parse(link)
	if err != nil {
		return false
	}

	return u1.Host == u2.Host
}
