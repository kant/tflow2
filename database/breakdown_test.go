package database

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBreakdownKeyString(t *testing.T) {
	assert := assert.New(t)

	// Empty Key
	key := BreakdownKey{}
	assert.Equal("", key.String())

	// Set one key
	key.Set("DstPort", "23")
	assert.Equal(key.Get("DstPort"), "23")
	assert.Equal("DstPort:23", key.String())

	// Set all keys
	for i := range breakdownLabels {
		key[i] = strconv.Itoa(i)
	}
	assert.Equal("Family:2,SrcAddr:3,DstAddr:4,Protocol:5,IntIn:6,IntOut:7,NextHop:8,SrcAsn:9,DstAsn:10,NextHopAsn:11,SrcPfx:12,DstPfx:13,SrcPort:14,DstPort:15", key.String())
}

func TestBreakdownFlags(t *testing.T) {
	assert := assert.New(t)

	// Defaults
	key := BreakdownFlags{}
	assert.False(key.DstAddr)

	// Enable all
	assert.NoError(key.Set([]string{"Family", "SrcAddr", "DstAddr", "Protocol", "IntIn", "IntOut", "NextHop", "SrcAsn", "DstAsn", "NextHopAsn", "SrcPfx", "DstPfx", "SrcPort", "DstPort"}))
	assert.True(key.DstAddr)

	// Invalid key
	assert.EqualError(key.Set([]string{"foobar"}), "invalid breakdown key: foobar")
}

func TestGetBreakdownLabels(t *testing.T) {
	assert := assert.New(t)

	labels := GetBreakdownLabels()
	assert.NotNil(labels)
	assert.Contains(labels, "SrcAddr")
}
