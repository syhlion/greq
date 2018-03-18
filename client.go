package greq

import (
	"context"
	"encoding/base64"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	requestwork "github.com/syhlion/requestwork.v2"
	httpstat "github.com/tcnksm/go-httpstat"
)

var traceTemplete = "[greq] Rrequest url:{{.Url}}\t Method:{{.Method}}\t Param:{{.Param}}\n" +
	"Response Body:{{.Body}}\n" +
	"{{.Time}}"

type Trace struct {
	Url    string
	Method string
	Body   string
	Param  string
	Time   string
}

//New return http client
func New(worker *requestwork.Worker, timeout time.Duration, debug bool) *Client {
	return &Client{
		worker:  worker,
		timeout: timeout,
		headers: make(map[string]string),
		lock:    &sync.RWMutex{},
		debug:   debug,
	}
}

//Client instance
type Client struct {
	worker  *requestwork.Worker
	timeout time.Duration
	headers map[string]string
	host    string
	lock    *sync.RWMutex
	debug   bool
}

//SetBasicAuth  set Basic auth
func (c *Client) SetBasicAuth(username, password string) *Client {
	auth := username + ":" + password
	hash := base64.StdEncoding.EncodeToString([]byte(auth))
	c.lock.Lock()
	defer c.lock.Unlock()
	c.headers["Authorization"] = "Basic " + hash
	return c
}

//SetHeader set http header
func (c *Client) SetHeader(key, value string) *Client {
	key = strings.Title(key)
	c.lock.Lock()
	defer c.lock.Unlock()
	c.headers[key] = value
	return c
}

//SetHost set host
func (c *Client) SetHost(host string) *Client {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.host = host
	return c
}

//Get http method get
func (c *Client) Get(url string, params url.Values) (data []byte, httpstatus int, err error) {
	if params != nil {
		url += "?" + params.Encode()
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	return c.resolveRequest(req, params, err)

}

//Post http method post
func (c *Client) Post(url string, params url.Values) (data []byte, httpstatus int, err error) {
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(params.Encode()))
	return c.resolveRequest(req, params, err)
}

//Put http method put
func (c *Client) Put(url string, params url.Values) (data []byte, httpstatus int, err error) {
	req, err := http.NewRequest(http.MethodPut, url, strings.NewReader(params.Encode()))
	return c.resolveRequest(req, params, err)
}

//Delete http method Delete
func (c *Client) Delete(url string, params url.Values) (data []byte, httpstatus int, err error) {
	req, err := http.NewRequest(http.MethodDelete, url, strings.NewReader(params.Encode()))
	return c.resolveRequest(req, params, err)
}

func (c *Client) resolveHeaders(req *http.Request) {
	c.lock.RLock()
	c.lock.RUnlock()
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}
	if c.host != "" {
		req.Host = c.host
	}
}

func (c *Client) resolveRequest(req *http.Request, params url.Values, e error) (data []byte, httpstatus int, err error) {
	var (
		body    []byte
		status  int
		endTime time.Time
		result  httpstat.Result
	)
	if c.debug {
		var stat Trace
		defer func() {
			stat.Param = params.Encode()
			stat.Url = req.URL.String()
			stat.Method = req.Method
			stat.Body = string(data)
			stat.Time = fmt.Sprintf("%+v", result)
			t := template.Must(template.New("trace templete").Parse(traceTemplete))
			t.Execute(os.Stdout, stat)
		}()
		sctx := httpstat.WithHTTPStat(req.Context(), &result)
		req = req.WithContext(sctx)
	}
	if e != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)

	defer cancel()
	c.resolveHeaders(req)

	switch req.Method {
	case "PUT", "POST", "DELETE":
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	}

	err = c.worker.Execute(ctx, req, func(resp *http.Response, err error) (er error) {
		if c.debug {
			defer func() {
				endTime = time.Now()
				result.End(endTime)
			}()
		}
		if err != nil {
			return err
		}
		var readErr error
		defer func() {
			resp.Body.Close()
		}()
		status = resp.StatusCode
		body, readErr = ioutil.ReadAll(resp.Body)
		if readErr != nil {
			return readErr
		}
		return
	})
	if err != nil {
		return
	}
	data = body
	httpstatus = status
	return

}
