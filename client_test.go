package acm

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

var ns = "您的命名空间"
var gn = "DEFAULT_GROUP"
var cli = New(
	Namespace(ns),
	GroupName(gn),
	Endpoint("acm.aliyun.com"),
	AuthCreds("您的ACCESS KEY", "您的SECRET KEY"),
)

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

		change, err := cli.Watch(dataId, content)
		expect := strings.Join([]string{dataId, gn, ns}, wordSeparator) + lineSeparator + "\n"
		if err != nil {
			t.Error(err)
		} else if string(change) != expect {
			t.Errorf("unexpected result, expect = %s, but = %s", expect, change)
		}

		finishSignal <- struct{}{}
	}()

	<-startSignal

	if err = cli.Remove(dataId); err != nil {
		t.Error(err)
	}

	<-finishSignal
}
