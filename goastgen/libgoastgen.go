package goastgen

import "C"
import (
	"encoding/json"
	"go/parser"
	"go/token"
	"log"
	"reflect"
	"unsafe"
)

//export ExternallyCalled
func ExternallyCalled() *C.char {
	result := "John"
	return C.CString(result)
}

//export Add
func Add(a int, b int) int {
	return a + b
}

/*
 It will parse given source code and generate AST in JSON format

 Parameters:
  filename: Filename used for generating AST metadata
  src: string, []byte, or io.Reader - Source code

 Returns:
  If given source is valid go source then it will generate AST in JSON format other will return "" string.
*/
func internalParseAstFromSource(filename string, src any) string {
	fset := token.NewFileSet()
	parsedAst, err := parser.ParseFile(fset, filename, src, 0)
	if err != nil {
		// TODO: convert this to just warning error log.
		log.Fatal(err)
	}
	result := serilizeToMap(parsedAst)
	return serilizeToJsonStr(result)
}

/*
 It will parse all the go files in given source folder location and generate AST in JSON format

 Parameters:
  file: absolute root directory path of source code

 Returns:
  If given directory contains valid go source code then it will generate AST in JSON format otherwise will return "" string.
*/
func internalParseAstFromDir(dir string) string {
	fset := token.NewFileSet()
	parsedAst, err := parser.ParseDir(fset, dir, nil, 0)
	if err != nil {
		// TODO: convert this to just warning error log.
		log.SetPrefix("[ERROR]")
		log.Println("Error while parsing source from source directory -> '", dir, ",")
		log.Print(err)
	}
	result := serilizeToMap(parsedAst)
	return serilizeToJsonStr(result)
}

/*
 It will parse the given file and generate AST in JSON format

 Parameters:
  file: absolute file path to be parsed

 Returns:
  If given file is a valid go code then it will generate AST in JSON format otherwise will return "" string.
*/
func internalParseAstFromFile(file string) string {
	fset := token.NewFileSet()
	// NOTE: Haven't explore much of mode parameter. Default value has been passed as 0
	parsedAst, err := parser.ParseFile(fset, file, nil, 0)
	if err != nil {
		log.SetPrefix("[ERROR]")
		log.Println("Error while parsing source file -> '", file, ",")
		log.Print(err)
		return ""
	} else {
		result := serilizeToMap(parsedAst)
		return serilizeToJsonStr(result)
	}
}

/*
 Independent function which handles serialisation of map[string]interface{} in to JSON

 Parameters:
  objectMap: Mostly it will be object of map[string]interface{}

 Returns:
  JSON string
*/
func serilizeToJsonStr(objectMap interface{}) string {
	jsonStr, err := json.MarshalIndent(objectMap, "", "  ")
	if err != nil {
		log.SetPrefix("[ERROR]")
		log.Println("Error while generating the AST JSON")
		log.Print(err)
	}
	return string(jsonStr)
}

/*
Process Map type objects. In order to process the contents of the map's value object.
If the value object is of type 'struct' then we are converting it to map[string]interface{} and using it.

Parameters:
 object: expects map[string] any

Returns:
 It returns and object of map[string]interface{} by converting any 'Struct' type value field to map
*/
func processMap(object interface{}) interface{} {
	value := reflect.ValueOf(object)
	objMap := make(map[string]interface{})
	for _, key := range value.MapKeys() {
		objValue := value.MapIndex(key)

		// If the map is created to accept valye of any type i.e. map[string]interface{}.
		// Then it's value's reflect.Kind is of type reflect.Interface.
		// We need to fetch original objects reflect.Value by calling .Elem() on it.
		if objValue.Kind() == reflect.Interface {
			objValue = objValue.Elem()
		}

		var ptrValue reflect.Value
		// Checking the reflect.Kind of value object and if its pointer
		// then fetching the reflect.Value of the object pointed to by this pointer
		if objValue.Kind() == reflect.Pointer {
			objValue = objValue.Elem()
			ptrValue = objValue
		}

		if objValue.IsValid() {
			switch objValue.Kind() {
			case reflect.String, reflect.Int, reflect.Bool, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				objMap[key.String()] = objValue.Interface()
			case reflect.Struct:
				objMap[key.String()] = processStruct(objValue.Interface(), ptrValue)
			default:
				log.SetPrefix("[WARNING]")
				log.Println(getLogPrefix(), objValue.Kind(), "- not handled")
			}
		}
	}
	return objMap
}

/*
 This will process the Array or Slice (Dynamic Array).
 It will identify the type/reflect.Kind of each array element and process the array element according.

 Parameters:
  object: []interface{} - expected to pass object of Array or Slice

 Returns:
  It will return []map[string]interface{}
*/
func processArrayOrSlice(object interface{}) interface{} {
	value := reflect.ValueOf(object)
	var nodeList []interface{}
	for j := 0; j < value.Len(); j++ {
		arrayElementValue := value.Index(j)
		elementKind := arrayElementValue.Kind()
		// If you create an array interface{} and assign pointer as elements into this array.
		// when we try to identify the reflect.Kind of such element it will be of type reflect.Interface.
		// In such case we need to call .elem() to fetch the original reflect.Value of the array element.
		// Refer test case - TestSimpleInterfaceWithArrayOfPointersType for the same.
		if elementKind == reflect.Interface {
			arrayElementValue = arrayElementValue.Elem()
			elementKind = arrayElementValue.Kind()
		}
		ptrValue := arrayElementValue

		switch elementKind {
		case reflect.String, reflect.Int, reflect.Bool, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			nodeList = append(nodeList, arrayElementValue.Interface())
		case reflect.Struct:
			nodeList = append(nodeList, processStruct(arrayElementValue.Interface(), ptrValue))
		case reflect.Map:
			nodeList = append(nodeList, processMap(arrayElementValue.Interface()))
		case reflect.Pointer:
			if arrayElementValue.Elem().IsValid() {
				arrayElementValuePtrKind := arrayElementValue.Elem().Kind()
				switch arrayElementValuePtrKind {
				case reflect.String, reflect.Int, reflect.Bool, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					nodeList = append(nodeList, arrayElementValue.Elem().Interface())
				case reflect.Struct:
					nodeList = append(nodeList, processStruct(arrayElementValue.Elem().Interface(), ptrValue))
				case reflect.Map:
					nodeList = append(nodeList, processMap(arrayElementValue.Elem().Interface()))
				default:
					log.SetPrefix("[WARNING]")
					log.Println(getLogPrefix(), arrayElementValuePtrKind, "- not handled for array pointer element")
				}
			}
		default:
			log.SetPrefix("[WARNING]")
			log.Println(getLogPrefix(), elementKind, "- not handled for array element")
		}
	}
	return nodeList
}

var nodeAddressMap = make(map[interface{}]interface{})

/*
 This will process object of 'struct' type and convert it into document / map[string]interface{}.
 It will process each field of this object, if it contains further child objects, arrays or maps.
 Then it will get those respective field objects processed through respective processors.
 e.g. if the field object is of type 'struct' then it will call function processStruct recursively

 Parameters:
  node: Object of struct

 Returns:
  It will return object of map[string]interface{} by converting all the child fields recursively into map

*/
func processStruct(node interface{}, objPtrValue reflect.Value) interface{} {
	objectMap := make(map[string]interface{})
	elementType := reflect.TypeOf(node)
	elementValueObj := reflect.ValueOf(node)

	process := true
	var objAddress uintptr
	if objPtrValue.Kind() == reflect.Pointer {
		ptr := unsafe.Pointer(objPtrValue.Pointer()) // Get the pointer address as an unsafe.Pointer
		objAddress = uintptr(ptr)                    // Convert unsafe.Pointer to uintptr
		_, ok := nodeAddressMap[objAddress]
		if ok {
			process = false
		}
	}
	objectMap["node_type"] = elementValueObj.Type().String()
	if process {
		if objPtrValue.Kind() == reflect.Pointer {
			nodeAddressMap[objAddress] = node
		}
		// We will iterate through each field process each field according to its reflect.Kind type.
		for i := 0; i < elementType.NumField(); i++ {
			field := elementType.Field(i)
			value := elementValueObj.Field(i)
			fieldKind := value.Type().Kind()

			// If object is defined with field type as interface{} and assigned with pointer value.
			// We need to first fetch the element from the interface
			if fieldKind == reflect.Interface {
				fieldKind = value.Elem().Kind()
				value = value.Elem()
			}
			var ptrValue reflect.Value

			if fieldKind == reflect.Pointer {
				// NOTE: This handles only one level of pointer. At this moment we don't expect to get pointer to pointer.
				// This will fetch the reflect.Kind of object pointed to by this field pointer
				fieldKind = value.Type().Elem().Kind()
				// This will fetch the reflect.Value of object pointed to by this field pointer.
				ptrValue = value
				value = value.Elem()
			}
			if value.IsValid() {
				switch fieldKind {
				case reflect.String, reflect.Int, reflect.Bool, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					objectMap[field.Name] = value.Interface()
				case reflect.Struct:
					objectMap[field.Name] = processStruct(value.Interface(), ptrValue)
				case reflect.Map:
					objectMap[field.Name] = processMap(value.Interface())
				case reflect.Array, reflect.Slice:
					objectMap[field.Name] = processArrayOrSlice(value.Interface())
				default:
					log.SetPrefix("[WARNING]")
					log.Println(getLogPrefix(), field.Name, "- of Kind ->", fieldKind, "- not handled")
				}
			}
		}
	}
	return objectMap
}

/*
 First step to convert the given object to Map, in order to export into JSON format.

 This function will check if the given passed object is of primitive, struct, map, array or slice (Dynamic array) type
 and process object accordingly to convert the same to map[string]interface

 In case the object itself is of primitive data type, it will not convert it to map, rather it will just return the same object as is.

 Parameters:
  node: any object

 Returns:
  possible return value types could be primitive type, map (map[string]interface{}) or slice ([]interface{})

*/
func serilizeToMap(node interface{}) interface{} {
	var elementType reflect.Type
	var elementValue reflect.Value
	var ptrValue reflect.Value
	nodeType := reflect.TypeOf(node)
	nodeValue := reflect.ValueOf(node)
	// If the first object itself is the pointer then get the underlying object 'Value' and process it.
	if nodeType.Kind() == reflect.Pointer {
		// NOTE: This handles only one level of pointer. At this moment we don't expect to get pointer to pointer.
		//This will get 'reflect.Value' object pointed to by this pointer.
		elementType = nodeType.Elem()
		//This will get 'reflect.Type' object pointed to by this pointer
		elementValue = nodeValue.Elem()
		ptrValue = nodeValue
	} else {
		elementType = nodeType
		elementValue = nodeValue
	}
	switch elementType.Kind() {
	case reflect.String, reflect.Int, reflect.Bool, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if elementValue.IsValid() {
			return elementValue.Interface()
		}
		return nil
	case reflect.Struct:
		if elementValue.IsValid() {
			return processStruct(elementValue.Interface(), ptrValue)
		}
		return nil
	case reflect.Map:
		if elementValue.IsValid() {
			return processMap(elementValue.Interface())
		}
		return nil
	case reflect.Array, reflect.Slice:
		if elementValue.IsValid() {
			return processArrayOrSlice(elementValue.Interface())
		}
		return nil
	default:
		log.SetPrefix("[WARNING]")
		log.Println(getLogPrefix(), elementType.Kind(), " - not handled")
		return elementValue.Interface()
	}
}

// build

//  go build -buildmode=c-shared -o lib-sample.dylib sample.go
