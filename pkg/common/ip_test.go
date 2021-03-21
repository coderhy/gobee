package common

import (
	"reflect"
	"testing"
)

//TestIP2Long 测试ip2long
func TestIP2long(t *testing.T) {
	type IPlist struct {
		input string
		want  int64
	}

	testMaps := map[string]IPlist{
		"yl_1": {
			input: "192.168.1.100",
			want:  3232235876,
		},
		"yl_2": {
			input: "8.8.8.8",
			want:  134744072,
		},
	}

	for name, val := range testMaps {
		t.Run(name, func(t *testing.T) {
			res := IP2long(val.input)
			if !reflect.DeepEqual(res, val.want) {
				t.Errorf("期望结果：%#v, 实际结果： %#v", val.want, res)
			}
		})
	}
}

func TestLong2IP(t *testing.T) {
	res := Long2IP(3232235876)
	val := "192.168.1.100"
	if val != res {
		t.Errorf("期望结果：%#v, 实际结果： %#v", val, res)
	}
}
