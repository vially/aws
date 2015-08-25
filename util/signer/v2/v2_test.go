package v2

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/service"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"
)

func buildSigner(serviceName string, region string, signTime time.Time, query url.Values) signer {
	endpoint := "https://" + serviceName + "." + region + ".amazonaws.com"
	req, _ := http.NewRequest("POST", endpoint, nil)
	req.URL.RawQuery = query.Encode()

	signer := signer{
		Request: req,
		Time:    signTime,
		Credentials: credentials.NewStaticCredentials(
			"AKIAJG6V72ZBDMWPMSWQ",
			"P9MvpCRxpwo2UexwMYbduHEoVcBPJQZO2GVsKNCD",
			""),
	}

	if os.Getenv("DEBUG") != "" {
		signer.Debug = 1
		signer.Logger = aws.NewDefaultLogger()
	}

	return signer
}

func TestSimpleSignRequest(t *testing.T) {
	values := url.Values{}
	values.Add("Action", "CreateDomain")
	values.Add("DomainName", "TestDomain-1437033376")
	values.Add("Version", "2009-04-15")

	timestamp := time.Date(2015, 7, 16, 7, 56, 16, 0, time.UTC)
	signer := buildSigner("sdb", "ap-southeast-2", timestamp, values)

	err := signer.sign()
	query := signer.Request.URL.Query()
	assert.Nil(t, err)
	assert.Equal(t, "u0v86smFkZhRcPjBFfhvoC5EGHTXrZiYBev5xlyW6Lw=", signer.signature)
	assert.Equal(t, 8, len(query))
	assert.Equal(t, "AKIAJG6V72ZBDMWPMSWQ", query.Get("AWSAccessKeyId"))
	assert.Equal(t, "2015-07-16T07:56:16Z", query.Get("Timestamp"))
	assert.Equal(t, "HmacSHA256", query.Get("SignatureMethod"))
	assert.Equal(t, "2", query.Get("SignatureVersion"))
	assert.Equal(t, "u0v86smFkZhRcPjBFfhvoC5EGHTXrZiYBev5xlyW6Lw=", query.Get("Signature"))
	assert.Equal(t, "CreateDomain", query.Get("Action"))
	assert.Equal(t, "TestDomain-1437033376", query.Get("DomainName"))
	assert.Equal(t, "2009-04-15", query.Get("Version"))
}

func TestMoreComplexSignRequest(t *testing.T) {
	query := make(url.Values)
	query.Add("Action", "PutAttributes")
	query.Add("DomainName", "TestDomain-1437041569")
	query.Add("Version", "2009-04-15")
	query.Add("Attribute.2.Name", "Attr2")
	query.Add("Attribute.2.Value", "Value2")
	query.Add("Attribute.2.Replace", "true")
	query.Add("Attribute.1.Name", "Attr1-%\\+ %")
	query.Add("Attribute.1.Value", " \tValue1 +!@#$%^&*(){}[]\"';:?/.>,<\x12\x00")
	query.Add("Attribute.1.Replace", "true")
	query.Add("ItemName", "Item 1")

	timestamp := time.Date(2015, 7, 16, 10, 12, 51, 0, time.UTC)
	signer := buildSigner("sdb", "ap-southeast-2", timestamp, query)

	err := signer.sign()
	assert.Nil(t, err)
	assert.Equal(t, "bXw9mWPcz59G5GFuLM5vnoH/E0cJ9ALb4mhaD0zQUgk=", signer.signature)
}

func TestAnonymousCredentials(t *testing.T) {
	s := service.New(&aws.Config{
		Credentials: credentials.AnonymousCredentials,
		Region:      aws.String("ap-southeast-2"),
	})
	r := s.NewRequest(
		&request.Operation{
			Name:       "PutAttributes",
			HTTPMethod: "POST",
			HTTPPath:   "/",
		},
		nil,
		nil,
	)
	r.Build()

	Sign(r)

	req := r.HTTPRequest
	req.ParseForm()

	assert.Empty(t, req.PostForm.Get("Signature"))
}
