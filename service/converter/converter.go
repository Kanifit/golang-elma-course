//Package converter для работы с json
package converter

import "encoding/json"

//Data сущность одного набора данных
type Data struct {
	Set   []int
	Shift int
}

//ParseDataSets парсит входные данные
func ParseDataSets(dataSet []byte) ([]Data, error) {
	var decodedDataSet []interface{}
	err := json.Unmarshal(dataSet, &decodedDataSet)
	if err != nil {
		return nil, err
	}

	var dataSets []Data
	for _, data := range decodedDataSet {
		if element, ok := data.([]interface{}); ok {
			var dataSet Data
			for _, setElement := range element[0].([]interface{}) {
				dataSet.Set = append(dataSet.Set, int(setElement.(float64)))
			}

			if len(element) > 1 {
				dataSet.Shift = int(element[1].(float64))
			}

			dataSets = append(dataSets, dataSet)
		}
	}

	return dataSets, nil
}
