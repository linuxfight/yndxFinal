package tasks

import (
	"orchestrator/internal/controllers/tasks/gen"
	"orchestrator/pkg/calc"
)

func getOperator(t calc.Task) gen.Operator {
	switch t.Operation {
	case "+":
		return gen.Operator_ADDICTION
	case "-":
		return gen.Operator_SUBTRACTION
	case "*":
		return gen.Operator_MULTIPLICATION
	case "/":
		return gen.Operator_DIVISION
	}

	panic("invalid operation")
}
