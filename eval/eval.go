package eval

import (
	"fmt"
	"strings"

	"github.com/Revolyssup/monkey/ast"
	"github.com/Revolyssup/monkey/obj"
)

//Mapping to Builtin Functions in monkey
var fns = map[string]*obj.Builtin{
	"len": &obj.Builtin{
		Fn: length,
	},
}

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
	case *ast.StringLiteral:
		{
			return &obj.String{Value: node.Value}
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
	case *ast.FunctionLiteral:
		{
			args := node.Params
			body := node.Body
			return &obj.Function{Args: args, Body: body, Env: env}
		}
	case *ast.FunctionCall:
		{
			//Create the function object
			fn := Eval(node.Function, env)
			if isError(fn) {
				return fn
			}

			//Create allt the argument objects
			args := evalExpressions(node.Arguments, env)

			if len(args) == 1 && isError(args[0]) {
				return args[0]
			}
			return execFunction(fn, args)
		}

	}
	return nil //handled by isError
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
	case left.DataType() == obj.STRING_OBJ && right.DataType() == obj.STRING_OBJ:
		{
			return evalString(op, left, right)
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

func evalString(op string, left obj.Object, right obj.Object) obj.Object {
	leftstr := left.(*obj.String).Value
	rightstr := right.(*obj.String).Value
	switch op {
	case "+":
		{
			total := leftstr + rightstr
			return &obj.String{Value: total}
		}
	case "-":
		{
			i := strings.Index(leftstr, rightstr)
			if i == -1 {
				return newErr("No right substring found in left string")
			}
			len := len(rightstr)
			total := leftstr[0:i] + leftstr[i+len:]
			return &obj.String{Value: total}
		}
	default:
		{
			return newErr("unknown operator: %s %s %s",
				left.DataType(), op, right.DataType())
		}
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
		//If its not user defined, check if its a builtin function
		bf, ok2 := fns[node.Value]

		if !ok2 {
			return newErr("Undefined variable: %s", node.Value)
		}
		return bf

	}
	return val
}

/****************/
//To evaluate a list of expressions into monkey objects.

func evalExpressions(node []ast.Expression, env *obj.Env) []obj.Object {
	exps := []obj.Object{}

	for _, exp := range node {
		evaluated := Eval(exp, env)
		if isError(evaluated) {
			return []obj.Object{evaluated}
		}
		exps = append(exps, evaluated)
	}
	return exps
}

//This function will do two things:
//1. It will pass the outer environments to the function such that if the function doesn't find a variable in its own environment, it checks in out env Object recursively
//2. It passes the arguments given to the functions into functions's Env object.
func extendFun(fn *obj.Function, args []obj.Object) *obj.Env {
	env := obj.NewEnclosedEnvironment(fn.Env)
	for i, param := range fn.Args {
		env.Set(param.Value, args[i])
	}
	return env
}

//Will be called after function has been executed and a Return object has been recieved.
func unwrapReturnValue(ob obj.Object) obj.Object {
	if returnValue, ok := ob.(*obj.Return); ok {
		return returnValue.Value
	}
	return ob
}

//Executing the function
func execFunction(fn obj.Object, args []obj.Object) obj.Object {
	function, ok := fn.(*obj.Function)
	if !ok {
		builin, ok2 := fn.(*obj.Builtin)
		if ok2 {
			return builin.Fn(args...)
		}
		return newErr("not a function: %s", fn.DataType())
	}
	newenv := extendFun(function, args)
	evaluated := Eval(function.Body, newenv)
	return unwrapReturnValue(evaluated)
}

/***Built in functions in Monkey*****/

func length(args ...obj.Object) obj.Object {
	if len(args) != 1 {
		return newErr("wrong number of arguments. got=%d, want=1", len(args))
	}
	s, ok := args[0].(*obj.String)
	if !ok {
		return &obj.Error{ErrMsg: "No string in arguments"}
	}
	return &obj.Integer{Value: int64(len(s.Value))}
}
