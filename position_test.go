package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPosition_Compare(t *testing.T) {
	tests := []struct {
		name      string
		positionA Position
		positionB Position
		result    int
	}{
		{
			name: "test equal",
			positionA: Position{
				"",
				4,
			},
			positionB: Position{
				"",
				4,
			},
			result: 0,
		},
		{
			name: "test equal",
			positionA: Position{
				"a",
				4,
			},
			positionB: Position{
				"a",
				4,
			},
			result: 0,
		},
		{
			name: "test equal",
			positionA: Position{
				"a",
				0,
			},
			positionB: Position{
				"a",
				0,
			},
			result: 0,
		},
		{
			name: "test bigger",
			positionA: Position{
				"a",
				2,
			},
			positionB: Position{
				"a",
				1,
			},
			result: 1,
		},
		{
			name: "test bigger",
			positionA: Position{
				"a",
				2,
			},
			positionB: Position{
				"",
				2,
			},
			result: 1,
		},
		{
			name: "test bigger",
			positionA: Position{
				"mysql-bin.000001",
				2,
			},
			positionB: Position{
				"aysql-bin.000001",
				2,
			},
			result: 1,
		},
		{
			name: "test bigger",
			positionA: Position{
				"mysql-bin.000002",
				2,
			},
			positionB: Position{
				"mysql-bin.000001",
				2,
			},
			result: 1,
		},
		{
			name: "test smaller",
			positionA: Position{
				"a",
				1,
			},
			positionB: Position{
				"a",
				2,
			},
			result: -1,
		},
		{
			name: "test smaller",
			positionA: Position{
				"",
				2,
			},
			positionB: Position{
				"a",
				2,
			},
			result: -1,
		},
		{
			name: "test smaller",
			positionB: Position{
				"mysql-bin.000001",
				2,
			},
			positionA: Position{
				"aysql-bin.000001",
				2,
			},
			result: -1,
		},
		{
			name: "test smaller",
			positionB: Position{
				"mysql-bin.000002",
				2,
			},
			positionA: Position{
				"mysql-bin.000001",
				2,
			},
			result: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := tt.positionA.Compare(tt.positionB)
			assert.Equal(t, tt.result, x)
		})
	}
}
