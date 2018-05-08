package util

import "testing"

func TestGetCookieDomain(t *testing.T) {
	host := "www.ts.com:8080"
	s := getCookieDomain(host)
	if s != "ts.com" {
		t.Log("s = ", s)
		t.Failed()
	}
}
