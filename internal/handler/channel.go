package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/VladPetriv/scanner_backend_api/internal/store/pg"
)

// GetChannelsCountHandler godoc
// @ID           get-channels-count
// @Summary      GetChannelsCount
// @Description  Handler will return channels count
// @Tags         channel
// @Produce      json
// @Success      200  {object}  object         "channels count"
// @Failure      404  {object}  lib.HttpError  "channels count not found"
// @Failure      500  {object}  lib.HttpError  "internal server error"
// @Router       /channel/count [get]
func (h *Handler) GetChannelsCountHandler(w http.ResponseWriter, r *http.Request) {
	count, err := h.service.Channel.GetChannelsCount()
	if err != nil {
		h.log.Error("get channels count error", zap.Error(err))

		if errors.Is(err, pg.ErrChannelsCountNotFound) {
			err = errors.Unwrap(err)

			h.WriteError(w, http.StatusNotFound, err.Error())

			return
		}

		h.WriteError(w, http.StatusInternalServerError, err.Error())

		return
	}

	h.WriteJSON(w, http.StatusOK, map[string]int{"count": count})
}

// GetChannelByNameHandler godoc
// @ID           get-channel-by-name
// @Summary      GetChannelByName
// @Description  Handler will return channel by name from url
// @Tags         channel
// @Produce      json
// @Param        name  path      string         true  "channel name"
// @Success      200   {object}  model.Channel  "channel by name"
// @Failure      400   {object}  lib.HttpError  "bad request"
// @Failure      404   {object}  lib.HttpError  "channel not found"
// @Failure      500   {object}  lib.HttpError  "internal server error"
// @Router       /channel/{name} [get]
func (h *Handler) GetChannelByNameHandler(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]

	channel, err := h.service.Channel.GetChannelByName(name)
	if err != nil {
		h.log.Error("get channel by name error", zap.String("name", name), zap.Error(err))

		if errors.Is(err, pg.ErrChannelNotFound) {
			err = errors.Unwrap(err)

			h.WriteError(w, http.StatusNotFound, err.Error())

			return
		}

		h.WriteError(w, http.StatusInternalServerError, err.Error())

		return
	}

	h.WriteJSON(w, http.StatusOK, channel)
}

// GetChannelsByPageHandler godoc
// @ID           get-channels-by-page
// @Summary      GetChannelsByPage
// @Description  Handler will return channels by page from query
// @Tags         channel
// @Produce      json
// @Param        page  query     integer        true  "page"
// @Success      200   {array}   model.Channel  "channels by page"
// @Failure      400   {object}  lib.HttpError  "bad request"
// @Failure      404   {object}  lib.HttpError  "channels not found"
// @Failure      500   {object}  lib.HttpError  "internal server error"
// @Router       /channel/ [get]
func (h *Handler) GetChannelsByPageHandler(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		h.log.Error("failed to get page", zap.Error(err))

		h.WriteError(w, http.StatusBadRequest, "page is not valid")

		return
	}

	channels, err := h.service.Channel.GetChannelsByPage(page)
	if err != nil {
		h.log.Error("get channels by page error", zap.String("page", strconv.Itoa(page)), zap.Error(err))

		if errors.Is(err, pg.ErrChannelsNotFound) {
			err = errors.Unwrap(err)

			h.WriteError(w, http.StatusNotFound, err.Error())

			return
		}

		h.WriteError(w, http.StatusInternalServerError, err.Error())

		return
	}

	h.WriteJSON(w, http.StatusOK, channels)
}
