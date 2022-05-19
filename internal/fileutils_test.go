package internal

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileutilsReadFileOrStdin(t *testing.T) {
	t.Parallel()

	t.Run("should read from stdin if file equals dash", func(t *testing.T) {
		t.Parallel()

		tmpfile, err := ioutil.TempFile(t.TempDir(), strings.ReplaceAll(t.Name(), "/", "_"))
		defer os.Remove(tmpfile.Name())
		assert.Nil(t, err)

		osStdin = tmpfile

		content := []byte("foobar")
		_, err = tmpfile.Write(content)
		assert.Nil(t, err)
		_, err = tmpfile.Seek(0, 0)
		assert.Nil(t, err)

		file := "-"
		actual, err := ReadFileOrStdin(file)
		assert.Nil(t, err)
		assert.Equal(t, "foobar", actual)
	})

	t.Run("should read from file otherwise", func(t *testing.T) {
		t.Parallel()

		content := "foobar"

		tmpfile, err := ioutil.TempFile(t.TempDir(), strings.ReplaceAll(t.Name(), "/", "_"))
		defer os.Remove(tmpfile.Name())
		assert.Nil(t, err)
		_, err = tmpfile.Write([]byte(content))
		assert.Nil(t, err)

		assert.NotEqual(t, "-", tmpfile.Name())
		actual, err := ReadFileOrStdin(tmpfile.Name())
		assert.Nil(t, err)
		assert.Equal(t, "foobar", actual)
	})

}
