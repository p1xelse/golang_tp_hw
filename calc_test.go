package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCalc(t *testing.T) {
	tests := map[string]struct {
		input   string
		result  float64
		needErr bool
	}{
		"Simple_test": {
			input:   "1 + 2 + 3 + 3 - 3",
			result:  6,
			needErr: false,
		},
		"Test_multiply": {
			input:   "1 * 2 * 3",
			result:  6,
			needErr: false,
		},
		"Test_devide": {
			input:   "1 * 6 / 2",
			result:  3,
			needErr: false,
		},
		"Test_brackets": {
			input:   "2 * (6 + 2)",
			result:  16,
			needErr: false,
		},
		"Test_more_brackets": {
			input:   "(2 + 1) * (6 + 2)",
			result:  24,
			needErr: false,
		},
		"Test_single_number": {
			input:   "3",
			result:  3,
			needErr: false,
		},
		"Test_invalid_expr": {
			input:   "(2 + 1)( * (6 + 2)",
			result:  0,
			needErr: true,
		},
		"Test_empty_expr": {
			input:   "",
			result:  0,
			needErr: true,
		},
		"Test_divide_by_zero": {
			input:   "23 / (2 - 2)",
			result:  0,
			needErr: true,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			got, err := Calc(test.input)
			expected := test.result

			if test.needErr {
				require.Error(t, err)
			} else {
				require.Equal(t, nil, err)
			}

			require.Equal(t, expected, got)
		})
	}
}
