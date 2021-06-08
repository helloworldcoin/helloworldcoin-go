package JsonUtil

import (
	"encoding/json"
	"fmt"
)

func toJson(object interface{}) string {
	jsonStu, err := json.Marshal(object)
	return jsonStu
}

func fromJson() {

}
