package aws

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

//golang = simple is beautiful
//golang the best language

const (
	timeFormatS3   = time.RFC1123Z
	subresourcesS3 = "acl,lifecycle,location,logging,notification,partNumber,policy,requestPayment,torrent,uploadId,uploads,versionId,versioning,versions,website"
)

type Credentials struct {
	AccessKeyID     string
	SecretAccessKey string
	SecurityToken   string `json:"Token"`
	Expiration      time.Time
}

func SignS3(request *http.Request, credentials Credentials) *http.Request {
	keys := credentials

	// Add the X-Amz-Security-Token header when using STS
	if keys.SecurityToken != "" {
		request.Header.Set("X-Amz-Security-Token", keys.SecurityToken)
	}

	prepareRequestS3(request)

	stringToSign := stringToSignS3(request)
	signature := signatureS3(stringToSign, keys)

	authHeader := "AWS " + keys.AccessKeyID + ":" + signature
	request.Header.Set("Authorization", authHeader)

	log.Println(request.Host)

	return request
}

func signatureS3(stringToSign string, keys Credentials) string {
	hashed := hmacSHA1([]byte(keys.SecretAccessKey), stringToSign)
	return base64.StdEncoding.EncodeToString(hashed)
}

func stringToSignS3(request *http.Request) string {
	str := request.Method + "\n"

	if request.Header.Get("Content-Md5") != "" {
		str += request.Header.Get("Content-Md5")
	} else {
		body := readAndReplaceBody(request)
		if len(body) > 0 {
			hashstr := hashMD5(body)
			str += hashstr
			request.Header.Add("Content-Md5", hashstr)
		}
	}
	str += "\n"

	str += request.Header.Get("Content-Type") + "\n"

	if request.Header.Get("Date") != "" {
		str += request.Header.Get("Date")
	} else {
		str += timestampS3()
	}

	str += "\n"

	canonicalHeaders := canonicalAmzHeadersS3(request)
	if canonicalHeaders != "" {
		str += canonicalHeaders
	}

	str += canonicalResourceS3(request)

	return str
}

func stringToSignS3Url(method string, expire time.Time, path string) string {
	return method + "\n\n\n" + timeToUnixEpochString(expire) + "\n" + path
}

func timeToUnixEpochString(t time.Time) string {
	return strconv.FormatInt(t.Unix(), 10)
}

func canonicalAmzHeadersS3(request *http.Request) string {
	var headers []string

	for header := range request.Header {
		standardized := strings.ToLower(strings.TrimSpace(header))
		if strings.HasPrefix(standardized, "x-amz-") {
			headers = append(headers, standardized)
		}
	}

	sort.Strings(headers)

	for i, header := range headers {
		headers[i] = header + ":" + strings.Replace(request.Header.Get(header), "\n", " ", -1)
	}

	if len(headers) > 0 {
		return strings.Join(headers, "\n") + "\n"
	} else {
		return ""
	}
}

func canonicalResourceS3(request *http.Request) string {
	res := ""

	if isS3VirtualHostedStyle(request) {
		bucketname := strings.Split(request.Host, ".")[0]
		res += "/" + bucketname
	}

	res += request.URL.Path

	for _, subres := range strings.Split(subresourcesS3, ",") {
		if strings.HasPrefix(request.URL.RawQuery, subres) {
			res += "?" + subres
		}
	}

	return res
}

func prepareRequestS3(request *http.Request) *http.Request {
	request.Header.Set("Date", timestampS3())
	if request.URL.Path == "" {
		request.URL.Path += "/"
	}
	return request
}

// Info: http://docs.aws.amazon.com/AmazonS3/latest/dev/VirtualHosting.html
func isS3VirtualHostedStyle(request *http.Request) bool {
	service, _ := serviceAndRegion(request.Host)
	return service == "s3" && strings.Count(request.Host, ".") == 3
}

func timestampS3() string {
	return time.Now().UTC().Format(timeFormatS3)
}

// serviceAndRegion parsers a hostname to find out which ones it is.
// http://docs.aws.amazon.com/general/latest/gr/rande.html
func serviceAndRegion(host string) (service string, region string) {
	// These are the defaults if the hostname doesn't suggest something else
	region = "us-east-1"
	service = "s3"

	parts := strings.Split(host, ".")
	if len(parts) == 4 {
		// Either service.region.amazonaws.com or virtual-host.region.amazonaws.com
		if parts[1] == "s3" {
			service = "s3"
		} else if strings.HasPrefix(parts[1], "s3-") {
			region = parts[1][3:]
			service = "s3"
		} else {
			service = parts[0]
			region = parts[1]
		}
	} else if len(parts) == 5 {
		service = parts[2]
		region = parts[1]
	} else {
		// Either service.amazonaws.com or s3-region.amazonaws.com
		if strings.HasPrefix(parts[0], "s3-") {
			region = parts[0][3:]
		} else {
			service = parts[0]
		}
	}

	if region == "external-1" {
		region = "us-east-1"
	}

	return
}

func hmacSHA256(key []byte, content string) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(content))
	return mac.Sum(nil)
}

func hmacSHA1(key []byte, content string) []byte {
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(content))
	return mac.Sum(nil)
}

func hashMD5(content []byte) string {
	h := md5.New()
	h.Write(content)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func readAndReplaceBody(request *http.Request) []byte {
	if request.Body == nil {
		return []byte{}
	}
	payload, _ := ioutil.ReadAll(request.Body)
	request.Body = ioutil.NopCloser(bytes.NewReader(payload))
	return payload
}
