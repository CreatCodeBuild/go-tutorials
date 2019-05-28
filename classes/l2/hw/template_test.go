package main

import (
	"testing"
	"net/http"
	"net/http/httptest"

	"time"

	"github.com/gorilla/mux"
)

func TestV1(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })        // 200
	r.HandleFunc("/a", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusCreated) })  // 201
	r.HandleFunc("/b", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusAccepted) }) // 202

	server := httptest.NewServer(r)
	defer server.Close()

	client := NewClient()
	client.SetBaseURL(server.URL)

	req := client.Request(http.MethodOptions, "/")
	stdReq, err := req.Final()
	require.NoError(t, err)

	res, err := client.Do(stdReq)
	require.NoError(t, err)
	require.Equal(t, res.StatusCode, http.StatusOK)
	compareReaders(t, res.Body, strings.NewReader(""))

	req2 := client.Request(http.MethodOptions, "/a")
	stdReq2, err := req2.Final()
	require.NoError(t, err)

	res2, err := client.Do(stdReq2)
	require.NoError(t, err)
	require.Equal(t, res2.StatusCode, http.StatusCreated)
	compareReaders(t, res2.Body, strings.NewReader(""))
}

func compareReaders(t *testing.T, a, b io.Reader) {
	b1, err := ioutil.ReadAll(a)
	if err != nil {
		t.FailNow()
	}
	b2, err := ioutil.ReadAll(b)
	if err != nil {
		t.FailNow()
	}
	require.Equal(t, b1, b2)
}