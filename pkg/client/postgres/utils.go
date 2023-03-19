package postgres

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"unicode"
)

// GenerateWhereClause : helper function to Generate where Clause
func GenerateWhereClause(query string, where map[string]interface{}) string {
	if len(where) > 0 {
		query += "WHERE "
		for key, value := range where {

			if strings.HasSuffix(key, "_ne") {
				query += NotEqualsTo(key, value)
			} else if strings.HasSuffix(key, "_lt") {
				query += LessThan(key, value)
			} else if strings.HasSuffix(key, "_gt") {
				query += GreaterThan(key, value)
			} else if strings.HasSuffix(key, "_lte") {
				query += LessThanEqualTo(key, value)
			} else if strings.HasSuffix(key, "_gte") {
				query += GreaterThanEqualTo(key, value)
			} else if strings.HasSuffix(key, "_null") {
				query += Null(key, value)
			} else if strings.HasSuffix(key, "_contains") {
				query += Contains(key, value)
			} else if strings.HasSuffix(key, "_in") {
				query += In(key, value)
			} else if strings.HasSuffix(key, "_nin") {
				query += NotIn(key, value)
			} else if strings.HasSuffix(key, "_lower") {
				query += Lower(key, value)
			} else {
				query += Equals(key, value)
			}
		}
		// remove last " AND " from query string if present
		query = strings.TrimSuffix(query, " AND ")
		//if strings.HasSuffix(query, " AND ") {
		//	query = query[:len(query)-5]
		//}
		//log.Println(query)
	}
	return query
}

func NotEqualsTo(key string, value interface{}) string {
	newKey := UpperRune(strings.TrimSuffix(key, "_ne"))

	var query string = fmt.Sprintf("%s != '%v' AND ", newKey, value)

	return query
}

func LessThan(key string, value interface{}) string {
	// less than clause
	newKey := UpperRune(strings.TrimSuffix(key, "_lt"))
	var query string = fmt.Sprintf("%s < '%v' AND ", newKey, value)

	return query
}

func GreaterThan(key string, value interface{}) string {
	// greater than
	newKey := UpperRune(strings.TrimSuffix(key, "_gt"))
	var query string = fmt.Sprintf("%s > '%v' AND ", newKey, value)
	return query
}

func LessThanEqualTo(key string, value interface{}) string {
	// less than and equal to
	newKey := UpperRune(strings.TrimSuffix(key, "_lte"))
	var query string = fmt.Sprintf("%s <= '%v' AND ", newKey, value)

	return query
}

func GreaterThanEqualTo(key string, value interface{}) string {
	// greater than and equal to
	newKey := UpperRune(strings.TrimSuffix(key, "_gte"))
	var query string = fmt.Sprintf("%s >= '%v' AND ", newKey, value)

	return query
}

func Null(key string, value interface{}) string {
	// contains _null
	newKey := UpperRune(key[:len(key)-5])

	var query string
	if value == true {
		query += fmt.Sprintf("%s IS NULL AND ", newKey)
	} else {
		query += fmt.Sprintf("%s IS NOT NULL AND ", newKey)
	}
	return query
}

func Contains(key string, value interface{}) string {
	// contains clause
	query := fmt.Sprintf("%s LIKE '%%%s%%' AND ", UpperRune(key[:len(key)-9]), value)
	return query
}

func In(key string, value interface{}) string {
	// In clause
	newKey := UpperRune(key[:len(key)-3])
	var query string
	if reflect.TypeOf(value).Kind() == reflect.String {
		query += newKey + " IN ('" + value.(string) + "') AND "
	} else if reflect.TypeOf(value).Kind() == reflect.Slice {
		switch v := value.(type) {
		case []string:
			query += newKey + " IN ('" + strings.Join(value.([]string), "','") + "') AND "
		case []int:
			// []int to []string
			var str []string
			for _, v := range value.([]int) {
				str = append(str, fmt.Sprintf("%v", v))
			}
			query += newKey + " IN ('" + strings.Join(str, "','") + "') AND "
		case []interface{}:
			var strList []string
			//query += newKey + " IN ('" + strings.Join(value.([]string), "','") + "') AND "
			for _, stuff := range value.([]interface{}) {
				switch v2 := stuff.(type) {
				case int, string:
					strList = append(strList, fmt.Sprintf("%v", v2))
				default:
					fmt.Printf("I don't know about type %T!\n", v)
				}
			}
			query += newKey + " IN ('" + strings.Join(strList, "','") + "') AND "
		default:
			fmt.Printf("I don't know about type %T!\n", v)
		}
	} else {
		log.Println("Invalid type for IN clause")
	}
	return query
}

func NotIn(key string, value interface{}) string {
	// Not in clause
	newKey := UpperRune(key[:len(key)-4])
	var query string
	if reflect.TypeOf(value).Kind() == reflect.String {
		query += newKey + " NOT IN ('" + value.(string) + "') AND "
	} else if reflect.TypeOf(value).Kind() == reflect.Slice {
		switch v := value.(type) {
		case []string:
			query += newKey + " NOT IN (" + strings.Join(value.([]string), ",") + ") AND "
		case []int:
			// []int to []string
			var str []string
			for _, v := range value.([]int) {
				str = append(str, fmt.Sprintf("%v", v))
			}
			query += newKey + " NOT IN (" + strings.Join(str, ",") + ") AND "
		case []interface{}:
			var strList []string
			//query += newKey + " IN ('" + strings.Join(value.([]string), "','") + "') AND "
			for _, stuff := range value.([]interface{}) {
				switch v2 := stuff.(type) {
				case int, string:
					strList = append(strList, fmt.Sprintf("%v", v2))
				default:
					fmt.Printf("I don't know about type %T!\n", v)
				}
			}
			query += newKey + " NOT IN ('" + strings.Join(strList, "','") + "') AND "
		default:
			log.Printf("I don't know about type %T!\n", v)
		}
	} else {
		log.Println("Invalid type for NOT IN clause")
	}
	return query
}

func Lower(key string, value interface{}) string {
	newKey := fmt.Sprintf(`lower("%v")=`, strings.TrimSuffix(key, "_lower"))
	// if value is string
	if reflect.TypeOf(value).Kind() == reflect.String {
		newKey += fmt.Sprintf("lower('%v')", value) + " AND "
	} else {
		newKey += fmt.Sprintf("lower(%v)", value) + " AND "
	}
	return newKey
}

func Equals(key string, value interface{}) string {
	query := UpperRune(key) + " = "
	// if value is string
	if reflect.TypeOf(value).Kind() == reflect.String {
		query += fmt.Sprintf("'%v'", value) + " AND "
	} else {
		query += fmt.Sprintf("%v", value) + " AND "
	}
	return query
}

func UpperRune(key string) string {
	switch key {
	case "user", "default":
		return fmt.Sprintf(`"%s"`, key)
	}

	runeArr := []rune(key)
	for i := 0; i < len(runeArr); i++ {
		if unicode.IsUpper(runeArr[i]) {
			stringLiteral := fmt.Sprintf(`"%s"`, key)
			return stringLiteral
		}
	}
	return key
}
