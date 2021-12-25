package auth

import (
	"fmt"
	"github.com/detecc/detecctor-v2/pkg/cache"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGenerateChatAuthenticationToken(t *testing.T) {
	require := require.New(t)
	chatId := "exampleChatId"

	GenerateChatAuthenticationToken(cache.Memory(), chatId)

	token, isFound := cache.Memory().Get(fmt.Sprintf("auth-token-%s", chatId))
	require.True(isFound)
	require.NotEmpty(token)
}
