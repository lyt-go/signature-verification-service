package service

import (
	"signatureservice/internal/config"
	"signatureservice/internal/model"
)

type signaturePolicy interface{ Allow(string) bool }

type PolicyVerifier struct {
	policy  signaturePolicy
	results map[string]bool
}

func NewPolicyVerifier(names []string) *PolicyVerifier {
	return &PolicyVerifier{
		policy:  config.LoadSignaturePolicy(names),
		results: make(map[string]bool),
	}
}

func (v *PolicyVerifier) Check(algorithm string) (bool, error) {
	if v.policy != nil && v.policy.Allow(algorithm) {
		v.results[algorithm] = true
		return true, nil
	}
	v.results[algorithm] = false
	return false, model.NewValidationError("algorithm", "algorithm is not allowed")
}

func (v *PolicyVerifier) Result(algorithm string) (bool, bool) {
	result, ok := v.results[algorithm]
	return result, ok
}
