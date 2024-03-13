package tests

import (
	"github.com/peter-mount/go-script/executor"
	"github.com/peter-mount/go-script/parser"
	_ "github.com/peter-mount/go-script/stdlib/fmt"
	"testing"
)

type forRangeIterator struct {
	Value int
	End   int
	Inc   int
}

func (i *forRangeIterator) HasNext() bool {
	return i.Value != i.End
}

func (i *forRangeIterator) Next() interface{} {
	if !i.HasNext() {
		panic("Next() without HasNext()")
	}
	r := i.Value
	i.Value = i.Value + i.Inc
	return r
}

func Test_forrange(t *testing.T) {

	tests := []struct {
		name           string
		script         string
		params         map[string]interface{}
		initialResult  interface{}
		expectedResult interface{}
	}{
		// ===============
		// maps
		//
		// Note these tests should return key "k4"
		// we can't look for the latter keys as the order can vary on each run
		// due to how maps work internally
		// ===============
		{
			name:   "map int",
			script: `main() { for k, v := range it { if k == "k4" result = v } }`,
			params: map[string]interface{}{
				"it": map[interface{}]interface{}{
					"k1": 1, "k2": 2, "k3": 3, "k4": 4,
					"k5": 5, "k6": 6, "k7": 7, "k8": 8,
				},
			},
			initialResult:  -1,
			expectedResult: 4,
		},
		{
			name:   "map int64",
			script: `main() { for k, v := range it { if k == "k4" result = v } }`,
			params: map[string]interface{}{
				"it": map[interface{}]interface{}{
					"k1": int64(1), "k2": int64(2), "k3": int64(3), "k4": int64(4),
					"k5": int64(5), "k6": int64(6), "k7": int64(7), "k8": int64(8),
				},
			},
			initialResult:  -1,
			expectedResult: int64(4),
		},
		{
			name:   "map float64",
			script: `main() { for k, v := range it { if k == "k4" result = v } }`,
			params: map[string]interface{}{
				"it": map[interface{}]interface{}{
					"k1": float64(1), "k2": float64(2), "k3": float64(3), "k4": float64(4),
					"k5": float64(5), "k6": float64(6), "k7": float64(7), "k8": float64(8),
				},
			},
			initialResult:  -1.0,
			expectedResult: 4.0,
		},
		{
			name:   "map float32",
			script: `main() { for k, v := range it { if k == "k4" result = v } }`,
			params: map[string]interface{}{
				"it": map[interface{}]interface{}{
					"k1": float32(1), "k2": float32(2), "k3": float32(3), "k4": float32(4),
					"k5": float32(5), "k6": float32(6), "k7": float32(7), "k8": float32(8),
				},
			},
			initialResult:  -1.0,
			expectedResult: float32(4.0),
		},
		{
			name:   "map string",
			script: `main() { for k, v := range it { if k == "k4" result = v } }`,
			params: map[string]interface{}{
				"it": map[interface{}]interface{}{
					"k1": "str1", "k2": "str2", "k3": "str3", "k4": "str4",
					"k5": "str5", "k6": "str6", "k7": "str7", "k8": "str8",
				},
			},
			initialResult:  "",
			expectedResult: "str4",
		},
		{
			name:   "map break",
			script: `main() { for k, v := range it { result = v if k == "k3" break } }`,
			params: map[string]interface{}{
				"it": map[interface{}]interface{}{
					"k1": "str1", "k2": "str2", "k3": "str3", "k4": "str4",
					"k5": "str5", "k6": "str6", "k7": "str7", "k8": "str8",
				},
			},
			initialResult:  "",
			expectedResult: "str3",
		},
		{
			name:   "map continue",
			script: `main() { for k, v := range it { if k != "k6" continue result = v } }`,
			params: map[string]interface{}{
				"it": map[interface{}]interface{}{
					"k1": "str1", "k2": "str2", "k3": "str3", "k4": "str4",
					"k5": "str5", "k6": "str6", "k7": "str7", "k8": "str8",
				},
			},
			initialResult:  "",
			expectedResult: "str6",
		},

		// ===============
		// arrays
		// ===============
		{
			name:   "array int",
			script: `main() { for i,v := range it { result = v } }`,
			params: map[string]interface{}{
				"it": []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			},
			initialResult:  -1,
			expectedResult: 9,
		},
		{
			name:   "array int64",
			script: `main() { for i,v := range it { result = v } }`,
			params: map[string]interface{}{
				"it": []int64{1, 2, 3, 4, 5, 6, 7, 8, 9},
			},
			initialResult:  -1,
			expectedResult: int64(9),
		},
		{
			name:   "array float64",
			script: `main() { for i,v := range it { result = v } }`,
			params: map[string]interface{}{
				"it": []float64{1, 2, 3, 4, 5, 6, 7, 8, 9},
			},
			initialResult:  -1.0,
			expectedResult: 9.0,
		},
		{
			name:   "array float32",
			script: `main() { for i,v := range it { result = v } }`,
			params: map[string]interface{}{
				"it": []float32{1, 2, 3, 4, 5, 6, 7, 8, 9},
			},
			initialResult:  -1.0,
			expectedResult: float32(9.0),
		},
		{
			name:   "array string",
			script: `main() { for i,v := range it { result = v } }`,
			params: map[string]interface{}{
				"it": []string{"str1", "str2", "str3", "str4", "str5", "str6", "str7", "str8", "str9"},
			},
			initialResult:  "",
			expectedResult: "str9",
		},

		// ===============
		// slices
		// ===============
		{
			name:   "slice int",
			script: `main() { for i,v := range it { result = v } }`,
			params: map[string]interface{}{
				"it": append([]int{}, 1, 2, 3, 4, 5, 6, 7, 8, 9),
			},
			initialResult:  -1,
			expectedResult: 9,
		},
		{
			name:   "slice int64",
			script: `main() { for i,v := range it { result = v } }`,
			params: map[string]interface{}{
				"it": append([]int64{1}, 2, 3, 4, 5, 6, 7, 8, 9),
			},
			initialResult:  -1,
			expectedResult: int64(9),
		},
		{
			name:   "slice float64",
			script: `main() { for i,v := range it { result = v } }`,
			params: map[string]interface{}{
				"it": append([]float64{}, 1, 2, 3, 4, 5, 6, 7, 8, 9),
			},
			initialResult:  -1.0,
			expectedResult: 9.0,
		},
		{
			name:   "slice float32",
			script: `main() { for i,v := range it { result = v } }`,
			params: map[string]interface{}{
				"it": append([]float32{}, 1, 2, 3, 4, 5, 6, 7, 8, 9),
			},
			initialResult:  -1.0,
			expectedResult: float32(9.0),
		},
		{
			name:   "slice string",
			script: `main() { for i,v := range it { result = v } }`,
			params: map[string]interface{}{
				"it": append([]string{}, "str1", "str2", "str3", "str4", "str5", "str6", "str7", "str8", "str9"),
			},
			initialResult:  "",
			expectedResult: "str9",
		},

		// ===============
		// slices with break
		// ===============
		{
			name:   "slice break int",
			script: `main() { for i,v := range it { result = v if i>5 break } }`,
			params: map[string]interface{}{
				"it": append([]int{}, 1, 2, 3, 4, 5, 6, 7, 8, 9),
			},
			initialResult:  -1,
			expectedResult: 7,
		},
		{
			name:   "slice break int64",
			script: `main() { for i,v := range it { result = v if i>5 break } }`,
			params: map[string]interface{}{
				"it": append([]int64{1}, 2, 3, 4, 5, 6, 7, 8, 9),
			},
			initialResult:  -1,
			expectedResult: int64(7),
		},
		{
			name:   "slice break float64",
			script: `main() { for i,v := range it { result = v if i>5 break } }`,
			params: map[string]interface{}{
				"it": append([]float64{}, 1, 2, 3, 4, 5, 6, 7, 8, 9),
			},
			initialResult:  -1.0,
			expectedResult: 7.0,
		},
		{
			name:   "slice break float32",
			script: `main() { for i,v := range it { result = v if i>5 break } }`,
			params: map[string]interface{}{
				"it": append([]float32{}, 1, 2, 3, 4, 5, 6, 7, 8, 9),
			},
			initialResult:  -1.0,
			expectedResult: float32(7.0),
		},
		{
			name:   "slice break string",
			script: `main() { for i,v := range it { result = v if i>5 break } }`,
			params: map[string]interface{}{
				"it": append([]string{}, "str1", "str2", "str3", "str4", "str5", "str6", "str7", "str8", "str9"),
			},
			initialResult:  "",
			expectedResult: "str7",
		},

		// ===============
		// slices with continue
		// ===============
		{
			name:   "slice continue int",
			script: `main() { for i,v := range it { if i>5 continue result = v } }`,
			params: map[string]interface{}{
				"it": append([]int{}, 1, 2, 3, 4, 5, 6, 7, 8, 9),
			},
			initialResult:  -1,
			expectedResult: 6,
		},
		{
			name:   "slice continue int64",
			script: `main() { for i,v := range it { if i>5 continue result = v } }`,
			params: map[string]interface{}{
				"it": append([]int64{1}, 2, 3, 4, 5, 6, 7, 8, 9),
			},
			initialResult:  -1,
			expectedResult: int64(6),
		},
		{
			name:   "slice continue float64",
			script: `main() { for i,v := range it { if i>5 continue result = v } }`,
			params: map[string]interface{}{
				"it": append([]float64{}, 1, 2, 3, 4, 5, 6, 7, 8, 9),
			},
			initialResult:  -1.0,
			expectedResult: 6.0,
		},
		{
			name:   "slice continue float32",
			script: `main() { for i,v := range it { if i>5 continue result = v } }`,
			params: map[string]interface{}{
				"it": append([]float32{}, 1, 2, 3, 4, 5, 6, 7, 8, 9),
			},
			initialResult:  -1.0,
			expectedResult: float32(6.0),
		},
		{
			name:   "slice continue string",
			script: `main() { for i,v := range it { if i>5 continue result = v } }`,
			params: map[string]interface{}{
				"it": append([]string{}, "str1", "str2", "str3", "str4", "str5", "str6", "str7", "str8", "str9"),
			},
			initialResult:  "",
			expectedResult: "str6",
		},

		// ===============
		// strings
		// ===============
		{
			name:   "string",
			script: `main() { for i,v := range it { result = v } }`,
			params: map[string]interface{}{
				"it": "hello world!",
			},
			initialResult:  0,
			expectedResult: uint8('!'),
		},

		// ===============
		// iterator
		// ===============
		{
			name:   "iterator ascending",
			script: `main() { for i,v := range it { result = v } }`,
			params: map[string]interface{}{
				"it": &forRangeIterator{Value: 0, End: 10, Inc: 1},
			},
			initialResult:  -1,
			expectedResult: 9,
		},
		{
			name:   "iterator descending",
			script: `main() { for i,v := range it { result = v } }`,
			params: map[string]interface{}{
				"it": &forRangeIterator{Value: 10, End: 0, Inc: -1},
			},
			initialResult:  11,
			expectedResult: 1,
		},
		{
			name:   "iterator break",
			script: `main() { for i,v := range it { if i>5 break result = v } }`,
			params: map[string]interface{}{
				"it": &forRangeIterator{Value: 10, End: 20, Inc: 1},
			},
			initialResult:  -1,
			expectedResult: 15,
		},
		{
			name:   "iterator continue",
			script: `main() { for i,v := range it { if i>5 continue result = v } }`,
			params: map[string]interface{}{
				"it": &forRangeIterator{Value: 10, End: 20, Inc: 1},
			},
			initialResult:  -1,
			expectedResult: 15,
		},

		// ===============
		// integer ranges
		// ===============
		{
			name:   "integer int",
			script: `main() { for i,_ := range it { result = i } }`,
			params: map[string]interface{}{
				"it": 10,
			},
			initialResult:  -1,
			expectedResult: 9,
		},
		{
			name:   "integer break int",
			script: `main() { for i,_ := range it { result = i if i>5 break } }`,
			params: map[string]interface{}{
				"it": 10,
			},
			initialResult:  -1,
			expectedResult: 6,
		},
		{
			name:   "integer int64",
			script: `main() { for i,v := range it { result = i } }`,
			params: map[string]interface{}{
				"it": int64(10),
			},
			initialResult:  -1,
			expectedResult: 9,
		},
		{
			name:   "integer break int64",
			script: `main() { for i,v := range it { result = i if i>5 break } }`,
			params: map[string]interface{}{
				"it": int64(10),
			},
			initialResult:  -1,
			expectedResult: 6,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			p, err := parser.New().ParseString(test.name, test.script)
			if err != nil {
				t.Fatal(err)
				return
			}

			exec, err := executor.New(p)
			if err != nil {
				t.Fatal(err)
				return
			}

			// Add each resource to the global scope
			globals := exec.GlobalScope()

			globals.Declare("result")
			globals.Set("result", test.initialResult)

			if test.params != nil {
				for k, v := range test.params {
					globals.Declare(k)
					globals.Set(k, v)
				}
			}

			// Ignore errors as we test for them
			err = exec.Run()
			if err != nil {
				t.Fatal(err)
				return
			}

			result, ok := globals.Get("result")
			if !ok {
				t.Errorf("result not returned")
			} else {
				if result != test.expectedResult {
					t.Errorf("expected %v %T got %v %T", test.expectedResult, test.expectedResult, result, result)
				}
			}

		})
	}

}
