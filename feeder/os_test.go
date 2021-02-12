package feeder

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOS_Feed(t *testing.T) {
	_ = os.Setenv("APP_NAME", "Config")
	_ = os.Setenv("APP_URL", "https://github.com/golobby/config")
	_ = os.Setenv("APP_VERSION", "2.0")
	_ = os.Setenv("APP_NONE", "")

	e := OS{Variables: []string{"APP_NAME", "APP_URL", "APP_VERSION", "APP_EMPTY", "APP_NONE"}}

	items, err := e.Feed()
	assert.NoError(t, err)
	assert.Equal(t, "Config", items["app.name"])
	assert.Equal(t, "https://github.com/golobby/config", items["app.url"])
	assert.Equal(t, "2.0", items["app.version"])
	assert.Equal(t, 3, len(items))
}

func TestOS_Feed_With_Strict_Mode(t *testing.T) {
	_ = os.Setenv("APP_NAME", "Config")

	e := OS{Variables: []string{"APP_NAME", "APP_NONE"}, Strict: true}

	items, err := e.Feed()
	assert.NoError(t, err)
	assert.Equal(t, "Config", items["app.name"])
	assert.Equal(t, "", items["app.none"])
	assert.Equal(t, 2, len(items))
}
