package store

import (
	"context"
)

// VerificationProbe 用于等待验签就绪，直到 ready 通道关闭或 ctx 被取消。
// 它不持有任何请求上下文：每次 Wait 只使用调用方传入的 ctx，
// 避免一次取消的请求污染后续请求。
type VerificationProbe struct{}

func (p *VerificationProbe) Wait(ctx context.Context, ready <-chan struct{}) error {
	select {
	case <-ready:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
