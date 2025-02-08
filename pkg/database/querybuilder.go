package db

import "strings"

type WhereCondition struct {
	Key       string `json:"key"`
	Condition string `json:"condition"`
	Value     string `json:"value"`
	Table     string `json:"table"`
	Joins     string `json:"joins"`
	//JsonKey        *string              `json:"json_key"`
	//GroupCondition *GroupWhereCondition `json:"group_condition"`
	//SubQuery       *SubQueryCondition   `json:"sub_query"`
}

func QueryBuilder(where []WhereCondition, typeQuery string) (string, string) {
	if len(where) == 0 {
		return "", ""

	}

	var queryParts []string
	Joins := ""
	for _, value := range where {
		if typeQuery == "JOIN" {
			condition := value.Table + "." + value.Key + " " + value.Condition + " " + value.Value + ""
			queryParts = append(queryParts, condition)
		} else {
			condition := value.Key + " " + value.Condition + " " + value.Value + ""
			queryParts = append(queryParts, condition)
		}
		Joins = Joins + value.Joins
	}

	query := Joins + "WHERE " + strings.Join(queryParts, " AND ")

	return query, Joins
}
