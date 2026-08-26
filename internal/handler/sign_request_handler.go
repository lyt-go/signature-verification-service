package handler

import (
	"net/http"

	"signatureservice/internal/model"
	"signatureservice/pkg/httpx"
)

func (s *Server) registerSignRequestRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/signrequests", s.createSignRequest)
	mux.HandleFunc("GET /api/signrequests", s.listSignRequests)
	mux.HandleFunc("GET /api/signrequests/{id}", s.getSignRequest)
	mux.HandleFunc("PUT /api/signrequests/{id}", s.updateSignRequest)
	mux.HandleFunc("DELETE /api/signrequests/{id}", s.deleteSignRequest)
	mux.HandleFunc("PATCH /api/signrequests/{id}/status", s.updateSignRequestStatus)
	mux.HandleFunc("POST /api/signrequests/{id}/process", s.processSignRequest)
}

type createSignRequestRequest struct {
	KeyPairID   string `json:"key_pair_id"`
	Algorithm   string `json:"algorithm"`
	PayloadHash string `json:"payload_hash"`
	RequestID   string `json:"request_id"`
}

func (s *Server) createSignRequest(w http.ResponseWriter, r *http.Request) {
	var req createSignRequestRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	sr, err := s.svc.CreateSignRequest(model.SignRequest{
		KeyPairID:   req.KeyPairID,
		Algorithm:   req.Algorithm,
		PayloadHash: req.PayloadHash,
		RequestID:   req.RequestID,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, sr)
}

func (s *Server) listSignRequests(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.SignRequestFilter{
		KeyPairID: r.URL.Query().Get("key_pair_id"),
		Algorithm: r.URL.Query().Get("algorithm"),
		Status:    r.URL.Query().Get("status"),
		RequestID: r.URL.Query().Get("request_id"),
	}
	items, total, err := s.svc.ListSignRequests(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getSignRequest(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	sr, err := s.svc.GetSignRequest(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sr)
}

type updateSignRequestRequest struct {
	Algorithm   string `json:"algorithm"`
	PayloadHash string `json:"payload_hash"`
	RequestID   string `json:"request_id"`
}

func (s *Server) updateSignRequest(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateSignRequestRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	sr, err := s.svc.UpdateSignRequest(id, model.SignRequest{
		Algorithm:   req.Algorithm,
		PayloadHash: req.PayloadHash,
		RequestID:   req.RequestID,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sr)
}

func (s *Server) deleteSignRequest(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteSignRequest(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) updateSignRequestStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateStatusRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	sr, err := s.svc.UpdateSignRequestStatus(id, req.Status)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sr)
}

func (s *Server) processSignRequest(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	sig, err := s.svc.ProcessSignRequest(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sig)
}
