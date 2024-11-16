package eventstore

import "errors"

var ErrDupEvent = errors.New("duplicate: event already exists")
