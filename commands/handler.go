package commands

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Coderx44/gg/api"
	"github.com/Coderx44/gg/apperrors"
	"github.com/Coderx44/gg/domain"
)

func HandleCommands(dt *Datastore) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		data := r.Context().Value("data").([]string)

		for _, i := range data {
			fmt.Println(i)
		}

		switch data[0] {
		case "SET":
			SetCommand(dt, data, w, r)
			return
		case "GET":
			GetCommand(dt, data, w, r)
			return
		case "QPUSH":
			QpushCommand(dt, data, w, r)
			return
		case "QPOP":
			QpopCommand(dt, data, w, r)
			return
		case "BQPOP":
			BQpopCommand(dt, data, w, r)
			return
		default:
			api.RespondWithJSON(w, http.StatusBadRequest, api.Response{
				Error: apperrors.ErrInvalidCommand.Error(),
			})
			return
		}

	})
}

func SetCommand(dt *Datastore, data []string, w http.ResponseWriter, r *http.Request) {
	if len(data) < 3 || data[0] != "SET" {
		api.RespondWithJSON(w, http.StatusBadRequest, api.Response{
			Error: apperrors.ErrInvalidCommand.Error(),
		})
		return
	}
	for i := len(data); i < 7; i++ {
		data = append(data, " ")
	}
	setCommand := domain.SetCommand{}
	setCommand.Key = data[1]
	setCommand.Value = data[2]
	setCommand.Expiry_time = -1
	if data[3] == "EX" {
		time, err := strconv.ParseInt(string(data[4]), 10, 64)
		if err == nil {
			setCommand.Expiry_time = int(time)
		}

		setCommand.Condition = data[5]
	} else {
		setCommand.Condition = data[3]
	}

	err := dt.Set(r.Context(), setCommand)
	if err != nil {
		if err == apperrors.ErrKeyExists {
			api.RespondWithJSON(w, http.StatusOK, api.Response{
				Message: apperrors.ErrKeyExists.Error(),
			})
			return
		}

		if err == apperrors.ErrKeyNotFound {
			api.RespondWithJSON(w, http.StatusOK, api.Response{
				Message: apperrors.ErrKeyNotFound.Error(),
			})
			return
		}

		api.RespondWithJSON(w, http.StatusInternalServerError, api.Response{
			Error: "Something went wrong!",
		})
		return
	}
	api.RespondWithJSON(w, http.StatusOK, api.Response{})
	return
}

func GetCommand(dt *Datastore, data []string, w http.ResponseWriter, r *http.Request) {

	if len(data) < 2 || data[0] != "GET" {
		api.RespondWithJSON(w, http.StatusBadRequest, api.Response{
			Error: apperrors.ErrInvalidCommand.Error(),
		})
		return
	}

	value, err := dt.Get(r.Context(), data[1])
	if err != nil {
		if err == apperrors.ErrKeyNotFound {
			api.RespondWithJSON(w, http.StatusBadRequest, api.Response{
				Error: apperrors.ErrKeyNotFound.Error(),
			})
			return
		}
	}

	api.RespondWithJSON(w, http.StatusOK, api.Response{
		Value: value,
	})
	return

}

func QpushCommand(dt *Datastore, data []string, w http.ResponseWriter, r *http.Request) {
	if len(data) < 3 {
		api.RespondWithJSON(w, http.StatusBadRequest, api.Response{
			Error: apperrors.ErrInvalidCommand.Error(),
		})
		return
	}
	list := data[2:]
	dt.Qpush(r.Context(), data[1], list)
	api.RespondWithJSON(w, http.StatusOK, api.Response{})

}

func QpopCommand(dt *Datastore, data []string, w http.ResponseWriter, r *http.Request) {
	if len(data) < 2 {
		api.RespondWithJSON(w, http.StatusBadRequest, api.Response{
			Error: apperrors.ErrInvalidCommand.Error(),
		})
		return
	}
	value := dt.Qpop(r.Context(), data[1])
	if len(value) == 0 {
		api.RespondWithJSON(w, http.StatusOK, api.Response{
			Error: "queue is empty",
		})
		return
	}
	api.RespondWithJSON(w, http.StatusOK, api.Response{
		Value: value,
	})

}

func BQpopCommand(dt *Datastore, data []string, w http.ResponseWriter, r *http.Request) {
	if len(data) < 3 {
		api.RespondWithJSON(w, http.StatusBadRequest, api.Response{
			Error: apperrors.ErrInvalidCommand.Error(),
		})
		return
	}
	timeout, err := strconv.ParseFloat(data[2], 64)
	if err != nil {
		api.RespondWithJSON(w, http.StatusOK, api.Response{
			Error: "Invalid time",
		})
		return
	}
	value, err := dt.Bqpop(r.Context(), data[1], timeout) // Wait up to 10 seconds
	if err != nil {
		if err == apperrors.ErrQueueEmpty || err == apperrors.ErrQueueTimeout {
			api.RespondWithJSON(w, http.StatusOK, api.Response{
				Error: "queue is empty",
			})
			return
		}
		api.RespondWithJSON(w, http.StatusInternalServerError, api.Response{
			Error: "Something went wrong!",
		})
		return
	}
	api.RespondWithJSON(w, http.StatusOK, api.Response{
		Value: value,
	})

}
