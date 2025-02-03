package modifiers

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"strings"
	"unicode/utf16"
)

func GetComparator(modifiers ...string) (ComparatorFunc, error) {
	return getComparator(Comparators, false, modifiers...)
}

func GetComparatorCaseSensitive(modifiers ...string) (ComparatorFunc, error) {
	return getComparator(ComparatorsCaseSensitive, true, modifiers...)
}

func getComparator(comparators map[string]Comparator, caseSensitive bool, modifiers ...string) (ComparatorFunc, error) {
	if len(modifiers) == 0 {
		if caseSensitive {
			return baseComparatorCaseSensitive{}.Bridges, nil
		}
		return baseComparator{}.Bridges, nil
	}

	// A valid sequence of modifiers is ([ValueModifier]*)[Comparator]?
	// If a comparator is specified, it must be in the last position and cannot be succeeded by any other modifiers
	// If no comparator is specified, the default comparator is used
	var valueModifiers []ValueModifier
	var comparator Comparator
	for i, modifier := range modifiers {
		comparatorModifier := comparators[modifier]
		valueModifier := ValueModifiers[modifier]
		switch {
		// Validate correctness
		case comparatorModifier == nil && valueModifier == nil:
			return nil, fmt.Errorf("unknown modifier %s", modifier)
		case i < len(modifiers)-1 && comparators[modifier] != nil:
			return nil, fmt.Errorf("comparator modifier %s must be the last modifier", modifier)

		// Build up list of modifiers
		case valueModifier != nil:
			valueModifiers = append(valueModifiers, valueModifier)
		case comparatorModifier != nil:
			comparator = comparatorModifier
		}
	}
	if comparator == nil {
		if caseSensitive {
			comparator = baseComparatorCaseSensitive{}
		} else {
			comparator = baseComparator{}
		}
	}

	return func(field, value any) (string, error) {
		var err error
		for _, modifier := range valueModifiers {
			value, err = modifier.Modify(value)
			if err != nil {
				return "", err
			}
		}

		return comparator.Bridges(field, value)
	}, nil
}

type Comparator interface {
	Bridges(field any, value any) (string, error)
}

type ComparatorFunc func(field, value any) (string, error)

// ValueModifier modifies the expected value before it is passed to the comparator.
// For example, the `base64` modifier converts the expected value to base64.
type ValueModifier interface {
	Modify(value any) (any, error)
}

var Comparators = map[string]Comparator{
	"contains":   contains{},
	"endswith":   endswith{},
	"startswith": startswith{},
	"re":         re{},
	"cidr":       cidr{},
	"gt":         gt{},
	"gte":        gte{},
	"lt":         lt{},
	"lte":        lte{},
}

var ComparatorsCaseSensitive = map[string]Comparator{
	"contains":   containsCS{},
	"endswith":   endswithCS{},
	"startswith": startswithCS{},
	"re":         re{},
	"cidr":       cidr{},
	"gt":         gt{},
	"gte":        gte{},
	"lt":         lt{},
	"lte":        lte{},
}

var ValueModifiers = map[string]ValueModifier{
	"base64": b64{},
	"wide":   wide{},
}

type baseComparator struct{}

func (baseComparator) Bridges(field, value any) (string, error) {
	switch {
	case field == nil && value == "null":
		return "", nil
	default:
		// The Sigma spec defines that by default comparisons are case-insensitive
		return fmt.Sprintf("%v=\"%v\"", strings.ToLower(coerceString(field)), strings.ToLower(EscapeBackslashes(coerceString(value)))), nil
	}
}

type contains struct{}

func (contains) Bridges(field, value any) (string, error) {
	return fmt.Sprintf("%v=\"*%v*\"", strings.ToLower(coerceString(field)), strings.ToLower(EscapeBackslashes(coerceString(value)))), nil
}

type endswith struct{}

func (endswith) Bridges(field, value any) (string, error) {
	return fmt.Sprintf("%v=\"*%v\"", strings.ToLower(coerceString(field)), strings.ToLower(EscapeBackslashes(coerceString(value)))), nil
}

type startswith struct{}

func (startswith) Bridges(field, value any) (string, error) {
	return fmt.Sprintf("%v=\"%v*\"", strings.ToLower(coerceString(field)), strings.ToLower(EscapeBackslashes(coerceString(value)))), nil
}

type baseComparatorCaseSensitive struct{}

func (baseComparatorCaseSensitive) Bridges(field, value any) (string, error) {
	switch {
	case field == nil && value == "null":
		return "", nil
	default:
		return fmt.Sprintf("%v=\"%v\"", strings.ToLower(coerceString(field)), EscapeBackslashes(coerceString(value))), nil
	}
}

type containsCS struct{}

func (containsCS) Bridges(field, value any) (string, error) {
	return fmt.Sprintf("%v=\"*%v*\"", strings.ToLower(coerceString(field)), EscapeBackslashes(coerceString(value))), nil
}

type endswithCS struct{}

func (endswithCS) Bridges(field, value any) (string, error) {
	return fmt.Sprintf("%v=\"*%v\"", strings.ToLower(coerceString(field)), EscapeBackslashes(coerceString(value))), nil
}

type startswithCS struct{}

func (startswithCS) Bridges(field, value any) (string, error) {
	return fmt.Sprintf("%v=\"%v*\"", strings.ToLower(coerceString(field)), EscapeBackslashes(coerceString(value))), nil
}

type re struct{}

func (re) Bridges(field any, value any) (string, error) {
	return fmt.Sprintf(" | regex %v=\"%v\"", strings.ToLower(coerceString(field)), EscapeBackslashes(coerceString(value))), nil
}

type cidr struct{}

func (cidr) Bridges(field any, value any) (string, error) {
	return fmt.Sprintf("%v=\"%v\"", strings.ToLower(coerceString(field)), coerceString(value)), nil
}

type gt struct{}

func (gt) Bridges(field any, value any) (string, error) {
	return fmt.Sprintf("%v > \"%v\"", strings.ToLower(coerceString(field)), coerceString(value)), nil
}

type gte struct{}

func (gte) Bridges(field any, value any) (string, error) {
	return fmt.Sprintf("%v >= \"%v\"", strings.ToLower(coerceString(field)), coerceString(value)), nil
}

type lt struct{}

func (lt) Bridges(field any, value any) (string, error) {
	return fmt.Sprintf("%v < \"%v\"", strings.ToLower(coerceString(field)), coerceString(value)), nil
}

type lte struct{}

func (lte) Bridges(field any, value any) (string, error) {
	return fmt.Sprintf("%v <= \"%v\"", strings.ToLower(coerceString(field)), coerceString(value)), nil
}

type b64 struct{}

func (b64) Modify(value any) (any, error) {
	return base64.StdEncoding.EncodeToString([]byte(coerceString(value))), nil
}

type wide struct{}

func (wide) Modify(value any) (any, error) {
	runes := utf16.Encode([]rune(coerceString(value)))
	bytes := make([]byte, 2*len(runes))
	for i, r := range runes {
		binary.LittleEndian.PutUint16(bytes[i*2:], r)
	}
	return coerceString(bytes), nil
}

func coerceString(v interface{}) string {
	switch vv := v.(type) {
	case string:
		return vv
	case []byte:
		return string(vv)
	default:
		return fmt.Sprint(vv)
	}
}

// EscapeBackslashes takes a string and doubles all backslashes (`\`).
func EscapeBackslashes(input string) string {
	var builder strings.Builder
	builder.Grow(len(input)) // Optimize memory allocation for large strings

	for _, char := range input {
		if char == '\\' {
			// Add an extra backslash for escaping
			builder.WriteRune('\\')
		}
		// Write the original character (escaped or not)
		builder.WriteRune(char)
	}

	return builder.String()
}
