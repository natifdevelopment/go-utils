package utils

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mcuadros/go-defaults"
	"gorm.io/gorm"
)

// PaginationLimit is the global max page size. Services should set this
// during initialization (e.g. from configs.PAGINATION_LIMIT).
var PaginationLimit int

type BasePagination struct {
	Page       int    `json:"page" default:"1"`
	PageSize   int    `json:"page_size" default:"10"`
	Sort       string `json:"sort"`
	TotalRows  int64  `json:"total_rows"`
	TotalPages int    `json:"total_pages"`
}

type PaginatedResponse struct {
	Status     bool           `json:"status"`
	Message    string         `json:"message"`
	Pagination BasePagination `json:"pagination"`
	Data       interface{}    `json:"data"`
}

type Pagination struct {
	Db         *gorm.DB
	Page       int         `json:"page" default:"1"`
	PageSize   int         `json:"page_size" default:"10"`
	Sort       string      `json:"sort" default:"-id"`
	TotalRows  int64       `json:"total_rows"`
	TotalPages int         `json:"total_pages"`
	Data       interface{} `json:"data"`
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination) GetLimit() int {
	if p.PageSize == 0 {
		p.PageSize = 10
	}
	if PaginationLimit > 0 && p.PageSize > PaginationLimit {
		p.PageSize = PaginationLimit
	}
	return p.PageSize
}

func (p *Pagination) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *Pagination) GetSort() []string {
	if p.Sort == "" {
		p.Sort = "ID"
	}
	sortFields := strings.Split(p.Sort, ",")
	var sortClauses []string

	for _, field := range sortFields {
		field = strings.TrimSpace(field)
		if field != "" {
			sortClause := field + " asc"
			if strings.HasPrefix(field, "-") {
				sortClause = strings.TrimPrefix(field, "-") + " desc"
			}
			sortClauses = append(sortClauses, strings.ToLower(sortClause))
		}
	}
	return sortClauses
}

func InitPagination(ctx *gin.Context, db *gorm.DB) *Pagination {
	pagination := Pagination{
		Db: db,
	}
	defaults.SetDefaults(&pagination)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	if page > 0 {
		pagination.Page = page
	}
	if pageSize > 0 {
		pagination.PageSize = pageSize
	}
	if PaginationLimit > 0 && pagination.PageSize > PaginationLimit {
		pagination.PageSize = PaginationLimit
	}
	return &pagination
}

func (p *Pagination) SortQuery(results interface{}) ([]string, error) {
	sliceType := reflect.TypeOf(results)

	if sliceType.Kind() == reflect.Ptr {
		sliceType = sliceType.Elem()
	}
	if sliceType.Kind() != reflect.Slice {
		return nil, errors.New("not a slice or pointer to slice, cannot proceed")
	}
	modelType := sliceType.Elem()
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}
	modelInstance := reflect.New(modelType).Elem().Interface()

	validFields, err := GetModelFields(modelInstance)
	if err != nil {
		return nil, err
	}
	for i, field := range validFields {
		validFields[i] = strings.ToLower(field)
	}
	var validSortClauses []string
	for _, sortClause := range p.GetSort() {
		parts := strings.Fields(sortClause)
		if len(parts) != 2 {
			continue
		}
		field, direction := parts[0], parts[1]
		field = strings.ToLower(field)
		if Contains(SliceStrToAny(validFields), field) {
			validSortClauses = append(validSortClauses, fmt.Sprintf("%s %s", field, direction))
		}
	}
	for _, sortClause := range validSortClauses {
		p.Db = p.Db.Order(sortClause)
	}
	return validSortClauses, nil
}

func (p *Pagination) PaginateQuery(results interface{}) {
	limit := p.GetLimit()
	page := p.GetPage()

	var totalRows int64
	p.Db.Model(results).Count(&totalRows)
	p.TotalRows = totalRows
	totalPages := int(totalRows/int64(limit)) + 1
	if totalRows <= 1 {
		totalPages = 1
	}
	p.TotalPages = totalPages
	offset := (page - 1) * limit
	sortClauses, err := p.SortQuery(results)
	if err != nil {
		return
	}
	p.Sort = strings.Join(sortClauses, ",")
	p.Db.Offset(offset).Limit(limit).Find(results)
	p.Data = results
}

func (p *Pagination) GetPaginatedResponse() *PaginatedResponse {
	status := true
	message := "Data Found."
	if p.TotalRows == 0 {
		status = false
		message = "Data Not Found."
	}
	return &PaginatedResponse{
		Status:  status,
		Message: message,
		Data:    p.Data,
		Pagination: BasePagination{
			Page:       p.Page,
			PageSize:   p.PageSize,
			Sort:       p.Sort,
			TotalRows:  p.TotalRows,
			TotalPages: p.TotalPages,
		},
	}
}

func Paginate(value interface{}, pagination *Pagination) func(db *gorm.DB) *gorm.DB {
	limit := pagination.GetLimit()
	var totalRows int64
	pagination.Db.Model(value).Count(&totalRows)
	pagination.TotalRows = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(limit)))
	pagination.TotalPages = totalPages
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(pagination.GetSort())
	}
}
