package auth

import (
	"context"
	"fmt"
	"reflect"
	"strconv"

	log "github.com/sirupsen/logrus"
)

// Changes user info given a userID string
// and a pointer to a struct containing the new info
func ChangeUserInfo(ctx context.Context, userID string, newUserInfo *UserInfo) error {
	if newUserInfo == nil {
		return nil
	}

	v := reflect.ValueOf(*newUserInfo)
	uInfoFields := make([]interface{}, 0)
	uInfoFieldTypes := make([]string, 0)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.IsNil() || v.Type().Field(i).Tag.Get("sql") == "" {
			continue
		}
		uInfoFields = append(uInfoFields, field.Elem().String())
		uInfoFieldTypes = append(uInfoFieldTypes, v.Type().Field(i).Tag.Get("sql"))
	}
	numFields := len(uInfoFields)
	if numFields == 0 {
		return nil
	}

	db := GetDBFromContext(ctx)
	query := "UPDATE users SET "

	for i, fieldType := range uInfoFieldTypes {
		query += fmt.Sprintf("%v=$%v", fieldType, i+1)
	}

	query += " WHERE id=$" + strconv.Itoa(numFields+1)

	log.Debug(query)

	uInfoFields = append(uInfoFields, userID)
	_, err := db.Exec(query, uInfoFields...)
	return err
}
