package handler

import (
	"net/http"

	"signatureservice/internal/model"
	"signatureservice/pkg/httpx"
)

func (s *Server) registerKeyPairRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/keypairs", s.createKeyPair)
	mux.HandleFunc("GET /api/keypairs", s.listKeyPairs)
	mux.HandleFunc("GET /api/keypairs/{id}", s.getKeyPair)
	mux.HandleFunc("PUT /api/keypairs/{id}", s.updateKeyPair)
	mux.HandleFunc("DELETE /api/keypairs/{id}", s.deleteKeyPair)
	mux.HandleFunc("PATCH /api/keypairs/{id}/status", s.updateKeyPairStatus)
}

type createKeyPairRequest struct {
	Name      string `json:"name"`
	Algorithm string `json:"algorithm"`
	KeySize   int    `json:"key_size"`
	PublicKey string `json:"public_key"`
}

func (s *Server) createKeyPair(w http.ResponseWriter, r *http.Request) {
	var req createKeyPairRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	kp, err := s.svc.CreateKeyPair(model.KeyPair{
		Name:      req.Name,
		Algorithm: req.Algorithm,
		KeySize:   req.KeySize,
		PublicKey: req.PublicKey,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, kp)
}

func (s *Server) listKeyPairs(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.KeyPairFilter{
		Name:      r.URL.Query().Get("name"),
		Algorithm: r.URL.Query().Get("algorithm"),
		Status:    r.URL.Query().Get("status"),
		Keyword:   r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListKeyPairs(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getKeyPair(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	kp, err := s.svc.GetKeyPair(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, kp)
}

type updateKeyPairRequest struct {
	Name      string `json:"name"`
	KeySize   int    `json:"key_size"`
	PublicKey string `json:"public_key"`
}

func (s *Server) updateKeyPair(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateKeyPairRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	kp, err := s.svc.UpdateKeyPair(id, model.KeyPair{
		Name:      req.Name,
		KeySize:   req.KeySize,
		PublicKey: req.PublicKey,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, kp)
}

func (s *Server) deleteKeyPair(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteKeyPair(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

type updateStatusRequest struct {
	Status string `json:"status"`
}

func (s *Server) updateKeyPairStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateStatusRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	kp, err := s.svc.UpdateKeyPairStatus(id, req.Status)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, kp)
}
