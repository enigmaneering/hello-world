package std

import (
	"fmt"
	"math/big"
	"strconv"

	"git.ignitelabs.net/janos/core/sys/atlas"
)

// Stringify converts the provided value into a string.
//
// NOTE: If the value does not satisfy Stringable, this will panic.
//
// See Stringable, StringableMany, Stringify, and StringifyMany
func Stringify(value any) string {
	if value == nil {
		return ""
	}

	var out string
	switch typed := value.(type) {
	case nil:
		return ""
	case string:
		return typed
	case fmt.Stringer:
		return typed.String()
	case *big.Int:
		return typed.Text(10)
	case *big.Float:
		return typed.Text('f', int(atlas.Precision))
	case *big.Rat:
		return typed.String()
	case complex64, complex128:
		return fmt.Sprintf("%v", typed)
	case float32:
		out = strconv.FormatFloat(float64(typed), 'f', -1, 32)
	case float64:
		out = strconv.FormatFloat(typed, 'f', -1, 64)
	case uint:
		out = strconv.FormatUint(uint64(typed), 10)
	case uint8:
		out = strconv.FormatUint(uint64(typed), 10)
	case uint16:
		out = strconv.FormatUint(uint64(typed), 10)
	case uint32:
		out = strconv.FormatUint(uint64(typed), 10)
	case uint64:
		out = strconv.FormatUint(typed, 10)
	case uintptr:
		out = strconv.FormatUint(uint64(typed), 10)
	case int:
		out = strconv.FormatInt(int64(typed), 10)
	case int8:
		out = strconv.FormatInt(int64(typed), 10)
	case int16:
		out = strconv.FormatInt(int64(typed), 10)
	case int32:
		out = strconv.FormatInt(int64(typed), 10)
	case int64:
		out = strconv.FormatInt(typed, 10)
	default:
		panic(fmt.Errorf("%T is not a stringable type", typed))
	}
	return out
}

// StringifyMany converts the provided values into a []string.
//
// NOTE: If the value does not satisfy Stringable, this will panic.
//
// See Stringable, StringableMany, Stringify, and StringifyMany
func StringifyMany(values ...any) []string {
	out := make([]string, len(values))
	for i, raw := range values {
		out[i] = Stringify(raw)
	}
	return out
}
