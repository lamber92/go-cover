package convert

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"

	"github.com/lamber92/go-cover/internal/metadata"
	"golang.org/x/tools/cover"
)

type converter struct {
	packages map[string]*metadata.Package
}

// statement metadata.Statement 的包装器
type statement struct {
	*metadata.Statement
	*StmtExtent
}

// convertProfile 转换 profile 文件内容
func (c *converter) convertProfile(packages packagesCache, p *cover.Profile) error {
	file, pkgPath, err := c.findFile(packages, p.FileName)
	if err != nil {
		return err
	}
	pkg := c.packages[pkgPath]
	if pkg == nil {
		pkg = &metadata.Package{Name: pkgPath}
		c.packages[pkgPath] = pkg
	}

	// 查找函数和语句范围；创建相应的 convert.Functions 和 convert.Statements，
	// 并保留一个单独的 convert.Statements 片段，以便将它们与 profile 匹配。
	extents, err := c.findFuncs(file)
	if err != nil {
		return err
	}
	var stmts []statement
	for _, fe := range extents {
		f := &metadata.Function{
			Name:      fe.name,
			File:      file,
			Start:     fe.startOffset,
			End:       fe.endOffset,
			StartLine: fe.startLine,
			EndLine:   fe.endLine,
		}
		for _, se := range fe.stmts {
			s := statement{
				Statement: &metadata.Statement{
					Start:     se.startOffset,
					End:       se.endOffset,
					StartLine: se.startLine,
					EndLine:   se.endLine,
				},
				StmtExtent: se,
			}
			f.Statements = append(f.Statements, s.Statement)
			stmts = append(stmts, s)
		}
		pkg.Functions = append(pkg.Functions, f)
	}
	// 对于文件中的每个配置文件块，找到它涵盖的语句并递增 Reached 字段。
	blocks := p.Blocks
	for _, s := range stmts {
		for i, b := range blocks {
			if b.StartLine > s.endLine || (b.StartLine == s.endLine && b.StartCol >= s.endCol) {
				// 超过语句末尾
				blocks = blocks[i:]
				break
			}
			if b.EndLine < s.startLine || (b.EndLine == s.startLine && b.EndCol <= s.startCol) {
				// 在语句开始之前
				continue
			}
			s.Reached += int64(b.Count)
			break
		}
	}
	return nil
}

// findFile 在 GOROOT、GOPATH 等中查找命名文件的位置。
func (c *converter) findFile(packages packagesCache, file string) (filename, pkgPath string, err error) {
	dir, file := filepath.Split(file)
	if dir != "" {
		dir = strings.TrimSuffix(dir, "/")
	}
	pkg, ok := packages[dir]
	if !ok {
		pkg, err = build.Import(dir, ".", build.FindOnly)
		if err != nil {
			return "", "", fmt.Errorf("can't find %q: %w", file, err)
		}
		packages[dir] = pkg
	}

	return filepath.Join(pkg.Dir, file), pkg.ImportPath, nil
}

// findFuncs 解析文件并返回一段 FuncExtent 描述符
func (c *converter) findFuncs(name string) ([]*FuncExtent, error) {
	fset := token.NewFileSet()
	parsedFile, err := parser.ParseFile(fset, name, nil, 0)
	if err != nil {
		return nil, err
	}
	visitor := &FuncVisitor{fset: fset}
	ast.Walk(visitor, parsedFile)
	return visitor.funcs, nil
}
