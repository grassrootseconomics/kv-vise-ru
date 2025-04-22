package api

import (
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/kamikazechaser/common/httputil"
	"github.com/uptrace/bunrouter"
)

type PhoneParam struct {
	Phone string `validate:"required"`
}

func (a *API) addressHandler(w http.ResponseWriter, req bunrouter.Request) error {
	r := PhoneParam{
		Phone: req.Param("phone"),
	}

	if err := a.validator.Validate(r); err != nil {
		return httputil.JSON(w, http.StatusBadRequest, ErrResponse{
			Ok:          false,
			Description: "Address validation failed",
		})
	}

	address, err := a.store.GetAddress(req.Context(), r.Phone)
	if err != nil {
		if err == pgx.ErrNoRows {
			return httputil.JSON(w, http.StatusNotFound, ErrResponse{
				Ok:          false,
				Description: "Address not found",
			})
		}
		return err
	}

	return httputil.JSON(w, http.StatusOK, OKResponse{
		Ok:          true,
		Description: "Address resolved",
		Result: map[string]any{
			"address": address,
		},
	})
}
