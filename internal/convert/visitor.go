package convert

import (
	"fmt"
	"go/ast"
	"go/token"
)

type extent struct {
	startOffset int
	startLine   int
	startCol    int
	endOffset   int
	endLine     int
	endCol      int
}

// StmtExtent 按文件和位置描述语句在源中的范围。
type StmtExtent extent

// FuncExtent 按文件和位置描述函数在源中的范围。
type FuncExtent struct {
	extent
	name  string
	stmts []*StmtExtent
}

// FuncVisitor 实现了为文件构建函数位置列表的访问者。
type FuncVisitor struct {
	fset  *token.FileSet
	funcs []*FuncExtent
}

// Visit 实现了 ast.Visitor 接口。
func (v *FuncVisitor) Visit(node ast.Node) ast.Visitor {
	var body *ast.BlockStmt
	var name string
	switch n := node.(type) {
	case *ast.FuncLit:
		body = n.Body
	case *ast.FuncDecl:
		body = n.Body
		name = v.functionName(n)
	}
	if body != nil {
		start := v.fset.Position(node.Pos())
		end := v.fset.Position(node.End())
		if name == "" {
			name = fmt.Sprintf("@%d:%d", start.Line, start.Column)
		}
		fe := &FuncExtent{
			name: name,
			extent: extent{
				startOffset: start.Offset,
				startLine:   start.Line,
				startCol:    start.Column,
				endOffset:   end.Offset,
				endLine:     end.Line,
				endCol:      end.Column,
			},
		}
		v.funcs = append(v.funcs, fe)
		sv := StmtVisitor{fset: v.fset, function: fe}
		sv.VisitStmt(body)
	}
	return v
}

func (v *FuncVisitor) functionName(f *ast.FuncDecl) string {
	name := f.Name.Name
	if f.Recv == nil {
		return name
	} else {
		// 函数名前面有“T”。如果有接收者，其中 T 是接收者的类型，如果它是指针则取消引用。
		return v.exprName(f.Recv.List[0].Type) + "." + name
	}
}

func (v *FuncVisitor) exprName(x ast.Expr) string {
	switch y := x.(type) {
	case *ast.StarExpr:
		return v.exprName(y.X)
	case *ast.IndexExpr:
		return fmt.Sprintf("%s[%s]", v.exprName(y.X), v.exprName(y.Index))
	case *ast.Ident:
		return y.Name
	default:
		return ""
	}
}

type StmtVisitor struct {
	fset     *token.FileSet
	function *FuncExtent
}

func (v *StmtVisitor) VisitStmt(s ast.Stmt) {
	var statements *[]ast.Stmt
	switch s := s.(type) {
	case *ast.BlockStmt:
		statements = &s.List
	case *ast.CaseClause:
		statements = &s.Body
	case *ast.CommClause:
		statements = &s.Body
	case *ast.ForStmt:
		if s.Init != nil {
			v.VisitStmt(s.Init)
		}
		if s.Post != nil {
			v.VisitStmt(s.Post)
		}
		v.VisitStmt(s.Body)
	case *ast.IfStmt:
		if s.Init != nil {
			v.VisitStmt(s.Init)
		}
		v.VisitStmt(s.Body)
		if s.Else != nil {
			// 从 go.tools/cmd/convert 复制的代码，用于处理“if x {} else if y {}
			const backupToElse = token.Pos(len("else ")) // AST 不记得 else 的位置。我们可以做出准确的预测。
			switch stmt := s.Else.(type) {
			case *ast.IfStmt:
				block := &ast.BlockStmt{
					Lbrace: stmt.If - backupToElse, // 所以被覆盖的部分看起来像是从“else”开始的。
					List:   []ast.Stmt{stmt},
					Rbrace: stmt.End(),
				}
				s.Else = block
			case *ast.BlockStmt:
				stmt.Lbrace -= backupToElse // 所以这个块看起来像是从“else”开始的。
			default:
				panic("unexpected node type in if")
			}
			v.VisitStmt(s.Else)
		}
	case *ast.LabeledStmt:
		v.VisitStmt(s.Stmt)
	case *ast.RangeStmt:
		v.VisitStmt(s.Body)
	case *ast.SelectStmt:
		v.VisitStmt(s.Body)
	case *ast.SwitchStmt:
		if s.Init != nil {
			v.VisitStmt(s.Init)
		}
		v.VisitStmt(s.Body)
	case *ast.TypeSwitchStmt:
		if s.Init != nil {
			v.VisitStmt(s.Init)
		}
		v.VisitStmt(s.Assign)
		v.VisitStmt(s.Body)
	}
	if statements == nil {
		return
	}
	for i := 0; i < len(*statements); i++ {
		s := (*statements)[i]
		switch s.(type) {
		case *ast.CaseClause, *ast.CommClause, *ast.BlockStmt:
			break
		default:
			start, end := v.fset.Position(s.Pos()), v.fset.Position(s.End())
			se := &StmtExtent{
				startOffset: start.Offset,
				startLine:   start.Line,
				startCol:    start.Column,
				endOffset:   end.Offset,
				endLine:     end.Line,
				endCol:      end.Column,
			}
			v.function.stmts = append(v.function.stmts, se)
		}
		v.VisitStmt(s)
	}
}
