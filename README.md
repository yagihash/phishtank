# phishtank
[![CircleCI](https://circleci.com/gh/yagihashoo/phishtank.svg?style=svg)](https://circleci.com/gh/yagihashoo/phishtank) [![codecov](https://codecov.io/gh/yagihashoo/phishtank/branch/master/graph/badge.svg)](https://codecov.io/gh/yagihashoo/phishtank)

```go
package main

import (
	"fmt"

	"github.com/yagihashoo/phishtank"
)

func main() {
	c := phishtank.New("YOUR API KEY")
	body, _ := c.CheckURL("https://example.com")
	fmt.Printf("%v", body.Results.InDatabase)
}
```
[Playground](https://play.golang.org/p/MzcsDqXqqMF)
