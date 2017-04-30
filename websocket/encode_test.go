package websocket

import "testing"



type testpair struct {
	input string
	output string
}

var tests = []testpair{
	{ "Hei, dette er en test", "Hei, dette er en test"},
	{"0xb00000000","0xb00000000"},
}

func Test_Encoding_Decoding(t *testing.T) {
	for _, pair := range tests {
		v := decode(encode(pair.input))
		if v != pair.output {
			t.Error(
				"For", pair.input,
				"expected", pair.output,
				"got", v,
			)
		}
	}
}
