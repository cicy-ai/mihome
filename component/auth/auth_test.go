package auth

import "testing"

func TestNewAuthenticatorRequiresDeclaredUser(t *testing.T) {
	a := NewAuthenticator([]AuthUser{{User: "agent-a", Pass: "ignored"}, {User: "agent-b", Pass: "ignored"}}, "")
	if a == nil {
		t.Fatal("expected authenticator")
	}
	if !a.Verify("agent-a", "agent-a") {
		t.Fatal("expected username-equals-password fallback to work")
	}
	if a.Verify("agent-a", "wrong") {
		t.Fatal("unexpected password accepted without api_token")
	}
	if a.Verify("unknown", "unknown") {
		t.Fatal("unknown user should not authenticate")
	}
}

func TestNewAuthenticatorAPITokenTakesPriority(t *testing.T) {
	a := NewAuthenticator([]AuthUser{{User: "agent-a", Pass: "ignored"}, {User: "agent-b", Pass: "ignored"}}, "shared-pass")
	if a == nil {
		t.Fatal("expected authenticator")
	}
	if !a.Verify("agent-a", "shared-pass") {
		t.Fatal("expected api_token/global password to work")
	}
	if a.Verify("agent-a", "agent-a") {
		t.Fatal("username fallback should not apply when api_token exists")
	}
	if a.Verify("unknown", "shared-pass") {
		t.Fatal("global password must not authenticate unknown users")
	}
}
