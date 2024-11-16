package sdk

import (
	"slices"
	"time"
)

func appendUnique[I comparable](arr []I, item ...I) []I {
	for _, item := range item {
		if slices.Contains(arr, item) {
			return arr
		}
		arr = append(arr, item)
	}
	return arr
}

func doThisNotMoreThanOnceAnHour(key string) (doItNow bool) {
	if _dtnmtoah == nil {
		go func() {
			_dtnmtoah = make(map[string]time.Time)
			for {
				time.Sleep(time.Minute * 10)
				_dtnmtoahLock.Lock()
				now := time.Now()
				for k, v := range _dtnmtoah {
					if v.Before(now) {
						delete(_dtnmtoah, k)
					}
				}
				_dtnmtoahLock.Unlock()
			}
		}()
	}

	_dtnmtoahLock.Lock()
	defer _dtnmtoahLock.Unlock()

	_, exists := _dtnmtoah[key]
	return !exists
}
