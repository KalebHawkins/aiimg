package cmd

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testApiKey = "tisbutatestkey"
)

func TestRootCmd(t *testing.T) {
	testcases := []struct {
		name      string
		apiKey    string
		useConfig bool
		useEnv    bool
		expErr    error
	}{
		{
			name:      "ApiErr",
			apiKey:    "",
			useConfig: false,
			useEnv:    false,
			expErr:    errNoApiKey,
		},
		{
			name:      "NoErrUseConfigFile",
			apiKey:    testApiKey,
			useConfig: true,
			useEnv:    false,
			expErr:    nil,
		},
		{
			name:      "NoErrUseEnvVar",
			apiKey:    testApiKey,
			useConfig: false,
			useEnv:    true,
			expErr:    nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.useConfig {
				rootCmd.SetArgs([]string{"--config", "../testdata/config.yaml"})
			}
			if tc.useEnv {
				os.Setenv("AIIMG_API_KEY", testApiKey)
			}

			err := rootCmd.Execute()
			if tc.expErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expErr, err)
			}

			if tc.expErr == nil {
				assert.NoError(t, err)
				assert.Equal(t, apiKey, testApiKey)
			}
		})
	}

}
