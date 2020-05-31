package transform

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegexp(t *testing.T) {
	t.Run("Test building transformer from JSON", func(t *testing.T) {
		rawJson := []byte(`
{
	"type":"regexp",
	"expression":"\\d{1}(\\d{8})\\d+",
	"submatch_number":1
}
`)

		fakeAction := TransformerFunc(func(body []byte) ([]byte, error) {
			require.Equal(t, []byte("23456789"), body)
			return body, nil
		})
		transformation, err := transformerFromJSON(rawJson, fakeAction)
		require.NoError(t, err)

		body := []byte("number 1234567890")

		body, err = transformation.Transform(body)
	})

	t.Run("Test one submatch", func(t *testing.T) {
		tr := &Regexp{
			Expression:     `number: \d(\d+)\d{4}`,
			SubmatchNumber: 1,
			Action: TransformerFunc(func(body []byte) ([]byte, error) {
				require.Equal(t, []byte("23456"), body)
				return body, nil
			}),
		}

		body := []byte("number: 1234567890")
		_, err := tr.Transform(body)
		require.NoError(t, err)
	})

	t.Run("Test multiple submatch", func(t *testing.T) {
		tr := &Regexp{
			Expression:     `number: (\d+)(\d{4})`,
			SubmatchNumber: 2,
			Action: TransformerFunc(func(body []byte) ([]byte, error) {
				require.Equal(t, []byte("7890"), body)
				return body, nil
			}),
		}

		body := []byte("number: 1234567890")
		_, err := tr.Transform(body)
		require.NoError(t, err)
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

	t.Run("Test transformation of multiple matches", func(t *testing.T) {
		tr := &Regexp{
			Expression:     `number: (\d+)(\d{4})`,
			SubmatchNumber: 2,
			Action: TransformerFunc(func(body []byte) ([]byte, error) {
				newBody := bytes.Repeat([]byte("x"), len(body))
				return newBody, nil
			}),
		}

		body := []byte("number: 12345 whatever number: 54321")
		want := []byte("number: 1xxxx whatever number: 5xxxx")
		got, err := tr.Transform(body)
		require.NoError(t, err)
		require.Equal(t, string(want), string(got))
	})
}
