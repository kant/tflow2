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
	assert.Equal("SrcAddr:6,DstAddr:7,Protocol:8,IntIn:9,IntOut:10,NextHop:11,SrcAsn:12,DstAsn:13,NextHopAsn:14,SrcPfx:15,DstPfx:16,SrcPort:17,DstPort:18", key.String())
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
