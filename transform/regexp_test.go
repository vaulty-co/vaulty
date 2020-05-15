package transform

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegexp(t *testing.T) {
	t.Run("Test one submatch", func(t *testing.T) {
		tr := &Regexp{
			Expression:     `number: \d(\d+)\d{4}`,
			SubmatchNumber: 1,
			Action: TransformerFunc(func(body []byte) ([]byte, error) {
				return []byte("xxxx"), nil
			}),
		}

		body := []byte("number: 4242424242424242")
		newBody, err := tr.Transform(body)
		require.NoError(t, err)
		require.Contains(t, string(newBody), "number: 4xxxx4242")
	})

	t.Run("Test multiple submatch", func(t *testing.T) {
		tr := &Regexp{
			Expression:     `number: (\d+)(\d{4})`,
			SubmatchNumber: 2,
			Action: TransformerFunc(func(body []byte) ([]byte, error) {
				return []byte("xxxx"), nil
			}),
		}

		body := []byte("number: 4242424242424242")
		newBody, err := tr.Transform(body)
		require.NoError(t, err)
		require.Contains(t, string(newBody), "number: 424242424242xxxx")
	})

	t.Run("Test no submatch", func(t *testing.T) {
		tr := &Regexp{
			Expression:     `number: (\d+)(\d{4})`,
			SubmatchNumber: 5,
			Action: TransformerFunc(func(body []byte) ([]byte, error) {
				return []byte("xxxx"), nil
			}),
		}

		body := []byte("number: 4242424242424242")
		newBody, err := tr.Transform(body)
		require.NoError(t, err)
		require.Contains(t, string(newBody), "number: 4242424242424242")

		tr2 := &Regexp{
			Expression:     `number: (\d+)(\d{4})`,
			SubmatchNumber: -1,
			Action: TransformerFunc(func(body []byte) ([]byte, error) {
				return []byte("xxxx"), nil
			}),
		}

		body = []byte("hello")
		newBody, err = tr2.Transform(body)
		require.NoError(t, err)
		require.Contains(t, string(newBody), "hello")
	})
}
