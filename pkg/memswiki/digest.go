package memswiki

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func digestPost(method, host, uri string, additionalHeaders map[string]string, payload io.Reader) ([]byte, error) {
	fmt.Println()
	req, err := http.NewRequest(method, host+uri, nil)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)},
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	log.Printf("getting digest key...\n\x1b[34m%#v\x1b[0m\n", resp)
	if resp.StatusCode != http.StatusUnauthorized {
		return nil, fmt.Errorf("Recieved status code '%v' auth skipped", resp.StatusCode)
	}

	digestParts := getDigestParts(resp)
	digestParts["uri"] = uri
	digestParts["method"] = method
	digestParts["username"] = username
	digestParts["password"] = password
	req, err = http.NewRequest(method, host+uri, payload)
	req.Header.Set("Authorization", getDigestAuthrization(digestParts))
	req.Header.Set("Content-Type", "application/json")
	// log.Printf("digestParts: \n\x1b[33m%s\x1b[0m\n", digestParts)
	// log.Printf("getDigestAuthrization(digestParts): \n\x1b[33m%s\x1b[0m\n", getDigestAuthrization(digestParts))
	for k, v := range additionalHeaders {
		req.Header.Set(k, v)
	}
	resp, err = client.Transport.RoundTrip(req)
	// resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	// log.Printf("req.Header:\n\x1b[33m%s\x1b[0m\n", req.Header)
	// log.Printf("resp.Request.Header:\n\x1b[33m%s\x1b[0m\n", resp.Request.Header)
	if resp.StatusCode != http.StatusOK {
		if err != nil {
			return nil, err
		}
		if resp.StatusCode >= 300 && resp.StatusCode < 400 {
			next, err := resp.Location()
			if err != nil {
				return nil, fmt.Errorf("ErrStatusContradiction")
			}
			log.Printf("Redirect:\n\x1b[33m%#v\x1b[0m\n", resp)

			host := next.Scheme + "://" + next.Host
			uri := next.Path + "?" + next.RawQuery
			log.Printf("next location:\nhost:\n\x1b[33m%s\x1b[0m\nuri:\n\x1b[33m%s\x1b[0m\nurl:\n\x1b[33m%s\x1b[0m\n", host, uri, host+uri)
			redirectRes, err := digestPost(http.MethodGet, host, uri, additionalHeaders, nil)
			if err != nil {
				return redirectRes, fmt.Errorf("failed to request: %w", err)
			}
			return redirectRes, nil
		}
		log.Printf("not OK:\n\x1b[31m%#v\x1b[0m\n", resp)
		log.Printf("not OK string(body):\n\x1b[31m%s\x1b[0m\n", string(body))
		return nil, fmt.Errorf("401 desu")
	}
	log.Printf("Success:\n\x1b[32m%#v\x1b[0m\n", resp)

	return body, nil
}

func getDigestParts(resp *http.Response) map[string]string {
	result := map[string]string{}
	if len(resp.Header["Www-Authenticate"]) > 0 {
		wantedHeaders := []string{"nonce", "realm", "qop"}
		responseHeaders := strings.Split(resp.Header["Www-Authenticate"][0], ",")
		for _, r := range responseHeaders {
			for _, w := range wantedHeaders {
				if strings.Contains(r, w) {
					result[w] = strings.Split(r, `"`)[1]
				}
			}
		}
	}
	return result
}

func getMD5(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func getCnonce() string {
	b := make([]byte, 8)
	io.ReadFull(rand.Reader, b)
	return fmt.Sprintf("%x", b)[:16]
}

func getDigestAuthrization(digestParts map[string]string) string {
	d := digestParts
	ha1 := getMD5(d["username"] + ":" + d["realm"] + ":" + d["password"])
	ha2 := getMD5(d["method"] + ":" + d["uri"])
	nonceCount := 00000001
	cnonce := getCnonce()
	response := getMD5(fmt.Sprintf("%s:%s:%v:%s:%s:%s", ha1, d["nonce"], nonceCount, cnonce, d["qop"], ha2))
	return fmt.Sprintf(
		`Digest username="%s", realm="%s", nonce="%s", uri="%s", cnonce="%s", nc="%v", qop="%s", response="%s"`,
		d["username"],
		d["realm"],
		d["nonce"],
		d["uri"],
		cnonce,
		nonceCount,
		d["qop"],
		response,
	)
}

type ErrDocNotFound struct{}

func (e *ErrDocNotFound) Error() string {
	return fmt.Sprintln("there is no element. Maybe already page is created")
}

func getOriginal(b []byte) (string, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(b))
	if err != nil {
		return "", err
	}

	return doc.Find("textarea[name='original']").Text(), nil
}

func getDigest(b []byte) (string, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(b))
	if err != nil {
		return "", err
	}

	dgst, ok := doc.Find("input[name='digest']").Attr("value")
	if !ok {
		return "", &ErrDocNotFound{}
	}

	return dgst, nil
}
