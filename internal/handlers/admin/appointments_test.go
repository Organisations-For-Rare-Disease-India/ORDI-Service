package admin

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNoofDays(t *testing.T) {
	days := noOfDays(time.Now())
	if days < 28 || days > 31 {
		t.Fatalf("invalid days:%d\t received", days)
	}
	t.Logf("no of days:%d\n", days)

}

func TestToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			c := &http.Cookie{
				Name:  t.Name(),
				Value: t.Name(),
			}
			http.SetCookie(w, c)

		}))
	defer server.Close()
	client := http.DefaultClient
	res, err := client.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	cookie := res.Header.Get("Set-Cookie")
	t.Logf("cookie:%s\n", cookie)

}
