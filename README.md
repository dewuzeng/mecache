# mecache

MeCache is an in-memory key:value store/cache.


```go
package main

import (
	"fmt"
	"github.com/somecodeio/mecache"
	"time"
)

func main() {

	c := mecache.New(5*time.Minute, 10*time.Minute)

	c.Set("foo", "bar", mecache.DefaultExpiration)
	if val, ok := c.Get("foo"); ok {
		fmt.Println(val)
	}

	c.SetDefault("num", 9527)
	if val, ok := c.Get("num"); ok {
		fmt.Println(val)
	}
}

```