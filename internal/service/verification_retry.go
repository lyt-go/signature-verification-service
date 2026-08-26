package service

import (
	"signatureservice/internal/model"
	"signatureservice/internal/store"
)

func RunVerificationRetry(st *store.VerificationJobStore, id string, delay func(func()), sideEffect func()) {
	first := &model.VerificationJob{ID: id, Version: 1, Status: "running"}
	_ = st.Save(first)
	sideEffect()
	delay(func() { first.Status = "running"; _ = st.Save(first) })
	retry := &model.VerificationJob{ID: id, Version: 2, Status: "succeeded"}
	sideEffect()
	_ = st.Save(retry)
}
