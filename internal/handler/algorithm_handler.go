package handler

import (
	"net/http"

	"signatureservice/internal/model"
	"signatureservice/pkg/httpx"
)

func (s *Server) registerAlgorithmRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/algorithms", s.createAlgorithm)
	mux.HandleFunc("GET /api/algorithms", s.listAlgorithms)
	mux.HandleFunc("GET /api/algorithms/{id}", s.getAlgorithm)
	mux.HandleFunc("PUT /api/algorithms/{id}", s.updateAlgorithm)
	mux.HandleFunc("DELETE /api/algorithms/{id}", s.deleteAlgorithm)
}

type createAlgorithmRequest struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	KeySize int    `json:"key_size"`
	Enabled bool   `json:"enabled"`
}

func (s *Server) createAlgorithm(w http.ResponseWriter, r *http.Request) {
	var req createAlgorithmRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	a, err := s.svc.CreateAlgorithm(model.Algorithm{
		Name:    req.Name,
		Type:    req.Type,
		KeySize: req.KeySize,
		Enabled: req.Enabled,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, a)
}

func (s *Server) listAlgorithms(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.AlgorithmFilter{
		Name: r.URL.Query().Get("name"),
		Type: r.URL.Query().Get("type"),
	}
	if v := r.URL.Query().Get("enabled"); v != "" {
		if v == "true" {
			b := true
			filter.Enabled = &b
		} else if v == "false" {
			b := false
			filter.Enabled = &b
		}
	}
	items, total, err := s.svc.ListAlgorithms(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getAlgorithm(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	a, err := s.svc.GetAlgorithm(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, a)
}

type updateAlgorithmRequest struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	KeySize int    `json:"key_size"`
	Enabled bool   `json:"enabled"`
}

func (s *Server) updateAlgorithm(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateAlgorithmRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	a, err := s.svc.UpdateAlgorithm(id, model.Algorithm{
		Name:    req.Name,
		Type:    req.Type,
		KeySize: req.KeySize,
		Enabled: req.Enabled,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, a)
}

func (s *Server) deleteAlgorithm(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteAlgorithm(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
