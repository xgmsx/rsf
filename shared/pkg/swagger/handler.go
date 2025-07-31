package swagger

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

type SwaggerHandler struct {
	prefix     string
	filename   string
	fileServer http.Handler
}

func NewSwaggerHandler(prefix, filename, apiPath string) *SwaggerHandler {
	// Создаем файловый сервер для указанной директории
	fileServer := http.FileServer(http.Dir(apiPath))

	return &SwaggerHandler{
		fileServer: fileServer,
		filename:   filename,
		prefix:     prefix,
	}
}

func (h *SwaggerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	// Убираем префикс из пути
	if strings.HasPrefix(path, h.prefix) {
		path = strings.TrimPrefix(path, h.prefix)
	}

	// Redirect для базового пути
	if path == "" || path == "/" {
		http.Redirect(w, r, h.prefix+"/swagger.html", http.StatusMovedPermanently)
		return
	}

	switch {
	case strings.HasSuffix(path, ".json"):
		h.serveSwaggerJSON(w, r, path)
	case strings.HasSuffix(path, ".html"):
		h.serveSwaggerHTML(w, r)
	default:
		http.NotFound(w, r)
	}
}

// serveSwaggerJSON отдает JSON файлы через стандартный файловый сервер
func (h *SwaggerHandler) serveSwaggerJSON(w http.ResponseWriter, r *http.Request, filePath string) {
	// Создаем новый запрос с обрезанным путем
	newReq := r.Clone(r.Context())
	newReq.URL.Path = filePath

	// Добавляем заголовки для JSON
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=3600")

	// Передаем управление стандартному файловому серверу
	h.fileServer.ServeHTTP(w, newReq)
}

// serveSwaggerHTML генерирует и отдает HTML страницу с Swagger UI
func (h *SwaggerHandler) serveSwaggerHTML(w http.ResponseWriter, r *http.Request) {
	html, err := h.generateSwaggerHTML()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating HTML: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	w.Write([]byte(html))
}

// generateSwaggerHTML генерирует HTML страницу Swagger UI
func (h *SwaggerHandler) generateSwaggerHTML() (string, error) {
	data := struct {
		SwaggerURL string
	}{
		SwaggerURL: h.prefix + "/" + h.filename, // Используем префикс + имя файла
	}

	t, err := template.New("swagger").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
