package service

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	cfg "github.com/MelnikovNA/noolingoswagger/internal/domain"
	"github.com/sirupsen/logrus"
)

type swagger struct {
	log    *logrus.Logger
	url    string
	uiPath string
}

func Swagger(config *cfg.Config, logger *logrus.Logger) error {

	addr := strings.Join([]string{
		config.Listen.Host,
		config.Listen.Ports.HTTP,
	}, ":")

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	if config.App.FilesPath == "" {
		config.App.FilesPath = filepath.Join(cwd, "ui")
	}

	s := &swagger{
		log:    logger,
		url:    config.App.URL,
		uiPath: config.App.FilesPath,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.router())

	return http.ListenAndServe(addr, mux)
}

func (s *swagger) router() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		url := req.URL.Path
		s.log.Debugf(`request path %s`, url)

		if url == "/swagger.json" {
			s.swaggerFile(w)
			return
		}

		fp := path.Clean(filepath.Join(s.uiPath, url))
		s.log.Debugf(`file path %s`, fp)
		http.ServeFile(w, req, fp)
	}
}

func (s *swagger) swaggerFile(w http.ResponseWriter) {
	swaggerPath := filepath.Join(s.uiPath, "swagger.json")
	file, err := os.Open(swaggerPath)

	if err != nil {
		s.log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	bytes, err := io.ReadAll(file)

	if err != nil {
		s.log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	swagger := make(map[string]json.RawMessage)
	err = json.Unmarshal(bytes, &swagger)

	if err != nil {
		s.log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	host := s.url
	hostJson, err := json.Marshal(host)

	if err != nil {
		s.log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	swagger["host"] = hostJson

	swaggerJson, err := json.Marshal(swagger)

	if err != nil {
		s.log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.Header().Add("Cache-Control", "no-cache")
	_, err = w.Write(swaggerJson)

	if err != nil {
		s.log.Error(err)
		return
	}
}
