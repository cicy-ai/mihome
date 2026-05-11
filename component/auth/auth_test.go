package auth

import "testing"

func TestNewAuthenticatorPerUserPassword(t *testing.T) {
	a := NewAuthenticator([]AuthUser{{User: "agent-a", Pass: "pass-a"}, {User: "agent-b", Pass: "pass-b"}}, "")
	if a == nil {
		t.Fatal("expected authenticator")
	}
	if !a.Verify("agent-a", "pass-a") {
		t.Fatal("expected per-user password to work")
	}
	if a.Verify("agent-a", "pass-b") {
		t.Fatal("unexpected password accepted")
	}
	if a.Verify("agent-c", "pass-a") {
		t.Fatal("unknown user should not authenticate")
	}
}

func TestNewAuthenticatorGlobalPasswordFallback(t *testing.T) {
	a := NewAuthenticator([]AuthUser{{User: "agent-a", Pass: "pass-a"}, {User: "agent-b", Pass: "pass-b"}}, "shared-pass")
	if a == nil {
		t.Fatal("expected authenticator")
	}
	if !a.Verify("agent-a", "pass-a") {
		t.Fatal("expected explicit per-user password to work")
	}
	if !a.Verify("agent-a", "shared-pass") {
		t.Fatal("expected global password fallback to work")
	}
	if !a.Verify("agent-b", "shared-pass") {
		t.Fatal("expected global password fallback to work for other declared users")
	}
	if a.Verify("unknown", "shared-pass") {
		t.Fatal("global password must not authenticate unknown users")
	}
}
