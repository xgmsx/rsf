package utils

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
)

func TestToPtr(t *testing.T) {
	var str string
	require.Equal(t, &str, ToPtr(str))

	str = gofakeit.Sentence(5)
	require.Equal(t, &str, ToPtr(str))
}
