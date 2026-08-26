package service

import (
	"signatureservice/internal/model"
)

type OverviewStats struct {
	KeyPairByAlgorithm  map[string]int `json:"key_pair_by_algorithm"`
	KeyPairByStatus     map[string]int `json:"key_pair_by_status"`
	SignRequestSuccess  float64        `json:"sign_request_success_rate"`
	VerifyPassRate      float64        `json:"verify_pass_rate"`
	VerifyFailRate      float64        `json:"verify_fail_rate"`
	CertificateByStatus map[string]int `json:"certificate_by_status"`
}

func (s *Service) OverviewStats() (*OverviewStats, error) {
	stats := &OverviewStats{
		KeyPairByAlgorithm:  make(map[string]int),
		KeyPairByStatus:     make(map[string]int),
		CertificateByStatus: make(map[string]int),
	}

	keyPairs := s.store.ListKeyPairs()
	for _, kp := range keyPairs {
		stats.KeyPairByAlgorithm[kp.Algorithm]++
		stats.KeyPairByStatus[kp.Status]++
	}

	signRequests := s.store.ListSignRequests()
	total := len(signRequests)
	signedCount := 0
	for _, sr := range signRequests {
		if sr.Status == model.SignRequestSigned {
			signedCount++
		}
	}
	if total > 0 {
		stats.SignRequestSuccess = float64(signedCount) / float64(total)
	}

	verifyRecords := s.store.ListVerifyRecords()
	verifyTotal := len(verifyRecords)
	passCount := 0
	for _, vr := range verifyRecords {
		if vr.Valid {
			passCount++
		}
	}
	if verifyTotal > 0 {
		stats.VerifyPassRate = float64(passCount) / float64(verifyTotal)
		stats.VerifyFailRate = float64(verifyTotal-passCount) / float64(verifyTotal)
	}

	certificates := s.store.ListCertificates()
	for _, c := range certificates {
		stats.CertificateByStatus[c.Status]++
	}

	return stats, nil
}
