package api

import (
	"fmt"
	"net/http"

	"github.com/grassrootseconomics/kv-vise-ru/pkg/data"
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

func (a *API) reverseAddressHandler(w http.ResponseWriter, req bunrouter.Request) error {
	r := struct {
		Address string `validate:"required,eth_addr_checksum"`
	}{
		Address: req.Param("address"),
	}

	if err := a.validator.Validate(r); err != nil {
		return httputil.JSON(w, http.StatusBadRequest, ErrResponse{
			Ok:          false,
			Description: "Address validation failed",
		})
	}

	sessionID, err := a.store.ReverseAddress(req.Context(), r.Address)
	if err != nil {
		if err == pgx.ErrNoRows {
			return httputil.JSON(w, http.StatusNotFound, ErrResponse{
				Ok:          false,
				Description: "Phone not found",
			})
		}
		return err
	}

	return httputil.JSON(w, http.StatusOK, OKResponse{
		Ok:          true,
		Description: "Phone resolved",
		Result: map[string]any{
			"phone": sessionID,
		},
	})
}

func (a *API) profileDetailsForSMSHandler(w http.ResponseWriter, req bunrouter.Request) error {
	r := PhoneParam{
		Phone: req.Param("phone"),
	}

	if err := a.validator.Validate(r); err != nil {
		return httputil.JSON(w, http.StatusBadRequest, ErrResponse{
			Ok:          false,
			Description: "Phone validation failed",
		})
	}

	details, err := a.store.GetProfileDetailsForSMS(req.Context(), r.Phone)
	if err != nil {
		if err == pgx.ErrNoRows {
			return httputil.JSON(w, http.StatusNotFound, ErrResponse{
				Ok:          false,
				Description: "Profile not found",
			})
		}
		return err
	}

	return httputil.JSON(w, http.StatusOK, OKResponse{
		Ok:          true,
		Description: "Profile details",
		Result: map[string]any{
			"publicKey":    details.PublicKey,
			"firstName":    details.FirstName,
			"familyName":   details.FamilyName,
			"languageCode": details.LanguageCode,
			"accountAlias": details.AccountAlias,
		},
	})
}

func (a *API) dumpAllHandler(w http.ResponseWriter, req bunrouter.Request) error {
	phone := req.Param("phone")
	if phone == "" {
		return httputil.JSON(w, http.StatusBadRequest, ErrResponse{
			Ok:          false,
			Description: "Phone is required",
		})
	}

	prefix := data.EncodeSessionID(phone)
	res, err := a.store.GetSessionData(req.Context(), prefix)
	if err != nil {
		return err
	}

	names := map[uint16]string{
		data.DATA_TRACKING_ID:            "DATA_TRACKING_ID",
		data.DATA_PUBLIC_KEY:             "DATA_PUBLIC_KEY",
		data.DATA_ACCOUNT_PIN:            "DATA_ACCOUNT_PIN",
		data.DATA_FIRST_NAME:             "DATA_FIRST_NAME",
		data.DATA_FAMILY_NAME:            "DATA_FAMILY_NAME",
		data.DATA_YOB:                    "DATA_YOB",
		data.DATA_LOCATION:               "DATA_LOCATION",
		data.DATA_GENDER:                 "DATA_GENDER",
		data.DATA_OFFERINGS:              "DATA_OFFERINGS",
		data.DATA_RECIPIENT:              "DATA_RECIPIENT",
		data.DATA_AMOUNT:                 "DATA_AMOUNT",
		data.DATA_TEMPORARY_VALUE:        "DATA_TEMPORARY_VALUE",
		data.DATA_ACTIVE_SYM:             "DATA_ACTIVE_SYM",
		data.DATA_ACTIVE_BAL:             "DATA_ACTIVE_BAL",
		data.DATA_BLOCKED_NUMBER:         "DATA_BLOCKED_NUMBER",
		data.DATA_PUBLIC_KEY_REVERSE:     "DATA_PUBLIC_KEY_REVERSE",
		data.DATA_ACTIVE_DECIMAL:         "DATA_ACTIVE_DECIMAL",
		data.DATA_ACTIVE_ADDRESS:         "DATA_ACTIVE_ADDRESS",
		data.DATA_INCORRECT_PIN_ATTEMPTS: "DATA_INCORRECT_PIN_ATTEMPTS",
		data.DATA_SELECTED_LANGUAGE_CODE: "DATA_SELECTED_LANGUAGE_CODE",
		data.DATA_INITIAL_LANGUAGE_CODE:  "DATA_INITIAL_LANGUAGE_CODE",
		data.DATA_ACCOUNT_ALIAS:          "DATA_ACCOUNT_ALIAS",
		data.DATA_VOUCHER_SYMBOLS:        "DATA_VOUCHER_SYMBOLS",
		data.DATA_VOUCHER_BALANCES:       "DATA_VOUCHER_BALANCES",
		data.DATA_VOUCHER_DECIMALS:       "DATA_VOUCHER_DECIMALS",
		data.DATA_VOUCHER_ADDRESSES:      "DATA_VOUCHER_ADDRESSES",
		data.DATA_TX_SENDERS:             "DATA_TX_SENDERS",
		data.DATA_TX_RECIPIENTS:          "DATA_TX_RECIPIENTS",
		data.DATA_TX_VALUES:              "DATA_TX_VALUES",
		data.DATA_TX_ADDRESSES:           "DATA_TX_ADDRESSES",
		data.DATA_TX_HASHES:              "DATA_TX_HASHES",
		data.DATA_TX_DATES:               "DATA_TX_DATES",
		data.DATA_TX_SYMBOLS:             "DATA_TX_SYMBOLS",
		data.DATA_TX_DECIMALS:            "DATA_TX_DECIMALS",
		data.DATA_TRANSACTIONS:           "DATA_TRANSACTIONS",
	}

	out := make(map[string][]string, len(res))
	for k, v := range res {
		name, ok := names[k]
		if !ok {
			name = fmt.Sprintf("DATA_%d", k)
		}
		out[name] = v
	}

	return httputil.JSON(w, http.StatusOK, OKResponse{
		Ok:          true,
		Description: "Dump successful",
		Result: map[string]any{
			"data": out,
		},
	})
}
