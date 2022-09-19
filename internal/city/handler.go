package city

import (
	"Attestation_work/internal/handlers"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type handler struct {
	Storage *Store
}

func NewHandler(s *Store) handlers.Handler {
	return &handler{
		Storage: s,
	}
}

const (
	cityUrl     = "/api/v1/city-create"
	cityUrlInfo = "/api/v1/city"
	cityUrlId   = "/api/v1/city/:id"
)

func (h *handler) Register(r *httprouter.Router) {

	r.POST(cityUrl, h.AddCityView)
	r.GET(cityUrlId, h.CityView)
	r.GET(cityUrlInfo, h.GetInfoCityView)
	r.DELETE(cityUrlId, h.CityView)
	r.PATCH(cityUrlId, h.CityView)
}

func (h *handler) CityView(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	var res []byte
	city, err := getCity(h, r, params)
	if err != nil {
		res, _ = json.Marshal(map[string]string{"status": err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}
	switch {
	case r.Method == "GET":
		res, _ = json.Marshal(city)
	case r.Method == "DELETE":
		h.Storage.mu.Lock()
		defer h.Storage.mu.Unlock()
		delete(h.Storage.Storage, city.Id)
		res, _ = json.Marshal(city)

	case r.Method == "PATCH":
		h.Storage.mu.Lock()
		defer h.Storage.mu.Unlock()

		response, err := changePopulationCity(city, r)
		if err != nil {
			res, _ = json.Marshal(map[string]string{"status": err.Error()})
			w.WriteHeader(http.StatusBadRequest)
			w.Write(res)
			return
		}
		h.Storage.Storage[city.Id] = *response
		res, _ = json.Marshal(response)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (h *handler) AddCityView(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	h.Storage.mu.Lock()
	defer h.Storage.mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	var res []byte
	switch {
	case r.Method == "POST":
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			res, _ = json.Marshal(map[string]string{"status": err.Error()})
			w.WriteHeader(http.StatusBadRequest)
			w.Write(res)
			return
		}
		defer r.Body.Close()
		var newCity City
		if err := json.Unmarshal(content, &newCity); err != nil {
			res, _ = json.Marshal(map[string]string{"status": err.Error()})
			w.WriteHeader(http.StatusBadRequest)
			w.Write(res)
			return

		}
		err = addCity(h, &newCity)
		if err != nil {
			res, _ = json.Marshal(map[string]string{"status": err.Error()})
			w.WriteHeader(http.StatusBadRequest)
			w.Write(res)
			return
		}
		res, _ := json.Marshal(newCity)
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	}
}

func (h *handler) GetInfoCityView(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	switch {
	case r.Method == "GET":
		w.Header().Set("Content-Type", "application/json")

		var response []City
		var err error
		var res []byte

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
				res, _ = json.Marshal(map[string]string{"status": err.Error()})
				w.WriteHeader(http.StatusBadRequest)
				w.Write(res)
				return
			}
		case foundationFromOk && foundationToOk:
			response, err = getCityListByFoundation(h, r, strings.Join(foundationFrom, ""), strings.Join(foundationTo, ""))
			if err != nil {
				res, _ = json.Marshal(map[string]string{"status": err.Error()})
				w.WriteHeader(http.StatusBadRequest)
				w.Write(res)
				return
			}
		default:
			res, _ = json.Marshal(map[string]string{"status": "ошибка в параметрах запроса"})
			w.WriteHeader(http.StatusNotFound)
			w.Write(res)
			return

		}
		w.WriteHeader(http.StatusOK)
		result, _ := json.Marshal(map[string][]City{"result": response})
		w.Write(result)
	}
}

func getCityListByRegion(h *handler, r *http.Request, value string) []City {
	result := make([]City, 0)
	for _, v := range h.Storage.Storage {
		if v.Region == value {
			result = append(result, v)
		}
	}
	return result
}

func getCityListByPopulation(h *handler, r *http.Request, FoundationFrom string, FoundationTo string) ([]City, error) {
	start, errStart := strconv.Atoi(FoundationFrom)
	end, errEnd := strconv.Atoi(FoundationTo)
	if errStart != nil || errEnd != nil {
		return nil, errors.New("ошибка в параметрах запроса")
	}
	result := make([]City, 0)
	for _, v := range h.Storage.Storage {
		if v.Population >= start && v.Population <= end {
			result = append(result, v)
		}
	}
	return result, nil
}

func getCityListByFoundation(h *handler, r *http.Request, PopulationFrom string, PopulationTo string) ([]City, error) {
	start, errStart := strconv.Atoi(PopulationFrom)
	end, errEnd := strconv.Atoi(PopulationTo)
	if errStart != nil || errEnd != nil {
		return nil, errors.New("ошибка в параметрах запроса")
	}
	result := make([]City, 0)
	for _, v := range h.Storage.Storage {
		if v.Foundation >= start && v.Foundation <= end {
			result = append(result, v)
		}
	}
	return result, nil
}

func getCityListByDistrict(h *handler, r *http.Request, value string) []City {
	result := make([]City, 0)
	for _, v := range h.Storage.Storage {
		if v.District == value {
			result = append(result, v)
		}
	}
	return result
}

func getCity(h *handler, r *http.Request, params httprouter.Params) (*City, error) {
	idParam := params.ByName("id")
	cityId, err := strconv.Atoi(idParam)
	if err != nil {
		return nil, errors.New("ошибка в параметрах запроса")
	}
	city, ok := h.Storage.Storage[cityId]
	if !ok {
		return nil, errors.New("город с таким id не найден")
	}
	return &city, nil
}

func changePopulationCity(s *City, r *http.Request) (*City, error) {
	populationParam := r.Header.Get("population")
	newPopulation, err := strconv.Atoi(populationParam)
	s.Population = newPopulation
	if err != nil {
		return nil, errors.New("ошибка в в параметрах запроса")
	}
	return s, nil
}

func addCity(h *handler, city *City) error {
	_, ok := h.Storage.Storage[city.Id]
	if ok {
		return errors.New("такой город уже есть в списке")
	}
	h.Storage.Storage[city.Id] = *city
	return nil
}
