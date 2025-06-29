package lib

import (
	"net/http/httptest"
	"testing"
)

func TestSetCookie(t *testing.T) {
	for _, tt := range []struct {
		name       string
		options    Options
		host       string
		cookieName string
	}{
		{
			name:       "basic",
			options:    Options{},
			host:       "",
			cookieName: "techaro.lol-anubis-auth",
		},
		{
			name:       "domain techaro.lol",
			options:    Options{CookieDomain: "techaro.lol"},
			host:       "",
			cookieName: "techaro.lol-anubis-auth-for-techaro.lol",
		},
		{
			name:       "dynamic cookie domain",
			options:    Options{CookieDynamicDomain: true},
			host:       "techaro.lol",
			cookieName: "techaro.lol-anubis-auth-for-techaro.lol",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			srv := spawnAnubis(t, tt.options)
			rw := httptest.NewRecorder()

			srv.SetCookie(rw, srv.cookieName, "test", "/", tt.host)

			resp := rw.Result()
			cookies := resp.Cookies()

			ckie := cookies[0]

			if ckie.Name != tt.cookieName {
				t.Errorf("wanted cookie named %q, got cookie named %q", tt.cookieName, ckie.Name)
			}
		})
	}
}

func TestClearCookie(t *testing.T) {
	srv := spawnAnubis(t, Options{})
	rw := httptest.NewRecorder()

	srv.ClearCookie(rw, srv.cookieName, "/", "localhost")

	resp := rw.Result()

	cookies := resp.Cookies()

	if len(cookies) != 1 {
		t.Errorf("wanted 1 cookie, got %d cookies", len(cookies))
	}

	ckie := cookies[0]

	if ckie.Name != "techaro.lol-anubis-auth" {
		t.Errorf("wanted cookie named %q, got cookie named %q", "techaro.lol-anubis-auth", ckie.Name)
	}

	if ckie.MaxAge != -1 {
		t.Errorf("wanted cookie max age of -1, got: %d", ckie.MaxAge)
	}
}

func TestClearCookieWithDomain(t *testing.T) {
	srv := spawnAnubis(t, Options{CookieDomain: "techaro.lol"})
	rw := httptest.NewRecorder()

	srv.ClearCookie(rw, srv.cookieName, "/", "locahost")

	resp := rw.Result()

	cookies := resp.Cookies()

	if len(cookies) != 1 {
		t.Errorf("wanted 1 cookie, got %d cookies", len(cookies))
	}

	ckie := cookies[0]

	if ckie.Name != srv.cookieName {
		t.Errorf("wanted cookie named %q, got cookie named %q", srv.cookieName, ckie.Name)
	}

	if ckie.MaxAge != -1 {
		t.Errorf("wanted cookie max age of -1, got: %d", ckie.MaxAge)
	}
}

func TestClearCookieWithDynamicDomain(t *testing.T) {
	srv := spawnAnubis(t, Options{CookieDynamicDomain: true})
	rw := httptest.NewRecorder()

	srv.ClearCookie(rw, srv.cookieName, "/", "xeiaso.net")

	resp := rw.Result()

	cookies := resp.Cookies()

	if len(cookies) != 1 {
		t.Errorf("wanted 1 cookie, got %d cookies", len(cookies))
	}

	ckie := cookies[0]

	if ckie.Name != "techaro.lol-anubis-auth-for-xeiaso.net" {
		t.Errorf("wanted cookie named %q, got cookie named %q", srv.cookieName, ckie.Name)
	}

	if ckie.MaxAge != -1 {
		t.Errorf("wanted cookie max age of -1, got: %d", ckie.MaxAge)
	}
}
