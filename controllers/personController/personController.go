package personController

import (
	"awesomeProject/models"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"net/http"
)

func GetAllPersons(w http.ResponseWriter, r *http.Request) {
	var db = r.Context().Value("db").(*gorm.DB)
	var persons []models.Person
	if result := db.Table("persons").Find(&persons); result.Error != nil {
		fmt.Println(result.Error)
	}

	personsJson, _ := json.Marshal(persons)

	w.WriteHeader(http.StatusOK)
	w.Write(personsJson)
}

func CreatePerson(w http.ResponseWriter, r *http.Request) {
	var db = r.Context().Value("db").(*gorm.DB)
	var personRequest models.PersonRequest
	err := json.NewDecoder(r.Body).Decode(&personRequest)
	if err != nil {
		fmt.Println(err)
	}
	var person = models.Person{
		Name: personRequest.Name,
		Age:  personRequest.Age,
	}
	if result := db.Table("persons").Create(&person); result.Error != nil {
		fmt.Println(result.Error)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Person was created successfully!"))
}

func DeletePerson(w http.ResponseWriter, r *http.Request) {
	//TODO
}
