package json

import (
	"strings"
)

func StringToJson(s string) *map[string]string {
	// This is a
	// naive implementation, but it works for now for json. update later
	inBody := false
	jsonMap := make(map[string]string)
	for _, line := range strings.Split(s, "\n") {
		if line == "\r" {
			break
		}
		lineTrimmed := strings.TrimLeft(line, " ")
		if string(lineTrimmed[0]) == "{" {
			inBody = true
		} else if string(lineTrimmed[0]) == "}" {
			break
		} else if inBody {
			body := strings.Split(string(lineTrimmed), ": ")
			trimVal := strings.Trim(body[1], ",")
			jsonMap[strings.Trim(body[0], "\"")] = strings.Trim(trimVal, "\"")
		}
	}
	return &jsonMap

}

func JsonToString(jsonMap *map[string]string) string {
	jsonString := "{\n"
	for key, value := range *jsonMap {
		jsonString += "\t\"" + key + "\": \"" + value + "\",\n"
	}
	jsonString += "}"
	return jsonString
}
