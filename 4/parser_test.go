package parser

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func TestURLParser(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*90)
	p := NewWebParser()
	parser := New(p, 2*time.Second, 2)
	urls := parser.Parse(ctx, "https://go.dev/")
	defer cancel()
	t.Logf("len(urls): %d", len(urls))
}

// WebParser реализует интерфейс Parser для веб-страниц.
type WebParser struct {
	client *http.Client // HTTP-клиент для запросов
}

// NewWebParser создаёт новый экземпляр WebParser.
func NewWebParser() *WebParser {
	return &WebParser{
		client: &http.Client{}, // Можно настроить Timeout и другие параметры
	}
}

// Parse извлекает все абсолютные ссылки с веб-страницы.
func (wp *WebParser) Parse(ctx context.Context, rawURL string) ([]string, error) {
	// Проверяем, что URL валиден
	baseURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	// Создаём HTTP-запрос с контекстом
	req, err := http.NewRequestWithContext(ctx, "GET", rawURL, nil)
	if err != nil {
		return nil, err
	}

	// Отправляем запрос
	resp, err := wp.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return nil, &url.Error{
			Op:  "GET",
			URL: rawURL,
			Err: err,
		}
	}

	// Парсим HTML с помощью goquery
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	// Извлекаем все ссылки
	var links []string
	doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}

		// Пропускаем пустые ссылки и якоря (#)
		if href == "" || strings.HasPrefix(href, "#") {
			return
		}

		// Преобразуем относительные ссылки в абсолютные
		parsedURL, err := baseURL.Parse(href)
		if err != nil {
			return // Пропускаем некорректные URL
		}

		// Нормализуем URL (убираем фрагменты и т. д.)
		parsedURL.Fragment = ""
		links = append(links, parsedURL.String())
	})

	return links, nil
}
