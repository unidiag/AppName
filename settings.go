package main

import "main/models"

func loadSettings() {
	settings = make(map[string]string)
	var rows []models.Setting
	db.Find(&rows)
	for _, s := range rows {
		settings[s.Key] = s.Value
	}
}

func setSetting(key, value string, desc ...string) {
	if settings == nil {
		settings = make(map[string]string)
	}
	settings[key] = value
	row := models.Setting{
		Key:   key,
		Value: value,
	}
	if len(desc) > 0 {
		row.Description = desc[0]
	}
	db.Save(&row)
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
