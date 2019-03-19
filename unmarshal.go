package config

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	// ErrNilPointer indicates the v parameter is a nil pointer
	ErrNilPointer = errors.New("destination parameter cannot be nil")
	// ErrCannotSetValue indicates the v parameter is not writable
	ErrCannotSetValue = errors.New("destination parameter must be writable")
	// ErrInvalidType indicates that the v parameter is not a pointer
	ErrInvalidType = errors.New("must be a pointer type")
	// ErrBadMap indicates the s parameter is not a comma separated list of key:value pairs
	ErrBadMap = errors.New("string for a map must be comma separated list of key:value")
	// ErrUnsupportedType indicates that Unmarshal doesn't support that type
	ErrUnsupportedType = errors.New("cannot unmarshal unsupported type")
)

// Unmarshal attempts to convert string 's' into a value of type 'v'.
// It infers the type from v itself, and uses the appropriate conversion routine.
// This only supports the common basic types, slices, and maps, but not structs.
// Supported types:
//
//	Type									Examples
//	-----------------------------------		---------------------------------------------
//	string                                  "", "some string"
//	bool                                    "0", "false", "F", "1", "true", "T"
//	time.Duration                           "1h22m33s", "6h"
//	int, int8, int16, int32, int64          "12", "-77"
//	uint, uint8, uint16, uint32, uint64     "42", "0xDEADBEEF"
//	float32, float64                        "345.71, "-3.14159"
//	slice                                   "super, duper", "12, 20, 36"
//	map                                     "good:26, bad:37, ugly:55", "1:true, 7:false"
func Unmarshal(s string, v interface{}) error {
	rv := reflect.ValueOf(v)
	return UnmarshalValue(s, rv)
}

// UnmarshalValue converts a string representation of a value and writes it into rv.
// The type is inferred on the type of rv.
func UnmarshalValue(s string, rv reflect.Value) error {
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return ErrNilPointer
		}
		// if rv is a pointer, dereference into its inner type
		rv = rv.Elem()
	}

	// skip over a value if it cannot be written
	if !rv.CanSet() {
		return ErrCannotSetValue
	}

	kind := rv.Kind()
	typ := rv.Type()

	switch {
	case kind == reflect.String:
		rv.SetString(s)

	case kind == reflect.Bool:
		b, err := strconv.ParseBool(s)
		if err != nil {
			return err
		}
		rv.SetBool(b)

		// this case needs to appear before the one for integers as they're both type int64
	case kind == reflect.Int64 && typ.String() == "time.Duration":
		d, err := time.ParseDuration(s)
		if err != nil {
			return err
		}
		rv.SetInt(int64(d))

	case kind == reflect.Int || kind == reflect.Int8 || kind == reflect.Int16 || kind == reflect.Int32 || kind == reflect.Int64:
		i, err := strconv.ParseInt(s, 0, typ.Bits())
		if err != nil {
			return err
		}
		rv.SetInt(i)

	case kind == reflect.Uint || kind == reflect.Uint8 || kind == reflect.Uint16 || kind == reflect.Uint32 || kind == reflect.Uint64:
		u, err := strconv.ParseUint(s, 0, typ.Bits())
		if err != nil {
			return err
		}
		rv.SetUint(u)

	case kind == reflect.Float32 || kind == reflect.Float64:
		f, err := strconv.ParseFloat(s, typ.Bits())
		if err != nil {
			return err
		}
		rv.SetFloat(f)

	case kind == reflect.Slice:
		sp := strings.Split(s, ",")
		slice := reflect.MakeSlice(typ, len(sp), len(sp))
		for i, t := range sp {
			c := strings.TrimSpace(t)
			err := UnmarshalValue(c, slice.Index(i))
			if err != nil {
				return err
			}
		}
		rv.Set(slice)

	case kind == reflect.Map:
		items := strings.Split(s, ",")
		m := reflect.MakeMap(typ)
		for _, pair := range items {
			kv := strings.Split(pair, ":")
			if len(kv)%2 != 0 {
				return ErrBadMap
			}
			k := reflect.New(typ.Key()).Elem()
			err := UnmarshalValue(strings.TrimSpace(kv[0]), k)
			if err != nil {
				return err
			}
			v := reflect.New(typ.Elem()).Elem()
			err = UnmarshalValue(strings.TrimSpace(kv[1]), v)
			if err != nil {
				return err
			}
			m.SetMapIndex(k, v)
		}
		rv.Set(m)

	default:
		return ErrUnsupportedType
	}

	return nil
}
