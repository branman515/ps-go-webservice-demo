package models

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Book struct {
	ID        int64    `json:"id"` //change name to lower case
	Title     string   `json:"title"`
	Published int      `json:"published,omitempty"`
	Pages     int      `json:"pages,omitempty,string"` // change return data type to string
	Genres    []string `json:"genres,omitempty"`       //string slice
	Rating    float32  `json:"rating,omitempty"`
}

type BookResponse struct {
	Book *Book `json:"book"`
}

type BooksResponse struct {
	Books *[]Book `json:"books"`
}

type ReadingListModel struct {
	Endpoint string
}

func (m *ReadingListModel) GetAll() (*[]Book, error) { //book slice (like a list)
	resp, err := http.Get(m.Endpoint) //sends get call to API
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // close connection at end

	if resp.StatusCode != http.StatusOK { //ensure we get a good response
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	data, err := io.ReadAll(io.Reader(resp.Body)) //read from the response
	if err != nil {
		return nil, err
	}

	var booksResp BooksResponse
	err = json.Unmarshal(data, &booksResp) //decode the message
	if err != nil {
		return nil, err
	}

	return booksResp.Books, nil
}

func (m *ReadingListModel) Get(id int64) (*Book, error) { //singular book returned
	url := fmt.Sprintf("%s/%d", m.Endpoint, id) //generate the endpoint/uri
	resp, err := http.Get(url)                  //sends get call to API
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // close connection at end

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	data, err := io.ReadAll(io.Reader(resp.Body)) //read from the response
	if err != nil {
		return nil, err
	}

	var bookResp BookResponse
	err = json.Unmarshal(data, &bookResp) //decode the message
	if err != nil {
		return nil, err
	}

	return bookResp.Book, nil
}
