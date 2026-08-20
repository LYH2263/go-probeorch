package probeorch_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"example.com/probeorch"
)

type closeTrackBody struct {
	io.ReadCloser
	closed *int32
}

func (c closeTrackBody) Close() error {
	atomic.AddInt32(c.closed, 1)
	return c.ReadCloser.Close()
}

func TestBug09_HTTPResponseBodyClosed(t *testing.T) {
	var closed int32
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte("ok"))
	}))
	defer ts.Close()

	base := ts.Client()
	rt := roundTripClose{base: base.Transport, closed: &closed}
	client := *base
	client.Transport = rt

	o := probeorch.New(probeorch.WithHTTPClient(&client))
	defer o.Close()
	id, err := o.Register(probeorch.Spec{
		Name: "http", Address: ts.URL, Kind: probeorch.KindHTTP,
		Interval: time.Second, Timeout: time.Second, Expect: 200,
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := o.RunOnce(id); err != nil {
		t.Fatal(err)
	}
	if atomic.LoadInt32(&closed) < 1 {
		t.Fatal("HTTP response Body was not Closed")
	}
}

type roundTripClose struct {
	base   http.RoundTripper
	closed *int32
}

func (r roundTripClose) RoundTrip(req *http.Request) (*http.Response, error) {
	if r.base == nil {
		r.base = http.DefaultTransport
	}
	resp, err := r.base.RoundTrip(req)
	if err != nil {
		return resp, err
	}
	resp.Body = closeTrackBody{ReadCloser: resp.Body, closed: r.closed}
	return resp, nil
}
