package hw

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"strings"
	"io/ioutil"
	"io"

	"github.com/stretchr/testify/require"
	"github.com/gorilla/mux"
)

func setup() *httptest.Server {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })        // 200
	r.HandleFunc("/a", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusCreated) })  // 201
	r.HandleFunc("/b/{id}", func(w http.ResponseWriter, r *http.Request) { 
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(mux.Vars(r)["id"]))
	}) // 202
	server := httptest.NewServer(r)
	return server
}

func TestV1(t *testing.T) {
	server := setup()
	defer server.Close()
	client := NewClient(server.URL) // You need to implement NewClient

	t.Run("Challenge 1", func(t2 *testing.T) {
		req := client.Request(http.MethodOptions, "/")
		res, err := client.Do(req)
		require.Equal(t, res.StatusCode, http.StatusOK)
		compareReaders(t, res.Body, strings.NewReader(""))
		require.NoError(t, err)
	})
	
	t.Run("Challenge 2", func(t2 *testing.T) {
		req2 := client.Request(http.MethodOptions, "/a")
		res2, err := client.Do(req2)
		require.Equal(t, res2.StatusCode, http.StatusCreated)
		compareReaders(t, res2.Body, strings.NewReader(""))
		require.NoError(t, err)
	})

	t.Run("Challenge 3", func(t2 *testing.T) {
		req2 := client.Request(http.MethodOptions, "/b/{id}").WithArgs(map[string]string{"id": "xxx"})
		res2, err := client.Do(req2)
		require.Equal(t, res2.StatusCode, http.StatusAccepted)
		compareReaders(t, res2.Body, strings.NewReader("xxx"))
		require.NoError(t, err)
	})
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