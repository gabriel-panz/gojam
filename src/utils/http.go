package utils

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	gojamErrors "github.com/gabriel-panz/gojam/errors"
	"github.com/gabriel-panz/gojam/types"
	"github.com/google/go-querystring/query"
)

func DecodeJsonBody[K interface{}](resp *http.Response) (*K, error) {
	p := new(K)
	err := json.NewDecoder(resp.Body).Decode(p)

	if err != nil {
		return nil, err
	}

	return p, nil
}

func EncodeQueryParameters(k interface{}) (string, error) {
	v, err := query.Values(k)
	if err != nil {
		return "", err
	}
	return v.Encode(), nil
}

func HandleHttpError(err error, logger *log.Logger, w http.ResponseWriter) {
	logger.Println(err)

	var code int

	switch err {
	case gojamErrors.ErrUnauthorized:
		code = http.StatusUnauthorized
	case gojamErrors.ErrNotFound:
		code = http.StatusNotFound
	default:
		code = http.StatusBadRequest
	}

	http.Error(w, err.Error(), code)
}

func GetAuthorizedUser(r *http.Request) types.Auth {
	return r.Context().Value(types.AuthorizedKey).(types.Auth)
}

func GetPagination(r *http.Request) (*types.Pagination, error) {
	var err error = nil
	var index int64 = 0
	pgStr := r.FormValue("page")
	if len(pgStr) > 0 {
		index, err = strconv.ParseInt(pgStr, 10, 32)
		if err != nil {
			return nil, errors.New("page number is not an integer")
		}
	}

	var size int64 = 0
	szStr := r.FormValue("size")
	if len(szStr) > 0 {
		size, err = strconv.ParseInt(szStr, 10, 32)
		if err != nil {
			return nil, errors.New("page size is not an integer")
		}
	}

	return &types.Pagination{
		PageIndex: int(index),
		PageSize:  int(size),
	}, nil
}
