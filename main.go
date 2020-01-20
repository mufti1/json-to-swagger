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
	"is_refund": true}`)
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
		case int, int16, int32, int64, int8, float32, float64, uint, uint16, uint32, uint64, uint8:
			dataType = "integer"
		case bool:
			dataType = "boolean"
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
		}

		fmt.Printf("%v:\n  type: %v\n", k, dataType)

		yamlSchema += fmt.Sprintf("%v:\n  type: %v\n", k, dataType)
	}

	clipboard.WriteAll(yamlSchema)
}

