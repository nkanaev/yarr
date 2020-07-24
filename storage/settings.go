package storage

import "encoding/json"

func settingsDefaults() map[string]interface{} {
	return map[string]interface{}{
		"filter":          "",
		"feed":            "",
		"feed_list_width": 300,
		"item_list_width": 300,
	}
}

func (s *Storage) GetSettingsValue(key string) interface{} {
	row := s.db.QueryRow(`select val from settings where key=?`, key)
	if row == nil {
		return settingsDefaults()[key]
	}
	var val []byte
	row.Scan(&val)
	var valDecoded interface{}
	if err := json.Unmarshal([]byte(val), &valDecoded); err != nil {
		s.log.Print(err)
		return nil
	}
	return valDecoded
}

func (s *Storage) GetSettings() map[string]interface{} {
	result := settingsDefaults()
	rows, err := s.db.Query(`select key, val from settings;`)
	if err != nil {
		s.log.Print(err)
		return result
	}
	for rows.Next() {
		var key string
		var val []byte
		var valDecoded interface{}

		rows.Scan(&key, &val)
		if err = json.Unmarshal([]byte(val), &valDecoded); err != nil {
			s.log.Print(err)
			continue
		}
		result[key] = valDecoded
	}
	return result
}

func (s *Storage) UpdateSettings(kv map[string]interface{}) bool {
	defaults := settingsDefaults()
	for key, val := range kv {
		if defaults[key] == nil {
			continue
		}
		valEncoded, err := json.Marshal(val)
		if err != nil {
			s.log.Print(err)
			return false
		}
		_, err = s.db.Exec(`
			insert into settings (key, val) values (?, ?)
			on conflict (key) do update set val=?`,
			key, valEncoded, valEncoded,
		)
		if err != nil {
			s.log.Print(err)
			return false
		}
	}
	return true
}
