package models

type NhicFormat struct {
	CardNumber  string `json:"cardNumber"`
	Name        string `json:"name"`
	IdNumber    string `json:"idNumber"`
	Birthday    string `json:"birthday"`
	Sex         string `json:"sex"`
	CardDate    string `json:"cardDate"`
	IsCardExist bool   `json:"isCardExist"`
}
