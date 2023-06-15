package lang

import (
	"strings"
	"testing"

	"github.com/Snare295/architect-lab3/painter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_parse_func(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedResult painter.Operation
	}{
		{
			name:           "Valid White Command",
			input:          "white\nupdate",
			expectedResult: painter.OperationFunc(painter.WhiteFill),
		},
		{
			name:           "Valid Green Command",
			input:          "green\nupdate",
			expectedResult: painter.OperationFunc(painter.GreenFill),
		},
		{
			name:           "Valid Reset Command",
			input:          "reset\nupdate",
			expectedResult: painter.OperationFunc(painter.Reset),
		},
	}
	parser := &Parser{}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expResults, err := parser.Parse(strings.NewReader(test.input))
			require.NoError(t, err)
			assert.IsType(t, test.expectedResult, expResults[0])
		})
	}
}

func Test_parse_struct(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedResult painter.Operation
	}{
		{
			name:           "Valid BgRect Command",
			input:          "bgrect 0.25 0.25 0.75 0.75\nupdate ",
			expectedResult: &painter.BackRect{X1: 200, Y1: 200, X2: 600, Y2: 600},
		},
		{
			name:           "Valid Figure Command",
			input:          "figure 0.5 0.5\nupdate",
			expectedResult: &painter.Figure{X: 400, Y: 400},
		},
		{
			name:           "Valid Move Command",
			input:          "move 0.3 0.3\nupdate",
			expectedResult: &painter.Move{X: 240, Y: 240},
		},
		{
			name:           "Valid Update Command",
			input:          "update",
			expectedResult: painter.UpdateOp,
		},
		{
			name:           "Invalid Command",
			input:          "invalidcommand\nupdate",
			expectedResult: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parser := &Parser{}
			expResults, err := parser.Parse(strings.NewReader(test.input))
			if test.expectedResult == nil {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.IsType(t, test.expectedResult, expResults[1])
				assert.Equal(t, test.expectedResult, expResults[1])
			}
		})
	}
}
