package dto

import (
	"effective-mobile-test/internal/entities"
)

type CreateSongRequest struct {
	Group string `json:"group" validate:"required"`
	Song  string `json:"song" validate:"required"`
}

type UpdateSongRequest struct {
	Group       string  `json:"group" validate:"required"`
	Song        string  `json:"song" validate:"required"`
	ReleaseDate *string `json:"releaseDate" db:"release_date"`
	Link        *string `json:"link" db:"link"`
	Text        *string `json:"text" db:"text"`
}

type DeleteSongRequest struct {
	Group string `json:"group" validate:"required"`
	Song  string `json:"song" validate:"required"`
}

type GetTextRequest struct {
	Group string `json:"group" validate:"required"`
	Song  string `json:"song" validate:"required"`
}

type GetTextResponse struct {
	Text string `json:"text"`
}

type GetSongRequest struct {
	Group string `schema:"group" validate:"required"`
	Song  string `schema:"song" validate:"required"`
}

type GetSongResponse struct {
	ReleaseDate *string `json:"releaseDate"`
	Link        *string `json:"link"`
	Text        *string `json:"text"`
}

func NewGetSongResponse(res *entities.Song) *GetSongResponse {
	return &GetSongResponse{
		ReleaseDate: res.ReleaseDate,
		Link:        res.Link,
		Text:        res.Text,
	}
}

type GetSongsListRequest struct {
	Group       string  `schema:"group" db:"group"`
	Song        string  `schema:"song" db:"song"`
	ReleaseDate *string `schema:"releaseDate" db:"release_date"`
	Link        *string `schema:"link" db:"link"`
	Text        *string `schema:"text" db:"text"`
}

type GetSongsListResponse struct {
	Group       string  `json:"group" db:"group"`
	Song        string  `json:"song" db:"song"`
	ReleaseDate *string `json:"releaseDate,omitempty" db:"release_date"`
	Link        *string `json:"link,omitempty" db:"link"`
	Text        *string `json:"text,omitempty" db:"text"`
}

func NewSongResponse(res *entities.Song) *GetSongsListResponse {
	return &GetSongsListResponse{
		Group:       res.Group,
		Song:        res.Song,
		ReleaseDate: res.ReleaseDate,
		Link:        res.Link,
		Text:        res.Text,
	}
}

func NewGetSongsListResponse(res []*entities.Song) []*GetSongsListResponse {
	var songs []*GetSongsListResponse
	for _, song := range res {
		songs = append(songs, NewSongResponse(song))
	}
	return songs
}
