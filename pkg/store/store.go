package store

type Store struct {
	Store *[]City
}

func GetStore() *Store{
	var store Store
	fileData := ReadCsvFile("./pkg/store/cities.csv")
	citiesList := createCitiesList(fileData)
	store.Store = &citiesList
	return &store
}

