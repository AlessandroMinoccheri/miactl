package sdk

import (
	"fmt"
	"testing"

	"github.com/davidebianchi/go-jsonclient"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("throws with empty options", func(t *testing.T) {
		tests := []struct {
			option Options
		}{
			{option: Options{}},
			{option: Options{Secret: "sid=asd", APIBaseURL: "base"}},
			{option: Options{APIBaseURL: "base", APICookie: "cookie"}},
			{option: Options{APICookie: "cookie", Secret: "sid=asd"}},
		}
		for _, test := range tests {
			client, err := New(test.option)
			require.EqualError(t, err, fmt.Sprintf("%s: client options are not correct", ErrCreateClient))
			require.Nil(t, client)
		}
	})

	t.Run("throws with wrong base url", func(t *testing.T) {
		client, err := New(Options{
			APIBaseURL: "wrong	",
			Secret:    "secret",
			APICookie: "sid=asd",
		})
		require.Error(t, err)
		require.Nil(t, client)
	})

	t.Run("correctly returns mia client", func(t *testing.T) {
		opts := Options{
			APIBaseURL: "http://my-url/path",
			Secret:     "my secret",
			APICookie:  "sid=asd",
		}
		client, err := New(opts)

		expectedJSONClient, _ := jsonclient.New(jsonclient.Options{
			BaseURL: opts.APIBaseURL,
			Headers: map[string]string{
				"secret": opts.Secret,
				"cookie": opts.APICookie,
			},
		})

		require.NoError(t, err, "new client error")
		require.Exactly(t, &MiaClient{
			Projects: &ProjectsClient{
				JSONClient: expectedJSONClient,
			},
		}, client)
	})
}