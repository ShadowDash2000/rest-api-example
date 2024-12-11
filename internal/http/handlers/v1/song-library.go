package handlers

import (
	"effective-mobile-test/internal/entities/dto"
	"effective-mobile-test/internal/http/middlewares/pagination"
	"effective-mobile-test/internal/http/response"
	"effective-mobile-test/internal/usecases"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
	"log/slog"
	"net/http"
)

type songLibrary struct {
	sluc *usecases.SongLibrary
	log  *slog.Logger
}

func newSongLibrary(sluc *usecases.SongLibrary, log *slog.Logger) *songLibrary {
	return &songLibrary{
		sluc: sluc,
		log:  log,
	}
}

// @Summary Song Library
// @Tags song-library
// @Description Create a song
// @ID create-song
// @Accept json
// @Produce json
// @Param input body dto.CreateSongRequest true "song info"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /v1/songs [post]
func (sl *songLibrary) create(w http.ResponseWriter, r *http.Request) {
	const fn = "http.handlers.songLibrary.create"

	sl.log = sl.log.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	var req dto.CreateSongRequest

	err := render.DecodeJSON(r.Body, &req)
	if err != nil {
		sl.log.Error("failed to decode request body", slog.String("error", err.Error()))

		response.RenderError(w, r, http.StatusBadRequest, "")

		return
	}

	sl.log.Info("request body decoded", slog.Any("request", req))

	if err = validator.New().Struct(req); err != nil {
		response.RenderError(w, r, http.StatusBadRequest, err.Error())

		return
	}

	if err = sl.sluc.Create(req.Group, req.Song); err != nil {
		sl.log.Error("failed to create song", slog.String("error", err.Error()))

		if errors.Is(err, usecases.ErrAlreadyExists) {
			response.RenderError(w, r, http.StatusBadRequest, "this song is already exists")

			return
		}

		response.RenderError(w, r, http.StatusInternalServerError, "internal error")

		return
	}

	response.RenderSuccess(w, r, http.StatusCreated, "")
}

// @Summary Song Library
// @Tags song-library
// @Description Get the lyrics of the song
// @ID get-song-lyrics
// @Accept json
// @Produce json
// @Param offset query int false "paginate through the song lyrics paragraphs"
// @Param group query string true "group name"
// @Param song query string true "song name"
// @Success 200 {object} dto.GetTextResponse
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /v1/songs/text [get]
func (sl *songLibrary) getText(w http.ResponseWriter, r *http.Request) {
	const fn = "http.handlers.songLibrary.getText"

	sl.log = sl.log.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	var req dto.GetTextRequest

	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err := decoder.Decode(&req, r.URL.Query())
	if err != nil {
		sl.log.Error("failed to decode request query", slog.String("error", err.Error()))

		response.RenderError(w, r, http.StatusBadRequest, "")

		return
	}

	sl.log.Info("request query decoded", slog.Any("request", req))

	if err = validator.New().Struct(req); err != nil {
		response.RenderError(w, r, http.StatusBadRequest, err.Error())

		return
	}

	textRes, err := sl.sluc.GetText(req.Group, req.Song, pagination.Get(r.Context()))
	if err != nil {
		sl.log.Error("failed to get text of the song", slog.String("error", err.Error()))

		if errors.Is(err, usecases.ErrNullFields) {
			response.RenderError(w, r, http.StatusBadRequest, "this song doesn't have text yet")

			return
		} else if errors.Is(err, usecases.ErrNoRowsAffected) {
			response.RenderError(w, r, http.StatusBadRequest, "song not found")

			return
		}

		response.RenderError(w, r, http.StatusInternalServerError, "internal error")

		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, textRes)
}

// @Summary Song Library
// @Tags song-library
// @Description Get the song info
// @ID get-song-info
// @Accept json
// @Produce json
// @Param group query string true "group name"
// @Param song query string true "song name"
// @Success 200 {object} dto.GetSongResponse
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /info [get]
func (sl *songLibrary) get(w http.ResponseWriter, r *http.Request) {
	const fn = "http.handlers.songLibrary.get"

	sl.log = sl.log.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	var req dto.GetSongRequest

	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err := decoder.Decode(&req, r.URL.Query())
	if err != nil {
		sl.log.Error("failed to decode request query", slog.String("error", err.Error()))

		response.RenderError(w, r, http.StatusBadRequest, "")

		return
	}

	sl.log.Info("request query decoded", slog.Any("request", req))

	if err = validator.New().Struct(req); err != nil {
		response.RenderError(w, r, http.StatusBadRequest, err.Error())

		return
	}

	textRes, err := sl.sluc.Get(req.Group, req.Song)
	if err != nil {
		sl.log.Error("failed to get song", slog.String("error", err.Error()))

		if errors.Is(err, usecases.ErrNoRowsAffected) {
			response.RenderError(w, r, http.StatusBadRequest, "song not found")

			return
		}

		response.RenderError(w, r, http.StatusInternalServerError, "internal error")

		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, textRes)
}

// @Summary Song Library
// @Tags song-library
// @Description Get a list of songs
// @ID get-songs-list
// @Accept json
// @Produce json
// @Param offset query int false "paginate through the songs list"
// @Param limit query int false "sets the list limit"
// @Param group query string false "group name"
// @Param song query string false " song name"
// @Param releaseDate query string false "release date" format(date)
// @Param link query string false "link"
// @Param text query string false "lyrics"
// @Success 200 {array} dto.GetSongsListResponse
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /v1/songs [get]
func (sl *songLibrary) getList(w http.ResponseWriter, r *http.Request) {
	const fn = "http.handlers.songLibrary.getList"

	sl.log = sl.log.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	var req dto.GetSongsListRequest

	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err := decoder.Decode(&req, r.URL.Query())
	if err != nil {
		sl.log.Error("failed to decode request query", slog.String("error", err.Error()))

		response.RenderError(w, r, http.StatusBadRequest, "")

		return
	}

	sl.log.Info("request query decoded", slog.Any("request", req))

	if err = validator.New().Struct(req); err != nil {
		response.RenderError(w, r, http.StatusBadRequest, err.Error())

		return
	}

	songs, err := sl.sluc.GetList(&req, pagination.Get(r.Context()))
	if err != nil {
		sl.log.Error("failed to get list of songs", slog.String("error", err.Error()))

		if errors.Is(err, usecases.ErrNoRowsAffected) {
			response.RenderError(w, r, http.StatusBadRequest, "songs not found")

			return
		}

		response.RenderError(w, r, http.StatusInternalServerError, "internal error")

		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, songs)
}

// @Summary Song Library
// @Tags song-library
// @Description Update a specific song
// @ID update-song
// @Accept json
// @Produce json
// @Param input body dto.UpdateSongRequest true "song info and the fields to update"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /v1/songs [put]
func (sl *songLibrary) update(w http.ResponseWriter, r *http.Request) {
	const fn = "http.handlers.songLibrary.update"

	sl.log = sl.log.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	var req dto.UpdateSongRequest

	err := render.DecodeJSON(r.Body, &req)
	if err != nil {
		sl.log.Error("failed to decode request body", slog.String("error", err.Error()))

		response.RenderError(w, r, http.StatusBadRequest, "")

		return
	}

	sl.log.Info("request body decoded", slog.Any("request", req))

	if err = validator.New().Struct(req); err != nil {
		response.RenderError(w, r, http.StatusBadRequest, err.Error())

		return
	}

	err = sl.sluc.Update(&req)
	if err != nil {
		sl.log.Error("failed to update song", slog.String("error", err.Error()))

		if errors.Is(err, usecases.ErrNoRowsAffected) {
			response.RenderError(w, r, http.StatusBadRequest, "song for update is not found")

			return
		}

		response.RenderError(w, r, http.StatusInternalServerError, "internal error")

		return
	}

	response.RenderSuccess(w, r, http.StatusOK, "")
}

// @Summary Song Library
// @Tags song-library
// @Description Delete a specific song
// @ID delete-song
// @Accept json
// @Produce json
// @Param input body dto.DeleteSongRequest true "song info"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /v1/songs [delete]
func (sl *songLibrary) delete(w http.ResponseWriter, r *http.Request) {
	const fn = "http.handlers.songLibrary.delete"

	sl.log = sl.log.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	var req dto.DeleteSongRequest

	err := render.DecodeJSON(r.Body, &req)
	if err != nil {
		sl.log.Error("failed to decode request body", slog.String("error", err.Error()))

		response.RenderError(w, r, http.StatusBadRequest, "")

		return
	}

	sl.log.Info("request body decoded", slog.Any("request", req))

	if err = validator.New().Struct(req); err != nil {
		response.RenderError(w, r, http.StatusBadRequest, err.Error())

		return
	}

	err = sl.sluc.Delete(req.Group, req.Song)
	if err != nil {
		sl.log.Error("failed to delete song", slog.String("error", err.Error()))

		if errors.Is(err, usecases.ErrNoRowsAffected) {
			response.RenderError(w, r, http.StatusBadRequest, "song for deletion is not found")

			return
		}

		response.RenderError(w, r, http.StatusInternalServerError, "internal error")

		return
	}

	response.RenderSuccess(w, r, http.StatusOK, "")
}
