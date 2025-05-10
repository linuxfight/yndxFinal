package calc

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestDecodeEncode(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    *Task
		expectError bool
		roundTrip   *Task
	}{
		{
			name:  "valid numeric task",
			input: "t1;5;3.5;+;8.5",
			expected: &Task{
				ID:        "t1",
				Arg1:      5.0,
				Arg2:      3.5,
				Operation: "+",
				Result:    8.5,
			},
			roundTrip: &Task{
				ID:        "t1",
				Arg1:      5.0,
				Arg2:      3.5,
				Operation: "+",
				Result:    8.5,
			},
		},
		{
			name:  "valid string arguments",
			input: "t2;foo;bar;concat;foobar",
			expected: &Task{
				ID:        "t2",
				Arg1:      "foo",
				Arg2:      "bar",
				Operation: "concat",
				Result:    "foobar",
			},
		},
		{
			name:        "invalid format - too few parts",
			input:       "t3;arg1;op",
			expectError: true,
		},
		{
			name:        "invalid format - empty string",
			input:       "",
			expectError: true,
		},
		{
			name:  "mixed numeric and string arguments",
			input: "t4;5.5;three;*;15.5",
			expected: &Task{
				ID:        "t4",
				Arg1:      5.5,
				Arg2:      "three",
				Operation: "*",
				Result:    15.5,
			},
		},
		{
			name: "integer float formatting",
			roundTrip: &Task{
				ID:        "t5",
				Arg1:      10.0,
				Arg2:      2.0,
				Operation: "/",
				Result:    5.0,
			},
		},
		{
			name:  "scientific notation",
			input: "t6;1e3;2e-1;*;200",
			expected: &Task{
				ID:        "t6",
				Arg1:      1000.0,
				Arg2:      0.2,
				Operation: "*",
				Result:    200.0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test Decode
			if tt.input != "" {
				task, err := Decode(tt.input)

				if tt.expectError {
					require.Error(t, err)
					if tt.expected == nil {
						return
					}
				} else {
					require.NoError(t, err)
					assert.Equal(t, tt.expected, task)
				}
			}

			// Test Encode and round trip
			if tt.roundTrip != nil {
				encoded := Encode(*tt.roundTrip)

				// Verify numeric formatting
				parts := strings.Split(encoded, ";")
				require.Len(t, parts, 5)

				// Check integer formatting
				assert.Equal(t, ToString(tt.roundTrip.Arg1), ToString(parts[1]))
				assert.Equal(t, ToString(tt.roundTrip.Result), ToString(parts[4]))

				// Test round trip decoding
				decoded, err := Decode(encoded)
				require.NoError(t, err)
				assert.Equal(t, tt.roundTrip, decoded)
			}
		})
	}
}

func TestEdgeCases(t *testing.T) {
	t.Run("empty fields", func(t *testing.T) {
		encoded := Encode(Task{
			ID:        "",
			Arg1:      "",
			Arg2:      nil,
			Operation: "",
			Result:    "",
		})
		assert.Equal(t, ";;<nil>;;", encoded)
	})

	t.Run("special characters", func(t *testing.T) {
		task := &Task{
			ID:        "t7",
			Arg1:      "foo;bar",
			Arg2:      "baz",
			Operation: "join",
			Result:    "foo;bar;baz",
		}
		encoded := Encode(*task)
		assert.Equal(t, "t7;foo;bar;baz;join;foo;bar;baz", encoded)

		// This should fail decode due to extra semicolons
		_, err := Decode(encoded)
		assert.Error(t, err)
	})

	t.Run("large numbers", func(t *testing.T) {
		task := &Task{
			ID:        "t8",
			Arg1:      1e18,
			Arg2:      2e-18,
			Operation: "*",
			Result:    2.0,
		}
		encoded := Encode(*task)
		assert.Equal(t, "t8;1000000000000000000;2e-18;*;2", encoded)
	})
}

func TestTypePreservation(t *testing.T) {
	tests := []struct {
		name  string
		task  Task
		check func(t *testing.T, decoded *Task)
	}{
		{
			name: "string arguments",
			task: Task{
				ID:        "t9",
				Arg1:      "foo",
				Arg2:      "bar",
				Operation: "concat",
				Result:    "foobar",
			},
			check: func(t *testing.T, decoded *Task) {
				assert.IsType(t, "", decoded.Arg1)
				assert.IsType(t, "", decoded.Arg2)
				assert.IsType(t, "", decoded.Result)
			},
		},
		{
			name: "mixed types",
			task: Task{
				ID:        "t10",
				Arg1:      5.0,
				Arg2:      "items",
				Operation: "count",
				Result:    "5 items",
			},
			check: func(t *testing.T, decoded *Task) {
				assert.IsType(t, 0.0, decoded.Arg1)
				assert.IsType(t, "", decoded.Arg2)
				assert.IsType(t, "", decoded.Result)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded := Encode(tt.task)
			decoded, err := Decode(encoded)
			require.NoError(t, err)
			tt.check(t, decoded)
		})
	}
}

func TestParseExpression(t *testing.T) {
	tests := []struct {
		name        string
		expression  string
		wantTasks   int
		wantErr     error
		checkResult func(t *testing.T, tasks []Task)
	}{
		// Original parsing tests
		{
			name:       "simple addition",
			expression: "2+3",
			wantTasks:  1,
			checkResult: func(t *testing.T, tasks []Task) {
				assertTask(t, tasks[0], "+", 2.0, 3.0)
			},
		},
		{
			name:       "operator precedence",
			expression: "2+3*4",
			wantTasks:  2,
			checkResult: func(t *testing.T, tasks []Task) {
				assertTask(t, tasks[0], "*", 3.0, 4.0)
				assertTask(t, tasks[1], "+", 2.0, tasks[0].ID)
			},
		},
		{
			name:       "parentheses",
			expression: "(2+3)*4",
			wantTasks:  2,
			checkResult: func(t *testing.T, tasks []Task) {
				assertTask(t, tasks[0], "+", 2.0, 3.0)
				assertTask(t, tasks[1], "*", tasks[0].ID, 4.0)
			},
		},
		{
			name:       "unary minus",
			expression: "-5+3",
			wantTasks:  2,
			checkResult: func(t *testing.T, tasks []Task) {
				assertTask(t, tasks[0], "-", 0.0, 5.0)
				assertTask(t, tasks[1], "+", tasks[0].ID, 3.0)
			},
		},
		{
			name:       "division by zero",
			expression: "5/0",
			wantTasks:  0,
			wantErr:    errDivisionByZero,
		},
		{
			name:       "floating point numbers",
			expression: "3.5*2",
			wantTasks:  1,
			checkResult: func(t *testing.T, tasks []Task) {
				assertTask(t, tasks[0], "*", 3.5, 2.0)
			},
		},
		{
			name:       "complex expression",
			expression: "((3+5)*2-4)/2",
			wantTasks:  4,
			checkResult: func(t *testing.T, tasks []Task) {
				assertTask(t, tasks[0], "+", 3.0, 5.0)
				assertTask(t, tasks[1], "*", tasks[0].ID, 2.0)
				assertTask(t, tasks[2], "-", tasks[1].ID, 4.0)
				assertTask(t, tasks[3], "/", tasks[2].ID, 2.0)
			},
		},

		// Validation error tests
		{
			name:       "empty string",
			expression: "",
			wantTasks:  0,
			wantErr:    errors.New("parsing error: 1:1: expected operand, found 'EOF'"), // Actually returns nil in original impl
		},
		{
			name:       "invalid space",
			expression: "1 + 2",
			wantTasks:  0,
			wantErr:    fmt.Errorf("invalid character: %c", ' '),
		},
		{
			name:       "invalid letter",
			expression: "a+3",
			wantTasks:  0,
			wantErr:    fmt.Errorf("invalid character: %c", 'a'),
		},
		{
			name:       "invalid comma",
			expression: "5,000",
			wantTasks:  0,
			wantErr:    fmt.Errorf("invalid character: %c", ','),
		},
		{
			name:       "invalid dollar sign",
			expression: "$100",
			wantTasks:  0,
			wantErr:    fmt.Errorf("invalid character: %c", '$'),
		},
		{
			name:       "multiple invalid characters",
			expression: "a$b√c",
			wantTasks:  0,
			wantErr:    fmt.Errorf("invalid character: %c", 'a'),
		},
		{
			name:       "invalid at end",
			expression: "123+4#",
			wantTasks:  0,
			wantErr:    fmt.Errorf("invalid character: %c", '#'),
		},
		{
			name:       "invalid underscore",
			expression: "var_1",
			wantTasks:  0,
			wantErr:    fmt.Errorf("invalid character: %c", 'v'),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tasks, err := ParseExpression(tt.expression)

			// Handle validation errors
			if tt.wantErr != nil {
				if err == nil {
					t.Fatal("expected error but got nil")
				}
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("unexpected error:\nwant: %q\ngot:  %q", tt.wantErr.Error(), err.Error())
				}
				return
			}

			// Handle parsing errors
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(tasks) != tt.wantTasks {
				t.Errorf("expected %d tasks, got %d", tt.wantTasks, len(tasks))
			}

			if tt.checkResult != nil {
				tt.checkResult(t, tasks)
			}
		})
	}
}

func assertTask(t *testing.T, task Task, op string, arg1, arg2 interface{}) {
	t.Helper()

	if task.Operation != op {
		t.Errorf("expected operation %q, got %q", op, task.Operation)
	}

	if !compareArgs(task.Arg1, arg1) {
		t.Errorf("arg1 mismatch: expected %v (%T), got %v (%T)", arg1, arg1, task.Arg1, task.Arg1)
	}

	if !compareArgs(task.Arg2, arg2) {
		t.Errorf("arg2 mismatch: expected %v (%T), got %v (%T)", arg2, arg2, task.Arg2, task.Arg2)
	}
}

func compareArgs(a, b interface{}) bool {
	switch v1 := a.(type) {
	case float64:
		if v2, ok := b.(float64); ok {
			return v1 == v2
		}
	case string:
		if v2, ok := b.(string); ok {
			return v2 == v1
		}
	}
	return false
}
