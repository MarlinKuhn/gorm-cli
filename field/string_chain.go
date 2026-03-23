package field

import "gorm.io/gorm/clause"

// StringChain represents a string field that allows chaining operations to build complex SQL expressions.
type StringChain struct {
	column any
}

// Query functions

// Eq creates an equality comparison expression (field = value).
func (s StringChain) Eq(value string) clause.Expression {
	return clause.Eq{Column: s.column, Value: value}
}

// EqExpr creates an equality comparison expression (field = expression).
func (s StringChain) EqExpr(expr any) clause.Expression {
	if col, ok := expr.(ColumnInterface); ok {
		expr = col.Column()
	}
	return clause.Eq{Column: s.column, Value: expr}
}

// Neq creates a not equal comparison expression (field != value).
func (s StringChain) Neq(value string) clause.Expression {
	return clause.Neq{Column: s.column, Value: value}
}

// NeqExpr creates a not equal comparison expression (field != expression).
func (s StringChain) NeqExpr(expr any) clause.Expression {
	if col, ok := expr.(ColumnInterface); ok {
		expr = col.Column()
	}
	return clause.Neq{Column: s.column, Value: expr}
}

// Gt creates a greater than comparison expression (field > value).
func (s StringChain) Gt(value string) clause.Expression {
	return clause.Gt{Column: s.column, Value: value}
}

// GtExpr creates a greater than comparison expression (field > expression).
func (s StringChain) GtExpr(expr any) clause.Expression {
	if col, ok := expr.(ColumnInterface); ok {
		expr = col.Column()
	}
	return clause.Gt{Column: s.column, Value: expr}
}

// Gte creates a greater than or equal comparison expression (field >= value).
func (s StringChain) Gte(value string) clause.Expression {
	return clause.Gte{Column: s.column, Value: value}
}

// GteExpr creates a greater than or equal comparison expression (field >= expression).
func (s StringChain) GteExpr(expr any) clause.Expression {
	if col, ok := expr.(ColumnInterface); ok {
		expr = col.Column()
	}
	return clause.Gte{Column: s.column, Value: expr}
}

// Lt creates a less than comparison expression (field < value).
func (s StringChain) Lt(value string) clause.Expression {
	return clause.Lt{Column: s.column, Value: value}
}

// LtExpr creates a less than comparison expression (field < expression).
func (s StringChain) LtExpr(expr any) clause.Expression {
	if col, ok := expr.(ColumnInterface); ok {
		expr = col.Column()
	}
	return clause.Lt{Column: s.column, Value: expr}
}

// Lte creates a less than or equal comparison expression (field <= value).
func (s StringChain) Lte(value string) clause.Expression {
	return clause.Lte{Column: s.column, Value: value}
}

// LteExpr creates a less than or equal comparison expression (field <= expression).
func (s StringChain) LteExpr(expr any) clause.Expression {
	if col, ok := expr.(ColumnInterface); ok {
		expr = col.Column()
	}
	return clause.Lte{Column: s.column, Value: expr}
}

// Like creates a LIKE pattern matching expression (field LIKE pattern).
func (s StringChain) Like(pattern string) clause.Expression {
	return clause.Like{Column: s.column, Value: pattern}
}

// NotLike creates a NOT LIKE pattern matching expression (field NOT LIKE pattern).
func (s StringChain) NotLike(pattern string) clause.Expression {
	return clause.Expr{SQL: "? NOT LIKE ?", Vars: []any{s.column, pattern}}
}

// ILike creates a case-insensitive LIKE pattern matching expression (field ILIKE pattern).
func (s StringChain) ILike(pattern string) clause.Expression {
	return clause.Expr{SQL: "? ILIKE ?", Vars: []any{s.column, pattern}}
}

// NotILike creates a case-insensitive NOT LIKE pattern matching expression (field NOT ILIKE pattern).
func (s StringChain) NotILike(pattern string) clause.Expression {
	return clause.Expr{SQL: "? NOT ILIKE ?", Vars: []any{s.column, pattern}}
}

// Regexp creates a regular expression matching expression (field REGEXP pattern).
func (s StringChain) Regexp(pattern string) clause.Expression {
	return clause.Expr{SQL: "? REGEXP ?", Vars: []any{s.column, pattern}}
}

// NotRegexp creates a regular expression not matching expression (field NOT REGEXP pattern).
func (s StringChain) NotRegexp(pattern string) clause.Expression {
	return clause.Expr{SQL: "? NOT REGEXP ?", Vars: []any{s.column, pattern}}
}

// In creates an IN comparison expression (field IN (values...)).
func (s StringChain) In(values ...string) clause.Expression {
	interfaceValues := make([]any, len(values))
	for i, v := range values {
		interfaceValues[i] = v
	}
	return clause.IN{Column: s.column, Values: interfaceValues}
}

// NotIn creates a NOT IN comparison expression (field NOT IN (values...)).
func (s StringChain) NotIn(values ...string) clause.Expression {
	interfaceValues := make([]any, len(values))
	for i, v := range values {
		interfaceValues[i] = v
	}
	return clause.Not(clause.IN{Column: s.column, Values: interfaceValues})
}

// IsNull creates a NULL check expression (field IS NULL).
func (s StringChain) IsNull() clause.Expression {
	return clause.Expr{SQL: "? IS NULL", Vars: []any{s.column}}
}

// IsNotNull creates a NOT NULL check expression (field IS NOT NULL).
func (s StringChain) IsNotNull() clause.Expression {
	return clause.Expr{SQL: "? IS NOT NULL", Vars: []any{s.column}}
}

// Concat creates a string concatenation expression.
func (s StringChain) Concat(value string) StringChain {
	return StringChain{column: clause.Expr{SQL: "CONCAT(?, ?)", Vars: []any{s.column, value}}}
}

// ConcatExpr creates a string concatenation expression with another expression.
func (s StringChain) ConcatExpr(expr any) StringChain {
	if col, ok := expr.(ColumnInterface); ok {
		expr = col.Column()
	}
	return StringChain{column: clause.Expr{SQL: "CONCAT(?, ?)", Vars: []any{s.column, expr}}}
}

// Length creates a string length expression.
func (s StringChain) Length() NumberChain[int] {
	return NumberChain[int]{column: clause.Expr{SQL: "LENGTH(?)", Vars: []any{s.column}}}
}

// Upper creates an uppercase conversion expression.
func (s StringChain) Upper() StringChain {
	return StringChain{column: clause.Expr{SQL: "UPPER(?)", Vars: []any{s.column}}}
}

// Lower creates a lowercase conversion expression.
func (s StringChain) Lower() StringChain {
	return StringChain{column: clause.Expr{SQL: "LOWER(?)", Vars: []any{s.column}}}
}

// Trim creates a whitespace trimming expression.
func (s StringChain) Trim() StringChain {
	return StringChain{column: clause.Expr{SQL: "TRIM(?)", Vars: []any{s.column}}}
}

// Left creates a left substring expression.
func (s StringChain) Left(length int) StringChain {
	return StringChain{column: clause.Expr{SQL: "LEFT(?, ?)", Vars: []any{s.column, length}}}
}

// Right creates a right substring expression.
func (s StringChain) Right(length int) StringChain {
	return StringChain{column: clause.Expr{SQL: "RIGHT(?, ?)", Vars: []any{s.column, length}}}
}

// Substring creates a substring expression.
func (s StringChain) Substring(start, length int) StringChain {
	return StringChain{column: clause.Expr{SQL: "SUBSTRING(?, ?, ?)", Vars: []any{s.column, start, length}}}
}
