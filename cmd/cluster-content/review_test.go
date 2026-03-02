package main

import "testing"

func TestValidateReviewListenAddr(t *testing.T) {
	if err := validateReviewListenAddr("127.0.0.1:8097", false); err != nil {
		t.Fatalf("localhost should pass: %v", err)
	}
	if err := validateReviewListenAddr("localhost:8097", false); err != nil {
		t.Fatalf("localhost name should pass: %v", err)
	}
	if err := validateReviewListenAddr("0.0.0.0:8097", false); err == nil {
		t.Fatal("expected non-local bind to fail without override")
	}
	if err := validateReviewListenAddr("0.0.0.0:8097", true); err != nil {
		t.Fatalf("override should allow non-local bind: %v", err)
	}
}
