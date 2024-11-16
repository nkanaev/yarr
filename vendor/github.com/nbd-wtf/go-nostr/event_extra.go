package nostr

// Deprecated: this was never a good idea, stop using.
func (evt *Event) SetExtra(key string, value any) {
	if evt.extra == nil {
		evt.extra = make(map[string]any)
	}
	evt.extra[key] = value
}

// Deprecated: this was never a good idea, stop using.
func (evt *Event) RemoveExtra(key string) {
	if evt.extra == nil {
		return
	}
	delete(evt.extra, key)
}

// Deprecated: this was never a good idea, stop using.
func (evt Event) GetExtra(key string) any {
	ival, _ := evt.extra[key]
	return ival
}

// Deprecated: this was never a good idea, stop using.
func (evt Event) GetExtraString(key string) string {
	ival, ok := evt.extra[key]
	if !ok {
		return ""
	}
	val, ok := ival.(string)
	if !ok {
		return ""
	}
	return val
}

// Deprecated: this was never a good idea, stop using.
func (evt Event) GetExtraNumber(key string) float64 {
	ival, ok := evt.extra[key]
	if !ok {
		return 0
	}

	switch val := ival.(type) {
	case float64:
		return val
	case int:
		return float64(val)
	case int64:
		return float64(val)
	}

	return 0
}

// Deprecated: this was never a good idea, stop using.
func (evt Event) GetExtraBoolean(key string) bool {
	ival, ok := evt.extra[key]
	if !ok {
		return false
	}
	val, ok := ival.(bool)
	if !ok {
		return false
	}
	return val
}
