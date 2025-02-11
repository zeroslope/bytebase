package mysql

// Framework code is generated by the generator.

import (
	"fmt"

	"github.com/pingcap/tidb/parser/ast"

	"github.com/bytebase/bytebase/plugin/advisor"
	"github.com/bytebase/bytebase/plugin/advisor/db"
)

var (
	_ advisor.Advisor = (*InsertUpdateNoOrderByAdvisor)(nil)
	_ ast.Visitor     = (*insertUpdateNoOrderByChecker)(nil)
)

func init() {
	advisor.Register(db.MySQL, advisor.MySQLInsertUpdateNoOrderBy, &InsertUpdateNoOrderByAdvisor{})
	advisor.Register(db.TiDB, advisor.MySQLInsertUpdateNoOrderBy, &InsertUpdateNoOrderByAdvisor{})
}

// InsertUpdateNoOrderByAdvisor is the advisor checking for no ORDER BY clause in INSERT/UPDATE statement.
type InsertUpdateNoOrderByAdvisor struct {
}

// Check checks for no ORDER BY clause in INSERT/UPDATE statement.
func (*InsertUpdateNoOrderByAdvisor) Check(ctx advisor.Context, statement string) ([]advisor.Advice, error) {
	stmtList, errAdvice := parseStatement(statement, ctx.Charset, ctx.Collation)
	if errAdvice != nil {
		return errAdvice, nil
	}

	level, err := advisor.NewStatusBySQLReviewRuleLevel(ctx.Rule.Level)
	if err != nil {
		return nil, err
	}
	checker := &insertUpdateNoOrderByChecker{
		level: level,
		title: string(ctx.Rule.Type),
	}

	for _, stmt := range stmtList {
		checker.text = stmt.Text()
		checker.line = stmt.OriginTextPosition()
		(stmt).Accept(checker)
	}

	if len(checker.adviceList) == 0 {
		checker.adviceList = append(checker.adviceList, advisor.Advice{
			Status:  advisor.Success,
			Code:    advisor.Ok,
			Title:   "OK",
			Content: "",
		})
	}
	return checker.adviceList, nil
}

type insertUpdateNoOrderByChecker struct {
	adviceList []advisor.Advice
	level      advisor.Status
	title      string
	text       string
	line       int
}

// Enter implements the ast.Visitor interface.
func (checker *insertUpdateNoOrderByChecker) Enter(in ast.Node) (ast.Node, bool) {
	code := advisor.Ok
	switch node := in.(type) {
	case *ast.UpdateStmt:
		if node.Order != nil {
			code = advisor.UpdateUseOrderBy
		}
	case *ast.InsertStmt:
		if useOrderBy(node) {
			code = advisor.InsertUseOrderBy
		}
	}

	if code != advisor.Ok {
		checker.adviceList = append(checker.adviceList, advisor.Advice{
			Status:  checker.level,
			Code:    code,
			Title:   checker.title,
			Content: fmt.Sprintf("ORDER BY clause is forbidden in INSERT and UPDATE statement, but \"%s\" uses", checker.text),
			Line:    checker.line,
		})
	}
	return in, false
}

// Leave implements the ast.Visitor interface.
func (*insertUpdateNoOrderByChecker) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

func useOrderBy(node *ast.InsertStmt) bool {
	if node.Select != nil {
		switch stmt := node.Select.(type) {
		case *ast.SelectStmt:
			return stmt.OrderBy != nil
		case *ast.SetOprStmt:
			return stmt.OrderBy != nil
		}
	}
	return false
}
