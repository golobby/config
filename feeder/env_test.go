package feeder

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnv_Feed_With_Invalid_File_Path_It_Should_Fail(t *testing.T) {
	e := Env{Path: "test/.404"}

	_, err := e.Feed()
	assert.Error(t, err)
}

func TestEnv_Feed_With_Empty_File_Path_It_Should_Fail(t *testing.T) {
	e := Env{}

	_, err := e.Feed()
	assert.Error(t, err)
}

func TestEnv_Feed_With_Empty_File_It_Should_Hold_No_Item(t *testing.T) {
	e := Env{Path: "test/.empty.env"}

	items, err := e.Feed()
	assert.NoError(t, err)
	assert.Empty(t, items)
}

func TestEnv_Feed_With_Buggy_File_It_Should_Fail(t *testing.T) {
	e := Env{Path: "test/.buggy.env"}

	_, err := e.Feed()
	assert.Error(t, err)
}

func TestEnv_Feed_It_Should_Read_The_Sample_Env_File(t *testing.T) {
	e := Env{Path: "test/.env"}

	items, err := e.Feed()
	l := len(items)

	assert.NoError(t, err)
	assert.Equalf(t, 10, l, "Expected %v got %v", 10, l)

	assert.Equal(t, "https://example.com", items["url"])
	assert.Equal(t, "127.0.0.1", items["db.host"])
	assert.Equal(t, "NewApp", items["db.name"])
	assert.Equal(t, "3306", items["db.port"])
	assert.Equal(t, "MySQL", items["db.type"])
	assert.Equal(t, "", items["app.name"])
	assert.Equal(t, "https://app.url", items["app.url"])
	assert.Equal(t, "true", items["debug"])
	assert.Equal(t, "#VALUE!", items["not.comment"])
	assert.Equal(t, "", items["name"])
}

func TestEnv_Feed_It_Should_Use_OS_When_Key_Is_Empty_In_Env_File(t *testing.T) {
	_ = os.Setenv("APP_NAME", "FROM_OS")

	e := Env{Path: "test/.env"}

	items, err := e.Feed()
	assert.NoError(t, err)

	assert.Equal(t, "FROM_OS", items["app.name"])
}

func TestEnv_Feed_It_Should_Not_Use_OS_When_The_Feature_Is_Disabled(t *testing.T) {
	_ = os.Setenv("APP_NAME", "FROM_OS")

	e := Env{Path: "test/.env", DisableOSVariables: true}

	items, err := e.Feed()
	assert.NoError(t, err)

	assert.Equal(t, "", items["app.name"])
}
