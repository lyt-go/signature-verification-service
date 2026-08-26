package handler

import (
	"net/http"

	"signatureservice/internal/model"
	"signatureservice/pkg/httpx"
)

func (s *Server) registerSignatureRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/signatures", s.createSignature)
	mux.HandleFunc("GET /api/signatures", s.listSignatures)
	mux.HandleFunc("GET /api/signatures/{id}", s.getSignature)
	mux.HandleFunc("PUT /api/signatures/{id}", s.updateSignature)
	mux.HandleFunc("DELETE /api/signatures/{id}", s.deleteSignature)
	mux.HandleFunc("POST /api/signatures/{id}/verify", s.verifySignature)
}

type createSignatureRequest struct {
	SignRequestID string `json:"sign_request_id"`
	KeyPairID     string `json:"key_pair_id"`
	Value         string `json:"value"`
	Algorithm     string `json:"algorithm"`
}

func (s *Server) createSignature(w http.ResponseWriter, r *http.Request) {
	var req createSignatureRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	sig, err := s.svc.CreateSignature(model.Signature{
		SignRequestID: req.SignRequestID,
		KeyPairID:     req.KeyPairID,
		Value:         req.Value,
		Algorithm:     req.Algorithm,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, sig)
}

func (s *Server) listSignatures(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.SignatureFilter{
		KeyPairID:     r.URL.Query().Get("key_pair_id"),
		Algorithm:     r.URL.Query().Get("algorithm"),
		SignRequestID: r.URL.Query().Get("sign_request_id"),
	}
	items, total, err := s.svc.ListSignatures(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getSignature(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	sig, err := s.svc.GetSignature(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sig)
}

type updateSignatureRequest struct {
	Value     string `json:"value"`
	Algorithm string `json:"algorithm"`
}

func (s *Server) updateSignature(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateSignatureRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	sig, err := s.svc.UpdateSignature(id, model.Signature{
		Value:     req.Value,
		Algorithm: req.Algorithm,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sig)
}

func (s *Server) deleteSignature(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteSignature(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

type verifySignatureRequest struct {
	PayloadHash string `json:"payload_hash"`
	Verifier    string `json:"verifier"`
}

func (s *Server) verifySignature(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req verifySignatureRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	vr, err := s.svc.VerifySignature(id, req.PayloadHash, req.Verifier)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, vr)
}
