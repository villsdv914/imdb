package utils

import "time"

type ResponseBody struct  {
	Id uint `json:"id,omitempty"`
	Title string `json:"title,omitempty"`
	ReleasedYear string `json:"released_year,omitempty"`
	Rating int `json:"rating,omitempty"`
	Genres string `json:"genres,omitempty"`
	CreatedAt time.Time `json:"CreatedAt,omitempty"`
	UpdatedAt time.Time `json:"UpdatedAt,omitempty"`
	DeletedAt time.Time `json:"DeletedAt,omitempty"`
}

type RequestBody struct  {
	Title string `json:"title"`
	ReleasedYear string `json:"released_year"`
	Rating int `json:"rating"`
	Genres string `json:"genres"`
}

type SearchBody struct  {
	SearchType string `json:"type"`
	Condition string `json:"condition"`
	TypeValue string `json:"type_value"`

}
