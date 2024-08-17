// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.
package receptor_sdk

import "encoding/json"

func ConfigToMap(configs []Config) (map[string]interface{}, error) {
	jsonData, err := json.Marshal(configs)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
