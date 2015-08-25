package ecommerce

import (
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/service"
	"github.com/aws/aws-sdk-go/aws/service/serviceinfo"
	"github.com/vially/aws/util/signer/v2"
	"net/url"
)

type ECommerce struct {
	*service.Service
	AssociateTag string
}

// New returns a new ECommerce client.
func New(associateTag string) *ECommerce {
	service := &service.Service{
		ServiceInfo: serviceinfo.ServiceInfo{
			Config:        defaults.DefaultConfig.WithEndpoint("webservices.amazon.com"),
			ServiceName:   "AWSECommerceService",
			APIVersion:    "2013-08-01",
			SigningRegion: "us-east-1",
		},
	}
	service.Initialize()

	// Handlers
	service.Handlers.Sign.PushBack(v2.Sign)
	service.Handlers.ValidateResponse.PushBack(validateItemLookupRequestResponseHandler)
	service.Handlers.Unmarshal.PushBack(unmarshalHandler)
	service.Handlers.UnmarshalError.PushBack(unmarshalItemLookupErrorHandler)

	return &ECommerce{service, associateTag}
}

// NewOperationRequest creates a new request for an ECommerce operation
func (e *ECommerce) NewOperationRequest(operation string, params url.Values, data interface{}) *request.Request {
	if params == nil {
		params = url.Values{}
	}

	params.Set("Operation", operation)
	params.Set("Service", e.ServiceInfo.ServiceName)
	params.Set("Version", e.ServiceInfo.APIVersion)
	params.Set("AssociateTag", e.AssociateTag)

	op := &request.Operation{
		Name:       operation,
		HTTPMethod: "GET",
		HTTPPath:   "/onca/xml?" + params.Encode(),
	}

	req := e.NewRequest(op, nil, data)

	return req
}
