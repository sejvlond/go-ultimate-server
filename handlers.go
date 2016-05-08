package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

const TEMPLATE_DIR = "./templates/"
const LAYOUT_FILE = TEMPLATE_DIR + "@layout"

var templates map[string]*template.Template

func init() {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}
	files, err := filepath.Glob(TEMPLATE_DIR + "*.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		name := filepath.Base(f)[:len(filepath.Base(f))-len(filepath.Ext(f))]
		templates[name] = template.Must(template.ParseFiles(LAYOUT_FILE, f))
	}
}

// ----------------------------------------------------------------------------

type HomepageHandler struct {
	lgr LOGGER
}

type HomepageData struct {
	UltimateAnswer int
}

func NewHomepageHandler(lgr LOGGER) (*HomepageHandler, error) {
	h := &HomepageHandler{
		lgr: lgr,
	}
	return h, nil
}

func (h *HomepageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.lgr.Warnf("Unallowed method %q", r.Method)
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}
	h.lgr.Infof("--> %v %v %v", r.Method, r.URL, r.RemoteAddr)
	if r.URL.Path != "/" {
		notFound(h.lgr, w, r)
		return
	}

	templates["Homepage"].Execute(w, HomepageData{
		UltimateAnswer: 42,
	})
	h.lgr.Infof("<-- 200 OK")
	return
}

// ----------------------------------------------------------------------------

type NotFoundData struct {
	Code    int
	Message string
	URI     string
}

func notFound(lgr LOGGER, w http.ResponseWriter, r *http.Request) {
	lgr.Warnf("<-- 404 Not Found")
	templates["NotFound"].Execute(w, NotFoundData{
		Code:    404,
		Message: "NotFound",
		URI:     r.URL.Path,
	})
	return
}
