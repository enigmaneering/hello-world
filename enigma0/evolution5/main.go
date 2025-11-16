package main

import (
	"fmt"
	"math/big"
	"strconv"

	"git.enigmaneering.net/hello-world/enigma0/solution0/evolution5/core/std"
	"git.ignitelabs.net/janos/core/sys/atlas"
)

func main() {
	fmt.Println(Stringify(new(big.Int)))
}

type Operation[T any] struct {
	Context std.Path
}

// Operation.WithBase() Operation
// Operation.WithPrecision() Operation

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
	case *big.Int:
		return typed.Text(10)
	case *big.Float:
		return typed.Text('f', int(atlas.Precision))
	case *big.Rat:
		return typed.String()
	// NOTE: fmt.Stringer exists EXACTLY here for a reason!
	// Switch cases are evaluated top to bottom, so below the fmt.Stringer case
	// exist the primitive types, and above it exist composite types which we
	// explicitly define the behavior of.  All others should fall into this
	// particular case because they define their own string functionality.

	// For instance: big.Float's String() function uses exponential notation, so we must override it with our own rules to standard notation.
	case fmt.Stringer:
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
