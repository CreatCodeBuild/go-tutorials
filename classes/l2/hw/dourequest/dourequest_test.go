package dourequest_test

import (
	"github.com/CreatCodeBuild/go-tutorials/classes/l2/hw/dourequest"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func setup() *httptest.Server {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })       // 200
	r.HandleFunc("/a", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusCreated) }) // 201
	r.HandleFunc("/b/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(mux.Vars(r)["id"]))
	}) // 202
	r.HandleFunc("/timeout", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second * time.Duration(1))
		w.WriteHeader(http.StatusAccepted)
	})

	r.HandleFunc("/query", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Duration(1) * time.Second)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(r.URL.Query().Get("test")))

	}).Methods("GET")

	r.HandleFunc("/body", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			w.Write([]byte(err.Error()))
		}
		w.Write(body)

	}).Methods("POST")
	server := httptest.NewServer(r)
	return server
}

func TestV1(t *testing.T) {
	server := setup()
	defer server.Close()

	dourequest.BaseURL = server.URL

	t.Run("Challenge 1", func(t2 *testing.T) {
		res, err := dourequest.NewRequest("/").
			Method(http.MethodOptions).Do()
		require.Equal(t, res.StatusCode, http.StatusOK)
		compareReaders(t, res.Body, strings.NewReader(""))
		require.NoError(t, err)
	})

	t.Run("Challenge 2", func(t2 *testing.T) {
		res2, err := dourequest.NewRequest("/a").
			Method(http.MethodOptions).Do()
		require.Equal(t, res2.StatusCode, http.StatusCreated)
		compareReaders(t, res2.Body, strings.NewReader(""))
		require.NoError(t, err)
	})

	t.Run("Challenge 3", func(t2 *testing.T) {
		res2, err := dourequest.NewRequest("/b/{id}").
			Method(http.MethodOptions).
			SetArgs(map[string]string{"id": "123"}).
			Arg("id", "xxx").
			Do()
		require.Equal(t, res2.StatusCode, http.StatusAccepted)
		compareReaders(t, res2.Body, strings.NewReader("xxx"))
		require.NoError(t, err)
	})

	t.Run("Challenge Timeout", func(t2 *testing.T) {
		res2, err := dourequest.NewRequest("/timeout").
			Method(http.MethodOptions).
			Timeout(500).
			RetryTimes(2).
			Do()
		require.Equal(t, res2.StatusCode, http.StatusAccepted)
		require.NoError(t, err)
	})

	t.Run("Challenge Query", func(t2 *testing.T) {
		res2, err := dourequest.Get("/query").
			Query(url.Values{"test": []string{"ok"}}).
			Timeout(500).
			RetryTimes(2).
			Do()
		require.Equal(t, res2.StatusCode, http.StatusOK)
		compareReaders(t, res2.Body, strings.NewReader("ok"))
		require.NoError(t, err)
	})

	t.Run("Challenge Body", func(t2 *testing.T) {
		res2, err := dourequest.NewRequest("/body").
			Method(http.MethodPost).
			Body([]byte("ok")).
			Do()
		require.Equal(t, res2.StatusCode, http.StatusOK)
		compareReaders(t, res2.Body, strings.NewReader("ok"))
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
