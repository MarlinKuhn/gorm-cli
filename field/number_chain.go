package field

import (
	"golang.org/x/exp/constraints"
	"gorm.io/gorm/clause"
)

// NumberChain represents a numeric field that supports both integer and float types, and allows for chaining operations to build complex SQL expressions.
type NumberChain[T constraints.Integer | constraints.Float] struct {
	column any
}

// Query functions

// Eq creates an equality comparison expression (field = value).
func (n NumberChain[T]) Eq(value T) clause.Expression {
	return clause.Eq{Column: n.column, Value: value}
}

// EqExpr creates an equality comparison expression (field = expression).
func (n NumberChain[T]) EqExpr(expr any) clause.Expression {
	if col, ok := expr.(ColumnInterface); ok {
		expr = col.Column()
	}
	return clause.Eq{Column: n.column, Value: expr}
}

// Neq creates a not equal comparison expression (field != value).
func (n NumberChain[T]) Neq(value T) clause.Expression {
	return clause.Neq{Column: n.column, Value: value}
}

// NeqExpr creates a not equal comparison expression (field != expression).
func (n NumberChain[T]) NeqExpr(expr any) clause.Expression {
	if col, ok := expr.(ColumnInterface); ok {
		expr = col.Column()
	}
	return clause.Neq{Column: n.column, Value: expr}
}

// Gt creates a greater than comparison expression (field > value).
func (n NumberChain[T]) Gt(value T) clause.Expression {
	return clause.Gt{Column: n.column, Value: value}
}

// GtExpr creates a greater than comparison expression (field > expression).
func (n NumberChain[T]) GtExpr(expr any) clause.Expression {
	if col, ok := expr.(ColumnInterface); ok {
		expr = col.Column()
	}
	return clause.Gt{Column: n.column, Value: expr}
}

// Gte creates a greater than or equal comparison expression (field >= value).
func (n NumberChain[T]) Gte(value T) clause.Expression {
	return clause.Gte{Column: n.column, Value: value}
}

// GteExpr creates a greater than or equal comparison expression (field >= expression).
func (n NumberChain[T]) GteExpr(expr any) clause.Expression {
	if col, ok := expr.(ColumnInterface); ok {
		expr = col.Column()
	}
	return clause.Gte{Column: n.column, Value: expr}
}

// Lt creates a less than comparison expression (field < value).
func (n NumberChain[T]) Lt(value T) clause.Expression {
	return clause.Lt{Column: n.column, Value: value}
}

// LtExpr creates a less than comparison expression (field < expression).
func (n NumberChain[T]) LtExpr(expr any) clause.Expression {
	if col, ok := expr.(ColumnInterface); ok {
		expr = col.Column()
	}
	return clause.Lt{Column: n.column, Value: expr}
}

// Lte creates a less than or equal comparison expression (field <= value).
func (n NumberChain[T]) Lte(value T) clause.Expression {
	return clause.Lte{Column: n.column, Value: value}
}

// LteExpr creates a less than or equal comparison expression (field <= expression).
func (n NumberChain[T]) LteExpr(expr any) clause.Expression {
	if col, ok := expr.(ColumnInterface); ok {
		expr = col.Column()
	}
	return clause.Lte{Column: n.column, Value: expr}
}

// Between creates a range comparison expression (field BETWEEN v1 AND v2).
func (n NumberChain[T]) Between(v1, v2 T) clause.Expression {
	return clause.And(
		clause.Gte{Column: n.column, Value: v1},
		clause.Lte{Column: n.column, Value: v2},
	)
}

// In creates an IN comparison expression (field IN (values...)).
func (n NumberChain[T]) In(values ...T) clause.Expression {
	interfaceValues := make([]any, len(values))
	for i, v := range values {
		interfaceValues[i] = v
	}
	return clause.IN{Column: n.column, Values: interfaceValues}
}

// NotIn creates a NOT IN comparison expression (field NOT IN (values...)).
func (n NumberChain[T]) NotIn(values ...T) clause.Expression {
	interfaceValues := make([]any, len(values))
	for i, v := range values {
		interfaceValues[i] = v
	}
	return clause.Not(clause.IN{Column: n.column, Values: interfaceValues})
}

// Add creates an addition expression (field + value).
func (n NumberChain[T]) Add(value T) NumberChain[T] {
	return NumberChain[T]{column: clause.Expr{SQL: "? + ?", Vars: []any{n.column, value}}}
}

// AddExpr creates an addition expression (field + expression).
func (n NumberChain[T]) AddExpr(expr any) NumberChain[T] {
	if col, ok := expr.(ColumnInterface); ok {
		expr = col.Column()
	}
	return NumberChain[T]{column: clause.Expr{SQL: "? + ?", Vars: []any{n.column, expr}}}
}

// Sub creates a subtraction expression (field - value).
func (n NumberChain[T]) Sub(value T) NumberChain[T] {
	return NumberChain[T]{column: clause.Expr{SQL: "? - ?", Vars: []any{n.column, value}}}
}

// SubExpr creates a subtraction expression (field - expression).
func (n NumberChain[T]) SubExpr(expr any) NumberChain[T] {
	if col, ok := expr.(ColumnInterface); ok {
		expr = col.Column()
	}
	return NumberChain[T]{column: clause.Expr{SQL: "? - ?", Vars: []any{n.column, expr}}}
}

// Mul creates a multiplication expression (field * value).
func (n NumberChain[T]) Mul(value T) NumberChain[T] {
	return NumberChain[T]{column: clause.Expr{SQL: "? * ?", Vars: []any{n.column, value}}}
}

// MulExpr creates a multiplication expression (field * expression).
func (n NumberChain[T]) MulExpr(expr any) NumberChain[T] {
	if col, ok := expr.(ColumnInterface); ok {
		expr = col.Column()
	}
	return NumberChain[T]{column: clause.Expr{SQL: "? * ?", Vars: []any{n.column, expr}}}
}

// Div creates a division expression (field / value).
func (n NumberChain[T]) Div(value T) NumberChain[T] {
	return NumberChain[T]{column: clause.Expr{SQL: "? / ?", Vars: []any{n.column, value}}}
}

// DivExpr creates a division expression (field / expression).
func (n NumberChain[T]) DivExpr(expr any) NumberChain[T] {
	if col, ok := expr.(ColumnInterface); ok {
		expr = col.Column()
	}
	return NumberChain[T]{column: clause.Expr{SQL: "? / ?", Vars: []any{n.column, expr}}}
}
