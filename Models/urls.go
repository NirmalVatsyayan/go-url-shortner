package models

import "time"

type UrlShortner struct {
	UserID    string    `json:"userId"`
	UrlEncodingKey int64 `json:"urlEncodingKey"`
	Url     string `json:"url"`
	UrlEncoded string `json:"urlEncoded"`
	CreatedOn time.Time `json:"createdOn"`
	ViewCount int `json:"viewCount"`
}

type Pagination struct {
	PrevUrl string `json:"prev_url"`
	NextUrl string `json:"next_url"`
	Count int `json:"count"`
}

type UrlWrapper struct {
	Pagination Pagination `json:"Page"`
	Urls []UrlShortner `json: "EncodedUrls"`
}

