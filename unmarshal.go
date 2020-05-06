package unmarshal

import (
	"encoding"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Unmarshal .
func Unmarshal(
	data reflect.Value,
	tags []string,
	unmarshaler *reflect.Type,
	check func(field string, tags []string) bool,
	setter func(field string, tags []string, typ reflect.Type, val reflect.Value) error,
) error {
	typ := data.Type()
	for typ.Kind() == reflect.Ptr {
		if data.IsNil() {
			data.Set(reflect.New(typ.Elem()))
		}
		typ = typ.Elem()
		data = data.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return fmt.Errorf("unmarshal data is not struct")
	}
	num := typ.NumField()
	for i := 0; i < num; i++ {
		field := typ.Field(i)
		var tagvals []string
		for _, tag := range tags {
			tagvals = append(tagvals, field.Tag.Get(tag))
		}
		ftyp := field.Type
		fval := data.Field(i)
		var value *reflect.Value
		if ftyp.Kind() == reflect.Ptr {
			for ftyp.Kind() == reflect.Ptr {
				if fval.IsNil() {
					fval.Set(reflect.New(ftyp.Elem()))
				}
				if unmarshaler != nil && ftyp.AssignableTo(*unmarshaler) && value == nil {
					val := fval
					value = &val
				}
				ftyp = ftyp.Elem()
				fval = fval.Elem()
			}
		} else {
			if unmarshaler != nil {
				if ftyp.AssignableTo(*unmarshaler) {
					val := fval
					value = &val
				} else if fval.Addr().Type().AssignableTo(*unmarshaler) {
					val := fval.Addr()
					value = &val
				}
			}
		}
		// TODO ...
		// if check != nil && !check(field.Name, tagvals) {
		// 	continue
		// }
		if value != nil {
			err := setter(field.Name, tagvals, value.Type(), *value)
			if err != nil {
				return err
			}
			continue
		}
		if ftyp != DurationType && ftyp.Kind() == reflect.Struct {
			err := Unmarshal(fval, tags, unmarshaler, check, setter)
			if err != nil {
				return err
			}
			continue
		}
		err := setter(field.Name, tagvals, ftyp, fval)
		if err != nil {
			return err
		}
	}
	return nil
}

// SkipEmpty .
func SkipEmpty(field string, tags []string) bool {
	for _, tag := range tags {
		if len(tag) <= 0 {
			return false
		}
	}
	return true
}

var (
	// TextUnmarshalerType .
	TextUnmarshalerType = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()
	// DurationType .
	DurationType = reflect.TypeOf((*time.Duration)(nil)).Elem()
)

// BindText .
func BindText(data reflect.Value, tag string, key func(string) string) error {
	return Unmarshal(
		data, []string{tag}, &TextUnmarshalerType, SkipEmpty,
		func(field string, tags []string, ftyp reflect.Type, fval reflect.Value) (err error) {
			text := key(tags[0])
			if len(text) <= 0 {
				return nil
			}
			if ftyp.AssignableTo(TextUnmarshalerType) {
				textUnmarshaler := fval.Interface().(encoding.TextUnmarshaler)
				err := textUnmarshaler.UnmarshalText([]byte(text))
				if err != nil {
					return fmt.Errorf("invalid %s value '%s' for %s, %s", tag, text, field, err)
				}
				return nil
			}
			if ftyp == DurationType {
				val, err := time.ParseDuration(text)
				if err != nil {
					return fmt.Errorf("invalid %s value '%s' for %s, %s", tag, text, field, err)
				}
				fval.Set(reflect.ValueOf(val))
				return nil
			}
			switch ftyp.Kind() {
			case reflect.Bool:
				val, err := strconv.ParseBool(text)
				if err != nil {
					return fmt.Errorf("invalid %s value '%s' for %s, %s", tag, text, field, err)
				}
				fval.Set(reflect.ValueOf(val))
			case reflect.Int:
				val, err := strconv.Atoi(text)
				if err != nil {
					return fmt.Errorf("invalid %s value '%s' for %s, %s", tag, text, field, err)
				}
				fval.Set(reflect.ValueOf(val))
				break
			case reflect.Int8:
				val, err := strconv.ParseInt(text, 10, 8)
				if err != nil {
					return fmt.Errorf("invalid %s value '%s' for %s, %s", tag, text, field, err)
				}
				fval.Set(reflect.ValueOf(int8(val)))
			case reflect.Int16:
				val, err := strconv.ParseInt(text, 10, 16)
				if err != nil {
					return fmt.Errorf("invalid %s value '%s' for %s, %s", tag, text, field, err)
				}
				fval.Set(reflect.ValueOf(int16(val)))
			case reflect.Int32:
				val, err := strconv.ParseInt(text, 10, 32)
				if err != nil {
					return fmt.Errorf("invalid %s value '%s' for %s, %s", tag, text, field, err)
				}
				fval.Set(reflect.ValueOf(int32(val)))
			case reflect.Int64:
				val, err := strconv.ParseInt(text, 10, 64)
				if err != nil {
					return fmt.Errorf("invalid %s value '%s' for %s, %s", tag, text, field, err)
				}
				fval.Set(reflect.ValueOf(int64(val)))
			case reflect.Uint:
				val, err := strconv.ParseUint(text, 10, 64)
				if err != nil {
					return fmt.Errorf("invalid %s value '%s' for %s, %s", tag, text, field, err)
				}
				fval.Set(reflect.ValueOf(uint64(val)))
			case reflect.Uint8:
				val, err := strconv.ParseUint(text, 10, 8)
				if err != nil {
					return fmt.Errorf("invalid %s value '%s' for %s, %s", tag, text, field, err)
				}
				fval.Set(reflect.ValueOf(uint8(val)))
			case reflect.Uint16:
				val, err := strconv.ParseUint(text, 10, 16)
				if err != nil {
					return fmt.Errorf("invalid %s value '%s' for %s, %s", tag, text, field, err)
				}
				fval.Set(reflect.ValueOf(uint16(val)))
			case reflect.Uint32:
				val, err := strconv.ParseUint(text, 10, 32)
				if err != nil {
					return fmt.Errorf("invalid %s value '%s' for %s, %s", tag, text, field, err)
				}
				fval.Set(reflect.ValueOf(uint32(val)))
			case reflect.Uint64:
				val, err := strconv.ParseUint(text, 10, 64)
				if err != nil {
					return fmt.Errorf("invalid %s value '%s' for %s, %s", tag, text, field, err)
				}
				fval.Set(reflect.ValueOf(val))
			case reflect.Float32:
				val, err := strconv.ParseFloat(text, 32)
				if err != nil {
					return fmt.Errorf("invalid %s value '%s' for %s, %s", tag, text, field, err)
				}
				fval.Set(reflect.ValueOf((uint32)(val)))
			case reflect.Float64:
				val, err := strconv.ParseFloat(text, 64)
				if err != nil {
					return fmt.Errorf("invalid %s value '%s' for %s, %s", tag, text, field, err)
				}
				fval.Set(reflect.ValueOf(val))
			case reflect.String:
				fval.Set(reflect.ValueOf(text))
			case reflect.Map:
				err := json.Unmarshal([]byte(text), fval.Addr().Interface())
				if err != nil {
					return fmt.Errorf("invalid %s value for %s, %s", tag, field, err)
				}
			case reflect.Slice:
				if fval.Type().Elem().Kind() == reflect.String {
					fval.Set(reflect.ValueOf(strings.Split(text, ",")))
				} else {
					err := json.Unmarshal([]byte("["+text+"]"), fval.Addr().Interface())
					if err != nil {
						return fmt.Errorf("invalid %s value for %s, %s", tag, field, err)
					}
				}
			default:
				return fmt.Errorf("not support bind type %v for %s", ftyp, field)
			}
			return nil
		})
}

// BindDefault .
func BindDefault(data interface{}) error {
	return BindText(reflect.ValueOf(data), "default", KeyValue)
}

// BindEnv .
func BindEnv(data interface{}) error {
	return BindText(reflect.ValueOf(data), "env", EnvValue)
}

// KeyValue .
func KeyValue(k string) string { return k }

// EnvValue .
func EnvValue(k string) string { return os.Getenv(k) }
