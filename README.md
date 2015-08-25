# aws
ECommerce API for aws-sdk-go

## Usage

```go
package main

import (
	"github.com/vially/aws/service/ecommerce"
	"log"
	"fmt"
	"os"
)

func main() {
	os.Setenv("AWS_ACCESS_KEY_ID", "MY_ACCESS_KEY")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "MY_SECRET_ACCESS_KEY")

	ec := ecommerce.New("my-associate-tag")

	responseGroups := []string{
		ecommerce.ResponseGroupSmall,
		ecommerce.ResponseGroupOffers,
	}

	lookupResponse, err := ec.ItemLookup([]string{"B00HCNH90W"}, responseGroups)
	if err != nil {
		log.Fatalln(err)
	}

	for _, item := range lookupResponse.Items {
		fmt.Println(item.ASIN, item.ItemAttributes.Title)
	}
}
```
