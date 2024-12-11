package usecases

import (
	"database/sql"
	"effective-mobile-test/internal/entities"
	"effective-mobile-test/internal/entities/dto"
	"effective-mobile-test/internal/http/middlewares/pagination"
	"errors"
	"fmt"
	"github.com/fatih/structs"
	"github.com/lib/pq"
	"log/slog"
	"reflect"
	"strings"
)

type SongLibraryRepo interface {
	Create(group, song string) error
	Get(group, song string) (*entities.Song, error)
	GetList(filter map[string]interface{}, pagination *pagination.Pagination) ([]*entities.Song, error)
	Update(group, song string, fields map[string]interface{}) error
	Delete(group, song string) error
}

type SongLibrary struct {
	repo SongLibraryRepo
	log  *slog.Logger
}

func NewSongLibrary(repo SongLibraryRepo, log *slog.Logger) *SongLibrary {
	return &SongLibrary{
		repo: repo,
		log:  log,
	}
}

func (sl *SongLibrary) Create(group, song string) error {
	const fn = "usecases.SongLibrary.Create"

	defer sl.log.With(
		slog.String("fn", fn),
	).Debug("", group, song)

	err := sl.repo.Create(group, song)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code.Name() == "unique_violation" {
				return fmt.Errorf("%s: %s", fn, ErrAlreadyExists)
			}
		}
		return fmt.Errorf("%s: %s", fn, err)
	}

	return nil
}

func (sl *SongLibrary) GetText(group, song string, pagination *pagination.Pagination) (*dto.GetTextResponse, error) {
	const fn = "usecases.SongLibrary.GetText"
	var paragraph string

	defer sl.log.With(
		slog.String("fn", fn),
	).Debug("",
		group,
		song,
		slog.Any("pagination", pagination),
	)

	songRes, err := sl.repo.Get(group, song)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", fn, ErrNoRowsAffected)
		}

		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	if songRes.Text == nil {
		return nil, fmt.Errorf("%s: %w", fn, ErrNullFields)
	}

	if pagination.Offset < 0 {
		pagination.Offset = 0
	}

	for key, text := range strings.Split(*songRes.Text, "\n\n") {
		if key == pagination.Offset {
			paragraph = text
			break
		}
	}

	return &dto.GetTextResponse{
		Text: paragraph,
	}, nil
}

func (sl *SongLibrary) Get(group, song string) (*dto.GetSongResponse, error) {
	const fn = "usecases.SongLibrary.Get"

	defer sl.log.With(
		slog.String("fn", fn),
	).Debug("", group, song)

	songRes, err := sl.repo.Get(group, song)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", fn, ErrNoRowsAffected)
		}

		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return dto.NewGetSongResponse(songRes), nil
}

func (sl *SongLibrary) GetList(filter *dto.GetSongsListRequest, pagination *pagination.Pagination) ([]*dto.GetSongsListResponse, error) {
	const fn = "usecases.SongLibrary.GetList"
	var filterMap = make(map[string]interface{})

	defer func(filterMap *map[string]interface{}) {
		sl.log.With(
			slog.String("fn", fn),
		).Debug("",
			slog.Any("filter", filter),
			slog.Any("filterMap", *filterMap),
			slog.Any("pagination", pagination),
		)
	}(&filterMap)

	s := structs.New(&filter)

	for _, field := range s.Fields() {
		if field.IsZero() {
			continue
		}

		var valRes any
		tag := field.Tag("db")
		val := reflect.ValueOf(field.Value())
		if val.Kind() == reflect.Ptr {
			if !val.IsNil() {
				valRes = val.Elem().Interface()
			} else {
				continue
			}
		} else {
			valRes = field.Value()
		}

		if tag != "" {
			filterMap[tag] = valRes
		}
	}

	songs, err := sl.repo.GetList(filterMap, pagination)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	if len(songs) == 0 {
		return nil, fmt.Errorf("%s: %w", fn, ErrNoRowsAffected)
	}

	return dto.NewGetSongsListResponse(songs), nil
}

func (sl *SongLibrary) Update(song *dto.UpdateSongRequest) error {
	const fn = "usecases.SongLibrary.Update"

	defer sl.log.With(
		slog.String("fn", fn),
	).Debug("", slog.Any("song", song))

	s := structs.New(&song)

	var fields = make(map[string]interface{})
	for _, field := range s.Fields() {
		if tag := field.Tag("db"); tag != "" {
			fields[tag] = field.Value()
		}
	}

	err := sl.repo.Update(song.Group, song.Song, fields)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%s: %w", fn, ErrNoRowsAffected)
		}
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (sl *SongLibrary) Delete(group, song string) error {
	const fn = "usecases.SongLibrary.Delete"

	defer sl.log.With(
		slog.String("fn", fn),
	).Debug("", group, song)

	err := sl.repo.Delete(group, song)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%s: %w", fn, ErrNoRowsAffected)
		}
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}
