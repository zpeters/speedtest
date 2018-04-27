package main

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	t.Parallel()
	assert.NotNil(t, viper.GetBool("debug"))
	assert.NotNil(t, viper.GetBool("quiet"))
	assert.NotNil(t, viper.GetBool("report"))
	assert.NotEmpty(t, viper.GetInt("numclosest"))
	assert.NotEmpty(t, viper.GetInt("numlatencytests"))
	assert.NotEmpty(t, viper.GetString("reportchar"))
	assert.NotEmpty(t, viper.GetString("algotype"))
	assert.NotEmpty(t, viper.GetInt("httptimeout"))
	assert.NotEmpty(t, viper.Get("dlsizes").([]int))
	assert.NotEmpty(t, viper.Get("ulsizes").([]int))
}
