package unmarshalflag

import (
	"fmt"
	"reflect"
	"time"

	"github.com/recallsong/unmarshal"
	"github.com/spf13/pflag"
)

var flagType = reflect.TypeOf((*pflag.Value)(nil)).Elem()

// BindFlag .
func BindFlag(flags *pflag.FlagSet, data interface{}) error {
	return BindFlagValue(flags, reflect.ValueOf(data))
}

// BindFlagValue .
func BindFlagValue(flags *pflag.FlagSet, data reflect.Value) error {
	return unmarshal.Unmarshal(
		data, []string{"flag", "desc"}, &flagType, nil,
		func(field string, tags []string, ftyp reflect.Type, fval reflect.Value) (err error) {
			if len(tags[0]) <= 0 {
				return nil
			}
			if ftyp.AssignableTo(flagType) {
				flags.Var(fval.Interface().(pflag.Value), tags[0], tags[1])
				return nil
			}
			if ftyp == unmarshal.DurationType {
				flags.DurationVar(fval.Addr().Interface().(*time.Duration), tags[0], fval.Interface().(time.Duration), tags[1])
				return nil
			}
			switch ftyp.Kind() {
			case reflect.Bool:
				flags.BoolVar(fval.Addr().Interface().(*bool), tags[0], fval.Bool(), tags[1])
			case reflect.Int:
				flags.IntVar(fval.Addr().Interface().(*int), tags[0], fval.Interface().(int), tags[1])
			case reflect.Int8:
				flags.Int8Var(fval.Addr().Interface().(*int8), tags[0], fval.Interface().(int8), tags[1])
			case reflect.Int16:
				flags.Int16Var(fval.Addr().Interface().(*int16), tags[0], fval.Interface().(int16), tags[1])
			case reflect.Int32:
				flags.Int32Var(fval.Addr().Interface().(*int32), tags[0], fval.Interface().(int32), tags[1])
			case reflect.Int64:
				flags.Int64Var(fval.Addr().Interface().(*int64), tags[0], fval.Interface().(int64), tags[1])
			case reflect.Uint:
				flags.UintVar(fval.Addr().Interface().(*uint), tags[0], fval.Interface().(uint), tags[1])
			case reflect.Uint8:
				flags.Uint8Var(fval.Addr().Interface().(*uint8), tags[0], fval.Interface().(uint8), tags[1])
			case reflect.Uint16:
				flags.Uint16Var(fval.Addr().Interface().(*uint16), tags[0], fval.Interface().(uint16), tags[1])
			case reflect.Uint32:
				flags.Uint32Var(fval.Addr().Interface().(*uint32), tags[0], fval.Interface().(uint32), tags[1])
			case reflect.Uint64:
				flags.Uint64Var(fval.Addr().Interface().(*uint64), tags[0], fval.Interface().(uint64), tags[1])
			case reflect.Float32:
				flags.Float32Var(fval.Addr().Interface().(*float32), tags[0], fval.Interface().(float32), tags[1])
			case reflect.Float64:
				flags.Float64Var(fval.Addr().Interface().(*float64), tags[0], fval.Interface().(float64), tags[1])
			case reflect.String:
				flags.StringVar(fval.Addr().Interface().(*string), tags[0], fval.Interface().(string), tags[1])
			case reflect.Slice:
				if fval.Type().Elem().Kind() == reflect.Bool {
					flags.BoolSliceVar(fval.Addr().Interface().(*[]bool), tags[0], fval.Interface().([]bool), tags[1])
				} else if fval.Type().Elem().Kind() == reflect.Int {
					flags.IntSliceVar(fval.Addr().Interface().(*[]int), tags[0], fval.Interface().([]int), tags[1])
				} else if fval.Type().Elem().Kind() == reflect.Int32 {
					flags.Int32SliceVar(fval.Addr().Interface().(*[]int32), tags[0], fval.Interface().([]int32), tags[1])
				} else if fval.Type().Elem().Kind() == reflect.Int64 {
					flags.Int64SliceVar(fval.Addr().Interface().(*[]int64), tags[0], fval.Interface().([]int64), tags[1])
				} else if fval.Type().Elem().Kind() == reflect.Uint {
					flags.UintSliceVar(fval.Addr().Interface().(*[]uint), tags[0], fval.Interface().([]uint), tags[1])
				} else if fval.Type().Elem().Kind() == reflect.Float32 {
					flags.Float32SliceVar(fval.Addr().Interface().(*[]float32), tags[0], fval.Interface().([]float32), tags[1])
				} else if fval.Type().Elem().Kind() == reflect.Float64 {
					flags.Float64SliceVar(fval.Addr().Interface().(*[]float64), tags[0], fval.Interface().([]float64), tags[1])
				} else if fval.Type().Elem().Kind() == reflect.String {
					flags.StringArrayVar(fval.Addr().Interface().(*[]string), tags[0], fval.Interface().([]string), tags[1])
				} else {
					return fmt.Errorf("not support bind type %v for %s", ftyp, field)
				}
			default:
				return fmt.Errorf("not support bind type %v for %s", ftyp, field)
			}
			return nil
		},
	)
}
