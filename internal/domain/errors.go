package domain

import "errors"

// ErrNotFound возвращается, когда подписка с указанным id не найдена.
var ErrNotFound = errors.New("subscription not found")
