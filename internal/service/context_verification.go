package service

import (
	"context"
	"signatureservice/internal/store"
)

func AwaitVerification(ctx context.Context, probe *store.VerificationProbe, ready <-chan struct{}) error {
	return probe.Wait(ctx, ready)
}

func AwaitDirectVerification(ctx context.Context, probe *store.VerificationProbe, ready <-chan struct{}) error {
	return probe.Wait(ctx, ready)
}
