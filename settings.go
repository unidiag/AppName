package main

import (
	"main/models"

	"github.com/spf13/cast"
)

func loadSettings() {
	settings = make(map[string]string)
	var rows []models.Setting
	db.Find(&rows)
	for _, s := range rows {
		settings[s.Key] = s.Value
	}
}

func setSetting(key, value string, another ...string) {
	if settings == nil {
		settings = make(map[string]string)
	}
	settings[key] = value
	row := models.Setting{
		Key:   key,
		Value: value,
	}
	if len(another) > 0 {
		row.Description = another[0]
		if len(another) == 2 {
			row.Position = cast.ToUint(another[1])
		}
	}

	if dbname != "" {
		db.Save(&row)
	}
}

func getSetting(key string, def ...string) string {
	val, ok := settings[key]
	if !ok || val == "" {
		if len(def) > 0 {
			return def[0]
		}
		return ""
	}
	return val
}
