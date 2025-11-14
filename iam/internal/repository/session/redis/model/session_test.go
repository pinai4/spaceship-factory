//go:build unit || !integration

package model

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUserNotificationMethods_MarshalBinary(t *testing.T) {
	orig := UserNotificationMethods{
		{Provider: "telegram", Target: "123"},
		{Provider: "email", Target: "test@example.com"},
	}

	data, err := orig.MarshalBinary()
	require.NoError(t, err)
	require.JSONEq(t, `[{"provider":"telegram","target":"123"},{"provider":"email","target":"test@example.com"}]`, string(data))
}

func TestUserNotificationMethods_UnmarshalText(t *testing.T) {
	data := `[{"provider":"telegram","target":"123"},{"provider":"email","target":"test@example.com"}]`

	var res UserNotificationMethods
	err := res.UnmarshalText([]byte(data))
	require.NoError(t, err)
	require.Equal(
		t,
		UserNotificationMethods{
			{Provider: "telegram", Target: "123"},
			{Provider: "email", Target: "test@example.com"},
		},
		res,
	)
}
