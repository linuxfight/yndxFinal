package utils

import (
	"errors"
	"fmt"
	"github.com/oklog/ulid/v2"
	"go/ast"
	"go/parser"
	"go/token"
	"orchestrator/internal/controllers/dto"
	"orchestrator/internal/controllers/tasksServer"
	"strconv"
	"strings"
)

type InternalTask struct {
	ID        string
	Arg1      interface{}
	Arg2      interface{}
	Operation string
	Result    interface{}
}

var (
	errDivisionByZero  = fmt.Errorf("division by zero")
	errUnsupportedNode = fmt.Errorf("unsupported node type")
)

func Decode(task string) (*InternalTask, error) {
	args := strings.Split(task, ";")
	if len(args) != 5 {
		return nil, errors.New("invalid string")
	}

	var arg1 interface{}
	arg1float, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		arg1 = args[1]
	} else {
		arg1 = arg1float
	}
	var arg2 interface{}
	arg2float, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		arg2 = args[2]
	} else {
		arg2 = arg2float
	}

	var res interface{}
	resfloat, err := strconv.ParseFloat(args[4], 64)
	if err != nil {
		res = args[4]
	} else {
		res = resfloat
	}

	return &InternalTask{
		ID:        args[0],
		Arg1:      arg1,
		Arg2:      arg2,
		Operation: args[3],
		Result:    res,
	}, nil
}

func Encode(t InternalTask) string {
	toString := func(v interface{}) string {
		switch val := v.(type) {
		case string:
			return val
		case float64:
			// Format without decimal if it's a whole number
			if val == float64(int(val)) {
				return fmt.Sprintf("%d", int(val))
			}
			return fmt.Sprintf("%g", val)
		default:
			return fmt.Sprintf("%v", val)
		}
	}

	return fmt.Sprintf("%s;%s;%s;%s;%s",
		toString(t.ID),
		toString(t.Arg1),
		toString(t.Arg2),
		toString(t.Operation),
		toString(t.Result),
	)
}

func GetOperator(t InternalTask) tasksServer.Operator {
	switch t.Operation {
	case "+":
		return tasksServer.Operator_ADDICTION
	case "-":
		return tasksServer.Operator_SUBTRACTION
	case "*":
		return tasksServer.Operator_MULTIPLICATION
	case "/":
		return tasksServer.Operator_DIVISION
	}

	panic("invalid operation")
}

// ParseExpression parses a mathematical expression into a sequence of tasksStorage
func ParseExpression(expression string) ([]InternalTask, error) {
	exprAst, err := parser.ParseExpr(expression)
	if err != nil {
		return nil, fmt.Errorf("parsing error: %w", err)
	}

	var tasks []InternalTask
	_, err = processNode(exprAst, &tasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

// processNode recursively processes AST nodes and creates tasksStorage
func processNode(node ast.Node, tasks *[]InternalTask) (interface{}, error) {
	switch n := node.(type) {
	case *ast.BinaryExpr:
		return processBinaryExpr(n, tasks)
	case *ast.UnaryExpr:
		return processUnaryExpr(n, tasks)
	case *ast.BasicLit:
		return processBasicLit(n)
	case *ast.ParenExpr:
		return processNode(n.X, tasks)
	default:
		return nil, errUnsupportedNode
	}
}

func processBinaryExpr(expr *ast.BinaryExpr, tasks *[]InternalTask) (interface{}, error) {
	left, err := processNode(expr.X, tasks)
	if err != nil {
		return nil, err
	}

	right, err := processNode(expr.Y, tasks)
	if err != nil {
		return nil, err
	}

	// Check for division by zero with literal values
	if expr.Op == token.QUO {
		if rval, ok := right.(float64); ok && rval == 0 {
			return nil, errDivisionByZero
		}
	}

	return createTask(tasks, left, right, expr.Op.String())
}

func processUnaryExpr(expr *ast.UnaryExpr, tasks *[]InternalTask) (interface{}, error) {
	if expr.Op != token.SUB {
		return nil, fmt.Errorf("unsupported unary operator: %v", expr.Op)
	}

	operand, err := processNode(expr.X, tasks)
	if err != nil {
		return nil, err
	}

	return createTask(tasks, 0.0, operand, token.SUB.String())
}

func processBasicLit(lit *ast.BasicLit) (float64, error) {
	switch lit.Kind {
	case token.INT, token.FLOAT:
		value, err := strconv.ParseFloat(lit.Value, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid number: %w", err)
		}
		return value, nil
	default:
		return 0, fmt.Errorf("unsupported literal type: %v", lit.Kind)
	}
}

func createTask(tasks *[]InternalTask, left, right interface{}, operation string) (string, error) {
	taskID := ulid.Make().String()
	*tasks = append(*tasks, InternalTask{
		ID:        taskID,
		Arg1:      left,
		Arg2:      right,
		Operation: operation,
		Result:    dto.Processing,
	})
	return taskID, nil
}
