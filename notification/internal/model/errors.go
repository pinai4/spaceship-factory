package model

import "errors"

var (
	ErrUserTelegramChatNotSpecified = errors.New("user telegram chat not specified")
	ErrUserTelegramChatInvalid      = errors.New("user telegram chat invalid")
)
