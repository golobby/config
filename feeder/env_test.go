package feeder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Load_It_Should_Return_Error_When_File_Not_Found(t *testing.T) {
	vs, err := Load("test/.404")
	assert.Nil(t, vs)
	assert.Error(t, err)
}

func Test_Load_It_Should_Read_Empty_File_And_Return_Empty_Map(t *testing.T) {
	vs, err := Load("test/.empty.env")
	assert.NoError(t, err)
	assert.Empty(t, vs)
}

func Test_Load_It_Should_Return_Error_When_File_Is_Invalid(t *testing.T) {
	vs, err := Load("test/.buggy.env")
	assert.Nil(t, vs)
	assert.Error(t, err)
}

func Test_Load_It_Should_Read_The_Sample_Env_File(t *testing.T) {
	vs, err := Load("test/.env")
	l := len(vs)

	assert.NoError(t, err)
	assert.Equalf(t, 10, l, "Expected %v got %v", 10, l)

	// Read
	assert.Equal(t, "127.0.0.1", vs["DB_HOST"])
	assert.Equal(t, "App", vs["DB_NAME"])
	assert.Equal(t, "3306", vs["DB_PORT"])
	assert.Equal(t, "MySQL", vs["DB_TYPE"])
	assert.Equal(t, "", vs["APP_NAME"])
	assert.Equal(t, "https://example.com", vs["APP_URL"])
	assert.Equal(t, "true", vs["DEBUG"])
	assert.Equal(t, "#VALUE!", vs["NOT_COMMENT"])
}