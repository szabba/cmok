package cmok

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"path"

	"github.com/szabba/cmok/auth"
)

const (
	_WWWAuthenticateHeader             = "WWW-Authenticate"
	_WWWAuthenticateHeaderWithoutRealm = "Basic"
	_WWWAuthenticateHeaderWithRealm    = "Basic realm=%q"
)

type Handler struct {
	realm        string
	authService  auth.Service
	accessPolicy AccessPolicy
	storage      Storage
}

var _ http.Handler = new(Handler)

func NewHandler(realm string, authService auth.Service, accessPolicy AccessPolicy, storage Storage) *Handler {
	return &Handler{
		realm:        realm,
		authService:  authService,
		accessPolicy: accessPolicy,
		storage:      storage,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("request %s %s", r.Method, r.URL.Path)
	user, ok := h.authService.Authenticate(r)
	if !ok {
		h.denyAccess(w)
		return
	}

	storage := h.accessPolicy.Protect(h.storage, user)

	switch r.Method {
	case http.MethodGet:
		h.handleGet(w, r, storage)
	case http.MethodPut:
		h.handlePut(w, r, storage)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleGet(w http.ResponseWriter, r *http.Request, storage Storage) {
	children, content, err := storage.Get(r.URL.Path)
	if err == ErrAccessDenied {
		h.denyAccess(w)

	} else if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)

	} else if children != nil {
		listTemplate.Execute(w, children)

	} else {
		defer content.Close()

		_, name := path.Split(r.URL.Path)
		h.forceDownload(w, name, content)
	}
}

func (h *Handler) handlePut(w http.ResponseWriter, r *http.Request, storage Storage) {
	err := storage.Set(r.URL.Path, r.Body)
	if err == ErrAccessDenied {
		h.denyAccess(w)

	} else if err != nil {
		w.WriteHeader(http.StatusConflict)
		fmt.Fprint(w, err)

	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) denyAccess(w http.ResponseWriter) {
	w.Header().Set(_WWWAuthenticateHeader, h.wwwAuthentiacateHeaderValue())
	http.Error(w, "access denied", http.StatusUnauthorized)
}

func (h *Handler) wwwAuthentiacateHeaderValue() string {
	if h.realm == "" {
		return _WWWAuthenticateHeaderWithoutRealm
	}
	return fmt.Sprintf(_WWWAuthenticateHeaderWithRealm, h.realm)
}

func (h *Handler) forceDownload(w http.ResponseWriter, name string, r io.Reader) {
	w.Header().Set("Content-Type", "application/force-download")
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", name))

	io.Copy(w, r)
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
