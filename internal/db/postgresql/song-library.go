package postgres

import (
	"database/sql"
	"effective-mobile-test/internal/entities"
	"effective-mobile-test/internal/http/middlewares/pagination"
	"fmt"
	"github.com/Masterminds/squirrel"
	"log/slog"
)

type SongLibrary struct {
	*DB
	stmtBuilder squirrel.StatementBuilderType
}

func NewSongLibrary(db *DB) *SongLibrary {
	stmtBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	return &SongLibrary{
		DB:          db,
		stmtBuilder: stmtBuilder,
	}
}

func buildPagination(queryBuilder squirrel.SelectBuilder, pagination *pagination.Pagination, defaultLimit int) squirrel.SelectBuilder {
	if pagination.Limit <= 0 && defaultLimit > 0 {
		pagination.Limit = defaultLimit
	}

	if pagination.Limit > 0 {
		queryBuilder = queryBuilder.Limit(uint64(pagination.Limit))

		if pagination.Offset > 0 {
			queryBuilder = queryBuilder.Offset(uint64(pagination.Offset))
		}
	}

	return queryBuilder
}

func (sl *SongLibrary) Create(group, song string) error {
	const fn = "sl.postgres.SongLibrary.Create"
	var query string

	defer sl.log.With(
		slog.String("fn", fn),
	).Debug("", slog.String("query", query))

	queryBuilder := sl.stmtBuilder.
		Insert("song_library").
		Columns(`"group"`, "song").
		Values(group, song)

	query, _, _ = queryBuilder.ToSql()

	_, err := queryBuilder.RunWith(sl.db).Exec()

	if err != nil {
		return fmt.Errorf("%w: %s", err, fn)
	}

	return nil
}

func (sl *SongLibrary) Get(group, song string) (*entities.Song, error) {
	const fn = "sl.postgres.SongLibrary.Get"
	var query string

	defer func(query *string) {
		sl.log.With(
			slog.String("fn", fn),
		).Debug("", slog.String("query", *query))
	}(&query)

	queryBuilder := sl.stmtBuilder.
		Select("text").
		From("song_library").
		Where(squirrel.Eq{`"group"`: group, `"song"`: song})

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	var songRes entities.Song
	err = sl.db.Get(&songRes, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &songRes, nil
}

func (sl *SongLibrary) GetList(filter map[string]interface{}, pagination *pagination.Pagination) ([]*entities.Song, error) {
	const fn = "sl.postgres.SongLibrary.GetList"
	var query string

	defer func(query *string) {
		sl.log.With(
			slog.String("fn", fn),
		).Debug("", slog.String("query", *query))
	}(&query)

	queryBuilder := sl.stmtBuilder.
		Select(`"group"`, "song", "release_date", "link", "text").
		From("song_library").
		OrderBy("id")

	queryBuilder = buildPagination(queryBuilder, pagination, 10)

	for key, value := range filter {
		queryBuilder = queryBuilder.Where(`"`+key+`" LIKE ?`, fmt.Sprint("%", value, "%"))
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	var songs []*entities.Song
	err = sl.db.Select(&songs, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return songs, nil
}

func (sl *SongLibrary) Update(group, song string, fields map[string]interface{}) error {
	const fn = "sl.postgres.SongLibrary.Update"
	var query string

	defer sl.log.With(
		slog.String("fn", fn),
	).Debug("", slog.String("query", query))

	queryBuilder := sl.stmtBuilder.
		Update("song_library").
		SetMap(fields).
		Where(squirrel.Eq{`"group"`: group, `"song"`: song})

	query, _, _ = queryBuilder.ToSql()

	res, err := queryBuilder.RunWith(sl.db).Exec()
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	if rows == 0 {
		return fmt.Errorf("%s: %w", fn, sql.ErrNoRows)
	}

	return nil
}

func (sl *SongLibrary) Delete(group, song string) error {
	const fn = "sl.postgres.SongLibrary.Delete"
	var query string

	defer func(query *string) {
		sl.log.With(
			slog.String("fn", fn),
		).Debug("", slog.String("query", *query))
	}(&query)

	queryBuilder := sl.stmtBuilder.
		Delete("song_library").
		Where(squirrel.Eq{`"group"`: group, `"song"`: song})

	query, _, _ = queryBuilder.ToSql()

	res, err := queryBuilder.RunWith(sl.db).Exec()
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	if rows == 0 {
		return fmt.Errorf("%s: %w", fn, sql.ErrNoRows)
	}

	return nil
}
