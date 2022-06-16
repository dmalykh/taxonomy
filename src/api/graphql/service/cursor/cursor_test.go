package cursor

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestEncodeDecode(t *testing.T) {

	type args struct {
		value   any
		entropy any
	}

	type testCase struct {
		name   string
		encode func() string
		need   any
		err    error
	}
	var timeExample = time.Now()
	tests := []testCase{
		{
			name: `Correct simple encoding and decoding string value`,
			encode: func() string {
				return Marshal("hello")
			},
			need: "hello",
		},
		{
			name: `Correct simple encoding and decoding long string value`,
			encode: func() string {
				return Marshal("OI*H(PG&fo6dc8yvp97g8P&F6oc8yvp9g87fo6DCYuviu9gp7f86odctuyviu9gp7f8o6cyuviu97gf8o6cuyvig9p7f8o6cyuvig7f8o6cyuvg97f8o6cyvg7f8o6cvy7g8f")
			},
			need: "OI*H(PG&fo6dc8yvp97g8P&F6oc8yvp9g87fo6DCYuviu9gp7f86odctuyviu9gp7f8o6cyuviu97gf8o6cuyvig9p7f8o6cyuvig7f8o6cyuvg97f8o6cyvg7f8o6cvy7g8f",
		},
		{
			name: `Correct simple encoding and decoding time value`,
			encode: func() string {
				return Marshal(timeExample)
			},
			need: timeExample,
		},
		{
			name: `Correct simple encoding and decoding int64 value`,
			encode: func() string {
				return Marshal(int64(433245362531))
			},
			need: int64(433245362531),
		},
		{
			name: `Correct simple encoding and decoding uint value`,
			encode: func() string {
				return Marshal(uint(433245362531))
			},
			need: uint(433245362531),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var encoded = tt.encode()
			assert.NotEmpty(t, encoded)

			switch tt.need.(type) {
			case string:
				var v string
				assert.NoError(t, Unmarshal(encoded, &v))
				assert.Equal(t, tt.need, v)
				break
			case int64:
				var v int64
				assert.NoError(t, Unmarshal(encoded, &v))
				assert.Equal(t, tt.need, v)
				break
			case uint:
				var v uint
				assert.NoError(t, Unmarshal(encoded, &v))
				assert.Equal(t, tt.need, v)
				break
			case time.Time:
				var v time.Time
				assert.NoError(t, Unmarshal(encoded, &v))
				assert.True(t, tt.need.(time.Time).Equal(v))
				//assert.Equal(t, tt.need.(time.Time).String(), v.String())
				break
			default:
				panic("Unknown type of tt.need")
				break
			}

		})
	}
}
