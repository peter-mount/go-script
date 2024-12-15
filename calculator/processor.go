package calculator

import "fmt"

var (
	Swap = swap{}
	Dup  = dup{}
	Drop = drop{}
	Over = over{}
	Rot  = rot{}
	Dump = dump{}
)

func (c *calculator) Process(instructions ...Instruction) error {
	for _, i := range instructions {
		if err := i.Invoke(c); err != nil {
			return err
		}
	}
	return nil
}

type Instruction interface {
	Invoke(c Calculator) error
}

type dump struct{}

func (d dump) Invoke(c Calculator) error {
	fmt.Println(c.Dump())
	return nil
}

func Push(v interface{}) Instruction { return push{v: v} }

type push struct {
	v interface{}
}

func (p push) Invoke(c Calculator) error {
	c.Push(p.v)
	return nil
}

type swap struct{}

func (p swap) Invoke(c Calculator) error {
	return c.Swap()
}

type dup struct{}

func (p dup) Invoke(c Calculator) error {
	return c.Dup()
}

type drop struct{}

func (p drop) Invoke(c Calculator) error {
	return c.Drop()
}

type over struct{}

func (p over) Invoke(c Calculator) error {
	return c.Over()
}

type rot struct{}

func (p rot) Invoke(c Calculator) error {
	return c.Rot()
}

func Op1(op string) Instruction { return op1{op: op} }

type op1 struct {
	op string
}

func (p op1) Invoke(c Calculator) error {
	return c.Op1(p.op)
}

func Op2(op string) Instruction { return op2{op: op} }

type op2 struct {
	op string
}

func (p op2) Invoke(c Calculator) error {
	return c.Op2(p.op)
}
