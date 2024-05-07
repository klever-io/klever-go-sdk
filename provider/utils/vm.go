package utils

import "strings"

func SplitTypes(fullType string) (string, string) {
	index := strings.Index(fullType, "<")
	if index == -1 {
		return "", fullType
	}

	wrapperType := fullType[:index]
	valueType := fullType[index+1:]
	valueType = valueType[:len(valueType)-1]

	return wrapperType, valueType
}

func SplitTupleTypes(tuple string) []string {
	var result []string
	var current int
	depth := 0
	start := 0

	for current < len(tuple) {
		switch tuple[current] {
		case '<':
			depth++
		case '>':
			depth--
		case ',':
			if depth == 0 {
				result = append(result, strings.TrimSpace(tuple[start:current]))
				start = current + 1
			}
		}
		current++
	}

	result = append(result, strings.TrimSpace(tuple[start:]))

	return result
}

func IsDynamicLengthType(t string) bool {
	typesMap := map[string]struct{}{
		ManagedBuffer:   {},
		TokenIdentifier: {},
		Bytes:           {},
		BoxedBytes:      {},
		String:          {},
		StrRef:          {},
		VecU8:           {},
		SliceU8:         {},
		BigInt:          {},
		BigUint:         {},
		BigFloat:        {},
	}

	_, exists := typesMap[t]
	return exists
}
