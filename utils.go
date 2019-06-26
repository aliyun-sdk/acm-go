package aliacm

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
)

func contentMd5(content []byte) string {
	hash := md5.New()
	hash.Write(content)
	return hex.EncodeToString(hash.Sum(nil))
}

func genSign(text, key string) (string, error) {
	hash := hmac.New(sha1.New, []byte(key))
	_, err := hash.Write([]byte(text))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(hash.Sum(nil)), nil
}

func httpGet(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Close = true
	req.Header.Set("Connection", "close")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http get error : url=%v , statusCode=%v", url, response.StatusCode)
	}
	return ioutil.ReadAll(response.Body)
}
