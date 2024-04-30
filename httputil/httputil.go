package httputil

import (
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

var (
	_client    *http.Client
	_cookieJar *cookiejar.Jar
)

func InitClient(proxys string) error {
	_cookieJar, _ = cookiejar.New(nil)
	if proxys != "" {
		proxyURL, err := url.Parse(proxys)
		if err != nil {
			return err
		}
		_client = &http.Client{
			Jar: _cookieJar,
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			},
		}
	} else {
		_client = &http.Client{
			Jar: _cookieJar,
		}
	}
	return nil
}

func HttpPost(urls string, data map[string]string) (string, error) {
	if _client == nil {
		return "", nil
	}
	formData := url.Values{}
	for key, value := range data {
		formData.Set(key, value)
	}
	resp, err := _client.PostForm(urls, formData)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func HttpGet(url string) (string, error) {
	if _client == nil {
		return "", nil
	}
	resp, err := _client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
