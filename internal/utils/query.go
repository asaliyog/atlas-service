package utils

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// QueryParams represents common query parameters for any list endpoint
type QueryParams struct {
	Page      int           `json:"page"`
	PageSize  int           `json:"pageSize"`
	SortBy    string        `json:"sortBy"`
	SortOrder string        `json:"sortOrder"`
	Filters   []QueryFilter `json:"filters"`
}

// QueryFilter represents a filter for any field
type QueryFilter struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}

// ParseQueryParams parses standard query parameters for any list endpoint
func ParseQueryParams(c *gin.Context) (QueryParams, error) {
	params := QueryParams{
		Page:      1,
		PageSize:  20,
		SortBy:    "createdAt",
		SortOrder: "asc",
		Filters:   []QueryFilter{},
	}

	// Parse page
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			params.Page = page
		} else {
			return params, fmt.Errorf("invalid page parameter: %s", pageStr)
		}
	}

	// Parse pageSize
	if pageSizeStr := c.Query("pageSize"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil && pageSize > 0 {
			if pageSize > 1000 {
				return params, fmt.Errorf("pageSize cannot exceed 1000")
			}
			params.PageSize = pageSize
		} else {
			return params, fmt.Errorf("invalid pageSize parameter: %s", pageSizeStr)
		}
	}

	// Parse sortBy (only allow one field)
	if sortBy := c.Query("sortBy"); sortBy != "" {
		sortFields := strings.Split(sortBy, ",")
		params.SortBy = strings.TrimSpace(sortFields[0])
	}

	// Parse sortOrder
	if sortOrder := c.Query("sortOrder"); sortOrder != "" {
		if sortOrder != "asc" && sortOrder != "desc" {
			return params, fmt.Errorf("sortOrder must be 'asc' or 'desc'")
		}
		params.SortOrder = sortOrder
	}

	// Parse filters using standard format: field=value, field_op=value
	params.Filters = ParseStandardFilters(c)

	return params, nil
}

// ParseStandardFilters parses filters in the format field=value (eq), field_op=value (for other ops)
func ParseStandardFilters(c *gin.Context) []QueryFilter {
	var filters []QueryFilter
	queryParams := c.Request.URL.Query()
	for key, values := range queryParams {
		if key == "page" || key == "pageSize" || key == "sortBy" || key == "sortOrder" {
			continue
		}
		if len(values) == 0 {
			continue
		}
		value := values[0]
		field := key
		operator := "eq"
		if idx := strings.LastIndex(key, "_"); idx != -1 {
			field = key[:idx]
			operator = key[idx+1:]
		}
		if !IsValidOperator(operator) {
			continue
		}
		convertedValue, err := ConvertFilterValue(value, operator)
		if err != nil {
			continue
		}
		filter := QueryFilter{
			Field:    field,
			Operator: operator,
			Value:    convertedValue,
		}
		filters = append(filters, filter)
	}
	return filters
}

// IsValidOperator checks if the operator is valid
func IsValidOperator(operator string) bool {
	validOperators := []string{
		"eq", "ne", "gt", "gte", "lt", "lte", 
		"contains", "starts_with", "ends_with", 
		"in", "not_in", "is_null", "is_not_null",
		"like", "ilike", "between",
	}
	
	for _, valid := range validOperators {
		if operator == valid {
			return true
		}
	}
	return false
}

// ConvertFilterValue converts the string value to appropriate type based on operator
func ConvertFilterValue(value, operator string) (interface{}, error) {
	switch operator {
	case "is_null", "is_not_null":
		return nil, nil
	case "in", "not_in":
		// Handle array values (comma-separated)
		if value == "" {
			return []string{}, nil
		}
		return strings.Split(value, ","), nil
	case "between":
		// Handle range values (comma-separated)
		if value == "" {
			return []string{}, nil
		}
		parts := strings.Split(value, ",")
		if len(parts) != 2 {
			return nil, fmt.Errorf("between operator requires exactly 2 values")
		}
		return parts, nil
	case "gt", "gte", "lt", "lte":
		// Try to convert to number
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal, nil
		}
		if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
			return floatVal, nil
		}
		// If not a number, treat as string
		return value, nil
	default:
		// For other operators, return as string
		return value, nil
	}
}

// ValidateQueryParams validates the query parameters
func ValidateQueryParams(params QueryParams) error {
	// Validate page
	if params.Page < 1 {
		return fmt.Errorf("page must be greater than 0")
	}

	// Validate page size
	if params.PageSize < 1 {
		return fmt.Errorf("pageSize must be greater than 0")
	}
	if params.PageSize > 1000 {
		return fmt.Errorf("pageSize cannot exceed 1000")
	}

	// Validate sort order
	if params.SortOrder != "" && params.SortOrder != "asc" && params.SortOrder != "desc" {
		return fmt.Errorf("sortOrder must be 'asc' or 'desc'")
	}

	// Validate filters
	for i, filter := range params.Filters {
		if err := ValidateFilter(filter); err != nil {
			return fmt.Errorf("filter %d: %w", i, err)
		}
	}

	return nil
}

// ValidateFilter validates a single filter
func ValidateFilter(filter QueryFilter) error {
	// Validate operator
	validOperators := map[string]bool{
		"eq":           true,
		"ne":           true,
		"lt":           true,
		"lte":          true,
		"gt":           true,
		"gte":          true,
		"in":           true,
		"not_in":       true,
		"contains":     true,
		"starts_with":  true,
		"ends_with":    true,
		"like":         true,
		"ilike":        true,
		"between":      true,
		"is_null":      true,
		"is_not_null":  true,
	}

	if !validOperators[filter.Operator] {
		return fmt.Errorf("unsupported operator: %s", filter.Operator)
	}

	// Validate field
	if filter.Field == "" {
		return fmt.Errorf("field cannot be empty")
	}

	// Validate value based on operator
	switch filter.Operator {
	case "in", "not_in":
		// Value should be a slice
		if filter.Value == nil {
			return fmt.Errorf("value cannot be nil for operator %s", filter.Operator)
		}
		v := reflect.ValueOf(filter.Value)
		if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
			return fmt.Errorf("value must be an array for operator %s", filter.Operator)
		}
		if v.Len() == 0 {
			return fmt.Errorf("value array cannot be empty for operator %s", filter.Operator)
		}
	case "between":
		// Value should be an array with exactly 2 elements
		if filter.Value == nil {
			return fmt.Errorf("value cannot be nil for operator %s", filter.Operator)
		}
		v := reflect.ValueOf(filter.Value)
		if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
			return fmt.Errorf("value must be an array for operator %s", filter.Operator)
		}
		if v.Len() != 2 {
			return fmt.Errorf("value array must have exactly 2 elements for operator %s", filter.Operator)
		}
	case "is_null", "is_not_null":
		// These operators don't need a value
		// Value can be nil or any value, it will be ignored
	case "contains", "starts_with", "ends_with", "like", "ilike":
		// Value should be a string
		if filter.Value == nil {
			return fmt.Errorf("value cannot be nil for operator %s", filter.Operator)
		}
		if _, ok := filter.Value.(string); !ok {
			return fmt.Errorf("value must be a string for operator %s", filter.Operator)
		}
	case "lt", "lte", "gt", "gte":
		// Value should be a comparable type (number, string, time)
		if filter.Value == nil {
			return fmt.Errorf("value cannot be nil for operator %s", filter.Operator)
		}
		if !IsComparableType(filter.Value) {
			return fmt.Errorf("value must be a comparable type (number, string, time) for operator %s", filter.Operator)
		}
	}

	return nil
}

// IsComparableType checks if a type is comparable for ordering operations
func IsComparableType(value interface{}) bool {
	switch value.(type) {
	case int, int8, int16, int32, int64,
		 uint, uint8, uint16, uint32, uint64,
		 float32, float64,
		 string, time.Time:
		return true
	}
	return false
}

// ApplyInMemoryFilters applies filters to any slice of structs
// data: slice of structs to filter
// filters: filters to apply
// fieldExtractor: function to extract field value from struct
func ApplyInMemoryFilters[T any](data []T, filters []QueryFilter, fieldExtractor func(T, string) interface{}) []T {
	filtered := data
	for _, filter := range filters {
		filtered = filterSliceByField(filtered, filter, fieldExtractor)
	}
	return filtered
}

// filterSliceByField filters a slice by a single filter
func filterSliceByField[T any](data []T, filter QueryFilter, fieldExtractor func(T, string) interface{}) []T {
	var out []T
	for _, item := range data {
		val := fieldExtractor(item, filter.Field)
		if matchFilter(val, filter.Operator, filter.Value) {
			out = append(out, item)
		}
	}
	return out
}

// matchFilter checks if a value matches a filter operator and value
func matchFilter(val interface{}, operator string, filterVal interface{}) bool {
	vs := fmt.Sprintf("%v", val)
	fs := fmt.Sprintf("%v", filterVal)
	switch operator {
	case "eq":
		return vs == fs
	case "ne":
		return vs != fs
	case "contains":
		return strings.Contains(vs, fs)
	case "starts_with":
		return strings.HasPrefix(vs, fs)
	case "ends_with":
		return strings.HasSuffix(vs, fs)
	case "in":
		if arr, ok := filterVal.([]string); ok {
			for _, v := range arr {
				if vs == v {
					return true
				}
			}
		}
		return false
	case "not_in":
		if arr, ok := filterVal.([]string); ok {
			for _, v := range arr {
				if vs == v {
					return false
				}
			}
			return true
		}
		return false
	case "is_null":
		return val == nil || vs == ""
	case "is_not_null":
		return val != nil && vs != ""
	default:
		return false
	}
}

// ApplyPagination applies pagination to any slice
func ApplyPagination[T any](data []T, page, pageSize int) []T {
	if len(data) == 0 {
		return data
	}
	
	start := (page - 1) * pageSize
	end := start + pageSize
	
	if start >= len(data) {
		return []T{}
	}
	
	if end > len(data) {
		end = len(data)
	}
	
	return data[start:end]
}

// CalculatePagination calculates pagination metadata
func CalculatePagination(totalItems, page, pageSize int) (int, int) {
	totalPages := totalItems / pageSize
	if totalItems%pageSize != 0 {
		totalPages++
	}
	return totalItems, totalPages
} 