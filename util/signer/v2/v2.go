package v2

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
)

const (
	signatureVersion = "2"
	signatureMethod  = "HmacSHA256"
	timeFormat       = "2006-01-02T15:04:05Z"
)

type signer struct {
	// Values that must be populated from the request
	Request      *http.Request
	Time         time.Time
	Credentials  *credentials.Credentials
	Debug        aws.LogLevelType
	Logger       aws.Logger
	stringToSign string
	signature    string
}

// Sign requests with signature version 2.
//
// Will sign the requests with the service config's Credentials object
// Signing is skipped if the credentials is the credentials.AnonymousCredentials
// object.
func Sign(req *request.Request) {
	// If the request does not need to be signed ignore the signing of the
	// request if the AnonymousCredentials object is used.
	if req.Service.Config.Credentials == credentials.AnonymousCredentials {
		return
	}

	v2 := signer{
		Request:     req.HTTPRequest,
		Time:        req.Time,
		Credentials: req.Service.Config.Credentials,
		Debug:       req.Service.Config.LogLevel.Value(),
		Logger:      req.Service.Config.Logger,
	}

	req.Error = v2.sign()
}

func (v2 *signer) sign() error {
	credValue, err := v2.Credentials.Get()
	if err != nil {
		return err
	}

	query := v2.Request.URL.Query()

	// Set new query parameters
	query.Set("AWSAccessKeyId", credValue.AccessKeyID)
	query.Set("SignatureVersion", signatureVersion)
	query.Set("SignatureMethod", signatureMethod)
	query.Set("Timestamp", v2.Time.UTC().Format(timeFormat))

	// in case this is a retry, ensure no signature present
	query.Del("Signature")

	method := v2.Request.Method
	host := v2.Request.URL.Host
	path := v2.Request.URL.Path
	if path == "" {
		path = "/"
	}

	// obtain all of the query keys and sort them
	queryKeys := make([]string, 0, len(query))
	for key := range query {
		queryKeys = append(queryKeys, key)
	}
	sort.Strings(queryKeys)

	// build URL-encoded query keys and values
	queryKeysAndValues := make([]string, len(queryKeys))
	for i, key := range queryKeys {
		k := strings.Replace(url.QueryEscape(key), "+", "%20", -1)
		v := strings.Replace(url.QueryEscape(query.Get(key)), "+", "%20", -1)
		queryKeysAndValues[i] = k + "=" + v
	}

	// build the canonical string for the V2 signature
	v2.stringToSign = strings.Join([]string{
		method,
		host,
		path,
		strings.Join(queryKeysAndValues, "&"),
	}, "\n")

	hash := hmac.New(sha256.New, []byte(credValue.SecretAccessKey))
	hash.Write([]byte(v2.stringToSign))
	v2.signature = base64.StdEncoding.EncodeToString(hash.Sum(nil))
	query.Set("Signature", v2.signature)
	v2.Request.URL.RawQuery = query.Encode()

	if v2.Debug.AtLeast(aws.LogDebug) {
		v2.logSigningInfo()
	}

	return nil
}

func (v2 *signer) logSigningInfo() {
	out := v2.Logger
	out.Log("---[ STRING TO SIGN ]--------------------------------\n")
	out.Log(v2.stringToSign)
	out.Log("---[ SIGNATURE ]-------------------------------------\n")
	out.Log(v2.signature)
	out.Log("-----------------------------------------------------\n")
}
