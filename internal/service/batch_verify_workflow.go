package service

import (
	"context"
	"fmt"

	"signatureservice/internal/store"
)

// VerifyBatch 依次对一批签名 ID 执行验签。
// 每条记录都在自己的 verifyOne 调用中获取并释放槽位，
// 避免 defer 在循环中累积导致槽位不释放（进而触发超时）。
func VerifyBatch(ctx context.Context, pool *store.SlotPool, signatureIDs []string) error {
	for _, id := range signatureIDs {
		if err := verifyOne(ctx, pool, id); err != nil {
			return err
		}
	}
	return nil
}

// verifyOne 处理单条签名：获取槽位后立即在函数返回时释放，
// 因此每轮循环的槽位都会在该轮结束时归还，不会跨迭代累积。
func verifyOne(ctx context.Context, pool *store.SlotPool, id string) error {
	lease, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer lease.Close()
	if id == "" {
		lease.Finish(false)
		return fmt.Errorf("signature id is empty")
	}
	lease.Finish(true)
	return nil
}
