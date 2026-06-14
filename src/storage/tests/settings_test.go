package sqlite

import (
	"reflect"
	"strings"
	"testing"

	"github.com/nkanaev/yarr/src/storage/model"
)

func TestSettingsDefaults(t *testing.T) {
	s := testDB()
	defer s.Close()

	settings := s.GetSettings()
	defaults := settingsDefaults()

	if !reflect.DeepEqual(settings, defaults) {
		t.Errorf("expected defaults %+v, got %+v", defaults, settings)
	}
}

func TestUpdateSettings(t *testing.T) {
	s := testDB()
	defer s.Close()

	params := model.UpdateSettingsParams{
		ThemeName:     ptr("night"),
		FeedListWidth: ptr(400),
		RefreshRate:   ptr(int64(15)),
	}

	if ok := s.UpdateSettings(params); !ok {
		t.Fatal("UpdateSettings failed")
	}

	settings := s.GetSettings()

	if settings.ThemeName != "night" {
		t.Errorf("expected theme_name night, got %s", settings.ThemeName)
	}
	if settings.FeedListWidth != 400 {
		t.Errorf("expected feed_list_width 400, got %d", settings.FeedListWidth)
	}
	if settings.RefreshRate != 15 {
		t.Errorf("expected refresh_rate 15, got %d", settings.RefreshRate)
	}
}

func TestGetSettings(t *testing.T) {
	s := testDB()
	defer s.Close()

	s.UpdateSettings(model.UpdateSettingsParams{Language: ptr("fr")})

	settings := s.GetSettings()
	if settings.Language != "fr" {
		t.Errorf("expected fr, got %v", settings.Language)
	}
	if settings.ThemeName != "light" {
		t.Errorf("expected light, got %v", settings.ThemeName)
	}
}

func TestSettingsExhaustive(t *testing.T) {
	s := testDB()
	defer s.Close()

	settingsType := reflect.TypeOf(model.Settings{})
	paramsType := reflect.TypeOf(model.UpdateSettingsParams{})
	
	settings := s.GetSettings()
	m := settings.Map()

	for i := 0; i < settingsType.NumField(); i++ {
		field := settingsType.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" {
			t.Errorf("Field %s missing json tag", field.Name)
			continue
		}
		// json tags might have options like "name,omitempty", take only the first part
		jsonKey := strings.Split(jsonTag, ",")[0]

		// 1. Check Map()
		if _, ok := m[jsonKey]; !ok {
			t.Errorf("Key %q (from field %s) missing from Settings.Map()", jsonKey, field.Name)
		}

		// 2. Check UpdateSettingsParams
		foundInParams := false
		for j := 0; j < paramsType.NumField(); j++ {
			pField := paramsType.Field(j)
			pJsonTag := strings.Split(pField.Tag.Get("json"), ",")[0]
			if pJsonTag == jsonKey {
				foundInParams = true
				// Also check it's a pointer
				if pField.Type.Kind() != reflect.Ptr {
					t.Errorf("Field %s in UpdateSettingsParams should be a pointer", pField.Name)
				}
				break
			}
		}
		if !foundInParams {
			t.Errorf("Key %q (from field %s) missing from UpdateSettingsParams", jsonKey, field.Name)
		}

		// 3. Test round-trip update
		// We'll create a new UpdateSettingsParams and set ONLY this field
		paramsValue := reflect.New(paramsType).Elem()
		for j := 0; j < paramsType.NumField(); j++ {
			pField := paramsType.Field(j)
			pJsonTag := strings.Split(pField.Tag.Get("json"), ",")[0]
			if pJsonTag == jsonKey {
				// Create a new value of the underlying type
				val := reflect.New(field.Type).Elem()
				switch field.Type.Kind() {
				case reflect.String:
					val.SetString("test_" + jsonKey)
				case reflect.Int, reflect.Int64:
					val.SetInt(42)
				case reflect.Bool:
					val.SetBool(false)
				}
				paramsValue.Field(j).Set(val.Addr())
				break
			}
		}
		
		if ok := s.UpdateSettings(paramsValue.Interface().(model.UpdateSettingsParams)); !ok {
			t.Errorf("UpdateSettings failed for %q", jsonKey)
		}
		
		updated := s.GetSettings()
		updatedValue := reflect.ValueOf(updated).Field(i)
		
		switch field.Type.Kind() {
		case reflect.String:
			if updatedValue.String() != "test_"+jsonKey {
				t.Errorf("Round-trip failed for %q: expected %q, got %q (check UpdateSettings/GetSettings switch)", jsonKey, "test_"+jsonKey, updatedValue.String())
			}
		case reflect.Int, reflect.Int64:
			if updatedValue.Int() != 42 {
				t.Errorf("Round-trip failed for %q: expected 42, got %d (check UpdateSettings/GetSettings switch)", jsonKey, updatedValue.Int())
			}
		case reflect.Bool:
			if updatedValue.Bool() != false {
				t.Errorf("Round-trip failed for %q: expected false, got %v (check UpdateSettings/GetSettings switch)", jsonKey, updatedValue.Bool())
			}
		}
	}
}
