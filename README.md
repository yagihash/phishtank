# phishtank
[![CircleCI](https://circleci.com/gh/yagihash/phishtank.svg?style=svg)](https://circleci.com/gh/yagihash/phishtank)
[![codecov](https://codecov.io/gh/yagihash/phishtank/branch/master/graph/badge.svg)](https://codecov.io/gh/yagihash/phishtank)
[![Go Report Card](https://goreportcard.com/badge/github.com/yagihash/phishtank)](https://goreportcard.com/report/github.com/yagihash/phishtank)

```go
package main

import (
	"fmt"

	"github.com/yagihash/phishtank"
)

func main() {
	c := phishtank.New("YOUR API KEY")
	body, _ := c.CheckURL("https://example.com")
	fmt.Printf("%v", body.Results.InDatabase)
}
```
[Playground](https://play.golang.org/p/z6JA3IP6fc9)
