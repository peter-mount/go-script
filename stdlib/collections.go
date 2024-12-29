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

// _mapContains implements mapContains(map,key) which returns true if the provided map contains a key.
//
// When multiple keys are provided then this returns true if any of those keys exist.
func _mapContains(e executor.Executor, call *script.CallFunc) error {
	a, err := executor.Args(e, call)

	if err == nil && len(a) < 2 {
		err = errors.Errorf(call.Pos, "mapContains(map,key)")
	}

	var ok bool
	var m map[string]interface{}
	if err == nil {
		m, ok = a[0].(map[string]interface{})
		if !ok {
			err = errors.Errorf(call.Parameters.Args[0].Pos, "expected map")
		}
	}

	if err == nil {
		var k string
		found := false
		for _, v := range a[1:] {
			k, err = calculator.GetString(v)
			if err != nil {
				break
			}
			// pass on the first key found
			_, found = m[k]
			if found {
				break
			}
		}
		if err == nil {
			e.Calculator().Push(found)
		}
	}

	return errors.Error(call.Pos, err)
}
