package store

import "os"

type Store struct {
	Storage *[]City
}

func GetStore() *Store{
	var store Store
	fileData := ReadCsvFile("./pkg/store/cities.csv")
	// citiesList := createCitiesList(fileData)
	store.Storage = createCitiesList(fileData)
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

func createCitiesList(data [][]string) *[]City {
	// var citiesList []City
	citiesList := make([]City, 0)
	for _, line := range data {
		var newCity City
		newCity.Id = toInt(line[0])
		newCity.Name = line[1]
		newCity.Region = line[2]
		newCity.District = line[3]
		newCity.Population = toInt(line[4])
		newCity.Foundation = toInt(line[5])
		citiesList = append(citiesList, newCity)
	}
	return &citiesList
}

func WriteToCsv(data *[]City) {
	f, err := os.Create("./pkg/store/cities2.csv")
	if err != nil {
		log.Println(err)
	}
	w := csv.NewWriter(f)
	defer f.Close()
	for _, line := range *data {
		var record []string
		record = append(record, strconv.Itoa(line.Id))
		record = append(record, line.Name)
		record = append(record, line.Region)
		record = append(record, line.District)
		record = append(record, strconv.Itoa(line.Population))
		record = append(record, strconv.Itoa(line.Foundation))
		w.Write(record)
	}
	w.Flush()
}