package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Coderx44/gg/api"
	"github.com/Coderx44/gg/apperrors"
)

func ParseCommands(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		req := map[string]string{}
		// data := [5]int{}
		err := json.NewDecoder(r.Body).Decode(&req)
		fmt.Println("req :", req)
		if err != nil {
			log.Printf("error while decoding %s", err.Error())
			api.RespondWithJSON(w, http.StatusInternalServerError, api.Response{
				Error: "Something went wrong !!",
			})
			return
		}
		data := []string{}
		value, ok := req["command"]
		if ok {
			data = strings.Split(strings.TrimSpace(value), " ")
		} else {
			api.RespondWithJSON(w, http.StatusBadRequest, api.Response{
				Error: apperrors.ErrInvalidCommand.Error(),
			})
			return
		}

		ctx := context.WithValue(r.Context(), "data", data)
		reqWithValueCtx := r.WithContext(ctx)

		next(w, reqWithValueCtx)
	}
}
