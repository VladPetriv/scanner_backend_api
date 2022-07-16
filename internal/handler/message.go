package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/VladPetriv/scanner_backend_api/internal/store/pg"
)

// GetMessagesCountHandler godoc
// @ID           get-messages-count
// @Summary      GetMessagesCount
// @Description  Handler will return messages count
// @Tags         message
// @Produce      json
// @Success      200         {object}  object         "messages count"
// @Failure      404         {object}  lib.HttpError  "messages count not found"
// @Failure      500         {object}  lib.HttpError  "internal server error"
// @Router       /message/count [get]
func (h *Handler) GetMessagesCountHandler(w http.ResponseWriter, r *http.Request) {
	count, err := h.service.Message.GetMessagesCount()
	if err != nil {
		h.log.Error("get messages count error", zap.Error(err))

		if errors.Is(err, pg.ErrMessagesCountNotFound) {
			err = errors.Unwrap(err)

			h.WriteError(w, http.StatusNotFound, err.Error())

			return
		}

		h.WriteError(w, http.StatusInternalServerError, err.Error())

		return
	}

	h.WriteJSON(w, http.StatusOK, map[string]int{"count": count})
}

// GetMessagesCountByChannelIDHandler godoc
// @ID           get-messages-by-channel-id-count
// @Summary      GetMessagesByChannelIDCount
// @Description  Handler will return messages count by channel id from url
// @Tags         message
// @Produce      json
// @Param        channel_id  path      integer        true  "channel id"
// @Success      200  {object}  object         "messages count"
// @Failure      400         {object}  lib.HttpError  "bad request"
// @Failure      404  {object}  lib.HttpError  "messages count not found"
// @Failure      500  {object}  lib.HttpError  "internal server error"
// @Router       /message/count/{channel_id} [get]
func (h *Handler) GetMessagesCountByChannelIDHandler(w http.ResponseWriter, r *http.Request) {
	channelID, err := strconv.Atoi(mux.Vars(r)["channel_id"])
	if err != nil {
		h.log.Error("get channel id from request error", zap.Error(err))

		h.WriteError(w, http.StatusBadRequest, "channel id is not valid")

		return
	}

	count, err := h.service.Message.GetMessagesCountByChannelID(channelID)
	if err != nil {
		h.log.Error("get messages count by channel id error", zap.String("id", strconv.Itoa(channelID)), zap.Error(err))

		if errors.Is(err, pg.ErrMessagesCountNotFound) {
			err = errors.Unwrap(err)

			h.WriteError(w, http.StatusNotFound, err.Error())

			return
		}

		h.WriteError(w, http.StatusInternalServerError, err.Error())

		return
	}

	h.WriteJSON(w, http.StatusOK, map[string]int{"count": count})
}

// GetFullMessagesByPageHandler godoc
// @ID           get-full-messages-by-page
// @Summary      GetFullMessagesByPage
// @Description  Handler will return full messages by page from query
// @Tags         message
// @Produce      json
// @Param        page  query     integer            true  "page"
// @Success      200   {array}   model.FullMessage  "full messages by page"
// @Failure      400         {object}  lib.HttpError      "bad request"
// @Failure      404         {object}  lib.HttpError      "full messages not found"
// @Failure      500         {object}  lib.HttpError      "internal server error"
// @Router       /message/ [get]
func (h *Handler) GetFullMessagesByPageHandler(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		h.log.Error("get page from query error", zap.Error(err))

		h.WriteError(w, http.StatusBadRequest, "page is not valid")

		return
	}

	messages, err := h.service.Message.GetFullMessagesByPage(page)
	if err != nil {
		h.log.Error("get messages by page error", zap.String("page", strconv.Itoa(page)), zap.Error(err))

		if errors.Is(err, pg.ErrFullMessagesNotFound) {
			err = errors.Unwrap(err)

			h.WriteError(w, http.StatusNotFound, err.Error())

			return
		}

		h.WriteError(w, http.StatusInternalServerError, err.Error())

		return
	}

	h.WriteJSON(w, http.StatusOK, messages)
}

// GetFullMessagesByChannelIDAndPageHandler godoc
// @ID           get-full-messages-by-page-and-channel-id
// @Summary      GetFullMessagesByChannelIDAndPage
// @Description  Handler will return full messages by page from query and channel id from url
// @Tags         message
// @Produce      json
// @Param        channel_id  path      integer            true  "channel id"
// @Param        page        query     integer            true  "page"
// @Success      200         {array}   model.FullMessage  "full messages by page and channel id"
// @Failure      400      {object}  lib.HttpError      "bad request"
// @Failure      404      {object}  lib.HttpError      "full messages not found"
// @Failure      500      {object}  lib.HttpError      "internal server error"
// @Router       /message/channel/{channel_id} [get]
func (h *Handler) GetFullMessagesByChannelIDAndPageHandler(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		h.log.Error("get page from query error", zap.Error(err))

		h.WriteError(w, http.StatusBadRequest, "page is not valid")

		return
	}

	channelID, err := strconv.Atoi(mux.Vars(r)["channel_id"])
	if err != nil {
		h.log.Error("get id from query error", zap.Error(err))

		h.WriteError(w, http.StatusBadRequest, "channel id is not valid")

		return
	}

	messages, err := h.service.Message.GetFullMessagesByChannelIDAndPage(channelID, page)
	if err != nil {
		h.log.Error("get full messages by page and channel id error", zap.String("id,page", fmt.Sprintf("%d,%d", channelID, page)))

		if errors.Is(err, pg.ErrFullMessagesNotFound) {
			err = errors.Unwrap(err)

			h.WriteError(w, http.StatusNotFound, err.Error())

			return
		}

		h.WriteError(w, http.StatusInternalServerError, err.Error())

		return
	}

	h.WriteJSON(w, http.StatusOK, messages)
}

// GetFullMessagesByUserIDHandler godoc
// @ID           get-full-messages-by-user-id
// @Summary      GetFullMessagesByUserID
// @Description  Handler will return full messages by user id from url
// @Tags         message
// @Produce      json
// @Param        user_id  path      integer            true  "user id"
// @Success      200      {array}   model.FullMessage  "full messages by user id"
// @Failure      400         {object}  lib.HttpError      "bad request"
// @Failure      404         {object}  lib.HttpError      "full messages not found"
// @Failure      500         {object}  lib.HttpError      "internal server error"
// @Router       /message/user/{user_id} [get]
func (h *Handler) GetFullMessagesByUserIDHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(mux.Vars(r)["user_id"])
	if err != nil {
		h.log.Error("get user id from url error", zap.Error(err))

		h.WriteError(w, http.StatusBadRequest, "user id is not valid")

		return
	}

	messages, err := h.service.Message.GetFullMessagesByUserID(userID)
	if err != nil {
		h.log.Error("get full messages by user id error", zap.String("id", strconv.Itoa(userID)), zap.Error(err))

		if errors.Is(err, pg.ErrFullMessagesNotFound) {
			err = errors.Unwrap(err)

			h.WriteError(w, http.StatusNotFound, err.Error())

			return
		}
		h.WriteError(w, http.StatusInternalServerError, err.Error())

		return
	}

	h.WriteJSON(w, http.StatusOK, messages)
}

// GetFullMessageByIDHandler godoc
// @ID           get-full-message-by-id
// @Summary      GetFullMessageByID
// @Description  Handler will return full message by id from url
// @Tags         message
// @Produce      json
// @Param        message_id  path      integer            true  "message id"
// @Success      200         {object}  model.FullMessage  "full message by user id"
// @Failure      400   {object}  lib.HttpError      "bad request"
// @Failure      404   {object}  lib.HttpError      "full messages not found"
// @Failure      500   {object}  lib.HttpError      "internal server error"
// @Router       /message/{message_id} [get]
func (h *Handler) GetFullMessageByIDHandler(w http.ResponseWriter, r *http.Request) {
	messageID, err := strconv.Atoi(mux.Vars(r)["message_id"])
	if err != nil {
		h.log.Error("get message id from url error", zap.Error(err))

		h.WriteError(w, http.StatusBadRequest, "message id is not valid")

		return
	}

	message, err := h.service.Message.GetFullMessageByID(messageID)
	if err != nil {
		h.log.Error("get message by id error", zap.String("id", strconv.Itoa(messageID)), zap.Error(err))

		if errors.Is(err, pg.ErrFullMessageNotFound) {
			err = errors.Unwrap(err)

			h.WriteError(w, http.StatusNotFound, err.Error())

			return
		}

		h.WriteError(w, http.StatusInternalServerError, err.Error())

		return
	}

	h.WriteJSON(w, http.StatusOK, message)
}
