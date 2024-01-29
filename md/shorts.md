---
runme:
  id: 01HJX1AHQCRQDVZW3535A91KH0
  version: v2.0
---

# Shorts

## Create psuedo-random numbers

```go {"id":"01HJX1B0J6PYEB2N554X773WN4"}
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	generator := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 10; i += 1 {
		fmt.Printf("%d ", generator.Intn(11))
	}
}
```