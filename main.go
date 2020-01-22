package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	clipboard "github.com/atotto/clipboard"
)

func main() {
	exampleData := []byte(`{"id": 123,
	"sku_id": 12232322,
	"retail_price": 110,
	"name": "Jacket A",
	"item_id": 1,
	"is_refund": true,
	"data":{
		"test": 1,
		"test-s": "a",
		"bool": true,
		"arr": [1, 2]
	}}`)
	var value interface{}

	json.Unmarshal(exampleData, &value)

	oneLevelJSON := value.(map[string]interface{})

	converToYAML(oneLevelJSON)
}

func converToYAML(oneLevelJSON map[string]interface{}) {
	var yamlSchema string

	for k, v := range oneLevelJSON {
		var dataType string

		switch v.(type) {
		case string:
			dataType = "string"
			yamlSchema += fmt.Sprintf("%v:\n  type: %v\n", k, dataType)
		case int, int16, int32, int64, int8, float32, float64, uint, uint16, uint32, uint64, uint8:
			dataType = "integer"
			yamlSchema += fmt.Sprintf("%v:\n  type: %v\n", k, dataType)
		case bool:
			dataType = "boolean"
			yamlSchema += fmt.Sprintf("%v:\n  type: %v\n", k, dataType)
		case []interface{}:
			temp := (v.([]interface{}))
			var ex string
			for _, v := range temp {
				switch v.(type) {
				case float64:
					ex = strconv.FormatFloat(v.(float64), 'f', 0, 64)
					break
				case string:
					ex = v.(string)
					break
				}
			}

			dataType = "array\n  items: {}\n  example:\n    - " + ex
			yamlSchema += fmt.Sprintf("%v:\n  type: %v\n", k, dataType)
		case map[string]interface{}:
			properties := (v.(map[string]interface{}))
			var dataProperties string
			for k, v := range properties {
				dataProperties += parse(k, v)
			}
			yamlSchema += fmt.Sprintf("%v:\n  type: object\n  properties:\n%v", k, dataProperties)
		}
	}

	fmt.Println(yamlSchema)

	clipboard.WriteAll(yamlSchema)
}

func parse(key interface{}, value interface{}) string {
	var yamlSchema string
	switch value.(type) {
	case string:
		dataType := "string"
		yamlSchema += fmt.Sprintf("    %v:\n      type: %v\n", key, dataType)
	case int, int16, int32, int64, int8, float32, float64, uint, uint16, uint32, uint64, uint8:
		dataType := "integer"
		yamlSchema += fmt.Sprintf("    %v:\n      type: %v\n", key, dataType)
	case bool:
		dataType := "boolean"
		yamlSchema += fmt.Sprintf("    %v:\n      type: %v\n", key, dataType)
	case []interface{}:
		temp := (value.([]interface{}))
		var ex string
		for _, v := range temp {
			switch v.(type) {
			case float64:
				ex = strconv.FormatFloat(v.(float64), 'f', 0, 64)
				break
			case string:
				ex = v.(string)
				break
			}
		}

		dataType := "array\n      items: {}\n      example:\n        - " + ex
		yamlSchema += fmt.Sprintf("    %v:\n      type: %v\n", key, dataType)
	case map[string]interface{}:
		properties := (value.(map[string]interface{}))
		var dataProperties string
		for k, v := range properties {
			dataProperties += parse(k, v)
		}
		yamlSchema += fmt.Sprintf("%v:\n  type: object\n  properties:\n%v", key, dataProperties)
	}

	return yamlSchema
}
