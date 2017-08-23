package frontend

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/taktv6/tflow2/database"
)

func TestTranslateQuery(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		Field         string
		Operand       string
		ExpectedField int
	}{
		{
			Field:         "Timestamp",
			Operand:       "1503432000",
			ExpectedField: database.FieldTimestamp,
		},
		{
			Field:         "Protocol",
			Operand:       "6",
			ExpectedField: database.FieldProtocol,
		},
		{
			Field:         "SrcAddr",
			Operand:       "1.2.3.4",
			ExpectedField: database.FieldSrcAddr,
		},
		{
			Field:         "SrcAs",
			Operand:       "5123",
			ExpectedField: database.FieldSrcAs,
		},
		{
			Field:         "SrcPfx",
			Operand:       "10.8.0.0/16",
			ExpectedField: database.FieldSrcPfx,
		},
	}

	for _, test := range tests {
		ext := QueryExt{
			Cond: ConditionsExt{ConditionExt{
				Field:   test.Field,
				Operand: test.Operand,
			}},
		}
		query, err := translateQuery(&ext)
		assert.NoError(err)
		assert.NotNil(query.Cond[0].Operand)
		assert.Equal(test.ExpectedField, query.Cond[0].Field)
	}

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
