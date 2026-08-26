package service

import (
	"errors"
	"signatureservice/internal/store"
	"signatureservice/pkg/logger"
)

func PersistVerifiedWithRetry(st *store.VerificationTxStore, publisher logger.VerificationPublisher, id string) error {
	for attempt := 0; attempt < 2; attempt++ {
		tx := st.Begin(id)
		if err := publisher.PublishVerification(id); err != nil {
			_ = tx.Rollback()
			return err
		}
		if err := tx.Commit(); err != nil {
			_ = tx.Rollback()
			if errors.Is(err, store.ErrTemporaryCommit) {
				continue
			}
			return err
		}
		return nil
	}
	return store.ErrTemporaryCommit
}
