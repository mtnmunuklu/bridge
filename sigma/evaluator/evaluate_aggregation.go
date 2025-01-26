package evaluator

import (
	"fmt"

	"github.com/mtnmunuklu/bridge/sigma"
)

// evaluateAggregationExpression evaluates an aggregation expression within a Sigma rule
func (rule RuleEvaluator) evaluateAggregationExpression(aggregation sigma.AggregationExpr) (string, error) {
	var aggregationResult string

	// Determine the type of aggregation expression
	switch agg := aggregation.(type) {
	case sigma.Near:
		return aggregationResult, fmt.Errorf("near isn't supported yet")

	case sigma.Comparison:
		// Evaluate the aggregation function
		aggregationResult, err := rule.evaluateAggregationFunc(agg.Func)
		if err != nil {
			return aggregationResult, err
		}

		// Return the aggregation result with the comparison operator and threshold
		return aggregationResult + " " + string(agg.Op) + " " + fmt.Sprintf("%d", int(agg.Threshold)), nil

	default:
		// Return an error if the aggregation expression is not recognized
		return aggregationResult, fmt.Errorf("unknown aggregation expression")
	}
}

// evaluateAggregationFunc evaluates the given aggregation function and returns the resulting query string.
func (rule RuleEvaluator) evaluateAggregationFunc(aggregation sigma.AggregationFunc) (string, error) {
	var result string
	switch agg := aggregation.(type) {
	case sigma.Count:
		// If the field is not specified, count all records
		if agg.Field == "" {
			// If there is a group by clause, add it to the select statement
			if agg.GroupedBy != "" {
				result = "| stats count by" + agg.GroupedBy + "| sort -count"
			} else {
				// Add the count function to the select statement
				result = "| sort -count"
			}
			return result, nil
		} else {
			// If the field is specified, count the number of records for each value of the field
			if len(rule.fieldmappings[agg.Field]) != 0 {
				agg.Field = rule.fieldmappings[agg.Field][0]
			}
			// If there is a group by clause, add it to the select statement
			if agg.GroupedBy != "" {
				if len(rule.fieldmappings[agg.GroupedBy]) != 0 {
					agg.GroupedBy = rule.fieldmappings[agg.GroupedBy][0]
				}
				result = "| stats count by" + agg.Field + "," + agg.GroupedBy + " | sort -count"
			} else {
				// Add the count function to the select statement
				result = "| stats count by" + agg.Field + " | sort -count"
			}
			return result, nil
		}

	case sigma.Average:
		// Compute the average of the specified field
		if len(rule.fieldmappings[agg.Field]) != 0 {
			agg.Field = rule.fieldmappings[agg.Field][0]
		}

		// If there is a group by clause, add it to the select statement
		if agg.GroupedBy != "" {
			if len(rule.fieldmappings[agg.GroupedBy]) != 0 {
				agg.GroupedBy = rule.fieldmappings[agg.GroupedBy][0]
			}
			result = "| stats avg(" + agg.Field + ") by " + agg.Field + "," + agg.GroupedBy + " AS average | sort -average"
		} else {
			// Add the average function to the select statement
			result = "| stats avg(" + agg.Field + ") by " + agg.Field + " AS average | sort -average"
		}
		return result, nil

	case sigma.Sum:
		// Compute the sum of the specified field
		if len(rule.fieldmappings[agg.Field]) != 0 {
			agg.Field = rule.fieldmappings[agg.Field][0]
		}

		// If there is a group by clause, add it to the select statement
		if agg.GroupedBy != "" {
			if len(rule.fieldmappings[agg.GroupedBy]) != 0 {
				agg.GroupedBy = rule.fieldmappings[agg.GroupedBy][0]
			}
			result += "| stats sum(" + agg.Field + ") by " + agg.Field + "," + agg.GroupedBy + " AS sum | sort -sum"
		} else {
			// Add the sum function to the select statement
			result += "| stats sum(" + agg.Field + ") by " + agg.Field + " AS sum | sort -sum"
		}

		return result, nil

	case sigma.Min:
		// If the aggregation function is a Min function, map the field to its equivalent in the data source.
		if len(rule.fieldmappings[agg.Field]) != 0 {
			agg.Field = rule.fieldmappings[agg.Field][0]
		}

		// If a group by clause is specified, add it to the query and map the field to its equivalent in the data source.
		if agg.GroupedBy != "" {
			if len(rule.fieldmappings[agg.GroupedBy]) != 0 {
				agg.GroupedBy = rule.fieldmappings[agg.GroupedBy][0]
			}
			result += "| stats min(" + agg.Field + ") by " + agg.Field + "," + agg.GroupedBy + " AS min | sort -min"
		} else {
			// Add the aggregation function to the query and set the having clause to filter by the minimum value of the field.
			result += "| stats min(" + agg.Field + ") by " + agg.Field + " AS min | sort -min"
		}
		return result, nil

	case sigma.Max:
		// If the aggregation function is a Max function, map the field to its equivalent in the data source.
		if len(rule.fieldmappings[agg.Field]) != 0 {
			agg.Field = rule.fieldmappings[agg.Field][0]
		}

		// If a group by clause is specified, add it to the query and map the field to its equivalent in the data source.
		if agg.GroupedBy != "" {
			if len(rule.fieldmappings[agg.GroupedBy]) != 0 {
				agg.GroupedBy = rule.fieldmappings[agg.GroupedBy][0]
			}
			result = "| stats max(" + agg.Field + ") by " + agg.Field + "," + agg.GroupedBy + " AS max | sort -max"
		} else {
			// Add the aggregation function to the query and set the having clause to filter by the maximum value of the field.
			result = "| stats max(" + agg.Field + ") by " + agg.Field + " AS max | sort -max"
		}
		return result, nil

	// If the aggregation function type is not supported, return an error.
	default:
		return result, fmt.Errorf("unsupported aggregation function")
	}
}
