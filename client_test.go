package acm

import (
	"bytes"
	"encoding/json"
	"flag"
	"testing"
	"time"
)

var cli *Client
var ns, ak, sk string

func init() {
	flag.StringVar(&ns, "namespace", "", "您的命名空间")
	flag.StringVar(&ak, "access_key", "", "您的ACCESS KEY")
	flag.StringVar(&sk, "secret_key", "", "您的SECRET KEY")
	flag.Parse()
	cli = New(Namespace(ns), Endpoint("acm.aliyun.com"), AuthCreds(ak, sk))
}

var dbCfg = map[string]interface{}{
	"host":   "192.168.1.1",
	"port":   3306,
	"user":   "root",
	"dbname": "test",
}

func TestClient(t *testing.T) {
	dataId := "test.test"
	content, _ := json.Marshal(dbCfg)
	err := cli.Write(dataId, content)
	if err != nil {
		t.Error(err)
	}

	// 等待100毫秒,确保后续能读到
	time.Sleep(time.Millisecond * 100)

	rc, err := cli.Read(dataId)
	if err != nil {
		t.Error(err)
	} else if !bytes.Equal(rc, content) {
		t.Errorf("unexpected result, expect = %s, but = %s", content, rc)
	}

	startSignal := make(chan struct{})
	finishSignal := make(chan struct{})

	go func() {
		startSignal <- struct{}{}

		changed, err := cli.Watch(dataId, content)
		if err != nil {
			t.Error(err)
		} else if changed == false {
			t.Errorf("unexpected result, expect = %v, but = %v", true, false)
		}

		finishSignal <- struct{}{}
	}()

	<-startSignal

	if err = cli.Remove(dataId); err != nil {
		t.Error(err)
	}

	<-finishSignal
}
