package sources

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestItParsesAnEnv(t *testing.T) {
	env := `FOO=FOO
	BAR="BAR"
	BAZ=1
	`
	secrets, err := NewSecretsFromEnv(
		strings.NewReader(env),
	)

	require.NoError(t, err)
	assert.Equal(t, "FOO", secrets["FOO"])
	assert.Equal(t, "BAR", secrets["BAR"])
	assert.Equal(t, "1", secrets["BAZ"])
}

func TestItWritesAnEnv(t *testing.T) {
	secrets := secretsMap{
		"BAR": "BAR",
		"FOO": "FOO",
		"BAZ": "1",
	}

	// The output should be sorted by key
	expected := `BAR="BAR"
BAZ="1"
FOO="FOO"
`

	var b strings.Builder
	err := secrets.ToEnv(&b)

	require.NoError(t, err)
	assert.Equal(t, expected, b.String())
}
