package auth

import "testing"

func TestVerifyPerUserPassword(t *testing.T) {
	a := NewAuthenticator([]AuthUser{{User: "agent-a", Pass: "per-user-pass"}}, "")
	if a == nil {
		t.Fatal("expected authenticator")
	}
	if !a.Verify("agent-a", "per-user-pass") {
		t.Fatal("per-user password should authenticate")
	}
	if a.Verify("agent-a", "wrong") {
		t.Fatal("wrong password must not authenticate")
	}
	if a.Verify("unknown", "per-user-pass") {
		t.Fatal("without globalPassword, unknown user cannot use someone else's password")
	}
}

func TestVerifyGlobalPasswordAllowsAnyUser(t *testing.T) {
	a := NewAuthenticator([]AuthUser{{User: "agent-a", Pass: "ignored-here"}}, "shared-pass")
	if a == nil {
		t.Fatal("expected authenticator")
	}
	// Per-user password still works when declared.
	if !a.Verify("agent-a", "ignored-here") {
		t.Fatal("declared per-user password should still work")
	}
	// globalPassword works for the declared user.
	if !a.Verify("agent-a", "shared-pass") {
		t.Fatal("declared user should authenticate via globalPassword")
	}
	// globalPassword also works for users NOT in the authentication list.
	if !a.Verify("agent-b-not-declared", "shared-pass") {
		t.Fatal("globalPassword must authenticate undeclared users too")
	}
	// Wrong password is still rejected for both declared and undeclared users.
	if a.Verify("agent-a", "wrong") {
		t.Fatal("wrong password must not authenticate declared user")
	}
	if a.Verify("agent-b-not-declared", "wrong") {
		t.Fatal("wrong password must not authenticate undeclared user")
	}
	// Empty user is always rejected.
	if a.Verify("", "shared-pass") {
		t.Fatal("empty user must not authenticate")
	}
}

func TestNewAuthenticatorGlobalOnly(t *testing.T) {
	a := NewAuthenticator(nil, "shared-pass")
	if a == nil {
		t.Fatal("expected authenticator with globalPassword only")
	}
	if !a.Verify("anyone", "shared-pass") {
		t.Fatal("global-only authenticator should accept any non-empty user with the global pass")
	}
	if a.Verify("anyone", "wrong") {
		t.Fatal("wrong password must not authenticate")
	}
}
