package cmok

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path"
)

type Handler struct {
	storage Storage
}

var _ http.Handler = new(Handler)

func NewHandler(storage Storage) *Handler {
	return &Handler{storage}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.handleGet(w, r)
	case "PUT":
		h.handlePut(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleGet(w http.ResponseWriter, r *http.Request) {
	entries, err := h.storage.List(r.URL.Path)
	if err == nil {
		listTemplate.Execute(w, entries)
		return
	}

	f, err := h.storage.Get(r.URL.Path)
	if err != nil {
		fmt.Fprintf(w, "cannot find %q", r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer f.Close()

	w.Header().Set("Content-Type", "application/force-download")
	w.Header().Set("Content-Transfer-Encoding", "binary")
	_, name := path.Split(r.URL.Path)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", name))

	io.Copy(w, f)
}

func (h *Handler) handlePut(w http.ResponseWriter, r *http.Request) {
	err := h.storage.Put(r.URL.Path, r.Body)
	if err != nil {
		fmt.Fprint(w, err)
		w.WriteHeader(http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusOK)
}

var listTemplate = template.Must(template.New("list").Parse(listTemplateText))

const listTemplateText = `
<html>
	<body>
		<ul>
		{{ range . }}
			<li>
				<a href="{{ .Path }}">{{ .Name }}</a>
			</li>
		{{ else }}
			<li>Nothing here!</li>
		{{ end }}
		</ul>
	</body>
</html>`