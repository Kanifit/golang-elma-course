//Package constructor позволяет работать с http запросами
package constructor

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
)

//GetNewRequest получает новый запрос
func GetNewRequest(path *url.URL, method string, body []byte, request *http.Request) (*http.Request, error) {
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

//GetResponse получает ответ на запрос
func GetResponse(client *http.Client, request *http.Request) ([]byte, error) {
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
