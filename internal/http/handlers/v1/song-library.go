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

func (sl *songLibrary) getText(w http.ResponseWriter, r *http.Request) {
	const fn = "http.handlers.songLibrary.getText"

	sl.log = sl.log.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	var req dto.GetTextRequest

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

func (sl *songLibrary) get(w http.ResponseWriter, r *http.Request) {
	const fn = "http.handlers.songLibrary.get"

	sl.log = sl.log.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	var req dto.GetSongRequest

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

func (sl *songLibrary) getList(w http.ResponseWriter, r *http.Request) {
	const fn = "http.handlers.songLibrary.getList"

	sl.log = sl.log.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	var req dto.GetSongsListRequest

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

	songs, err := sl.sluc.GetList(&req, pagination.Get(r.Context()))
	if err != nil {
		sl.log.Error("failed to get list of songs", slog.String("error", err.Error()))

		response.RenderError(w, r, http.StatusInternalServerError, "internal error")

		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, songs)
}

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
