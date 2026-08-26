package service

import "testing"

func TestDefaultPolicyRejectsWithoutPanicAndRecordsResult(t *testing.T) {
	defer func() {
		if recovered := recover(); recovered != nil {
			t.Fatalf("default policy path panicked instead of returning a rejection: %v", recovered)
		}
	}()
	verifier := NewPolicyVerifier(nil)
	allowed, err := verifier.Check("rsa")
	if err == nil || allowed {
		t.Fatalf("an unconfigured policy should reject rsa with a validation error")
	}
	if result, ok := verifier.Result("rsa"); !ok || result {
		t.Fatalf("the rejected result should be retained as false, got result=%v ok=%v", result, ok)
	}
}
