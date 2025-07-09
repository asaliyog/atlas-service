package utils

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"golang-service/internal/config"
	"golang-service/internal/models"
)

// ApplyFilters applies filters to a slice of VMs based on the filter configuration
func ApplyFilters(vms []models.VM, filters []config.FilterParam) []models.VM {
	if len(filters) == 0 {
		return vms
	}

	var filteredVMs []models.VM

	for _, vm := range vms {
		include := true
		for _, filter := range filters {
			if !applyFilter(vm, filter) {
				include = false
				break
			}
		}
		if include {
			filteredVMs = append(filteredVMs, vm)
		}
	}

	return filteredVMs
}

// applyFilter applies a single filter to a VM
func applyFilter(vm models.VM, filter config.FilterParam) bool {
	// Get the field value using reflection
	fieldValue := GetFieldValue(vm, filter.Field)
	if fieldValue == nil {
		// Handle null operators
		if filter.Operator == config.OperatorIsNull {
			return true
		}
		if filter.Operator == config.OperatorIsNotNull {
			return false
		}
		return false
	}

	// Handle null operators
	if filter.Operator == config.OperatorIsNull {
		return false
	}
	if filter.Operator == config.OperatorIsNotNull {
		return true
	}

	// Apply the filter based on the operator
	switch filter.Operator {
	case config.OperatorEquals:
		return equals(fieldValue, filter.Value)
	case config.OperatorNotEquals:
		return !equals(fieldValue, filter.Value)
	case config.OperatorContains:
		return contains(fieldValue, filter.Value)
	case config.OperatorStartsWith:
		return startsWith(fieldValue, filter.Value)
	case config.OperatorEndsWith:
		return endsWith(fieldValue, filter.Value)
	case config.OperatorLike:
		return like(fieldValue, filter.Value)
	case config.OperatorILike:
		return ilike(fieldValue, filter.Value)
	case config.OperatorIn:
		return in(fieldValue, filter.Value)
	case config.OperatorNotIn:
		return !in(fieldValue, filter.Value)
	case config.OperatorGreaterThan:
		return greaterThan(fieldValue, filter.Value)
	case config.OperatorGreaterEqual:
		return greaterEqual(fieldValue, filter.Value)
	case config.OperatorLessThan:
		return lessThan(fieldValue, filter.Value)
	case config.OperatorLessEqual:
		return lessEqual(fieldValue, filter.Value)
	case config.OperatorBetween:
		return between(fieldValue, filter.Value)
	default:
		return false
	}
}

// GetFieldValue gets the value of a field from a VM using reflection
func GetFieldValue(vm models.VM, fieldName string) interface{} {
	// Handle nested fields (e.g., "environment.id")
	if strings.Contains(fieldName, ".") {
		return getNestedFieldValue(vm, fieldName)
	}

	v := reflect.ValueOf(vm)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fieldType := t.Field(i)
		fieldValue := v.Field(i)

		// Check JSON tag first, then field name
		jsonTag := fieldType.Tag.Get("json")
		if jsonTag != "" {
			jsonField := strings.Split(jsonTag, ",")[0]
			if jsonField == fieldName {
				return fieldValue.Interface()
			}
		}

		// Check field name
		if strings.EqualFold(fieldType.Name, fieldName) {
			return fieldValue.Interface()
		}
	}

	return nil
}

// getNestedFieldValue gets the value of a nested field (e.g., "environment.id")
func getNestedFieldValue(vm models.VM, fieldName string) interface{} {
	parts := strings.Split(fieldName, ".")
	if len(parts) != 2 {
		return nil
	}

	// Get the parent field value
	parentValue := GetFieldValue(vm, parts[0])
	if parentValue == nil {
		return nil
	}

	// Handle pointer types
	parentReflect := reflect.ValueOf(parentValue)
	if parentReflect.Kind() == reflect.Ptr {
		if parentReflect.IsNil() {
			return nil
		}
		parentReflect = parentReflect.Elem()
	}

	// Get the nested field value
	if parentReflect.Kind() == reflect.Struct {
		parentType := parentReflect.Type()
		for i := 0; i < parentReflect.NumField(); i++ {
			fieldType := parentType.Field(i)
			fieldValue := parentReflect.Field(i)

			// Check JSON tag first, then field name
			jsonTag := fieldType.Tag.Get("json")
			if jsonTag != "" {
				jsonField := strings.Split(jsonTag, ",")[0]
				if jsonField == parts[1] {
					return fieldValue.Interface()
				}
			}

			// Check field name
			if strings.EqualFold(fieldType.Name, parts[1]) {
				return fieldValue.Interface()
			}
		}
	}

	return nil
}

// equals compares two values for equality
func equals(fieldValue interface{}, filterValue string) bool {
	fieldStr := fmt.Sprintf("%v", fieldValue)
	return strings.EqualFold(fieldStr, filterValue)
}

// contains checks if a string contains another string
func contains(fieldValue interface{}, filterValue string) bool {
	fieldStr := strings.ToLower(fmt.Sprintf("%v", fieldValue))
	return strings.Contains(fieldStr, strings.ToLower(filterValue))
}

// startsWith checks if a string starts with another string
func startsWith(fieldValue interface{}, filterValue string) bool {
	fieldStr := strings.ToLower(fmt.Sprintf("%v", fieldValue))
	return strings.HasPrefix(fieldStr, strings.ToLower(filterValue))
}

// endsWith checks if a string ends with another string
func endsWith(fieldValue interface{}, filterValue string) bool {
	fieldStr := strings.ToLower(fmt.Sprintf("%v", fieldValue))
	return strings.HasSuffix(fieldStr, strings.ToLower(filterValue))
}

// like performs simple pattern matching (supports % wildcard)
func like(fieldValue interface{}, filterValue string) bool {
	fieldStr := fmt.Sprintf("%v", fieldValue)
	pattern := strings.ReplaceAll(filterValue, "%", ".*")
	return strings.Contains(fieldStr, strings.ReplaceAll(pattern, ".*", ""))
}

// ilike performs case-insensitive pattern matching
func ilike(fieldValue interface{}, filterValue string) bool {
	fieldStr := strings.ToLower(fmt.Sprintf("%v", fieldValue))
	pattern := strings.ToLower(strings.ReplaceAll(filterValue, "%", ".*"))
	return strings.Contains(fieldStr, strings.ReplaceAll(pattern, ".*", ""))
}

// in checks if a value is in a comma-separated list
func in(fieldValue interface{}, filterValue string) bool {
	fieldStr := strings.ToLower(fmt.Sprintf("%v", fieldValue))
	values := strings.Split(filterValue, ",")
	for _, value := range values {
		if strings.TrimSpace(value) == fieldStr {
			return true
		}
	}
	return false
}

// greaterThan compares two values for greater than
func greaterThan(fieldValue interface{}, filterValue string) bool {
	return compareValues(fieldValue, filterValue) > 0
}

// greaterEqual compares two values for greater than or equal
func greaterEqual(fieldValue interface{}, filterValue string) bool {
	return compareValues(fieldValue, filterValue) >= 0
}

// lessThan compares two values for less than
func lessThan(fieldValue interface{}, filterValue string) bool {
	return compareValues(fieldValue, filterValue) < 0
}

// lessEqual compares two values for less than or equal
func lessEqual(fieldValue interface{}, filterValue string) bool {
	return compareValues(fieldValue, filterValue) <= 0
}

// between checks if a value is between two values (comma-separated)
func between(fieldValue interface{}, filterValue string) bool {
	parts := strings.Split(filterValue, ",")
	if len(parts) != 2 {
		return false
	}
	
	lower := strings.TrimSpace(parts[0])
	upper := strings.TrimSpace(parts[1])
	
	return compareValues(fieldValue, lower) >= 0 && compareValues(fieldValue, upper) <= 0
}

// compareValues compares two values, handling different types
func compareValues(fieldValue interface{}, filterValue string) int {
	// Handle time.Time
	if t, ok := fieldValue.(time.Time); ok {
		// Try to parse the filter value as time
		if filterTime, err := time.Parse(time.RFC3339, filterValue); err == nil {
			if t.Before(filterTime) {
				return -1
			} else if t.After(filterTime) {
				return 1
			} else {
				return 0
			}
		}
		// Try alternative date formats
		if filterTime, err := time.Parse("2006-01-02", filterValue); err == nil {
			if t.Before(filterTime) {
				return -1
			} else if t.After(filterTime) {
				return 1
			} else {
				return 0
			}
		}
		return 0
	}

	// Handle numeric types
	switch v := fieldValue.(type) {
	case int, int8, int16, int32, int64:
		if filterInt, err := strconv.ParseInt(filterValue, 10, 64); err == nil {
			fieldInt := reflect.ValueOf(v).Int()
			if fieldInt < filterInt {
				return -1
			} else if fieldInt > filterInt {
				return 1
			} else {
				return 0
			}
		}
	case uint, uint8, uint16, uint32, uint64:
		if filterUint, err := strconv.ParseUint(filterValue, 10, 64); err == nil {
			fieldUint := reflect.ValueOf(v).Uint()
			if fieldUint < filterUint {
				return -1
			} else if fieldUint > filterUint {
				return 1
			} else {
				return 0
			}
		}
	case float32, float64:
		if filterFloat, err := strconv.ParseFloat(filterValue, 64); err == nil {
			fieldFloat := reflect.ValueOf(v).Float()
			if fieldFloat < filterFloat {
				return -1
			} else if fieldFloat > filterFloat {
				return 1
			} else {
				return 0
			}
		}
	}

	// Handle string comparison
	fieldStr := fmt.Sprintf("%v", fieldValue)
	return strings.Compare(fieldStr, filterValue)
} 