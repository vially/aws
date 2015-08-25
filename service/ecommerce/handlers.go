package ecommerce

import (
	"bytes"
	"encoding/xml"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"io"
	"io/ioutil"
)

type xmlItemLookupResponse struct {
	XMLName   xml.Name    `xml:"ItemLookupResponse"`
	RequestID string      `xml:"OperationRequest>RequestId"`
	Valid     string      `xml:"Items>Request>IsValid"`
	Errors    []*xmlError `xml:"Items>Request>Errors>Error"`
}

type xmlError struct {
	Code    string
	Message string
}

func drainBody(b io.ReadCloser) (out *bytes.Buffer, err error) {
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(b); err != nil {
		return nil, err
	}
	if err = b.Close(); err != nil {
		return nil, err
	}
	return &buf, nil
}

var validateItemLookupRequestResponseHandler = func(r *request.Request) {
	if r.HTTPResponse.StatusCode == 0 || r.HTTPResponse.StatusCode >= 300 {
		// this may be replaced by an UnmarshalError handler
		r.Error = awserr.New("UnknownError", "unknown error", nil)
	} else if r.HTTPResponse.StatusCode == 200 {
		buf, err := drainBody(r.HTTPResponse.Body)
		if err != nil { // failed to read the response body, skip
			r.Error = awserr.New("IOError", "failed to read response body while validating response", err)
			return
		}

		// Reset body for subsequent reads
		r.HTTPResponse.Body = ioutil.NopCloser(bytes.NewReader(buf.Bytes()))

		resp := &xmlItemLookupResponse{}
		err = xml.NewDecoder(buf).Decode(resp)
		if err != nil && err != io.EOF {
			r.Error = awserr.New("SerializationError", "failed to decode item lookup XML response", err)
		} else if resp.Valid == "False" {
			r.Error = awserr.New("RequestValidationError", "invalid item lookup request", nil)
		}
	}
}

func unmarshalItemLookupErrorHandler(r *request.Request) {
	defer r.HTTPResponse.Body.Close()

	resp := &xmlItemLookupResponse{}
	err := xml.NewDecoder(r.HTTPResponse.Body).Decode(resp)
	if err != nil && err != io.EOF {
		r.Error = awserr.New("SerializationError", "failed to decode item lookup XML error response", err)
	} else {
		if len(resp.Errors) == 0 {
			r.Error = awserr.New("UnknownError", "unknown item lookup request error", nil)
		} else {
			r.Error = awserr.NewRequestFailure(
				awserr.New(resp.Errors[0].Code, resp.Errors[0].Message, nil),
				r.HTTPResponse.StatusCode,
				resp.RequestID,
			)
		}
	}
}

func unmarshalHandler(r *request.Request) {
	defer r.HTTPResponse.Body.Close()
	err := xml.NewDecoder(r.HTTPResponse.Body).Decode(r.Data)
	if err != nil && err != io.EOF {
		r.Error = awserr.New("SerializationError", "failed to decode item lookup REST XML response", err)
	}
}
