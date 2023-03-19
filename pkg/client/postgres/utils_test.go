package postgres

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateWhereClause(t *testing.T) {
	query := "SELECT * from test "

	tests := []struct {
		name          string
		whereClause   map[string]interface{}
		expectedQuery string
	}{
		{
			name: "Test Case - Decent Case Greater Than",
			whereClause: map[string]interface{}{
				"id_gt": 5,
			},
			expectedQuery: fmt.Sprintf("%sWHERE id > '5'", query),
		},
		{
			name: "Test Case - Capital Case Greater Than",
			whereClause: map[string]interface{}{
				"Caps_gt": 5,
			},
			expectedQuery: fmt.Sprintf(`%sWHERE "Caps" > '5'`, query),
		},
		{
			name: "Test Case - Decent Case Equal",
			whereClause: map[string]interface{}{
				"Caps": 5,
			},
			expectedQuery: fmt.Sprintf(`%sWHERE "Caps" = 5`, query),
		},
		{
			name: "Test Case - Capital Case Lower Equal",
			whereClause: map[string]interface{}{
				"Caps_lower": "doctor",
			},
			expectedQuery: fmt.Sprintf(`%sWHERE lower("Caps")=lower('doctor')`, query),
		},
		{
			name: "Test Case - Decent Case In",
			whereClause: map[string]interface{}{
				"test_in": []string{"chicken", "wings"},
			},
			expectedQuery: fmt.Sprintf(`%sWHERE test IN ('chicken','wings')`, query),
		},
		{
			name:          "Test Case Empty Where Clause",
			whereClause:   map[string]interface{}{},
			expectedQuery: query,
		},
		{
			name:          "Test Case Nil Where Clause",
			whereClause:   nil,
			expectedQuery: query,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formattedQuery := GenerateWhereClause(query, tt.whereClause)
			assert.Equal(t, tt.expectedQuery, formattedQuery)
		})
	}
}
