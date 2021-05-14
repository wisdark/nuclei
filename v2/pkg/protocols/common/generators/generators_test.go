package generators

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSniperGenerator(t *testing.T) {
	usernames := []string{"admin", "password"}
	moreUsernames := []string{"login", "test"}

	generator, err := New(map[string]interface{}{"username": usernames, "aliases": moreUsernames}, Sniper, "")
	require.Nil(t, err, "could not create generator")

	iterator := generator.NewIterator()
	count := 0
	for {
		_, ok := iterator.Value()
		if !ok {
			break
		}
		count++
	}
	require.Equal(t, len(usernames)+len(moreUsernames), count, "could not get correct sniper counts")
}

func TestPitchforkGenerator(t *testing.T) {
	usernames := []string{"admin", "token"}
	passwords := []string{"admin", "password"}

	generator, err := New(map[string]interface{}{"username": usernames, "password": passwords}, PitchFork, "")
	require.Nil(t, err, "could not create generator")

	iterator := generator.NewIterator()
	count := 0
	for {
		value, ok := iterator.Value()
		if !ok {
			break
		}
		count++
		require.Contains(t, usernames, value["username"], "Could not get correct pitchfork username")
		require.Contains(t, passwords, value["password"], "Could not get correct pitchfork password")
	}
	require.Equal(t, len(passwords), count, "could not get correct pitchfork counts")
}

func TestClusterbombGenerator(t *testing.T) {
	usernames := []string{"admin"}
	passwords := []string{"admin", "password", "token"}

	generator, err := New(map[string]interface{}{"username": usernames, "password": passwords}, ClusterBomb, "")
	require.Nil(t, err, "could not create generator")

	iterator := generator.NewIterator()
	count := 0
	for {
		value, ok := iterator.Value()
		if !ok {
			break
		}
		count++
		require.Contains(t, usernames, value["username"], "Could not get correct clusterbomb username")
		require.Contains(t, passwords, value["password"], "Could not get correct clusterbomb password")
	}
	require.Equal(t, 3, count, "could not get correct clusterbomb counts")

	iterator.Reset()
	count = 0
	for {
		value, ok := iterator.Value()
		if !ok {
			break
		}
		count++
		require.Contains(t, usernames, value["username"], "Could not get correct clusterbomb username")
		require.Contains(t, passwords, value["password"], "Could not get correct clusterbomb password")
	}
	require.Equal(t, 3, count, "could not get correct clusterbomb counts")
}
