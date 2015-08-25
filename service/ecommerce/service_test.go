package ecommerce

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestServiceConfig(t *testing.T) {
	svc := New("example-tag")
	assert.Equal(t, svc.APIVersion, "2013-08-01")
	assert.Equal(t, svc.ServiceName, "AWSECommerceService")
	assert.Equal(t, svc.Endpoint, "https://webservices.amazon.com")
	assert.Equal(t, svc.associateTag, "example-tag")
	assert.NotEmpty(t, svc.SigningRegion)
}

func TestServiceRequestParams(t *testing.T) {
	svc := New("example-tag")
	req := svc.NewOperationRequest("DoSomething", nil, nil)
	query := req.HTTPRequest.URL.Query()
	assert.Equal(t, query.Get("Operation"), "DoSomething")
	assert.Equal(t, query.Get("AssociateTag"), "example-tag")
}
