package ecommerce

import (
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/client/metadata"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/vially/aws/util/signer/v2"
	"net/url"
)

// A ServiceName is the name of the service the client will make API calls to.
const ServiceName = "AWSECommerceService"

type ECommerce struct {
	*client.Client
	AssociateTag string
}

// New returns a new ECommerce client.
func New(associateTag string) *ECommerce {
	sess := session.New()
	service := &ECommerce{
		Client: client.New(
			*sess.Config,
			metadata.ClientInfo{
				ServiceName:   ServiceName,
				SigningRegion: "us-east-1",
				Endpoint:      "https://webservices.amazon.com",
				APIVersion:    "2013-08-01",
			},
			sess.Handlers,
		),
		AssociateTag: associateTag,
	}

	// Handlers
	service.Handlers.Sign.PushBack(v2.Sign)
	service.Handlers.ValidateResponse.PushBack(validateItemLookupRequestResponseHandler)
	service.Handlers.Unmarshal.PushBack(unmarshalHandler)
	service.Handlers.UnmarshalError.PushBack(unmarshalItemLookupErrorHandler)

	return service
}

// NewOperationRequest creates a new request for an ECommerce operation
func (e *ECommerce) NewOperationRequest(operation string, params url.Values, data interface{}) *request.Request {
	if params == nil {
		params = url.Values{}
	}

	params.Set("Operation", operation)
	params.Set("Service", e.ServiceName)
	params.Set("Version", e.APIVersion)
	params.Set("AssociateTag", e.AssociateTag)

	op := &request.Operation{
		Name:       operation,
		HTTPMethod: "GET",
		HTTPPath:   "/onca/xml?" + params.Encode(),
	}

	req := e.NewRequest(op, nil, data)

	return req
}
