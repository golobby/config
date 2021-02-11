package feeder

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOS_Feed_It_Should_Return_Non_Empty_OS_Variables(t *testing.T) {
	_ = os.Setenv("APP_NAME", "Config")
	_ = os.Setenv("APP_URL", "https://github.com/golobby/config")
	_ = os.Setenv("APP_VERSION", "2.0")

	e := OS{Keys: []string{"APP_NAME", "APP_URL", "APP_VERSION", "APP_EMPTY"}}

	items, err := e.Feed()
	assert.NoError(t, err)
	assert.Equal(t, "Config", items["app.name"])
	assert.Equal(t, "https://github.com/golobby/config", items["app.url"])
	assert.Equal(t, "2.0", items["app.version"])
	assert.Equal(t, 3, len(items))
}
