package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ikasamt/rd/pkg/redmine"
)

// resolveCustomFields は "name=value" または "id=value" 形式のカスタムフィールド指定を解決する。
// キーが数値ならIDとして直接使用し、文字列なら名前からIDを解決する。
func resolveCustomFields(client *redmine.Client, fields []string) ([]redmine.CustomFieldValue, error) {
	var result []redmine.CustomFieldValue
	var needNameResolve []struct {
		name  string
		value string
	}

	for _, field := range fields {
		parts := strings.SplitN(field, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid field format '%s': expected name=value", field)
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		if id, err := strconv.Atoi(key); err == nil {
			result = append(result, redmine.CustomFieldValue{ID: id, Value: val})
		} else {
			needNameResolve = append(needNameResolve, struct {
				name  string
				value string
			}{key, val})
		}
	}

	if len(needNameResolve) > 0 {
		for _, nr := range needNameResolve {
			cf, err := client.FindCustomFieldByName(nr.name)
			if err != nil {
				return nil, fmt.Errorf("failed to resolve custom field '%s': %w", nr.name, err)
			}
			result = append(result, redmine.CustomFieldValue{ID: cf.ID, Value: nr.value})
		}
	}

	return result, nil
}
