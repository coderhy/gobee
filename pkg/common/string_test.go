package common

import (
	"fmt"
	"reflect"
	"testing"
)

// 一个test demo

//TestTrim 测试去空格
func TestTrim(t *testing.T) {
	type str struct {
		input string
		want  string
	}

	tests := map[string]str{
		"ltrim": {
			input: " aaa",
			want:  "aaa",
		},
		"rtrim": {
			input: "bbb ",
			want:  "bbb",
		},
		"mid_trim": {
			input: "c cc",
			want:  "ccc",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got := Trim(test.input, "\t")
			if got != test.want {
				t.Errorf("excepted:%#v, got:%#v", test.want, got)
			}
		})
	}
}

//TestStrtoupper 测试大写
func TestStrtoupper(t *testing.T) {
	type ValidRule struct {
		input string
		want  string
	}

	testMaps := map[string]ValidRule{
		"asciiChar": {
			input: "abcd",
			want:  "ABCD",
		},
		"number": {
			input: "123456abc",
			want:  "123456ABC",
		},
		"zhChar": {
			input: "汉字",
			want:  "汉字",
		},
	}

	for name, testMap := range testMaps {
		t.Run(name, func(t *testing.T) {
			res := Strtoupper(testMap.input)
			if !reflect.DeepEqual(res, testMap.want) {
				t.Logf("excepted:%v, got:%v", testMap.want, res)
				t.Errorf("excepted:%v, got:%v", testMap.want, res)
			}
		})
	}
}

//BenchmarkStrtoupper 基准测试
func BenchmarkStrtoupper(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Strtoupper("asdfasdf")
	}
}

// ExampleStrtoupper 示例函数  go test -run Example
func ExampleStrtoupper() {
	fmt.Println(Strtoupper("asdfasdfsa"))
	fmt.Println(Strtoupper("asdf23433"))

	// Output:
	// ASDFASDFSA
	// ASDF23433
}
