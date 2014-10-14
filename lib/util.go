package lib

import (
	"encoding"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
)

func debug(v ...interface{}) {
	if os.Getenv("DEBUG") != "" {
		log.Println(v...)
	}
}

var (
	printXMLStructOffset = 2
	textMarshalerType    = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()
)

func PrintXMLStruct(s interface{}, indent ...int) {
	sv := reflect.ValueOf(s)
	st := sv.Type()
	if sv.Kind() != reflect.Struct {
		return
	}

	idt := 0
	if len(indent) > 0 {
		idt = indent[0]
	}
	sp := strings.Repeat(" ", idt*printXMLStructOffset)

	numField := sv.NumField()
	for i := 0; i < numField; i++ {
		fv := sv.Field(i)
		fname := st.Field(i).Name
		if fname == "XMLName" {
			continue
		}
		if fv.Kind() == reflect.Ptr && fv.IsNil() {
			continue
		}
		if fv.CanInterface() && fv.Type().Implements(textMarshalerType) {
			text, _ := fv.Interface().(encoding.TextMarshaler).MarshalText()
			fmt.Printf("%s%s: %s\n", sp, fname, string(text))
			continue
		}
		switch fv.Kind() {
		case reflect.Bool:
			var bstr string
			if fv.Bool() {
				bstr = "true"
			} else {
				bstr = "false"
			}
			fmt.Printf("%s%s: %s\n", sp, fname, bstr)
		case reflect.Float32, reflect.Float64:
			fmt.Printf("%s%s: %f\n", sp, fname, fv.Float())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fmt.Printf("%s%s: %d\n", sp, fname, fv.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			fmt.Printf("%s%s: %d\n", sp, fname, fv.Uint())
		case reflect.String:
			fmt.Printf("%s%s: %s\n", sp, fname, fv.String())
		case reflect.Struct:
			fmt.Printf("%s%s:\n", sp, fname)
			if fv.CanInterface() {
				PrintXMLStruct(fv.Interface(), idt+1)
			}
		case reflect.Array, reflect.Slice:
			numArray := fv.Len()
			for j := 0; j < numArray; j++ {
				fmt.Printf("%s%s:\n", sp, fname)
				elem := fv.Index(j)
				if elem.CanInterface() {
					PrintXMLStruct(elem.Interface(), idt+1)
				}
			}
		case reflect.Map:
			fmt.Printf("%s%s:\n", sp, fname)
			for _, k := range fv.MapKeys() {
				mapsp := strings.Repeat(" ", (idt+1)*printXMLStructOffset)
				fmt.Printf("%s%s:\n", mapsp, k)
				elem := fv.MapIndex(k)
				if elem.CanInterface() {
					PrintXMLStruct(elem.Interface(), idt+2)
				}
			}
		case reflect.Ptr:
			elem := fv.Elem()
			switch elem.Kind() {
			case reflect.Bool:
				var bstr string
				if elem.Bool() {
					bstr = "true"
				} else {
					bstr = "false"
				}
				fmt.Printf("%s%s: %s\n", sp, fname, bstr)
			case reflect.Float32, reflect.Float64:
				fmt.Printf("%s%s: %f\n", sp, fname, elem.Float())
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				fmt.Printf("%s%s: %d\n", sp, fname, elem.Int())
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				fmt.Printf("%s%s: %d\n", sp, fname, elem.Uint())
			case reflect.String:
				fmt.Printf("%s%s: %s\n", sp, fname, elem.String())
			case reflect.Struct:
				fmt.Printf("%s%s:\n", sp, fname)
				if elem.CanInterface() {
					PrintXMLStruct(elem.Interface(), idt+1)
				}
			}
		}
	}
}
