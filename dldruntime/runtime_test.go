package dldruntime

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestStdioVersionAndShutdown(t *testing.T) {
	rt, err := New(Config{Store: NewMemoryStore()})
	if err != nil {
		t.Fatal(err)
	}
	in := strings.NewReader("{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"Version.Get\",\"params\":{}}\n{\"jsonrpc\":\"2.0\",\"id\":2,\"method\":\"Meta.Shutdown\",\"params\":{}}\n")
	var out bytes.Buffer
	if err := ServeStdio(context.Background(), in, &out, rt); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 responses, got %d: %s", len(lines), out.String())
	}
	if !strings.Contains(lines[0], "versionString") {
		t.Fatalf("Version.Get response missing versionString: %s", lines[0])
	}
}

func TestSyntheticOAuthAndFetchNoRealProviderTraffic(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/oauth/token", func(w http.ResponseWriter, req *http.Request) {
		if err := req.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if req.Form.Get("code") != "synthetic-code" || req.Form.Get("code_verifier") != "synthetic-verifier" {
			t.Fatalf("unexpected OAuth params")
		}
		_, _ = io.WriteString(w, `{"key":{"id":1,"userId":42,"key":"synthetic-api-key"},"cookie":{"itchio":"synthetic-cookie"}}`)
	})
	mux.HandleFunc("/profile", func(w http.ResponseWriter, req *http.Request) {
		if req.Header.Get("Authorization") != "synthetic-api-key" {
			t.Fatalf("missing synthetic API authorization")
		}
		_, _ = io.WriteString(w, `{"user":{"id":42,"username":"fixture","displayName":"Fixture"}}`)
	})
	mux.HandleFunc("/profile/games", func(w http.ResponseWriter, req *http.Request) {
		_, _ = io.WriteString(w, `{"games":[{"id":7,"title":"Fixture Game","type":"default","classification":"game","published":true,"platforms":{"windows":"amd64"}}]}`)
	})
	mux.HandleFunc("/games/7", func(w http.ResponseWriter, req *http.Request) {
		_, _ = io.WriteString(w, `{"game":{"id":7,"title":"Fixture Game","type":"default","classification":"game","platforms":{"windows":"amd64"}}}`)
	})
	mux.HandleFunc("/games/7/uploads", func(w http.ResponseWriter, req *http.Request) {
		_, _ = io.WriteString(w, `{"uploads":[{"id":9,"storage":"hosted","filename":"fixture.zip","size":100,"type":"default","platforms":{"windows":"amd64"}},{"id":10,"storage":"hosted","filename":"fixture.7z","size":100,"type":"default","platforms":{"windows":"amd64"}}]}`)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	rt, err := New(Config{Store: NewMemoryStore(), APIBaseURL: srv.URL, HTTPClient: srv.Client(), StorageRoot: t.TempDir(), StagingRoot: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}

	loginRaw, _ := json.Marshal(ProfileLoginWithOAuthCodeParams{Code: "synthetic-code", CodeVerifier: "synthetic-verifier", RedirectURI: "https://127.0.0.1/callback", ClientID: "synthetic-client"})
	loginAny, _, err := rt.Handle(context.Background(), "Profile.LoginWithOAuthCode", loginRaw)
	if err != nil {
		t.Fatal(err)
	}
	login := loginAny.(*ProfileLoginWithOAuthCodeResult)
	if login.Profile.ID != 42 {
		t.Fatalf("unexpected profile id: %d", login.Profile.ID)
	}

	gamesRaw := []byte(`{"profileId":42}`)
	gamesAny, _, err := rt.Handle(context.Background(), "Fetch.ProfileGames", gamesRaw)
	if err != nil {
		t.Fatal(err)
	}
	games := gamesAny.(*FetchProfileGamesResult)
	if len(games.Items) != 1 || games.Items[0].Game.ID != 7 {
		t.Fatalf("unexpected games: %+v", games)
	}

	uploadsAny, _, err := rt.Handle(context.Background(), "Install.GetUploads", []byte(`{"gameId":7,"profileId":42}`))
	if err != nil {
		t.Fatal(err)
	}
	uploads := uploadsAny.(*InstallGetUploadsResult)
	if len(uploads.Uploads) != 1 || uploads.Uploads[0].ID != 9 || len(uploads.IncompatibleUploads) != 1 {
		t.Fatalf("unexpected compatibility split: %+v", uploads)
	}
}

func TestNoCredentialInGenericRemoteError(t *testing.T) {
	store := NewMemoryStore()
	_ = store.SaveProfile(StoredProfile{ID: 1, APIKey: "top-secret", User: &User{ID: 1}})
	rt, err := New(Config{Store: store, APIBaseURL: "http://127.0.0.1:1", HTTPClient: &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return nil, &url.Error{Op: "Get", URL: req.URL.String(), Err: io.EOF}
	})}})
	if err != nil {
		t.Fatal(err)
	}
	err = redactError(&url.Error{Op: "Get", URL: "https://example.invalid/?api_key=top-secret", Err: io.EOF})
	if strings.Contains(err.Error(), "top-secret") {
		t.Fatal("credential leaked through error")
	}
	_ = rt
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
