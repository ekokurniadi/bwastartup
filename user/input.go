package user

import "encoding/json"

type RegisterUserInput struct {
	Name       string `json:"name" binding:"required"`
	Occupation string `json:"occupation" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type CheckEmailInput struct {
	Email string `json:"email" binding:"required,email"`
}

type SearchStruct struct {
	Regex bool   `json:"regex"`
	Value string `json:"value"`
}
type Column struct {
	Data       json.Number  `json:"data"`
	Name       string       `json:"name"`
	Orderable  bool         `json:"orderable"`
	Search     SearchStruct `json:"search"`
	Searchable bool         `json:"searchable"`
}
type Order struct {
	Column int    `json:"column"`
	Dir    string `json:"dir"`
}
type DTJson struct {
	Columms []Column     `json:"columns"`
	Draw    int          `json:"draw"`
	Length  int          `json:"length"`
	Orders  []Order      `json:"order"`
	Search  SearchStruct `json:"search"`
	Start   int          `json:"start"`
	Total   int          `json:"recordsFiltered"`
}
