package service

import (
	"fmt"
	"signatureservice/internal/model"
	"signatureservice/internal/store"
)

func BuildCertificateBundle(cache *store.BundleCache, id string, shouldPanic bool) (bundle *store.CertificateBundle, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			// 签名解码失败：把 panic 收敛为普通错误，服务不崩；
			// 同时不把残缺包交回调用方。
			bundle = nil
			err = fmt.Errorf("bundle build failed: %v", recovered)
		}
	}()
	bundle = &store.CertificateBundle{ID: id}
	bundle.Certificate = &model.Certificate{ID: "cert-" + id, KeyPairID: "kp", Subject: "subject", Issuer: "issuer", SerialNumber: id}
	if shouldPanic {
		panic("signature decoder rejected input")
	}
	bundle.Signature = &model.Signature{ID: "sig-" + id, SignRequestID: "request", KeyPairID: "kp", Value: "value", Algorithm: "rsa"}
	// 只在组合包构建完整、校验通过后才写入缓存，避免解码失败时把只含部分内容的包缓存下来。
	if err = store.ValidateBundle(bundle); err != nil {
		bundle = nil
		return nil, err
	}
	cache.Put(bundle)
	return bundle, nil
}
