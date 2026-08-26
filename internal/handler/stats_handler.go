package handler

import (
	"net/http"

	"signatureservice/pkg/httpx"
)

func (s *Server) registerStatsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/stats/overview", s.overviewStats)
}

func (s *Server) overviewStats(w http.ResponseWriter, r *http.Request) {
	stats, err := s.svc.OverviewStats()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, stats)
}
