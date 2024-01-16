package cursor_test

import (
	"testing"
	"time"

	"github.com/dmalykh/taxonomy/api/graphql/service/cursor"
	"github.com/stretchr/testify/assert"
)

func TestEncodeDecode(t *testing.T) {
	t.Parallel()

	timeExample := time.Now()

	tests := []struct {
		name   string
		encode func() string
		need   any
		err    error
	}{
		{
			name: `Correct simple encoding and decoding string value`,
			encode: func() string {
				return cursor.Marshal("hello")
			},
			need: "hello",
		},
		{
			name: `Correct simple encoding and decoding long string value`,
			encode: func() string {
				return cursor.Marshal(
					"OI*H(PG&g8P&F6oc8yvp9g87fo6DCYuviu9gp7f86odctuyviu9gp7f8o6cyuviu97gf8o6cuyvig9p7f8o6cyuvig7f8o6cyuvg97f8o6cyvg7f8o6cvy7g8f")
			},
			need: "OI*H(PG&g8P&F6oc8yvp9g87fo6DCYuviu9gp7f86odctuyviu9gp7f8o6cyuviu97gf8o6cuyvig9p7f8o6cyuvig7f8o6cyuvg97f8o6cyvg7f8o6cvy7g8f",
		},
		{
			name: `Correct simple encoding and decoding time value`,
			encode: func() string {
				return cursor.Marshal(timeExample)
			},
			need: timeExample,
		},
		{
			name: `Correct simple encoding and decoding int64 value`,
			encode: func() string {
				return cursor.Marshal(int64(433245362531))
			},
			need: int64(433245362531),
		},
		{
			name: `Correct simple encoding and decoding uint value`,
			encode: func() string {
				return cursor.Marshal(uint(433245362531))
			},
			need: uint(433245362531),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			encoded := tt.encode()
			assert.NotEmpty(t, encoded)

			switch val := tt.need.(type) {
			case string:
				var v string
				assert.NoError(t, cursor.Unmarshal(encoded, &v))
				assert.Equal(t, val, v)

				break
			case int64:
				var v int64
				assert.NoError(t, cursor.Unmarshal(encoded, &v))
				assert.Equal(t, val, v)

				break
			case uint:
				var v uint
				assert.NoError(t, cursor.Unmarshal(encoded, &v))
				assert.Equal(t, val, v)

				break
			case time.Time:
				var v time.Time
				assert.NoError(t, cursor.Unmarshal(encoded, &v))
				assert.True(t, val.Equal(v))

				break
			default:
				panic("Unknown type of tt.need")
			}
		})
	}
}
