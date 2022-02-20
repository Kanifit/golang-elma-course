package main

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"golang-elma-course/cyclic_rotation"
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

const SolutionServerPath = "http://116.203.203.76:3000"

const UserName = "Kanifit"

type Data struct {
	Set   []int
	Shift int
}

type TaskResults struct {
	Payload [][]interface{} `json:"payload"`
	Results [][]int         `json:"results"`
}

type Solution struct {
	UserName string      `json:"user_name"`
	Task     string      `json:"task"`
	Results  TaskResults `json:"results"`
}

func router() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/task/Циклическая ротация", func(writer http.ResponseWriter, request *http.Request) {

		client := &http.Client{Timeout: time.Minute}

		dataSets, err := getDataSets(client, request)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		data, err := parseDataSets(dataSets)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		solution := getTaskSolution(data)

		result, err := getResult(solution, client, request)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		_, _ = writer.Write(result)

	})

	return router
}

func getDataSets(client *http.Client, request *http.Request) ([]byte, error) {
	path, _ := url.Parse(SolutionServerPath + "/tasks/Циклическая ротация")
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

			dataSet.Shift = int(element[1].(float64))
			dataSets = append(dataSets, dataSet)
		}
	}

	return dataSets, nil
}

func getTaskSolution(dataSets []Data) Solution {
	var solution Solution
	solution.UserName = UserName
	solution.Task = "Циклическая ротация"

	var taskResults TaskResults
	for _, dataSet := range dataSets {
		var taskPayload []interface{}
		taskResults.Payload = append(taskResults.Payload, append(taskPayload, dataSet.Set, dataSet.Shift))
		taskResults.Results = append(taskResults.Results, cyclic_rotation.Solution(dataSet.Set, dataSet.Shift))
	}
	solution.Results = taskResults

	return solution
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

	return response, nil
}
