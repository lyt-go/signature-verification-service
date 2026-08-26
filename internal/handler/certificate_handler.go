package handler

import (
	"net/http"

	"signatureservice/internal/model"
	"signatureservice/pkg/httpx"
)

func (s *Server) registerCertificateRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/certificates", s.createCertificate)
	mux.HandleFunc("GET /api/certificates", s.listCertificates)
	mux.HandleFunc("GET /api/certificates/{id}", s.getCertificate)
	mux.HandleFunc("PUT /api/certificates/{id}", s.updateCertificate)
	mux.HandleFunc("DELETE /api/certificates/{id}", s.deleteCertificate)
	mux.HandleFunc("PATCH /api/certificates/{id}/status", s.updateCertificateStatus)
}

type createCertificateRequest struct {
	KeyPairID    string `json:"key_pair_id"`
	Subject      string `json:"subject"`
	Issuer       string `json:"issuer"`
	SerialNumber string `json:"serial_number"`
}

func (s *Server) createCertificate(w http.ResponseWriter, r *http.Request) {
	var req createCertificateRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	c, err := s.svc.CreateCertificate(model.Certificate{
		KeyPairID:    req.KeyPairID,
		Subject:      req.Subject,
		Issuer:       req.Issuer,
		SerialNumber: req.SerialNumber,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, c)
}

func (s *Server) listCertificates(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.CertificateFilter{
		KeyPairID: r.URL.Query().Get("key_pair_id"),
		Status:    r.URL.Query().Get("status"),
		Subject:   r.URL.Query().Get("subject"),
		Issuer:    r.URL.Query().Get("issuer"),
	}
	items, total, err := s.svc.ListCertificates(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getCertificate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	c, err := s.svc.GetCertificate(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}

type updateCertificateRequest struct {
	Subject      string `json:"subject"`
	Issuer       string `json:"issuer"`
	SerialNumber string `json:"serial_number"`
}

func (s *Server) updateCertificate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateCertificateRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	c, err := s.svc.UpdateCertificate(id, model.Certificate{
		Subject:      req.Subject,
		Issuer:       req.Issuer,
		SerialNumber: req.SerialNumber,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}

func (s *Server) deleteCertificate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteCertificate(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) updateCertificateStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateStatusRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	c, err := s.svc.UpdateCertificateStatus(id, req.Status)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}
