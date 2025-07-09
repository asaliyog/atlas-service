package config

import (
	"fmt"
	"strings"
	"time"
)

// FilterOperator represents a filter operator
type FilterOperator string

const (
	// String operators
	OperatorEquals        FilterOperator = "eq"
	OperatorNotEquals     FilterOperator = "ne"
	OperatorContains      FilterOperator = "contains"
	OperatorStartsWith    FilterOperator = "starts_with"
	OperatorEndsWith      FilterOperator = "ends_with"
	OperatorLike          FilterOperator = "like"
	OperatorILike         FilterOperator = "ilike"
	OperatorIn            FilterOperator = "in"
	OperatorNotIn         FilterOperator = "not_in"
	OperatorIsNull        FilterOperator = "is_null"
	OperatorIsNotNull     FilterOperator = "is_not_null"

	// Numeric/Date operators
	OperatorGreaterThan   FilterOperator = "gt"
	OperatorGreaterEqual  FilterOperator = "gte"
	OperatorLessThan      FilterOperator = "lt"
	OperatorLessEqual     FilterOperator = "lte"
	OperatorBetween       FilterOperator = "between"
)

// FieldType represents the data type of a field
type FieldType string

const (
	FieldTypeString  FieldType = "string"
	FieldTypeInt     FieldType = "int"
	FieldTypeFloat   FieldType = "float"
	FieldTypeBool    FieldType = "bool"
	FieldTypeDate    FieldType = "date"
	FieldTypeArray   FieldType = "array"
)

// FieldConfig defines the configuration for a filterable field
type FieldConfig struct {
	Type     FieldType        `json:"type"`
	Operators []FilterOperator `json:"operators"`
	Required bool             `json:"required"`
}

// FilterConfig defines the filter configuration for an endpoint
type FilterConfig struct {
	Fields map[string]FieldConfig `json:"fields"`
}

// VMsFilterConfig returns the filter configuration for the VMs endpoint
func VMsFilterConfig() FilterConfig {
	return FilterConfig{
		Fields: map[string]FieldConfig{
			"id": {
				Type:      FieldTypeString,
				Operators: []FilterOperator{OperatorEquals, OperatorNotEquals, OperatorContains, OperatorStartsWith, OperatorEndsWith, OperatorLike, OperatorILike, OperatorIn, OperatorNotIn, OperatorIsNull, OperatorIsNotNull},
			},
			"name": {
				Type:      FieldTypeString,
				Operators: []FilterOperator{OperatorEquals, OperatorNotEquals, OperatorContains, OperatorStartsWith, OperatorEndsWith, OperatorLike, OperatorILike, OperatorIn, OperatorNotIn, OperatorIsNull, OperatorIsNotNull},
			},
			"cloudType": {
				Type:      FieldTypeString,
				Operators: []FilterOperator{OperatorEquals, OperatorNotEquals, OperatorIn, OperatorNotIn, OperatorIsNull, OperatorIsNotNull},
			},
			"status": {
				Type:      FieldTypeString,
				Operators: []FilterOperator{OperatorEquals, OperatorNotEquals, OperatorIn, OperatorNotIn, OperatorIsNull, OperatorIsNotNull},
			},
			"cloudAccountId": {
				Type:      FieldTypeString,
				Operators: []FilterOperator{OperatorEquals, OperatorNotEquals, OperatorContains, OperatorStartsWith, OperatorEndsWith, OperatorLike, OperatorILike, OperatorIn, OperatorNotIn, OperatorIsNull, OperatorIsNotNull},
			},
			"location": {
				Type:      FieldTypeString,
				Operators: []FilterOperator{OperatorEquals, OperatorNotEquals, OperatorContains, OperatorStartsWith, OperatorEndsWith, OperatorLike, OperatorILike, OperatorIn, OperatorNotIn, OperatorIsNull, OperatorIsNotNull},
			},
			"instanceType": {
				Type:      FieldTypeString,
				Operators: []FilterOperator{OperatorEquals, OperatorNotEquals, OperatorContains, OperatorStartsWith, OperatorEndsWith, OperatorLike, OperatorILike, OperatorIn, OperatorNotIn, OperatorIsNull, OperatorIsNotNull},
			},
			"env": {
				Type:      FieldTypeString,
				Operators: []FilterOperator{OperatorEquals, OperatorNotEquals, OperatorIn, OperatorNotIn, OperatorIsNull, OperatorIsNotNull},
			},
			"environment": {
				Type:      FieldTypeString,
				Operators: []FilterOperator{OperatorEquals, OperatorNotEquals, OperatorIn, OperatorNotIn, OperatorIsNull, OperatorIsNotNull},
			},
			"environment.id": {
				Type:      FieldTypeString,
				Operators: []FilterOperator{OperatorEquals, OperatorNotEquals, OperatorIn, OperatorNotIn, OperatorIsNull, OperatorIsNotNull},
			},
			"environment.name": {
				Type:      FieldTypeString,
				Operators: []FilterOperator{OperatorEquals, OperatorNotEquals, OperatorContains, OperatorStartsWith, OperatorEndsWith, OperatorLike, OperatorILike, OperatorIn, OperatorNotIn, OperatorIsNull, OperatorIsNotNull},
			},
			"createdAt": {
				Type:      FieldTypeDate,
				Operators: []FilterOperator{OperatorEquals, OperatorNotEquals, OperatorGreaterThan, OperatorGreaterEqual, OperatorLessThan, OperatorLessEqual, OperatorBetween, OperatorIsNull, OperatorIsNotNull},
			},
			"updatedAt": {
				Type:      FieldTypeDate,
				Operators: []FilterOperator{OperatorEquals, OperatorNotEquals, OperatorGreaterThan, OperatorGreaterEqual, OperatorLessThan, OperatorLessEqual, OperatorBetween, OperatorIsNull, OperatorIsNotNull},
			},
		},
	}
}

// ValidateFilter validates a filter against the configuration
func (fc *FilterConfig) ValidateFilter(field, operator, value string) error {
	// Check if field exists
	fieldConfig, exists := fc.Fields[field]
	if !exists {
		return fmt.Errorf("field '%s' is not allowed for filtering", field)
	}

	// Check if operator is allowed for this field
	op := FilterOperator(operator)
	allowed := false
	for _, allowedOp := range fieldConfig.Operators {
		if op == allowedOp {
			allowed = true
			break
		}
	}
	if !allowed {
		return fmt.Errorf("operator '%s' is not allowed for field '%s' of type '%s'", operator, field, fieldConfig.Type)
	}

	// Validate value based on field type and operator
	return fc.validateValue(field, fieldConfig, op, value)
}

// validateValue validates the value based on field type and operator
func (fc *FilterConfig) validateValue(field string, fieldConfig FieldConfig, operator FilterOperator, value string) error {
	// Skip value validation for null operators
	if operator == OperatorIsNull || operator == OperatorIsNotNull {
		return nil
	}

	switch fieldConfig.Type {
	case FieldTypeString:
		return fc.validateStringValue(field, operator, value)
	case FieldTypeInt:
		return fc.validateIntValue(field, operator, value)
	case FieldTypeFloat:
		return fc.validateFloatValue(field, operator, value)
	case FieldTypeBool:
		return fc.validateBoolValue(field, operator, value)
	case FieldTypeDate:
		return fc.validateDateValue(field, operator, value)
	case FieldTypeArray:
		return fc.validateArrayValue(field, operator, value)
	default:
		return fmt.Errorf("unsupported field type '%s' for field '%s'", fieldConfig.Type, field)
	}
}

// validateStringValue validates string field values
func (fc *FilterConfig) validateStringValue(field string, operator FilterOperator, value string) error {
	if value == "" && operator != OperatorIsNull && operator != OperatorIsNotNull {
		return fmt.Errorf("value cannot be empty for field '%s' with operator '%s'", field, operator)
	}

	switch operator {
	case OperatorIn, OperatorNotIn:
		// For array operators, value should be comma-separated
		if value == "" {
			return fmt.Errorf("value cannot be empty for array operator '%s' on field '%s'", operator, field)
		}
		// Additional validation could be added here for array format
	}

	return nil
}

// validateIntValue validates integer field values
func (fc *FilterConfig) validateIntValue(field string, operator FilterOperator, value string) error {
	if value == "" {
		return fmt.Errorf("value cannot be empty for numeric field '%s' with operator '%s'", field, operator)
	}

	// Try to parse as int
	var testInt int
	if _, err := fmt.Sscanf(value, "%d", &testInt); err != nil {
		return fmt.Errorf("value '%s' is not a valid integer for field '%s'", value, field)
	}

	return nil
}

// validateFloatValue validates float field values
func (fc *FilterConfig) validateFloatValue(field string, operator FilterOperator, value string) error {
	if value == "" {
		return fmt.Errorf("value cannot be empty for numeric field '%s' with operator '%s'", field, operator)
	}

	// Try to parse as float
	var testFloat float64
	if _, err := fmt.Sscanf(value, "%f", &testFloat); err != nil {
		return fmt.Errorf("value '%s' is not a valid number for field '%s'", value, field)
	}

	return nil
}

// validateBoolValue validates boolean field values
func (fc *FilterConfig) validateBoolValue(field string, operator FilterOperator, value string) error {
	if value == "" {
		return fmt.Errorf("value cannot be empty for boolean field '%s' with operator '%s'", field, operator)
	}

	// Check if value is a valid boolean
	value = strings.ToLower(value)
	if value != "true" && value != "false" && value != "1" && value != "0" {
		return fmt.Errorf("value '%s' is not a valid boolean for field '%s'", value, field)
	}

	return nil
}

// validateDateValue validates date field values
func (fc *FilterConfig) validateDateValue(field string, operator FilterOperator, value string) error {
	if value == "" {
		return fmt.Errorf("value cannot be empty for date field '%s' with operator '%s'", field, operator)
	}

	// Try to parse as date (ISO 8601 format)
	if _, err := time.Parse(time.RFC3339, value); err != nil {
		// Try alternative formats
		if _, err := time.Parse("2006-01-02", value); err != nil {
			if _, err := time.Parse("2006-01-02T15:04:05Z", value); err != nil {
				return fmt.Errorf("value '%s' is not a valid date for field '%s'. Expected format: RFC3339 or YYYY-MM-DD", value, field)
			}
		}
	}

	return nil
}

// validateArrayValue validates array field values
func (fc *FilterConfig) validateArrayValue(field string, operator FilterOperator, value string) error {
	if value == "" {
		return fmt.Errorf("value cannot be empty for array field '%s' with operator '%s'", field, operator)
	}

	// For array fields, value should be comma-separated
	if !strings.Contains(value, ",") {
		return fmt.Errorf("value '%s' for array field '%s' should be comma-separated", value, field)
	}

	return nil
}

// ParseQueryParams parses query parameters in the format field_operator=value
func (fc *FilterConfig) ParseQueryParams(queryParams map[string][]string) ([]FilterParam, error) {
	var filters []FilterParam

	for key, values := range queryParams {
		// Skip non-filter parameters
		if key == "page" || key == "pageSize" || key == "sortBy" || key == "sortOrder" || key == "env" {
			continue
		}

		// Parse field_operator format
		parts := strings.SplitN(key, "_", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid filter parameter format: '%s'. Expected format: field_operator", key)
		}

		field := parts[0]
		operator := parts[1]

		// Get the first value (we don't support multiple values for the same filter)
		if len(values) == 0 {
			return nil, fmt.Errorf("no value provided for filter parameter '%s'", key)
		}
		value := values[0]

		// Validate the filter
		if err := fc.ValidateFilter(field, operator, value); err != nil {
			return nil, err
		}

		filters = append(filters, FilterParam{
			Field:    field,
			Operator: FilterOperator(operator),
			Value:    value,
		})
	}

	return filters, nil
}

// FilterParam represents a parsed filter parameter
type FilterParam struct {
	Field    string         `json:"field"`
	Operator FilterOperator `json:"operator"`
	Value    string         `json:"value"`
} 