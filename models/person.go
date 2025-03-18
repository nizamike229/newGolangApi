package models

type Person struct {
	Name string
	Age  int
	Id   int `gorm:"unique;primaryKey;autoIncrement"`
}
