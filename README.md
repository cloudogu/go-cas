# CAS Client library

CAS provides a http package compatible client implementation for use with
securing http frontends in golang.

```
import "gopkg.in/cas.v2"
```

If you are using go modules, get the library by running:

```
go get gopkg.in/cas.v2@v2.4.0
```

## Examples and Documentation

Installing goes like this (have your `go.mod` file ready):
```bash
go get github.com/cloudogu/go-cas/v2
```

Documentation is available at: https://pkg.go.dev/gopkg.in/cas.v2
Examples are included in the documentation but are also available in the
`_examples` directory.

It is possible to override the log-out detection with a custom one by providing a function

```go
package yourpackage

import (
	"net/http"

	cas "github.com/cloudogu/go-cas/v2"
)

func yourCode(next http.Handler) {
	browserClient := cas.NewClient(&cas.Options{
		//...
		IsLogoutRequest: func(r *http.Request) bool {
			return false // inspect request and react here but be sure to re-add any request bodies that you've read
		},
	})
	browserClient.CreateHandler(next)
}
```