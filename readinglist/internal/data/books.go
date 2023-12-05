package data

import (
	"time"
)

type Book struct {
	ID        int64     `json:"id"` //change name to lower case
	CreatedAt time.Time `json:"-"`  //hide the field in json marshalling
	Title     string    `json:"title"`
	Published int       `json:"published,omitempty"`
	Pages     int       `json:"pages,omitempty,string"` // change return data type to string
	Genres    []string  `json:"genres,omitempty"`       //string slice
	Rating    float32   `json:"rating,omitempty"`
	Version   int32     `json:"-"`
}
