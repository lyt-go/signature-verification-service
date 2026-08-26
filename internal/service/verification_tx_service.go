package service

import (
	"errors"
	"signatureservice/internal/store"
	"signatureservice/pkg/logger"
)

// PersistVerifiedWithRetry 持久化验签记录并在提交成功后发出通知。
// 只有 Commit 成功才表示记录真正落库，因此通知必须在提交成功之后发出，
// 否则遇到临时提交错误重试时会重复发送成功通知并产生重复记录。
func PersistVerifiedWithRetry(st *store.VerificationTxStore, publisher logger.VerificationPublisher, id string) error {
	for attempt := 0; attempt < 2; attempt++ {
		tx := st.Begin(id)
		if err := tx.Commit(); err != nil {
			_ = tx.Rollback()
			if errors.Is(err, store.ErrTemporaryCommit) {
				continue
			}
			return err
		}
		// 记录已成功落库，之后再发出对外的成功通知。
		// 若通知失败，记录本身已存在；由调用方决定后续处理，这里不重试通知。
		if err := publisher.PublishVerification(id); err != nil {
			return err
		}
		return nil
	}
	return store.ErrTemporaryCommit
}
