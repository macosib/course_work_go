package handlers

import (
	"Attestation_work/pkg/store"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type Handler struct {
	storage *store.Store
}

func GetHandler() Handler {
	var handler Handler
	handler.storage = store.GetStore()
	return handler

}

func (h *Handler) CityView(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "GET":
		fmt.Println(h.storage)
		w.Header().Set("Content-Type", "application/json")
		response, _, err := getCity(h, r)
		if err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}
		w.WriteHeader(http.StatusOK)
		res, _ := json.Marshal(response)
		w.Write(res)
	case r.Method == "DELETE":
		w.Header().Set("Content-Type", "application/json")
		response, index, err := getCity(h, r)
		if err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}
		remove(*h.storage.Store, index)
		w.WriteHeader(http.StatusOK)
		res, _ := json.Marshal(response)
		w.Write(res)
	case r.Method == "PATCH":
		w.Header().Set("Content-Type", "application/json")
		response, err := changePopulationCity(h, r)
		if err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}
		res, _ := json.Marshal(response)
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	default:
		res, _ := json.Marshal(&ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found."})
		w.Write(res)

	}
}

func (h *Handler) GetInfoCityView(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "GET":
		w.Header().Set("Content-Type", "application/json")
		var response []store.City
		var err error

		value, ok := r.URL.Query()["Region"]
		if ok {
			response = getCityListByRegion(h, r, strings.Join(value, ""))
		}

		value, ok = r.URL.Query()["District"]
		if ok {
			response = getCityListByDistrict(h, r, strings.Join(value, ""))
		}

		populationFrom, populationFromOk := r.URL.Query()["PopulationFrom"]
		populationTo, populationToOk := r.URL.Query()["PopulationTo"]
		if populationFromOk && populationToOk {
			response, err = getCityListByPopulation(h, r, strings.Join(populationFrom, ""), strings.Join(populationTo, ""))
			if err != nil {
				render.Render(w, r, ErrInvalidRequest(err))
				return
			}
		}

		foundationFrom, foundationFromOk := r.URL.Query()["FoundationFrom"]
		foundationTo, foundationToOk := r.URL.Query()["FoundationTo"]
		if foundationFromOk && foundationToOk {
			response, err = getCityListByFoundation(h, r, strings.Join(foundationFrom, ""), strings.Join(foundationTo, ""))
			if err != nil {
				render.Render(w, r, ErrInvalidRequest(err))
				return
			}
		}

		if response == nil {
			render.Render(w, r, ErrInvalidRequest(errors.New("Ошибка в параметрах запроса")))
			return
		}

		w.WriteHeader(http.StatusOK)
		res, _ := json.Marshal(response)
		w.Write(res)
	default:
		res, _ := json.Marshal(&ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found."})
		w.Write(res)
	}
}

func getCityListByRegion(h *Handler, r *http.Request, value string) []store.City {
	result := make([]store.City, 0)
	for _, v := range *h.storage.Store {
		if v.Region == value {
			result = append(result, v)
		}
	}
	return result
}

func getCityListByPopulation(h *Handler, r *http.Request, FoundationFrom string, FoundationTo string) ([]store.City, error) {
	start, err := strconv.Atoi(FoundationFrom)
	if err != nil {
		return nil, err
	}
	end, err := strconv.Atoi(FoundationTo)
	if err != nil {
		return nil, err
	}
	result := make([]store.City, 0)
	for _, v := range *h.storage.Store {
		if v.Population >= start && v.Population <= end {
			result = append(result, v)
		}
	}
	return result, nil
}

func getCityListByFoundation(h *Handler, r *http.Request, PopulationFrom string, PopulationTo string) ([]store.City, error) {
	start, err := strconv.Atoi(PopulationFrom)
	if err != nil {
		return nil, err
	}
	end, err := strconv.Atoi(PopulationTo)
	if err != nil {
		return nil, err
	}
	result := make([]store.City, 0)
	for _, v := range *h.storage.Store {
		if v.Foundation >= start && v.Foundation <= end {
			result = append(result, v)
		}
	}
	return result, nil
}

func getCityListByDistrict(h *Handler, r *http.Request, value string) []store.City {
	result := make([]store.City, 0)
	for _, v := range *h.storage.Store {
		if v.District == value {
			result = append(result, v)
		}
	}
	return result
}

func (h *Handler) AddCityView(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "POST":
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}
		defer r.Body.Close()
		var data store.City
		if err := json.Unmarshal(content, &data); err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}
		newCity := addCity(h, r, &data)
		w.WriteHeader(http.StatusOK)
		res, _ := json.Marshal(newCity)
		w.Write(res)
	default:
		res, _ := json.Marshal(&ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found."})
		w.Write(res)

	}
}

func addCity(h *Handler, r *http.Request, city *store.City) map[string]string {
	result := make(map[string]string)
	_, index, _ := findCity(h, city.Id)
	if index != -1 {
		result["error"] = "Такой город уже есть в списке!"
		return result
	}
	newStore := append(*h.storage.Store, *city)
	h.storage.Store = &newStore
	result["status"] = "Город успешно добавлен!"
	return result
}

func getCity(h *Handler, r *http.Request) (*store.City, int, error) {
	idParam := chi.URLParam(r, "Id")
	cityId, err := strconv.Atoi(idParam)
	if err != nil {
		return nil, -1, err
	}
	city, index, err := findCity(h, cityId)
	if err != nil {
		return nil, -1, err
	}
	return city, index, nil
}

func changePopulationCity(h *Handler, r *http.Request) (*store.City, error) {
	idParam := chi.URLParam(r, "Id")
	populationParam := r.Header.Get("population")
	cityId, err := strconv.Atoi(idParam)
	if err != nil {
		return nil, err
	}
	newPopulation, err := strconv.Atoi(populationParam)
	if err != nil {
		return nil, err
	}
	for index, v := range *h.storage.Store {
		if v.Id == cityId {
			v.Population = newPopulation
			remove(*h.storage.Store, index)
			newStore := append(*h.storage.Store, v)
			h.storage.Store = &newStore
			return &v, nil
		}
	}
	return nil, err
}

func findCity(h *Handler, cityId int) (*store.City, int, error) {
	for index, v := range *h.storage.Store {
		if v.Id == cityId {
			return &v, index, nil
		}
	}
	return nil, -1, errors.New("Город с таким id не найден")
}

// func remove(slice []store.City, s int) []store.City {
// 	return append(slice[:s], slice[s+1:]...)
// }

func remove(s []store.City, i int) []store.City {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
