package handlers

import (
	"Attestation_work/pkg/store"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type Handler struct {
	Storage *store.Store
}

func GetHandler() *Handler {
	var handler Handler
	handler.Storage = store.GetStore()
	return &handler

}

func (h *Handler) CityView(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var res []byte
	switch {
	case r.Method == "GET":
		index, err := getCity(h, r)
		if err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}
		response := h.Storage.Storage[index]
		res, _ = json.Marshal(response)
	case r.Method == "DELETE":
		index, err := getCity(h, r)
		if err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}
		response := h.Storage.Storage[index]
		remove(h.Storage.Storage, index)
		res, _ = json.Marshal(response)
	case r.Method == "PATCH":
		response := changePopulationCity(h, r)
		res, _ = json.Marshal(response)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (h *Handler) GetInfoCityView(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "GET":
		w.Header().Set("Content-Type", "application/json")

		var response []store.City
		var err error

		region, regionOk := r.URL.Query()["Region"]
		district, districtOk := r.URL.Query()["District"]
		populationFrom, populationFromOk := r.URL.Query()["PopulationFrom"]
		populationTo, populationToOk := r.URL.Query()["PopulationTo"]
		foundationFrom, foundationFromOk := r.URL.Query()["FoundationFrom"]
		foundationTo, foundationToOk := r.URL.Query()["FoundationTo"]

		switch true {
		case regionOk:
			response = getCityListByRegion(h, r, strings.Join(region, ""))
		case districtOk:
			response = getCityListByDistrict(h, r, strings.Join(district, ""))
		case populationFromOk && populationToOk:
			response, err = getCityListByPopulation(h, r, strings.Join(populationFrom, ""), strings.Join(populationTo, ""))
			if err != nil {
				render.Render(w, r, ErrInvalidRequest(err))
				return
			}
		case foundationFromOk && foundationToOk:
			response, err = getCityListByFoundation(h, r, strings.Join(foundationFrom, ""), strings.Join(foundationTo, ""))
			if err != nil {
				render.Render(w, r, ErrInvalidRequest(err))
				return
			}
		default:
			render.Render(w, r, ErrInvalidRequest(errors.New("ошибка в параметрах запроса")))
			return
		}
		w.WriteHeader(http.StatusOK)
		res, _ := json.Marshal(map[string][]store.City{"result": response})
		w.Write(res)
	default:
		res, _ := json.Marshal(&ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found."})
		w.Write(res)
	}
}

func getCityListByRegion(h *Handler, r *http.Request, value string) []store.City {
	result := make([]store.City, 0)
	for _, v := range h.Storage.Storage {
		if v.Region == value {
			result = append(result, v)
		}
	}
	return result
}

func getCityListByPopulation(h *Handler, r *http.Request, FoundationFrom string, FoundationTo string) ([]store.City, error) {
	start, errStart := strconv.Atoi(FoundationFrom)
	end, errEnd := strconv.Atoi(FoundationTo)
	if errStart != nil || errEnd != nil {
		return nil, errors.New("ошибка в параметрах запроса")
	}
	result := make([]store.City, 0)
	for _, v := range h.Storage.Storage {
		if v.Population >= start && v.Population <= end {
			result = append(result, v)
		}
	}
	return result, nil
}

func getCityListByFoundation(h *Handler, r *http.Request, PopulationFrom string, PopulationTo string) ([]store.City, error) {
	start, errStart := strconv.Atoi(PopulationFrom)
	end, errEnd := strconv.Atoi(PopulationTo)
	if errStart != nil || errEnd != nil {
		return nil, errors.New("ошибка в параметрах запроса")
	}
	result := make([]store.City, 0)
	for _, v := range h.Storage.Storage {
		if v.Foundation >= start && v.Foundation <= end {
			result = append(result, v)
		}
	}
	return result, nil
}

func getCityListByDistrict(h *Handler, r *http.Request, value string) []store.City {
	result := make([]store.City, 0)
	for _, v := range h.Storage.Storage {
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
	index, _ := findCity(h, city.Id)
	if index != -1 {
		result["status"] = "Такой город уже есть в списке!"
		return result
	}
	h.Storage.Storage = append(h.Storage.Storage, *city)
	result["status"] = "Город успешно добавлен!"
	return result
}

func getCity(h *Handler, r *http.Request) (int, error) {
	idParam := chi.URLParam(r, "Id")
	cityId, err := strconv.Atoi(idParam)
	if err != nil {
		return -1, err
	}
	index, err := findCity(h, cityId)
	if err != nil {
		return -1, err
	}
	return index, nil
}

func changePopulationCity(h *Handler, r *http.Request) map[string]string {
	result := make(map[string]string)
	idParam := chi.URLParam(r, "Id")
	populationParam := r.Header.Get("population")
	cityId, errId := strconv.Atoi(idParam)
	newPopulation, errPop := strconv.Atoi(populationParam)
	if errId != nil || errPop != nil {
		result["error"] = "Ошибка в в параметрах запроса!"
		return result
	}
	for index, v := range h.Storage.Storage {
		if v.Id == cityId {
			h.Storage.Storage[index].Population = newPopulation
			result["status"] = "Население успешно изменено!"
			return result
		}
	}
	result["status"] = "Указанный город не найден!"
	return result
}

func findCity(h *Handler, cityId int) (int, error) {
	for index, v := range h.Storage.Storage {
		if v.Id == cityId {
			return index, nil
		}
	}
	return -1, errors.New("город с таким id не найден")
}

func remove(s []store.City, i int) []store.City {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
