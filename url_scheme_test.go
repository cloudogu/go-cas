package cas

import (
	"testing"
	"net/url"
)

func TestDefaultURLScheme(t *testing.T) {
	url, _ := url.Parse("https://cas.org/cas")
	scheme := NewDefaultURLScheme(url)

	u, err := scheme.Login()
	assertUrl(t, "/cas/login", u, err)
	u, err = scheme.Logout()
	assertUrl(t, "/cas/logout", u, err)
	u, err = scheme.Validate()
	assertUrl(t, "/cas/validate", u, err)
	u, err = scheme.ServiceValidate()
	assertUrl(t, "/cas/serviceValidate", u, err)
}

func assertUrl(t *testing.T, expected string, u *url.URL, err error) {
	if err != nil {
		t.Fatalf("returned error")
	}

	if expected != u.Path {
		t.Errorf("%s should be equal to %s", u.Path, expected)
	}
}