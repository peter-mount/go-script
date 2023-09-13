package stdlib

import (
	"context"
	"github.com/peter-mount/go-script/calculator"
	error2 "github.com/peter-mount/go-script/error"
	"github.com/peter-mount/go-script/executor"
	"github.com/peter-mount/go-script/script"
)

func _map(e executor.Executor, call *script.CallFunc, ctx context.Context) error {
	a, err := executor.Args(e, call, ctx)
	if err != nil {
		return error2.Error(call.Pos, err)
	}

	m := make(map[string]interface{})

	for i, v := range a {
		kv, ok := calculator.GetKeyValue(v)
		if !ok {
			return error2.Errorf(call.Parameters.Args[i].Pos, "expected Key:Value pair")
		}

		m[kv.Key()] = kv.Value()
	}

	e.Calculator().Push(m)
	return nil
}
