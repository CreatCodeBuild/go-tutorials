package main

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"strings"
	"io/ioutil"
	"io"
	"os"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/gorilla/mux"
)


type r struct {
	i int
	next chan struct{}
}

func (r *r) Read(p []byte) (n int, err error) {
	r.i++
	p[0] = []byte("xyz")[r.i % 3]
	<-r.next
	return 1, nil
}

type logger struct {}

func (l *logger) Log(r io.Reader) error {
	_, err := io.Copy(os.Stderr, r)
	if err != nil {
		return err
	}
	io.Copy(os.Stderr, strings.NewReader("\n"))
	return err
}

func Test1(t *testing.T) {
	l := logger{}
	c := make(chan struct{})

	go func() {
		for {
			c <- struct{}{}
			time.Sleep(time.Second)
		}
	}()
	l.Log(&r{next:c})
}

func setup() *httptest.Server {
	l := logger{}
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := l.Log(r.Body) // ?
		if err != nil {
			// ??
		}
		w.WriteHeader(http.StatusOK)
	})        // 200

	r.HandleFunc("/a", func(w http.ResponseWriter, r *http.Request) { 
		w.WriteHeader(http.StatusCreated)
	})  // 201
	
	r.HandleFunc("/b/{id}", func(w http.ResponseWriter, r *http.Request) { 
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(mux.Vars(r)["id"]))
	}) // 202


	server := httptest.NewServer(r)
	return server
}

// func TestV1(t *testing.T) {
// 	server := setup()
// 	defer server.Close()
// 	client := NewClient(server.URL) // You need to implement NewClient

// 	t.Run("Challenge 1", func(t2 *testing.T) {
// 		req := client.Request(http.MethodOptions, "/")
// 		res, err := client.Do(req)
// 		require.Equal(t, res.StatusCode, http.StatusOK)
// 		compareReaders(t, res.Body, strings.NewReader(""))
// 		require.NoError(t, err)
// 	})
	
// 	t.Run("Challenge 2", func(t2 *testing.T) {
// 		req2 := client.Request(http.MethodOptions, "/a")
// 		res2, err := client.Do(req2)
// 		require.Equal(t, res2.StatusCode, http.StatusCreated)
// 		compareReaders(t, res2.Body, strings.NewReader(""))
// 		require.NoError(t, err)
// 	})

// 	t.Run("Challenge 3", func(t2 *testing.T) {
// 		req2 := client.Request(http.MethodOptions, "/b/{id}").WithArgs(map[string]string{"id": "xxx"})
// 		res2, err := client.Do(req2)
// 		require.Equal(t, res2.StatusCode, http.StatusAccepted)
// 		compareReaders(t, res2.Body, strings.NewReader("xxx"))
// 		require.NoError(t, err)
// 	})
// }

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
