package controller

import (
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/riskiramdan/efishery/golang/internal/concurrency"
	"github.com/riskiramdan/efishery/golang/internal/data"
	"github.com/riskiramdan/efishery/golang/internal/http/response"
	"github.com/riskiramdan/efishery/golang/internal/types"
	u "github.com/riskiramdan/efishery/golang/util"
)

// ConcurrencyList ..
type ConcurrencyList struct {
	Data  []*concurrency.ConResult `json:"data"`
	Total int                      `json:"total"`
}

// ConcurrencyController represents the user controller
type ConcurrencyController struct {
	concurrencyService concurrency.ServiceInterface
	utility            *u.Utility
	redisManager       *redis.Client
}

const (
	layoutTIMEISO = "2006-01-02"
)

// GetListPrices for getting list data prices
func (a *ConcurrencyController) GetListPrices(w http.ResponseWriter, r *http.Request) {
	result, cnt, err := a.concurrencyService.ListDataPrice(r.Context(), &r.Header)
	if err != nil {
		response.Error(w, "Internal Server Error", http.StatusInternalServerError, *err)
		return
	}
	response.JSON(w, http.StatusOK, ConcurrencyList{
		Data:  result,
		Total: cnt,
	})
}

// GetAggregation for aggregate resources data
func (a *ConcurrencyController) GetAggregation(w http.ResponseWriter, r *http.Request) {
	var err *types.Error
	var errConversion error
	queryValues := r.URL.Query()

	from := time.Time{}
	if queryValues.Get("dateFrom") != "" {
		from, errConversion = time.Parse(layoutTIMEISO, queryValues.Get("dateFrom"))
		if errConversion != nil {
			err = &types.Error{
				Path:    ".UserController->ListUser()",
				Message: errConversion.Error(),
				Error:   errConversion,
				Type:    "golang-error",
			}
			response.Error(w, "Bad Request", http.StatusBadRequest, *err)
			return
		}
	}

	to := time.Time{}
	if queryValues.Get("dateTo") != "" {
		to, errConversion = time.Parse(layoutTIMEISO, queryValues.Get("dateTo"))
		if errConversion != nil {
			err = &types.Error{
				Path:    ".UserController->ListUser()",
				Message: errConversion.Error(),
				Error:   errConversion,
				Type:    "golang-error",
			}
			response.Error(w, "Bad Request", http.StatusBadRequest, *err)
			return
		}
	}

	areaProvinsi := queryValues.Get("areaProvinsi")

	result, cnt, err := a.concurrencyService.AggregateDataPrice(r.Context(), &r.Header, &concurrency.ParamsAgg{
		DateFrom:     from,
		DateTo:       to,
		AreaProvinsi: areaProvinsi,
	})
	if cnt == 0 {
		err = &types.Error{
			Path:    ".UserController->ListUser()",
			Message: data.ErrNotFound.Error(),
			Error:   data.ErrNotFound,
			Type:    "golang-error",
		}
		response.Error(w, data.ErrNotFound.Error(), http.StatusNotFound, *err)
		return
	}
	if err != nil {
		err.Path = ".UserController->CreateUser()" + err.Path
		response.Error(w, "Internal Server Error", http.StatusInternalServerError, *err)
		return
	}
	response.JSON(w, http.StatusOK, result)
}

// NewConcurrencyController creates a new user controller
func NewConcurrencyController(
	concurrencyService concurrency.ServiceInterface,
	utility *u.Utility,
	redisManager *redis.Client,
) *ConcurrencyController {
	return &ConcurrencyController{
		concurrencyService: concurrencyService,
		utility:            utility,
		redisManager:       redisManager,
	}
}
