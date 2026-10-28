package auth

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/wrestle-R/Spotify-Tui/internal/config"
	"github.com/zalando/go-keyring"
)

func TestPKCEAndCallbackState(t *testing.T) {
	if got := Challenge("dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk"); got != "E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM" {
		t.Fatal("RFC PKCE vector mismatch")
	}
	var delivered atomic.Int32
	h := CallbackHandler("/callback", "expected", func(code string, err error) {
		if code != "valid" || err != nil {
			t.Error("unexpected result")
		}
		delivered.Add(1)
	})
	for _, path := range []string{"/callback?state=wrong&code=valid", "/different?state=expected&code=valid", "/callback?state=expected"} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest("GET", path, nil))
		if w.Code < 400 {
			t.Fatal("invalid callback accepted")
		}
	}
	if delivered.Load() != 0 {
		t.Fatal("invalid callback delivered a code")
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/callback?state=expected&code=valid", nil))
	if w.Code != 200 || delivered.Load() != 1 {
		t.Fatal("valid callback rejected")
	}
}

func TestLoginRoundTripWithoutRealCredentials(t *testing.T) {
	keyring.MockInit()
	listener, e := net.Listen("tcp", "127.0.0.1:0")
	if e != nil {
		t.Fatal(e)
	}
	address := listener.Addr().String()
	listener.Close()
	var challenge string
	tokens := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		if r.Form.Get("code") != "fixture-code" || r.Form.Get("grant_type") != "authorization_code" || Challenge(r.Form.Get("code_verifier")) != challenge {
			t.Error("invalid token exchange")
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"access_token":"fixture-access","refresh_token":"fixture-refresh","expires_in":3600}`)
	}))
	defer tokens.Close()
	m := New(config.Config{ClientID: strings.Repeat("a", 32), RedirectURI: "http://" + address + "/callback"}, config.Paths{Data: t.TempDir()})
	m.TokenURL = tokens.URL
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.Login(ctx, func(link string) error {
		u, e := url.Parse(link)
		if e != nil {
			return e
		}
		q := u.Query()
		challenge = q.Get("code_challenge")
		if q.Get("code_challenge_method") != "S256" {
			t.Error("PKCE not enabled")
		}
		resp, e := http.Get(q.Get("redirect_uri") + "?state=" + url.QueryEscape(q.Get("state")) + "&code=fixture-code")
		if e != nil {
			return e
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			t.Error("callback not accepted")
		}
		return nil
	}, io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	token, err := m.AccessToken(ctx)
	if err != nil || token != "fixture-access" {
		t.Fatal("token was not persisted")
	}
}

func TestRefreshIsCoordinatedAndLate401DoesNotInvalidateNewToken(t *testing.T) {
	keyring.MockInit()
	var refreshes atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		refreshes.Add(1)
		fmt.Fprint(w, `{"access_token":"new-access","expires_in":3600}`)
	}))
	defer server.Close()
	m := New(config.Config{ClientID: "fixture"}, config.Paths{Data: t.TempDir()})
	m.TokenURL = server.URL
	if err := m.Save(Token{AccessToken: "expired-access", RefreshToken: "fixture-refresh", ExpiresAt: time.Now().Add(-time.Hour)}); err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	for i := 0; i < 12; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s, e := m.AccessToken(context.Background())
			if e != nil || s != "new-access" {
				t.Error("refresh failed")
			}
		}()
	}
	wg.Wait()
	if refreshes.Load() != 1 {
		t.Fatalf("expected one refresh, got %d", refreshes.Load())
	}
	m.InvalidateIf("expired-access")
	_, err := m.AccessToken(context.Background())
	if err != nil || refreshes.Load() != 1 {
		t.Fatal("late 401 invalidated refreshed credentials")
	}
}
