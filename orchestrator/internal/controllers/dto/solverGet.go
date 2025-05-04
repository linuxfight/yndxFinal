package dto

import "orchestrator/internal/db/expressions"

type ListAllExpressionsResponse struct {
	Expressions []Expression `json:"expressions"`
}

type GetByIdExpressionResponse struct {
	Expression Expression `json:"expression"`
}

type Expression struct {
	Id     string  `json:"id" example:"928b303f-cfcc-46f4-ae24-aabb72bbb7d9"`
	Result float64 `json:"result"`
	Status string  `json:"status" example:"DONE"`
}

func NewExpression(expr expressions.Expression) Expression {
	res := Expression{
		Id:     expr.ID,
		Result: expr.Res,
	}

	if !expr.Finished {
		res.Status = Processing
	} else {
		if !expr.Error {
			res.Status = Done
		} else {
			res.Status = Failed
		}
	}

	return res
}
