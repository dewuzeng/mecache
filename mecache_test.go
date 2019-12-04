package mecache

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var c *MeCache

func setup() {
	assertCacheImplementation()
	c = New(2*time.Second, 4*time.Second)
}

func teardown() {

}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func BenchmarkMeCache_Set(b *testing.B) {
	key := "bKey"
	val := "bVal"
	for n := 0; n < b.N; n++ {
		key = fmt.Sprintf("%s%d", key, n)
		val = fmt.Sprintf("%s%d", val, n)
		c.SetDefault(key, val)
	}
}

func TestMeCache_Set(t *testing.T) {
	ass := assert.New(t)
	c.Set("key", "value", DefaultExpiration)
	val, found := c.Get("key")
	ass.True(found)
	ass.Equal(val, "value")
	time.Sleep(5 * time.Second)
	val, found = c.Get("key")
	ass.False(found)
	ass.Nil(val)
}

func TestMeCache_Get(t *testing.T) {
	ass := assert.New(t)
	c.SetDefault("keyDefault", "valueDefault")
	val, found := c.Get("keyDefault")
	ass.True(found)
	ass.Equal(val, "valueDefault")
	time.Sleep(2 * time.Second)
	c.DeleteExpired()
	val, found = c.Get("keyDefault")
	ass.False(found)
	ass.Nil(val)
}
