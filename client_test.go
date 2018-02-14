package greq

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	requestwork "github.com/syhlion/requestwork.v2"
)

var worker *requestwork.Worker

func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("method error"))
		return
	}
	a := r.FormValue("key")

	if a != "TEST_HELLO" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("param error ,request:" + a))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))
	return

}
func getHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("method error"))
		return
	}

	a := r.FormValue("key")
	if a != "TEST_HELLO" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("param error ,request:" + a))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))
	return

}
func putHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("method error"))
		return
	}
	a := r.FormValue("key")
	if a != "TEST_HELLO" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("param error ,request:" + a))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))
	return

}
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("method error"))
		return
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("read body error: " + err.Error()))
		return
	}
	v, err := url.ParseQuery(string(b))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("parser query error: " + err.Error()))
		return
	}
	a := v.Get("key")
	if a != "TEST_HELLO" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("param error ,request:" + a))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))
	return

}

func basicAuthHandler(w http.ResponseWriter, r *http.Request) {
	basicAuthPrefix := "Basic "
	auth := r.Header.Get("Authorization")
	if strings.HasPrefix(auth, basicAuthPrefix) {
		payload, err := base64.StdEncoding.DecodeString(auth[len(basicAuthPrefix):])
		if err == nil {
			pair := bytes.SplitN(payload, []byte(":"), 2)
			if len(pair) == 2 && bytes.Equal(pair[0], []byte("scott")) && bytes.Equal(pair[1], []byte("fine")) {
				w.Write([]byte("success"))
				return
			}
		}
	}
	w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
	w.WriteHeader(http.StatusUnauthorized)
	return
}

func TestAddBasicAuth(t *testing.T) {
	worker = requestwork.New(10)
	ts := httptest.NewServer(http.HandlerFunc(basicAuthHandler))
	defer ts.Close()
	client := New(worker, 15*time.Second)
	client.SetBasicAuth("scott", "fine")
	data, s, err := client.Get(ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	if s != http.StatusOK {
		t.Fatalf("status fatal:%d ,body:%s", s, string(data))
	}
	if string(data) != "success" {
		t.Fatal("body fatal :", string(data))
	}

	//test err response
	client.SetBasicAuth("scott", "nofine")
	data, s, err = client.Get(ts.URL, nil)
	if s == http.StatusOK {
		t.Fatalf("status fatal:%d ,body:%s", s, string(data))
	}
	if string(data) == "success" {
		t.Fatal("body fatal :", string(data))
	}
}

func TestGet(t *testing.T) {
	worker = requestwork.New(10)
	ts := httptest.NewServer(http.HandlerFunc(getHandler))
	defer ts.Close()

	client := New(worker, 15*time.Second)
	v := url.Values{}
	v.Set("key", "TEST_HELLO")
	data, s, err := client.Get(ts.URL, v)
	if err != nil {
		t.Fatal(err)
	}
	if s != http.StatusOK {
		t.Fatalf("status fatal:%d ,body:%s", s, string(data))
	}
	if string(data) != "success" {
		t.Fatal("body fatal :", string(data))
	}

}
func TestPost(t *testing.T) {
	worker = requestwork.New(10)
	ts := httptest.NewServer(http.HandlerFunc(postHandler))
	defer ts.Close()
	client := New(worker, 15*time.Second)
	v := url.Values{}
	v.Set("key", "TEST_HELLO")
	data, s, err := client.Post(ts.URL, v)
	if err != nil {
		t.Fatal(err)
	}
	if s != http.StatusOK {
		t.Fatalf("status fatal:%d ,body:%s", s, string(data))
	}
	if string(data) != "success" {
		t.Fatal("body fatal :", string(data))
	}
}
func TestPut(t *testing.T) {
	worker = requestwork.New(10)
	ts := httptest.NewServer(http.HandlerFunc(putHandler))
	defer ts.Close()
	client := New(worker, 15*time.Second)
	v := url.Values{}
	v.Set("key", "TEST_HELLO")
	data, s, err := client.Put(ts.URL, v)
	if err != nil {
		t.Fatal(err)
	}
	if s != http.StatusOK {
		t.Fatalf("status fatal:%d ,body:%s", s, string(data))
	}
	if string(data) != "success" {
		t.Fatal("body fatal :", string(data))
	}
}
func TestDelete(t *testing.T) {
	worker = requestwork.New(10)
	ts := httptest.NewServer(http.HandlerFunc(deleteHandler))
	defer ts.Close()
	client := New(worker, 15*time.Second)
	v := url.Values{}
	v.Set("key", "TEST_HELLO")
	data, s, err := client.Delete(ts.URL, v)
	if err != nil {
		t.Fatal(err)
	}
	if s != http.StatusOK {
		t.Fatalf("status fatal:%d ,body:%s", s, string(data))
	}
	if string(data) != "success" {
		t.Fatal("body fatal :", string(data))
	}
}
