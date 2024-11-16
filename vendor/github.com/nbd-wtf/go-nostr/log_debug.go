//go:build debug

package nostr

import (
	"encoding/json"
	"fmt"
)

func debugLogf(str string, args ...any) {
	// this is such that we don't modify the actual args that may be used outside of this function
	printableArgs := make([]any, len(args))

	for i, v := range args {
		printableArgs[i] = stringify(v)
	}

	DebugLogger.Printf(str, printableArgs...)
}

func stringify(anything any) any {
	switch v := anything.(type) {
	case []any:
		// this is such that we don't modify the actual values that may be used outside of this function
		printableValues := make([]any, len(v))
		for i, subv := range v {
			printableValues[i] = stringify(subv)
		}
		return printableValues
	case []json.RawMessage:
		j, _ := json.Marshal(v)
		return string(j)
	case []byte:
		return string(v)
	case fmt.Stringer:
		return v.String()
	default:
		return v
	}
}
