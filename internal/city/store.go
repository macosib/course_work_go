package city

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"sync"
)

type Store struct {
	mu sync.Mutex
	Storage map[int]City
}

func NewStore() *Store {
	var store Store
	store.Storage = make(map[int]City)
	fileData := store.ReadCsvFile("./internal/city/cities.csv")
	store.createCitiesList(fileData)
	return &store
}

func (s *Store) ReadCsvFile(fileName string) [][]string {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	data, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func toInt(str string) int {
	number, err := strconv.Atoi(str)
	if err != nil {
		log.Println(err)
	}
	return number
}

func (s *Store) createCitiesList(data [][]string) {
	for _, line := range data {
		var newCity City
		newCity.Id = toInt(line[0])
		newCity.Name = line[1]
		newCity.Region = line[2]
		newCity.District = line[3]
		newCity.Population = toInt(line[4])
		newCity.Foundation = toInt(line[5])
		s.Storage[newCity.Id] = newCity
	}
}

func WriteToCsv(s *Store) {
	f, err := os.Create("./internal/city/cities.csv")
	if err != nil {
		log.Println(err)
	}
	w := csv.NewWriter(f)
	defer f.Close()
	for _, value := range s.Storage {
		var record []string
		record = append(record, strconv.Itoa(value.Id))
		record = append(record, value.Name)
		record = append(record, value.Region)
		record = append(record, value.District)
		record = append(record, strconv.Itoa(value.Population))
		record = append(record, strconv.Itoa(value.Foundation))
		w.Write(record)
	}
	w.Flush()
}
