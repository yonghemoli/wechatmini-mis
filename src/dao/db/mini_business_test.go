package db

import "testing"

func TestBusinessStatusTransitions(t *testing.T) {
	if !validDemandTransition("PENDING_CONTACT", "CONTACTED") || validDemandTransition("PENDING_CONTACT", "MATCHING") {
		t.Fatal("demand transition rules are invalid")
	}
	if !validResumeTransition("VERIFYING", "APPROVED") || validResumeTransition("PENDING_CONTACT", "APPROVED") {
		t.Fatal("resume transition rules are invalid")
	}
}
