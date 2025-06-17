package main


import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)
func TestServeHTTP(t *testing.T) {
	for name, test := range map[string]struct {
		RequestURL         string
		ExpectedStatusCode int
		ExpectedHeader     http.Header
		ExpectedbodyString string
	}{
		"Request Custom Configuration Settings": {
			RequestURL:         "/custom_config_settings",
			ExpectedStatusCode: http.StatusOK,
			ExpectedHeader:     http.Header{"Content-Type": []string{"application/json"}},
			ExpectedbodyString: `{"QuestionPort":"5000","QuestionServerAddress":"127.0.0.1"}`,
		},
	} {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			plugin := &Plugin{}
			plugin.initializeAPI()
			cfg := plugin.getConfiguration()
			cfg.QuestionServerAddress = "127.0.0.1"
			cfg.QuestionPort = "5000"

			// This is required in order to get the expected result
			plugin.setConfiguration(cfg)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", test.RequestURL, nil)
			plugin.ServeHTTP(nil, w, r)

			result := w.Result()
			require.NotNil(t, result)
			defer result.Body.Close()

			bodyBytes, err := io.ReadAll(result.Body)
			require.Nil(t, err)
			bodyString := string(bodyBytes)

			assert.Equal(test.ExpectedbodyString, bodyString)
			assert.Equal(test.ExpectedStatusCode, result.StatusCode)
			assert.Equal(test.ExpectedHeader, result.Header)
		})
	}
}
