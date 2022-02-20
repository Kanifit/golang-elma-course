//Package router управление запросами
package router

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"golang-elma-course/service/task"
	"net/http"
	"time"
)

//Router маршрутизатор запросов к серверу
func Router() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	client := &http.Client{Timeout: time.Minute}

	router.Get("/tasks", func(writer http.ResponseWriter, request *http.Request) {

		tasks := []string{task.CyclicRotation, task.UniqueElement, task.SequenceCheck, task.SearchElement}
		results := map[string]string{}

		for _, taskName := range tasks {
			result, err := task.Solve(taskName, client, request)
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
		case task.CyclicRotation:
			taskName = task.CyclicRotation
		case task.UniqueElement:
			taskName = task.UniqueElement
		case task.SequenceCheck:
			taskName = task.SequenceCheck
		case task.SearchElement:
			taskName = task.SearchElement
		default:
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		result, err := task.Solve(taskName, client, request)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		_, _ = writer.Write(result)

	})

	return router
}
