package acm

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	wordSeparator = string([]byte{37, 48, 50})
	lineSeparator = string([]byte{37, 48, 49})
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Client ACM客户端
type Client struct {
	options Options
	servers []server
}

func (a *Client) Read(dataId string) ([]byte, error) {
	url := a.buildUrl(getConfig) +
		"?tenant=" + a.options.namespace +
		"&dataId=" + dataId +
		"&group=" + a.options.groupName
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return a.doRequest(req)
}

func (a *Client) Write(dataId string, content []byte) error {
	arg := fmt.Sprintf(
		"tenant=%s&dataId=%s&group=%s&content=%s",
		a.options.namespace, dataId,
		a.options.groupName, content,
	)
	url := a.buildUrl(setConfig)
	req, err := http.NewRequest("POST", url, strings.NewReader(arg))
	if err != nil {
		return err
	}
	res, err := a.doRequest(req)
	if err != nil {
		return err
	} else if string(res[:2]) != "OK" {
		return fmt.Errorf("acm client: write config failed, response = %v", res)
	}
	return nil
}

func (a *Client) Remove(dataId string) error {
	arg := fmt.Sprintf(
		"tenant=%s&dataId=%s&group=%s",
		a.options.namespace, dataId, a.options.groupName,
	)
	url := a.buildUrl(delConfig)
	req, err := http.NewRequest("POST", url, strings.NewReader(arg))
	if err != nil {
		return err
	}
	res, err := a.doRequest(req)
	if err != nil {
		return err
	} else if string(res[:2]) != "OK" {
		return fmt.Errorf("acm client: remove config failed, response = %s", res)
	}
	return nil
}

func (a *Client) Watch(dataId string, content []byte) ([]byte, error) {
	url := a.buildUrl(getConfig)
	arg := strings.Join([]string{dataId, a.options.groupName, contentMd5(content), a.options.namespace}, wordSeparator)
	req, err := http.NewRequest("POST", url, strings.NewReader("Probe-Modify-Request="+arg+lineSeparator))
	if err != nil {
		return nil, err
	}
	req.Header.Set("longPullingTimeout", "30000")
	return a.doRequest(req)
}

func (a *Client) buildUrl(path apiPath) string {
	idx := rand.Intn(len(a.servers))
	return path.URL(a.servers[idx].host, a.servers[idx].port)
}

func (a *Client) doRequest(r *http.Request) ([]byte, error) {
	err := a.addHeader(r)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("acm client: response error, code = %d, message = %s", resp.StatusCode, body)
	}
	return body, nil
}

func (a *Client) addHeader(r *http.Request) error {
	mtime := strconv.FormatInt(time.Now().Unix()*1000, 10)
	sign, err := genSign(
		fmt.Sprintf("%s+%s+%s", a.options.namespace, a.options.groupName, mtime),
		a.options.authCreds.secretKey,
	)
	if err != nil {
		return err
	}
	r.Header.Set("timeStamp", mtime)
	r.Header.Set("Spas-Signature", sign)
	r.Header.Set("Spas-AccessKey", a.options.authCreds.accessKey)
	if strings.ToUpper(r.Method) == "POST" {
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	return nil
}

func New(fns ...Option) *Client {
	a := &Client{
		options: newOptions(fns),
	}
	// 通过Endpoint查询服务IP列表, 以便后面能够通过IP发起请求
	url := getServer.URL(a.options.endpoint, 8080)
	if body, err := httpGet(url); err != nil {
		panic(fmt.Sprintf("acm client: get server list failed, err = %v", err))
	} else {
		servers := strings.Split(string(body), "\n")
		for _, server := range servers {
			server = strings.TrimSpace(server)
			if server == "" {
				continue
			}
			hps := strings.Split(server, ":")
			if len(hps) == 1 {
				a.servers = append(a.servers, newServer(hps[0], 8080))
			} else {
				port, _ := strconv.Atoi(hps[1])
				a.servers = append(a.servers, newServer(hps[0], uint16(port)))
			}
		}
	}
	return a
}
