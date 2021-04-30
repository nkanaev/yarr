package env

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

const sep = "_"

func Fill(prefix string, v interface{}) error {
	ind := reflect.Indirect(reflect.ValueOf(v))
	if reflect.ValueOf(v).Kind() != reflect.Ptr || ind.Kind() != reflect.Struct {
		return fmt.Errorf("only the pointer to a struct is supported")
	}
	if prefix == "" {
		prefix = ind.Type().Name()
	}
	for i := 0; i < ind.NumField(); i++ {
		f := ind.Type().Field(i)
		name := f.Name
		envName, exist := f.Tag.Lookup("env")
		if exist {
			name = envName
		}
		p := strings.ToUpper(prefix + sep + name)
		err := parse(p, ind.Field(i), f)
		if err != nil {
			return err
		}
	}
	return nil
}

func parse(prefix string, f reflect.Value, sf reflect.StructField) error {
	ev, exist := os.LookupEnv(prefix)
	if !exist {
		return nil
	}
	log.Printf("key: %s, env:%s\n", prefix, ev)
	switch f.Kind() {
	case reflect.String:
		f.SetString(ev)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		iv, err := strconv.ParseInt(ev, 10, f.Type().Bits())
		if err != nil {
			return fmt.Errorf("%s:%s", prefix, err)
		}
		f.SetInt(iv)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uiv, err := strconv.ParseUint(ev, 10, f.Type().Bits())
		if err != nil {
			return fmt.Errorf("%s:%s", prefix, err)
		}
		f.SetUint(uiv)
	case reflect.Float32, reflect.Float64:
		floatValue, err := strconv.ParseFloat(ev, f.Type().Bits())
		if err != nil {
			return fmt.Errorf("%s:%s", prefix, err)
		}
		f.SetFloat(floatValue)
	case reflect.Bool:
		if ev == "" {
			f.SetBool(false)
		} else {
			b, err := strconv.ParseBool(ev)
			if err != nil {
				return fmt.Errorf("%s:%s", prefix, err)
			}
			f.SetBool(b)
		}
	default:
		return fmt.Errorf("kind %s not supported", f.Kind())
	}
	return nil
}
