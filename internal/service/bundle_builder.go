package service

import (
	"fmt"
	"signatureservice/internal/model"
	"signatureservice/internal/store"
)

func BuildCertificateBundle(cache *store.BundleCache, id string, shouldPanic bool) (bundle *store.CertificateBundle, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("bundle build failed: %v", recovered)
		}
	}()
	bundle = &store.CertificateBundle{ID: id}
	cache.Put(bundle)
	bundle.Certificate = &model.Certificate{ID: "cert-" + id, KeyPairID: "kp", Subject: "subject", Issuer: "issuer", SerialNumber: id}
	if shouldPanic {
		panic("signature decoder rejected input")
	}
	bundle.Signature = &model.Signature{ID: "sig-" + id, SignRequestID: "request", KeyPairID: "kp", Value: "value", Algorithm: "rsa"}
	cache.Put(bundle)
	return bundle, nil
}
