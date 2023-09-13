package stdlib

import (
	"github.com/peter-mount/go-script/calculator"
	"github.com/peter-mount/go-script/errors"
	"github.com/peter-mount/go-script/executor"
	"github.com/peter-mount/go-script/script"
)

func _map(e executor.Executor, call *script.CallFunc) error {
	a, err := executor.Args(e, call)
	if err != nil {
		return errors.Error(call.Pos, err)
	}

	m := make(map[string]interface{})

	for i, v := range a {
		kv, ok := calculator.GetKeyValue(v)
		if !ok {
			return errors.Errorf(call.Parameters.Args[i].Pos, "expected Key:Value pair")
		}

		m[kv.Key()] = kv.Value()
	}

	e.Calculator().Push(m)
	return nil
}
