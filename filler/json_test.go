package filler

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJson_Fill_With_Valid_JSON(t *testing.T) {
	type config struct {
		Name    string
		Version float64
		Numbers []int
		Users []struct{
			Name string
			Year int
			Address struct{
				Country string
				State string
				City string
			}
		}
	}

	c := &config{}

	j := Json{"test/config.json"}
	err := j.Fill(c)
	assert.NoError(t, err)

	assert.Equal(t, "MyAppUsingConfig", c.Name)
	assert.Equal(t, 3.14, c.Version)
	assert.Equal(t, []int{1, 2, 3}, c.Numbers)
	assert.Equal(t, "Milad Rahimi", c.Users[0].Name)
	assert.Equal(t, "Amirreza Askarpour", c.Users[1].Name)
}
