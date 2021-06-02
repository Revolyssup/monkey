package eval

import (
	"github.com/Revolyssup/monkey/ast"
	"github.com/Revolyssup/monkey/obj"
)

//It take in the AST ,starting from the root node. And depending on the type of Node, calls other functions which evaluate and then call Eval recursively.
//Because all data type in AST implement Node interface ,ergo this works
func Eval(node ast.Node) obj.Object {
	switch node := node.(type) {
	//If it is the root node
	case *ast.Program:
		{
			return evalStatements(node.Statements)
		}
	// For each statement
	case *ast.ExpressionStatement:
		{
			return Eval(node.Expression)
		}

		//For different types of expressions
	case *ast.IntegerLiteral:
		{
			return &obj.Integer{Value: node.Value}
		}
	case *ast.Boolean:
		{
			return &obj.Boolean{Value: node.Value}
		}
	}
	return nil
}

func evalStatements(stmts []ast.Statement) obj.Object {
	var result obj.Object

	for _, stmt := range stmts {
		result = Eval(stmt)
	}
	return result // It actually consists of the result of last evaluated statement.
}
