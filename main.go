package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	clipboard "github.com/atotto/clipboard"
	"github.com/gorilla/mux"
)

type jsonRequest struct {
	JSON interface{} `json:"json_schema"`
}

func main() {
	routes := mux.NewRouter()
	server := os.Getenv("PORT")

	routes.HandleFunc("/convert", convertToSwagger).Methods("POST", "OPTIONS")
	routes.HandleFunc("/ping", ping).Methods("GET")

	log.Printf("server running on %v", server)
	err := http.ListenAndServe(":"+server, routes)
	if err != nil {
		log.Fatalf("Unable to run http server: %v", err)
	}

	log.Println("Stopping API Service...")
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")

	fmt.Fprintf(w, "pong")
}

func convertToSwagger(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")

	decoder := json.NewDecoder(r.Body)
	jsonSchema := jsonRequest{}

	err := decoder.Decode(&jsonSchema)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, err.Error())
		fmt.Println(err.Error())
		return
	}

	yamlSchema := converToYAML(jsonSchema.JSON.(map[string]interface{}))
	fmt.Fprintf(w, yamlSchema)
	fmt.Println("success")
}

func converToYAML(oneLevelJSON map[string]interface{}) string {
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
					dataType = "array\n  items: {}\n  example:\n    - " + ex
					yamlSchema += fmt.Sprintf("%v:\n  type: %v\n", k, dataType)
					break
				case string:
					ex = v.(string)
					dataType = "array\n  items: {}\n  example:\n    - " + ex
					yamlSchema += fmt.Sprintf("%v:\n  type: %v\n", k, dataType)
					break
				case map[string]interface{}:
					properties := (v.(map[string]interface{}))
					var dataProperties string
					for k, v := range properties {
						dataProperties += parse(k, v, 4, 6)
					}
					yamlSchema += fmt.Sprintf("%v:\n  type: object\n  properties:\n%v", k, dataProperties)
					break
				}
			}
		case map[string]interface{}:
			properties := (v.(map[string]interface{}))
			var dataProperties string
			for k, v := range properties {
				dataProperties += parse(k, v, 4, 6)
			}
			yamlSchema += fmt.Sprintf("%v:\n  type: object\n  properties:\n%v", k, dataProperties)
		}
	}

	// fmt.Println(yamlSchema)

	clipboard.WriteAll(yamlSchema)

	return yamlSchema
}

func parse(key interface{}, value interface{}, indentationOne int, indentationTwo int) string {
	var firstIndentation int = indentationOne
	var secondIndentation int = indentationTwo
	var yamlSchema string
	switch value.(type) {
	case string:
		dataType := "string"
		yamlSchema += fmt.Sprintf("%*s%v:\n%*stype: %v\n", firstIndentation, "", key, secondIndentation, "", dataType)
		firstIndentation += 2
		secondIndentation += 2
	case int, int16, int32, int64, int8, float32, float64, uint, uint16, uint32, uint64, uint8:
		dataType := "integer"
		yamlSchema += fmt.Sprintf("%*s%v:\n%*stype: %v\n", firstIndentation, "", key, secondIndentation, "", dataType)
		firstIndentation += 2
		secondIndentation += 2
	case bool:
		dataType := "boolean"
		yamlSchema += fmt.Sprintf("%*s%v:\n%*stype: %v\n", firstIndentation, "", key, secondIndentation, "", dataType)
		firstIndentation += 2
		secondIndentation += 2
	case []interface{}:
		temp := (value.([]interface{}))
		var ex string
		for _, v := range temp {
			switch v.(type) {
			case float64:
				ex = strconv.FormatFloat(v.(float64), 'f', 0, 64)
				dataType := fmt.Sprintf("array\n%*sitems: {}\n%*sexample:\n%*s  - %v", secondIndentation, "", secondIndentation, "", secondIndentation, "", ex)
				yamlSchema += fmt.Sprintf("%*s%v:\n%*stype: %v\n", firstIndentation, "", key, secondIndentation, "", dataType)
				firstIndentation += 2
				secondIndentation += 2
				break
			case string:
				ex = v.(string)
				dataType := fmt.Sprintf("array\n%*sitems: {}\n%*sexample:\n%*s  - %v", secondIndentation, "", secondIndentation, "", secondIndentation, "", ex)
				yamlSchema += fmt.Sprintf("%*s%v:\n%*stype: %v\n", firstIndentation, "", key, secondIndentation, "", dataType)
				firstIndentation += 2
				secondIndentation += 2
				break
			case map[string]interface{}:
				properties := (v.(map[string]interface{}))
				var dataProperties string
				for k, v := range properties {
					dataProperties += parse(k, v, firstIndentation+4, secondIndentation+4)
				}
				yamlSchema += fmt.Sprintf("%*s%v:\n%*stype: object\n%*sproperties:\n%v", firstIndentation, "", key, secondIndentation, "", secondIndentation, "", dataProperties)
				firstIndentation += 2
				secondIndentation += 2
				break
			}
		}
	case map[string]interface{}:
		properties := (value.(map[string]interface{}))
		var dataProperties string
		for k, v := range properties {
			dataProperties += parse(k, v, firstIndentation+4, secondIndentation+4)
		}
		yamlSchema += fmt.Sprintf("%*s%v:\n%*stype: object\n%*sproperties:\n%v", firstIndentation, "", key, secondIndentation, "", secondIndentation, "", dataProperties)
		firstIndentation += 2
		secondIndentation += 2
	}

	return yamlSchema
}
