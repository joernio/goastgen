package goastgen

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMapWithNilPointerCheck(t *testing.T) {
	var nilStr *string
	var nilObj *Phone
	mapType := make(map[string]interface{})
	mapType["first"] = "first value"
	mapType["second"] = nilStr
	mapType["third"] = nilObj

	result := processMap(mapType)

	expectedResult := make(map[string]interface{})
	expectedResult["first"] = "first value"

	assert.Equal(t, expectedResult, result, "It should return map with valid values only")

}

func TestArrayOfPointerOfMapOfObjectPointerType(t *testing.T) {
	first := Phone{PhoneNo: "1234567890", Type: "Home"}
	second := Phone{PhoneNo: "0987654321", Type: "Office"}
	third := Phone{PhoneNo: "1234567891", Type: "Home1"}
	forth := Phone{PhoneNo: "1987654321", Type: "Office1"}
	firstMap := make(map[string]*Phone)
	firstMap["fmfirst"] = &first
	firstMap["fmsecond"] = &second
	secondMap := make(map[string]*Phone)
	secondMap["smfirst"] = &third
	secondMap["smsecond"] = &forth
	array := [2]*map[string]*Phone{&firstMap, &secondMap}
	result := processArrayOrSlice(array)
	firstPhone := make(map[string]interface{})
	firstPhone["PhoneNo"] = "1234567890"
	firstPhone["Type"] = "Home"
	firstPhone["node_type"] = "goastgen.Phone"
	secondPhone := make(map[string]interface{})
	secondPhone["PhoneNo"] = "0987654321"
	secondPhone["Type"] = "Office"
	secondPhone["node_type"] = "goastgen.Phone"
	thirdPhone := make(map[string]interface{})
	thirdPhone["PhoneNo"] = "1234567891"
	thirdPhone["Type"] = "Home1"
	thirdPhone["node_type"] = "goastgen.Phone"
	forthPhone := make(map[string]interface{})
	forthPhone["PhoneNo"] = "1987654321"
	forthPhone["Type"] = "Office1"
	forthPhone["node_type"] = "goastgen.Phone"
	firstExpectedMap := make(map[string]interface{})
	firstExpectedMap["fmfirst"] = firstPhone
	firstExpectedMap["fmsecond"] = secondPhone
	secondExpectedMap := make(map[string]interface{})
	secondExpectedMap["smfirst"] = thirdPhone
	secondExpectedMap["smsecond"] = forthPhone
	expectedResult := []interface{}{firstExpectedMap, secondExpectedMap}

	assert.Equal(t, expectedResult, result, "Array of Map of simple string should match with expected result")
}

func TestArrayOfPointerOfMapOfPrimitivesType(t *testing.T) {
	firstMap := make(map[string]string)
	firstMap["fmfirst"] = "fmfirstvalue"
	firstMap["fmsecond"] = "fmsecondvalue"
	secondMap := make(map[string]string)
	secondMap["smfirst"] = "smfirstvalue"
	secondMap["smsecond"] = "smsecondvalue"
	array := [2]*map[string]string{&firstMap, &secondMap}
	result := processArrayOrSlice(array)

	firstExpectedMap := make(map[string]interface{})
	firstExpectedMap["fmfirst"] = "fmfirstvalue"
	firstExpectedMap["fmsecond"] = "fmsecondvalue"
	secondExpectedMap := make(map[string]interface{})
	secondExpectedMap["smfirst"] = "smfirstvalue"
	secondExpectedMap["smsecond"] = "smsecondvalue"
	expectedResult := []interface{}{firstExpectedMap, secondExpectedMap}

	assert.Equal(t, expectedResult, result, "Array of Map of simple string should match with expected result")
}

func TestArrayOfMapOfPrimitivesType(t *testing.T) {
	firstMap := make(map[string]string)
	firstMap["fmfirst"] = "fmfirstvalue"
	firstMap["fmsecond"] = "fmsecondvalue"
	secondMap := make(map[string]string)
	secondMap["smfirst"] = "smfirstvalue"
	secondMap["smsecond"] = "smsecondvalue"
	array := [2]map[string]string{firstMap, secondMap}
	result := processArrayOrSlice(array)

	firstExpectedMap := make(map[string]interface{})
	firstExpectedMap["fmfirst"] = "fmfirstvalue"
	firstExpectedMap["fmsecond"] = "fmsecondvalue"
	secondExpectedMap := make(map[string]interface{})
	secondExpectedMap["smfirst"] = "smfirstvalue"
	secondExpectedMap["smsecond"] = "smsecondvalue"
	expectedResult := []interface{}{firstExpectedMap, secondExpectedMap}

	assert.Equal(t, expectedResult, result, "Array of Map of simple string should match with expected result")

}

func TestMapObjPtrType(t *testing.T) {
	first := Phone{PhoneNo: "1234567890", Type: "Home"}
	second := Phone{PhoneNo: "0987654321", Type: "Office"}
	phones := make(map[string]*Phone)
	phones["first"] = &first
	phones["second"] = &second

	mapType := MapObjPtrType{Id: 90, Phones: phones}
	result := serilizeToMap(mapType)
	expectedResult := make(map[string]interface{})
	expectedResult["Id"] = 90
	expectedResult["node_type"] = "goastgen.MapObjPtrType"
	expectedPhones := make(map[string]interface{})
	firstPhone := make(map[string]interface{})
	firstPhone["PhoneNo"] = "1234567890"
	firstPhone["Type"] = "Home"
	firstPhone["node_type"] = "goastgen.Phone"
	secondPhone := make(map[string]interface{})
	secondPhone["PhoneNo"] = "0987654321"
	secondPhone["Type"] = "Office"
	secondPhone["node_type"] = "goastgen.Phone"
	expectedPhones["first"] = firstPhone
	expectedPhones["second"] = secondPhone
	expectedResult["Phones"] = expectedPhones

	assert.Equal(t, expectedResult, result, "Map with Object Pointer type result Map should match with expected result Map")

	jsonResult := serilizeToJsonStr(result)
	expectedJsonResult := "{\n  \"Id\": 90,\n  \"Phones\": {\n    \"first\": {\n      \"PhoneNo\": \"1234567890\",\n      \"Type\": \"Home\",\n      \"node_type\": \"goastgen.Phone\"\n    },\n    \"second\": {\n      \"PhoneNo\": \"0987654321\",\n      \"Type\": \"Office\",\n      \"node_type\": \"goastgen.Phone\"\n    }\n  },\n  \"node_type\": \"goastgen.MapObjPtrType\"\n}"
	assert.Equal(t, expectedJsonResult, jsonResult, "Map with Object type result json should match with expected result")
}

func TestMapStrPtrType(t *testing.T) {
	first := "firstvalue"
	second := "secondvalue"
	names := make(map[string]*string)
	names["firstname"] = &first
	names["secondname"] = &second
	mapType := MapStrPtrType{Id: 30, Names: names}
	result := serilizeToMap(mapType)
	expectedResult := make(map[string]interface{})
	expectedResult["Id"] = 30
	expectedResult["node_type"] = "goastgen.MapStrPtrType"
	expectedNames := make(map[string]interface{})
	expectedNames["firstname"] = "firstvalue"
	expectedNames["secondname"] = "secondvalue"
	expectedResult["Names"] = expectedNames
	assert.Equal(t, expectedResult, result, "Map String pointer type result Map should match with expected result Map")

	jsonResult := serilizeToJsonStr(result)
	expectedJsonResult := "{\n  \"Id\": 30,\n  \"Names\": {\n    \"firstname\": \"firstvalue\",\n    \"secondname\": \"secondvalue\"\n  },\n  \"node_type\": \"goastgen.MapStrPtrType\"\n}"
	assert.Equal(t, expectedJsonResult, jsonResult, "Map with String pointer type result json should match with expected result")
}

func TestMapObjType(t *testing.T) {
	phones := make(map[string]Phone)
	phones["first"] = Phone{PhoneNo: "1234567890", Type: "Home"}
	phones["second"] = Phone{PhoneNo: "0987654321", Type: "Office"}

	mapType := MapObjType{Id: 90, Phones: phones}
	result := serilizeToMap(mapType)
	expectedResult := make(map[string]interface{})
	expectedResult["Id"] = 90
	expectedResult["node_type"] = "goastgen.MapObjType"
	expectedPhones := make(map[string]interface{})
	firstPhone := make(map[string]interface{})
	firstPhone["PhoneNo"] = "1234567890"
	firstPhone["Type"] = "Home"
	firstPhone["node_type"] = "goastgen.Phone"
	secondPhone := make(map[string]interface{})
	secondPhone["PhoneNo"] = "0987654321"
	secondPhone["Type"] = "Office"
	secondPhone["node_type"] = "goastgen.Phone"
	expectedPhones["first"] = firstPhone
	expectedPhones["second"] = secondPhone
	expectedResult["Phones"] = expectedPhones

	assert.Equal(t, expectedResult, result, "Map with Object type result Map should match with expected result Map")

	jsonResult := serilizeToJsonStr(result)
	expectedJsonResult := "{\n  \"Id\": 90,\n  \"Phones\": {\n    \"first\": {\n      \"PhoneNo\": \"1234567890\",\n      \"Type\": \"Home\",\n      \"node_type\": \"goastgen.Phone\"\n    },\n    \"second\": {\n      \"PhoneNo\": \"0987654321\",\n      \"Type\": \"Office\",\n      \"node_type\": \"goastgen.Phone\"\n    }\n  },\n  \"node_type\": \"goastgen.MapObjType\"\n}"
	assert.Equal(t, expectedJsonResult, jsonResult, "Map with Object type result json should match with expected result")
}

func TestMapIntType(t *testing.T) {
	names := make(map[string]int)
	names["firstname"] = 1000
	names["secondname"] = 2000
	mapType := MapIntType{Id: 30, Names: names}
	result := serilizeToMap(mapType)
	expectedResult := make(map[string]interface{})
	expectedResult["Id"] = 30
	expectedResult["node_type"] = "goastgen.MapIntType"
	expectedNames := make(map[string]interface{})
	expectedNames["firstname"] = 1000
	expectedNames["secondname"] = 2000
	expectedResult["Names"] = expectedNames
	assert.Equal(t, expectedResult, result, "Simple Map type result Map should match with expected result Map")

	jsonResult := serilizeToJsonStr(result)
	expectedJsonResult := "{\n  \"Id\": 30,\n  \"Names\": {\n    \"firstname\": 1000,\n    \"secondname\": 2000\n  },\n  \"node_type\": \"goastgen.MapIntType\"\n}"
	assert.Equal(t, expectedJsonResult, jsonResult, "Simple Map type result json should match with expected result")
}

func TestMapType(t *testing.T) {
	names := make(map[string]string)
	names["firstname"] = "firstvalue"
	names["secondname"] = "secondvalue"
	mapType := MapType{Id: 30, Names: names}
	result := serilizeToMap(mapType)
	expectedResult := make(map[string]interface{})
	expectedResult["Id"] = 30
	expectedResult["node_type"] = "goastgen.MapType"
	expectedNames := make(map[string]interface{})
	expectedNames["firstname"] = "firstvalue"
	expectedNames["secondname"] = "secondvalue"
	expectedResult["Names"] = expectedNames
	assert.Equal(t, expectedResult, result, "Simple Map type result Map should match with expected result Map")
}

func TestSimpleMapType(t *testing.T) {
	phone1 := Phone{PhoneNo: "1234567890", Type: "Home"}
	phone2 := Phone{PhoneNo: "0987654321", Type: "Office"}

	mapType := make(map[string]*Phone)
	mapType["first"] = &phone1
	mapType["second"] = &phone2

	result := serilizeToMap(mapType)
	expectedResult := make(map[string]interface{})
	firstPhone := make(map[string]interface{})
	firstPhone["PhoneNo"] = "1234567890"
	firstPhone["Type"] = "Home"
	firstPhone["node_type"] = "goastgen.Phone"
	secondPhone := make(map[string]interface{})
	secondPhone["PhoneNo"] = "0987654321"
	secondPhone["Type"] = "Office"
	secondPhone["node_type"] = "goastgen.Phone"
	expectedResult["first"] = firstPhone
	expectedResult["second"] = secondPhone

	assert.Equal(t, expectedResult, result, "Map type with object pointer values should match with expected results")
}
