package handler

import (
	"net/http"
	"strconv"

	"signatureservice/internal/model"
	"signatureservice/pkg/httpx"
)

func (s *Server) registerVerifyRecordRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/verifyrecords", s.createVerifyRecord)
	mux.HandleFunc("GET /api/verifyrecords", s.listVerifyRecords)
	mux.HandleFunc("GET /api/verifyrecords/{id}", s.getVerifyRecord)
	mux.HandleFunc("PUT /api/verifyrecords/{id}", s.updateVerifyRecord)
	mux.HandleFunc("DELETE /api/verifyrecords/{id}", s.deleteVerifyRecord)
	mux.HandleFunc("POST /api/verifyrecords/batch", s.batchCreateVerifyRecords)
}

type createVerifyRecordRequest struct {
	SignatureID string `json:"signature_id"`
	PayloadHash string `json:"payload_hash"`
	Valid       bool   `json:"valid"`
	Verifier    string `json:"verifier"`
}

func (s *Server) createVerifyRecord(w http.ResponseWriter, r *http.Request) {
	var req createVerifyRecordRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	vr, err := s.svc.CreateVerifyRecord(model.VerifyRecord{
		SignatureID: req.SignatureID,
		PayloadHash: req.PayloadHash,
		Valid:       req.Valid,
		Verifier:    req.Verifier,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, vr)
}

func (s *Server) listVerifyRecords(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.VerifyRecordFilter{
		SignatureID: r.URL.Query().Get("signature_id"),
		Verifier:    r.URL.Query().Get("verifier"),
	}
	if v := r.URL.Query().Get("valid"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			filter.Valid = &b
		}
	}
	items, total, err := s.svc.ListVerifyRecords(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getVerifyRecord(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	vr, err := s.svc.GetVerifyRecord(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, vr)
}

type updateVerifyRecordRequest struct {
	PayloadHash string `json:"payload_hash"`
	Verifier    string `json:"verifier"`
}

func (s *Server) updateVerifyRecord(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateVerifyRecordRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	vr, err := s.svc.UpdateVerifyRecord(id, model.VerifyRecord{
		PayloadHash: req.PayloadHash,
		Verifier:    req.Verifier,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, vr)
}

func (s *Server) deleteVerifyRecord(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteVerifyRecord(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

type batchCreateVerifyRecordsRequest struct {
	Records []createVerifyRecordRequest `json:"records"`
}

func (s *Server) batchCreateVerifyRecords(w http.ResponseWriter, r *http.Request) {
	var req batchCreateVerifyRecordsRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	inputs := make([]model.VerifyRecord, len(req.Records))
	for i, rec := range req.Records {
		inputs[i] = model.VerifyRecord{
			SignatureID: rec.SignatureID,
			PayloadHash: rec.PayloadHash,
			Valid:       rec.Valid,
			Verifier:    rec.Verifier,
		}
	}
	if err := s.svc.BatchCreateVerifyRecords(inputs); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, map[string]int{"created": len(inputs)})
}
