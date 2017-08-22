package frontend

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/taktv6/tflow2/database"
)

func TestTranslateQuery(t *testing.T) {
	assert := assert.New(t)

	ext := QueryExt{
		Cond: ConditionsExt{ConditionExt{
			Operator: database.OpEqual,
			Field:    "Timestamp",
			Operand:  "1503432000",
		}},
	}
	query, err := translateQuery(&ext)
	assert.NoError(err)
	assert.NotNil(query)
}

func TestTranslateQueryInvalid(t *testing.T) {
	assert := assert.New(t)

	ext := QueryExt{
		Cond: ConditionsExt{ConditionExt{
			Operator: database.OpEqual,
			Field:    "Unknown",
		}},
	}
	query, err := translateQuery(&ext)
	assert.EqualError(err, "unknown field: Unknown")
	assert.Nil(query)
}
