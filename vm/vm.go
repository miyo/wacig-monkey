package vm

import (
	"fmt"
	"monkey/code"
	"monkey/compiler"
	"monkey/object"
)

const StackSize = 2048
const GlobalSize = 65536

var True = &object.Boolean{Value: true}
var False = &object.Boolean{Value: false}
var Null = &object.Null{}

type VM struct {
	constants    []object.Object
	instructions code.Instructions
	stack        []object.Object
	sp           int
	globals      []object.Object
}

func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		instructions: bytecode.Instructions,
		constants:    bytecode.Constants,
		stack:        make([]object.Object, StackSize),
		sp:           0,
		globals:      make([]object.Object, GlobalSize),
	}
}

func NewWithState(bytecode *compiler.Bytecode, g []object.Object) *VM {
	vm := New(bytecode)
	vm.globals = g
	return vm
}

func (vm *VM) StackTop() object.Object {
	if vm.sp == 0 {
		return nil
	}
	return vm.stack[vm.sp-1]
}

func (vm *VM) Run() error {
	for ip := 0; ip < len(vm.instructions); ip++ {
		op := code.Opcode(vm.instructions[ip])
		switch op {
		case code.OpNull:
			err := vm.push(Null)
			if err != nil {
				return err
			}
		case code.OpConstant:
			constIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2
			err := vm.push(vm.constants[constIndex])
			if err != nil {
				return err
			}
		case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv, code.OpEqual, code.OpNotEqual, code.OpGreaterThan:
			err := vm.executeBinaryOperation(op)
			if err != nil {
				return err
			}
		case code.OpTrue:
			err := vm.push(True)
			if err != nil {
				return err
			}
		case code.OpFalse:
			err := vm.push(False)
			if err != nil {
				return err
			}
		case code.OpMinus:
			err := vm.executeMinusOperation()
			if err != nil {
				return err
			}
		case code.OpBang:
			err := vm.executeBangOperation()
			if err != nil {
				return err
			}
		case code.OpPop:
			vm.pop()
		case code.OpJump:
			pos := int(code.ReadUint16(vm.instructions[ip+1:]))
			ip = pos - 1
		case code.OpJumpNotTruthy:
			pos := int(code.ReadUint16(vm.instructions[ip+1:]))
			ip += 2
			condition := vm.pop()
			if !isTruthy(condition) {
				ip = pos - 1
			}
		case code.OpSetGlobal:
			globalIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2
			vm.globals[globalIndex] = vm.pop()
		case code.OpGetGlobal:
			globalIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2
			err := vm.push(vm.globals[globalIndex])
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (vm *VM) push(o object.Object) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}

	vm.stack[vm.sp] = o
	vm.sp++
	return nil
}

func (vm *VM) pop() object.Object {
	o := vm.stack[vm.sp-1]
	vm.sp--
	return o
}

func (vm *VM) LastPoppedStackElem() object.Object {
	return vm.stack[vm.sp]
}

func (vm *VM) executeBinaryOperation(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()
	rightT := right.Type()
	leftT := left.Type()
	if leftT == object.INTEGER_OBJ && rightT == object.INTEGER_OBJ {
		return vm.executeBinaryIntegerOperation(op, left, right)
	} else if leftT == object.BOOLEAN_OBJ && rightT == object.BOOLEAN_OBJ {
		return vm.executeBinaryBooleanOperation(op, left, right)
	}
	return fmt.Errorf("unsupported types for binary operation: %s %s", leftT, rightT)
}

func (vm *VM) executeBinaryIntegerOperation(op code.Opcode, left object.Object, right object.Object) error {
	rightValue := right.(*object.Integer).Value
	leftValue := left.(*object.Integer).Value
	switch op {
	case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv:
		var result int64
		switch op {
		case code.OpAdd:
			result = leftValue + rightValue
		case code.OpSub:
			result = leftValue - rightValue
		case code.OpMul:
			result = leftValue * rightValue
		case code.OpDiv:
			result = leftValue / rightValue
		}
		err := vm.push(&object.Integer{Value: result})
		return err
	case code.OpEqual, code.OpNotEqual, code.OpGreaterThan:
		var result bool
		switch op {
		case code.OpEqual:
			result = leftValue == rightValue
		case code.OpNotEqual:
			result = leftValue != rightValue
		case code.OpGreaterThan:
			result = leftValue > rightValue
		}
		if result {
			return vm.push(True)
		} else {
			return vm.push(False)
		}
	default:
		return fmt.Errorf("unknown intger operator: %d", op)
	}
}

func (vm *VM) executeBinaryBooleanOperation(op code.Opcode, left object.Object, right object.Object) error {
	rightValue := right.(*object.Boolean).Value
	leftValue := left.(*object.Boolean).Value
	var result bool
	switch op {
	case code.OpEqual:
		result = leftValue == rightValue
	case code.OpNotEqual:
		result = leftValue != rightValue
	default:
		return fmt.Errorf("unknown boolean operator: %d", op)
	}
	if result {
		return vm.push(True)
	} else {
		return vm.push(False)
	}
}

func (vm *VM) executeMinusOperation() error {
	v := vm.pop()
	vv, ok := v.(*object.Integer)
	if !ok {
		return fmt.Errorf("unsupported types for minus operation: %s", v.Type())
	}
	vm.push(&object.Integer{Value: -vv.Value})
	return nil
}

func (vm *VM) executeBangOperation() error {
	v := vm.pop()
	switch v := v.(type) {
	case *object.Boolean:
		if v.Value {
			vm.push(False)
		} else {
			vm.push(True)
		}
	case *object.Null:
		vm.push(True)
	default:
		vm.push(False)
	}
	return nil
}

func isTruthy(obj object.Object) bool {
	switch obj := obj.(type) {
	case *object.Boolean:
		return obj.Value
	case *object.Null:
		return false
	default:
		return true
	}
}
