package concurrency

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/riskiramdan/efishery/golang/internal/data"
	"github.com/riskiramdan/efishery/golang/internal/hosts"
	"github.com/riskiramdan/efishery/golang/internal/types"

	"github.com/patrickmn/go-cache"
	"github.com/thedevsaddam/gojsonq"
)

// ConStorage ..
type ConStorage struct {
	UUID          string `json:"uuid"`
	Komoditas     string `json:"komoditas"`
	AreaProvinsi  string `json:"area_provinsi"`
	AreaKota      string `json:"area_kota"`
	Size          string `json:"size"`
	Price         string `json:"price"`
	TanggalParsed string `json:"tgl_parsed"`
	TimeStamp     string `json:"timestamp"`
}

// ConResult ..
type ConResult struct {
	UUID          string    `json:"uuid"`
	Komoditas     string    `json:"komoditas"`
	AreaProvinsi  string    `json:"area_provinsi"`
	AreaKota      string    `json:"area_kota"`
	Size          int       `json:"size"`
	Price         float64   `json:"price"`
	PriceUSD      float64   `json:"priceUSD"`
	TanggalParsed time.Time `json:"tgl_parsed"`
	TimeStamp     int       `json:"timestamp"`
}

// AggResponse ..
type AggResponse struct {
	Min     float64 `json:"min"`
	Max     float64 `json:"max"`
	Median  float64 `json:"median"`
	Average float64 `json:"average"`
}

// ConIDRUSD  ..
type ConIDRUSD struct {
	UsdIdr *float64 `json:"USD_IDR"`
}

// ParamsAgg ..
type ParamsAgg struct {
	AreaProvinsi string    `json:"areaProvinsi"`
	DateFrom     time.Time `json:"dateFrom"`
	DateTo       time.Time `json:"dateTo"`
}

// Prices ..
type Prices struct {
	Price float64
}

// ServiceInterface represents the user service interface
type ServiceInterface interface {
	ListDataPrice(ctx context.Context, header *http.Header) ([]*ConResult, int, *types.Error)
	AggregateDataPrice(ctx context.Context, header *http.Header, params *ParamsAgg) (interface{}, int, *types.Error)
}

// Service is the domain logic implementation of user Service interface
type Service struct {
	httpManager  *hosts.HTTPManager
	redisManager *redis.Client
}

// ListDataPrice ..
func (s *Service) ListDataPrice(ctx context.Context, header *http.Header) ([]*ConResult, int, *types.Error) {
	key := "CONVERTION"
	res := []*ConResult{}
	if getRedisCache(ctx, s.redisManager, key) != "" {
		if err := json.Unmarshal([]byte(getRedisCache(ctx, s.redisManager, key)), &res); err != nil {
			fmt.Println("get")
			return nil, 0, &types.Error{
				Path:    ".ConcurrencyService->ListDataPrice()",
				Message: err.Error(),
				Error:   err,
				Type:    "validation-error",
			}
		}
		return res, len(res), nil
	}

	result, err := getListData(s.httpManager, header)
	if err != nil {
		return nil, 0, &types.Error{
			Path:    ".ConcurrencyService->ListDataPrice()",
			Message: err.Error(),
			Error:   err,
			Type:    "validation-error",
		}
	}

	conv, err := getUSD(ctx, s.redisManager, s.httpManager, header)
	if err != nil {
		return nil, 0, &types.Error{
			Path:    ".ConcurrencyService->ListDataPrice()",
			Message: err.Error(),
			Error:   err,
			Type:    "validation-error",
		}
	}

	if len(result) > 0 {
		for i := 0; i < len(result); i++ {
			if result[i].UUID != "" {
				a := &ConResult{
					UUID:          result[i].UUID,
					Komoditas:     result[i].Komoditas,
					AreaProvinsi:  result[i].AreaProvinsi,
					PriceUSD:      convertUSD(result[i].Price, *conv.UsdIdr),
					Price:         convertFloat(result[i].Price),
					AreaKota:      result[i].AreaKota,
					Size:          convertInt(result[i].Size),
					TanggalParsed: convertDate(result[i].TanggalParsed[:19] + "+0700"),
					TimeStamp:     convertInt(result[i].TimeStamp),
				}
				res = append(res, a)
			}
		}
	}

	setRedisCache(ctx, s.redisManager, key, res)

	return res, len(res), nil
}

// AggregateDataPrice ..
func (s *Service) AggregateDataPrice(ctx context.Context, header *http.Header, params *ParamsAgg) (interface{}, int, *types.Error) {
	dataConversion, lendata, err := s.ListDataPrice(ctx, header)
	if err != nil {
		err.Path = ".ConcurrencyService->AggregateDataPrice()" + err.Path
		return nil, 0, err
	}
	if lendata == 0 {
		return nil, 0, &types.Error{
			Path:    ".ConcurrencyService->AggregateDataPrice()",
			Message: data.ErrNotFound.Error(),
			Error:   data.ErrNotFound,
			Type:    "validation-error",
		}
	}

	stringConversion, errConvert := jsonString(dataConversion)
	if errConvert != nil {
		return nil, 0, &types.Error{
			Path:    ".ConcurrencyService->AggregateDataPrice()",
			Message: errConvert.Error(),
			Error:   errConvert,
			Type:    "validation-error",
		}
	}

	var from, to int64
	var areaProvinsi = ""

	jsonq := gojsonq.New()
	query := jsonq.FromString(stringConversion)
	if !params.DateFrom.IsZero() {
		from = params.DateFrom.UnixNano() / int64(time.Millisecond)
		query = query.Where("timestamp", ">", from)
	}
	if !params.DateTo.IsZero() {
		to = params.DateTo.UnixNano() / int64(time.Millisecond)
		query = query.Where("timestamp", "<", to)
	}
	if params.AreaProvinsi != "" {
		areaProvinsi = params.AreaProvinsi
		query = query.Where("area_provinsi", "contains", areaProvinsi)
	}
	cnt := query.Count()
	return map[string]interface{}{
		"price": &AggResponse{
			Average: query.Avg("price"),
			Min:     query.Min("price"),
			Max:     query.Max("price"),
			Median:  findMedian(query.Select("price").Get()),
		}, "total": cnt,
	}, cnt, nil
}

func getListData(hm *hosts.HTTPManager, header *http.Header) ([]*ConStorage, error) {
	result := []*ConStorage{}
	url := "https://stein.efishery.com/v1/storages/5e1edf521073e315924ceab4/list"
	data, err := hm.HTTPGet(url, *header)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func getUSD(ctx context.Context, r *redis.Client, hm *hosts.HTTPManager, header *http.Header) (*ConIDRUSD, error) {
	result := &ConIDRUSD{}
	key := "USD_IDR"
	if getRedisCache(ctx, r, key) != "" {
		if err := json.Unmarshal([]byte(getRedisCache(ctx, r, key)), &result); err != nil {
			return nil, err
		}
		fmt.Println("get from redis idr")
		return result, nil
	}
	url := "https://free.currconv.com/api/v7/convert?q=USD_IDR&compact=ultra&apiKey=bbc99fbf718c1a23be93"
	data, err := hm.HTTPGet(url, *header)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	fmt.Println("get from hit api idr")
	setRedisCache(ctx, r, key, result)
	return result, err
}

func convertUSD(price string, priceUSD float64) float64 {
	return math.Abs(convertFloat(price) / priceUSD)
}

func convertFloat(str string) float64 {
	p, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0
	}
	return p
}

func convertInt(str string) int {
	p, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return p
}

func convertDate(str string) time.Time {
	layout := "2006-01-02T15:04:05-0700"
	t, err := time.Parse(layout, str)
	if err != nil {
		return time.Time{}
	}
	return t
}

func setCache(c *cache.Cache, key string, res interface{}) {
	c.Set(key, res, time.Hour*12)
}

func getCache(c *cache.Cache, key string) (interface{}, bool) {
	result, found := c.Get(key)
	return result, found
}

func setRedisCache(ctx context.Context, r *redis.Client, key string, value interface{}) error {
	result, err := json.Marshal(&value)
	if err != nil {
		return err
	}
	if err := r.Set(ctx, key, result, time.Hour*12).Err(); err != nil {
		return err
	}
	return nil
}

func getRedisCache(ctx context.Context, r *redis.Client, key string) string {
	return r.Get(ctx, key).Val()
}

func jsonString(data interface{}) (string, error) {
	a, errType := json.Marshal(data)
	if errType != nil {
		return "", nil
	}
	return string(a), nil
}

func findMedian(data interface{}) float64 {
	x, errMarshall := json.Marshal(data)
	if errMarshall != nil {
		return 0
	}
	prices := []*Prices{}
	if errMarshall := json.Unmarshal(x, &prices); errMarshall != nil {
		return 0
	}
	sort.SliceStable(prices[:], func(i, j int) bool {
		return prices[i].Price < prices[j].Price
	})
	if len(prices) == 0 {
		return 0
	}
	res := prices[len(prices)/2].Price
	return res
}

// NewService creates a new concurrency AppService
func NewService(
	httpManager *hosts.HTTPManager,
	redisManager *redis.Client,
) *Service {
	return &Service{
		httpManager:  httpManager,
		redisManager: redisManager,
	}
}
