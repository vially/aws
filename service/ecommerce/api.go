package ecommerce

import (
	"net/url"
	"strings"
)

func (e *ECommerce) ItemLookup(itemIds, responseGroups []string) (*ItemLookupResponse, error) {
	params := url.Values{}
	params.Set("ItemId", strings.Join(itemIds, ","))
	if responseGroups != nil {
		params.Set("ResponseGroup", strings.Join(responseGroups, ","))
	}

	out := &ItemLookupResponse{}
	err := e.NewOperationRequest("ItemLookup", params, out).Send()

	return out, err
}
