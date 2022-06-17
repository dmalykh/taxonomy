package cursor

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"time"
)

func Marshal[V string | time.Time | uint | int64](value V) string {
	buf := new(bytes.Buffer)

	if err := gob.NewEncoder(buf).Encode(value); err != nil {
		// @TODO log. Impossible situation, because we have full control on type for value
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes())
}

// Unmarshal decoder value string to cursor.
func Unmarshal(value string, v any) error {
	decoded, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return fmt.Errorf(`error to unmarshal cursor %q %w`, value, err)
	}

	if err := gob.NewDecoder(bytes.NewBuffer(decoded)).Decode(v); err != nil {
		return fmt.Errorf(`error to decode cursor %q %w`, value, err)
	}

	return nil
}
