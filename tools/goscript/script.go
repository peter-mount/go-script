package goscript

import (
	"flag"
	"fmt"
	"github.com/peter-mount/go-script/executor"
	"github.com/peter-mount/go-script/parser"
	"math"
)

type Script struct {
}

func (b *Script) Run() error {
	p := parser.New()

	for _, fileName := range flag.Args() {
		s, err := p.ParseFile(fileName)
		if err != nil {
			return err
		}

		exec, err := executor.New(s)
		if err != nil {
			return err
		}

		globals := exec.GlobalScope()

		globals.Declare("testMap")
		globals.Set("testMap",
			map[string]interface{}{
				"a": 3.1415926,
				"b": map[string]interface{}{
					"c": 42,
					"d": &TestStruct{
						C: math.E,
					},
				},
				"42": "Answer to life",
			},
		)

		globals.Declare("testStruct")
		globals.Set("testStruct", &TestStruct{
			A: TestStruct2{C: 10},
			B: &TestStruct{
				A: TestStruct2{
					C: 15,
				},
				B: nil,
				C: 20,
			},
			C: 42,
		})

		globals.Declare("testSlice")
		globals.Set("testSlice", []string{
			"entry1",
			"entry2",
			"entry3",
			"entry4",
		})

		globals.Declare("testSlice2")
		globals.Set("testSlice2", []interface{}{
			TestStruct2{C: 10},
			&TestStruct{
				A: TestStruct2{
					C: 15,
				},
				B: nil,
				C: 20,
			},
			TestStruct2{C: 42},
		})

		for i := 0; i < 3; i++ {
			n := fmt.Sprintf("testcl%d", i)
			globals.Declare(n)
			globals.Set(n, &TestClosable{Id: n})
		}

		err = exec.Run()
		if err != nil {
			return err
		}
	}

	return nil
}

type TestStruct struct {
	A TestStruct2
	B *TestStruct
	C float64
}

type TestStruct2 struct {
	B *TestStruct
	C float64
}

type TestClosable struct {
	Id string
}

func (t *TestClosable) Close() {
	fmt.Printf("**** Close %q ****\n", t.Id)
}
