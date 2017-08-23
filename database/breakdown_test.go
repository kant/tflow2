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
	assert.Equal("SrcAddr:2,DstAddr:3,Protocol:4,IntIn:5,IntOut:6,NextHop:7,SrcAsn:8,DstAsn:9,NextHopAsn:10,SrcPfx:11,DstPfx:12,SrcPort:13,DstPort:14", key.String())
}

func TestBreakdownFlags(t *testing.T) {
	assert := assert.New(t)

	// Defaults
	key := BreakdownFlags{}
	assert.False(key.DstAddr)

	// Enable all
	assert.NoError(key.Set([]string{"Router", "Family", "SrcAddr", "DstAddr", "Protocol", "IntIn", "IntOut", "NextHop", "SrcAsn", "DstAsn", "NextHopAsn", "SrcPfx", "DstPfx", "SrcPort", "DstPort"}))
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
