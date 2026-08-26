package service

import (
	"signatureservice/internal/model"
	"signatureservice/internal/store"
)

// RunVerificationRetry 模拟一次验签任务的失败重试。
//
// 时序问题：第一次尝试会注册一个延迟回调（更新状态），重试成功后若该回调迟到，
// 旧实现会用过期的 running 状态覆盖已落库的 succeeded，并且对外动作被执行了两次。
//
// 修复策略：
//   - 对外动作（sideEffect）仅在重试成功时执行一次，第一次尝试不再重复触发。
//   - 延迟回调改为基于版本号的 CAS 写入（SaveIfVersion）。
//     任务版本已随重试前进到 v2 时，基于 v1 的回调被视为过期回调丢弃，绝不覆盖 succeeded。
func RunVerificationRetry(st *store.VerificationJobStore, id string, delay func(func()), sideEffect func()) {
	first := &model.VerificationJob{ID: id, Version: 1, Status: "running"}
	_ = st.Save(first)
	delay(func() {
		// 仅当任务仍停留在本次尝试（v1）时才回写 running；
		// 若已重试到更新版本，本次回调即为迟到回调，直接丢弃。
		first.Status = "running"
		_, _ = st.SaveIfVersion(first, first.Version)
	})
	retry := &model.VerificationJob{ID: id, Version: 2, Status: "succeeded"}
	sideEffect()
	_ = st.Save(retry)
}
