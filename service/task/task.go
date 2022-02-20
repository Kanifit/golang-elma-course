//Package task работа с задачей
package task

import (
	"encoding/json"
	"errors"
	"golang-elma-course/service/converter"
	"golang-elma-course/service/http/constructor"
	"golang-elma-course/service/solver/cyclic_rotation"
	"golang-elma-course/service/solver/search_element"
	"golang-elma-course/service/solver/sequence_check"
	"golang-elma-course/service/solver/unique_element"
	"net/http"
	"net/url"
)

const (
	solutionServerPath = "http://116.203.203.76:3000"
	userName           = "Kanifit"
	CyclicRotation     = "Циклическая ротация"
	UniqueElement      = "Чудные вхождения в массив"
	SequenceCheck      = "Проверка последовательности"
	SearchElement      = "Поиск отсутствующего элемента"
)

//Solution сущность данных решения
type Solution struct {
	UserName string  `json:"user_name"`
	Task     string  `json:"task"`
	Results  Results `json:"results"`
}

//Results сущность результата решения задачи
type Results struct {
	Payload [][]interface{} `json:"payload"`
	Results []interface{}   `json:"results"`
}

//getDataSets получает наборы данных
func getDataSets(taskName string, client *http.Client, request *http.Request) ([]byte, error) {
	path, _ := url.Parse(solutionServerPath + "/tasks/" + taskName)
	newRequest, err := constructor.GetNewRequest(path, http.MethodGet, []byte{}, request)
	if err != nil {
		return nil, err
	}

	response, err := constructor.GetResponse(client, newRequest)
	if err != nil {
		return nil, err
	}

	return response, nil
}

//getSolution получает решение задачи
func getSolution(taskName string, dataSets []converter.Data) (Solution, error) {
	var solution Solution
	solution.UserName = userName
	solution.Task = taskName

	var taskResults Results
	for _, dataSet := range dataSets {
		var taskPayload []interface{}
		taskPayload = append(taskPayload, dataSet.Set)
		if dataSet.Shift > 0 {
			taskPayload = append(taskPayload, dataSet.Shift)
		}

		taskResults.Payload = append(taskResults.Payload, taskPayload)

		switch taskName {
		case CyclicRotation:
			taskResults.Results = append(taskResults.Results, cyclic_rotation.Solution(dataSet.Set, dataSet.Shift))
		case UniqueElement:
			taskResults.Results = append(taskResults.Results, unique_element.Solution(dataSet.Set))
		case SequenceCheck:
			taskResults.Results = append(taskResults.Results, sequence_check.Solution(dataSet.Set))
		case SearchElement:
			taskResults.Results = append(taskResults.Results, search_element.Solution(dataSet.Set))
		default:
			return solution, errors.New("task name is not found")
		}
	}
	solution.Results = taskResults

	return solution, nil
}

//getResult получает результат проверки на сервисе "Solution"
func getResult(solution Solution, client *http.Client, request *http.Request) ([]byte, error) {
	raw, err := json.Marshal(solution)
	if err != nil {
		return nil, err
	}

	path, _ := url.Parse(solutionServerPath + "/tasks/solution")

	newRequest, err := constructor.GetNewRequest(path, http.MethodPost, raw, request)
	if err != nil {
		return nil, err
	}

	response, err := constructor.GetResponse(client, newRequest)
	if err != nil {
		return nil, err
	}

	return response, nil
}

//Solve решает конкретную задачу
func Solve(taskName string, client *http.Client, request *http.Request) ([]byte, error) {
	dataSets, err := getDataSets(taskName, client, request)
	if err != nil {
		return nil, err
	}

	data, err := converter.ParseDataSets(dataSets)
	if err != nil {
		return nil, err
	}

	solution, err := getSolution(taskName, data)
	if err != nil {
		return nil, err
	}

	result, err := getResult(solution, client, request)
	if err != nil {
		return nil, err
	}

	return result, nil
}
