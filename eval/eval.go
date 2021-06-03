package eval

import (
	"fmt"

	"github.com/Revolyssup/monkey/ast"
	"github.com/Revolyssup/monkey/obj"
)

//Because different "true" are not different so creating new instance everytime a bool instance is created is a waste of space. SO we point all booleans of one type
//to a single instance. Same for null
var (
	TRUE  = &obj.Boolean{Value: true}
	FALSE = &obj.Boolean{Value: false}
	NULL  = &obj.Null{}
)

//It take in the AST ,starting from the root node. And depending on the type of Node, calls other functions which evaluate and then call Eval recursively.
//Because all data type in AST implement Node interface ,ergo this works
func Eval(node ast.Node, env *obj.Env) obj.Object {
	switch node := node.(type) {
	//If it is the root node
	case *ast.Program:
		{
			return evalStatements(node.Statements, env)
		}
	// For each statement
	case *ast.ExpressionStatement:
		{
			return Eval(node.Expression, env)
		}

		//For different types of expressions
	case *ast.IntegerLiteral:
		{
			return &obj.Integer{Value: node.Value}
		}
	case *ast.Boolean:
		{
			return returnSingleBooleanInstance(node.Value)
		}
		//Evaluating prefix expressions
	case *ast.PrefixExpression:
		{
			evalRight := Eval(node.RightExpression, env)
			return evalPrefixExpression(node.Operator, evalRight)
		}
	case *ast.InfixExpression:
		{
			evalLeft := Eval(node.LeftExpression, env)
			evalRight := Eval(node.RightExpression, env)
			return evalInfixExpression(node.Operator, evalLeft, evalRight)
		}
	case *ast.BlockStatement:
		{
			return evalBlockStatement(node, env)
		}
	case *ast.IfExpression:
		{
			return evalIfExpression(node, env)
		}
	case *ast.ReturnStatement:
		{
			val := Eval(node.ReturnValue, env)
			return &obj.Return{Value: val}
		}
	case *ast.Identifier:
		{

			return evalIdentifiers(node, env)
		}
	case *ast.LetStatement:
		{
			val := Eval(node.Value, env)
			if isError(val) {
				return val
			}
			env.Set(node.Name.Value, val)
		}
	}
	return nil
}

func returnSingleBooleanInstance(input bool) *obj.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}
func evalStatements(stmts []ast.Statement, env *obj.Env) obj.Object {
	var result obj.Object

	for _, stmt := range stmts {
		result = Eval(stmt, env)
		//if we encounter a return statement,we have to take that result and just exit that scope.
		if rs, ok := result.(*obj.Return); ok {
			return rs.Value
		}

		if err, ok := result.(*obj.Error); ok {
			return err
		}
	}
	return result // It actually consists of the result of last evaluated statement.
}

//PREFIX
/****************/
func evalPrefixExpression(op string, right obj.Object) obj.Object {
	switch op {
	case "!":
		{
			return evalBangOperator(right)
		}
	case "-":
		{
			if right.DataType() != obj.INTEGER_OBJ { //Would be changed into switch case later. Based on datatype different methods will be called
				return newErr("unknown operator: %s%s", op, right.DataType())
			}
			rightint := right.(*obj.Integer)
			return evalMinusOperator(rightint)
		}

	default:
		{
			return newErr("unknown operator: %s %s", op, right.DataType())
		}
	}

}

//for different prefix operators.
func evalBangOperator(right obj.Object) obj.Object {
	switch right {
	case TRUE:
		{
			return FALSE
		}
	case FALSE:
		{
			return TRUE
		}
	case NULL:
		{
			return TRUE
		}
	default:
		{
			return FALSE
		}
	}
}

func evalMinusOperator(right *obj.Integer) obj.Object {
	value := right.Value
	ans := &obj.Integer{Value: -value}
	return ans
}

/*********/
//INFIX
func evalInfixExpression(op string, left obj.Object, right obj.Object) obj.Object {
	switch {
	//If we have integers on either side
	case left.DataType() == obj.INTEGER_OBJ && right.DataType() == obj.INTEGER_OBJ:
		{
			return evalInteger(op, left, right)
		}
	// Directly compare the objects and return boolean. As all booleans are pointing to a single instance, it swiftly works for booleans
	case op == "==":
		{
			return returnSingleBooleanInstance(left == right)
		}
	case op == "!=":
		{
			return returnSingleBooleanInstance(left != right)
		}
	case left.DataType() != right.DataType():
		{
			return newErr("type mismatch: %s %s %s", left.DataType(), op, right.DataType())
		}
	default:
		{
			return newErr("unknown operator: %s %s %s", left.DataType(), op, right.DataType())
		}
	}

}

func evalInteger(op string, left obj.Object, right obj.Object) obj.Object {
	leftVal := left.(*obj.Integer).Value
	rightVal := right.(*obj.Integer).Value

	switch op {
	case "+":
		{
			return &obj.Integer{Value: leftVal + rightVal}
		}
	case "-":
		{
			return &obj.Integer{Value: leftVal - rightVal}
		}
	case "*":
		{
			return &obj.Integer{Value: leftVal * rightVal}
		}
	case "/":
		{
			return &obj.Integer{Value: leftVal / rightVal}
		}
	case "<":
		{
			return returnSingleBooleanInstance(leftVal < rightVal)
		}
	case ">":
		{
			return returnSingleBooleanInstance(leftVal > rightVal)
		}
	case "==":
		{
			return returnSingleBooleanInstance(leftVal == rightVal)
		}
	case "!=":
		{
			return returnSingleBooleanInstance(leftVal != rightVal)
		}
	default:
		return newErr("unknown operator: %s %s %s",
			left.DataType(), op, right.DataType())
	}

}

/***********/
//IF

func evalIfExpression(node *ast.IfExpression, env *obj.Env) obj.Object {
	cond := Eval(node.Condition, env)
	if isTruthy((cond)) {
		return Eval(node.MainStmt, env)
	} else if node.AltStmt != nil {
		return Eval(node.AltStmt, env)
	}
	return NULL
}

func isTruthy(object obj.Object) bool {
	switch object {
	case NULL:
		{
			return false
		}
	case TRUE:
		{
			return true
		}
	case FALSE:
		{
			return false
		}
	default:
		{
			return true
		}
	}

}

/***************/
//BLOCK

func evalBlockStatement(block *ast.BlockStatement, env *obj.Env) obj.Object {
	var result obj.Object
	for _, stmt := range block.Stmts {
		result = Eval(stmt, env)

		if result != nil && result.DataType() == obj.RETURN_OBJ || result.DataType() == obj.ERROR_OBJ {
			return result
		}
	}
	return result
}

/***************/
//To return specific Error
func newErr(f string, a ...interface{}) *obj.Error {
	return &obj.Error{ErrMsg: fmt.Sprintf(f, a...)}
}

func isError(ob obj.Object) bool {
	if ob == nil {
		return false
	}
	return ob.DataType() == obj.ERROR_OBJ //If it is not nil, it has to be an error object
}

/***********/
//Identifiers
func evalIdentifiers(node *ast.Identifier, env *obj.Env) obj.Object {
	val, ok := env.Get(node.Value)
	if !ok {
		return newErr("Undefined variable: %s", node.Value)
	}
	return val
}
