package main

import (
	"net/http"
	"strings"
	"testing"
)

func TestSessionStoreCookieOptions(t *testing.T) {
	t.Setenv("SESSION_SECRET", strings.Repeat("k", minSessionSecretLength))

	t.Setenv("COOKIE_SECURE", "")
	o := newSessionStore().Options
	if !o.HttpOnly || o.SameSite != http.SameSiteLaxMode || o.Secure || o.Path != "/" {
		t.Fatalf("%+v", o)
	}

	t.Setenv("COOKIE_SECURE", "1")
	if !newSessionStore().Options.Secure {
		t.Fatal("COOKIE_SECURE=1 must set the Secure flag")
	}
}
