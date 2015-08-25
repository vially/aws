package ecommerce

// The returned result of the corresponding request.
type ItemLookupResponse struct {
	Items []*Item `xml:"Items>Item"`
}

type Item struct {
	ASIN           string
	ParentASIN     string
	ItemAttributes *ItemAttributes
	OfferSummary   *OfferSummary
	Offers         []*Offer	`xml:"Offers>Offer"`
}

type ItemAttributes struct {
	Manufacturer string
	ProductGroup string
	Title        string
}

type OfferSummary struct {
	LowestNewPrice  int64 `xml:"LowestNewPrice>Amount"`
	LowestUsedPrice int64 `xml:"LowestUsedPrice>Amount"`
}

type Offer struct {
	Condition   string `xml:"OfferAttributes>Condition"`
	Price       int64  `xml:"OfferListing>Price>Amount"`
	AmountSaved int64  `xml:"OfferListing>AmountSaved>Amount"`
}
