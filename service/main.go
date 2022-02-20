package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"golang-elma-course/cyclic_rotation"
	"golang-elma-course/search_element"
	"golang-elma-course/sequence_check"
	"golang-elma-course/unique_element"
	"io"
	"net/http"
	"net/url"
	"time"
)

func main() {

	err := http.ListenAndServe(":8080", router())
	if err != nil {
		panic(err)
	}

}

const (
	SolutionServerPath = "http://116.203.203.76:3000"
	UserName           = "Kanifit"
	cyclicRotation     = "Циклическая ротация"
	uniqueElement      = "Чудные вхождения в массив"
	sequenceCheck      = "Проверка последовательности"
	searchElement      = "Поиск отсутствующего элемента"
)

type Data struct {
	Set   []int
	Shift int
}

type TaskResults struct {
	Payload [][]interface{} `json:"payload"`
	Results []interface{}   `json:"results"`
}

type Solution struct {
	UserName string      `json:"user_name"`
	Task     string      `json:"task"`
	Results  TaskResults `json:"results"`
}

func router() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	client := &http.Client{Timeout: time.Minute}

	router.Get("/tasks", func(writer http.ResponseWriter, request *http.Request) {

		tasks := []string{cyclicRotation, uniqueElement, sequenceCheck, searchElement}
		results := map[string]string{}

		for _, taskName := range tasks {
			result, err := solveTask(taskName, client, request)
			if err != nil {
				writer.WriteHeader(http.StatusBadRequest)
				return
			}
			results[taskName] = string(result)
		}

		res, err := json.Marshal(results)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		_, _ = writer.Write(res)
	})

	router.Get("/task/{taskName}", func(writer http.ResponseWriter, request *http.Request) {

		var taskName string
		switch chi.URLParam(request, "taskName") {
		case cyclicRotation:
			taskName = cyclicRotation
		case uniqueElement:
			taskName = uniqueElement
		case sequenceCheck:
			taskName = sequenceCheck
		case searchElement:
			taskName = searchElement
		default:
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		result, err := solveTask(taskName, client, request)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		_, _ = writer.Write(result)

	})

	return router
}

func getDataSets(taskName string, client *http.Client, request *http.Request) ([]byte, error) {
	path, _ := url.Parse(SolutionServerPath + "/tasks/" + taskName)
	newRequest, err := getNewRequest(path, http.MethodGet, []byte{}, request)
	if err != nil {
		return nil, err
	}

	response, err := getResponse(client, newRequest)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func getNewRequest(path *url.URL, method string, body []byte, request *http.Request) (*http.Request, error) {
	newRequest, err := http.NewRequestWithContext(
		request.Context(),
		method,
		path.String(),
		bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	return newRequest, nil
}

func getResponse(client *http.Client, request *http.Request) ([]byte, error) {
	requestResponse, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	response, err := io.ReadAll(requestResponse.Body)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(requestResponse.Body)

	return response, nil
}

func parseDataSets(dataSet []byte) ([]Data, error) {
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

func getTaskSolution(taskName string, dataSets []Data) (Solution, error) {
	var solution Solution
	solution.UserName = UserName
	solution.Task = taskName

	var taskResults TaskResults
	for _, dataSet := range dataSets {
		var taskPayload []interface{}
		taskPayload = append(taskPayload, dataSet.Set)
		if dataSet.Shift > 0 {
			taskPayload = append(taskPayload, dataSet.Shift)
		}

		taskResults.Payload = append(taskResults.Payload, taskPayload)

		switch taskName {
		case cyclicRotation:
			taskResults.Results = append(taskResults.Results, cyclic_rotation.Solution(dataSet.Set, dataSet.Shift))
		case uniqueElement:
			taskResults.Results = append(taskResults.Results, unique_element.Solution(dataSet.Set))
		case sequenceCheck:
			taskResults.Results = append(taskResults.Results, sequence_check.Solution(dataSet.Set))
		case searchElement:
			taskResults.Results = append(taskResults.Results, search_element.Solution(dataSet.Set))
		default:
			return solution, errors.New("task name is not found")
		}
	}
	solution.Results = taskResults

	return solution, nil
}

func getResult(solution Solution, client *http.Client, request *http.Request) ([]byte, error) {
	raw, err := json.Marshal(solution)
	if err != nil {
		return nil, err
	}

	path, _ := url.Parse(SolutionServerPath + "/tasks/solution")

	newRequest, err := getNewRequest(path, http.MethodPost, raw, request)
	if err != nil {
		return nil, err
	}

	response, err := getResponse(client, newRequest)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func solveTask(taskName string, client *http.Client, request *http.Request) ([]byte, error) {
	dataSets, err := getDataSets(taskName, client, request)
	if err != nil {
		return nil, err
	}

	data, err := parseDataSets(dataSets)
	if err != nil {
		return nil, err
	}

	solution, err := getTaskSolution(taskName, data)
	if err != nil {
		return nil, err
	}

	result, err := getResult(solution, client, request)
	if err != nil {
		return nil, err
	}

	return result, nil
}
