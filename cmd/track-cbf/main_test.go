package main

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/arylatt/go-cbf/api"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockClient struct {
	mock.Mock
}

func (m *MockClient) Get(ev, cat string) (*api.Response, error) {
	args := m.Called(ev, cat)
	return args.Get(0).(*api.Response), args.Error(1)
}

func resetVars() {
	viper.Reset()

	viper.Set(FlagCategory, "category")
	viper.Set(FlagTimeFormat, DefaultFileTimeFormat)
	viper.Set(FlagOutputDir, "/output-dir/")
	viper.Set(FlagLogLevel, "disabled")
	viper.Set(CfgRunTime, time.Now())
}

func TestMain(m *testing.M) {
	fs = afero.NewMemMapFs()
	log.Logger = log.Level(zerolog.Disabled)

	resetVars()

	os.Exit(m.Run())
}

func TestFilename(t *testing.T) {
	tests := []struct {
		time     time.Time
		category string
		expected string
	}{
		{
			time.Unix(1684262624, 0),
			"test-cat",
			"2023-05-16T19-43-44_test-cat.json",
		},
		{
			time.Unix(1684262674, 0),
			"cat-test",
			"2023-05-16T19-44-34_cat-test.json",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, filename(test.time, test.category))
	}
}

func TestWrite(t *testing.T) {
	tests := []struct {
		data   *api.Response
		assert func(t *testing.T, err error)
	}{
		{
			&api.Response{
				Timestamp: "Tue May 16 2023 20:16:35 BST",
			},
			func(t *testing.T, err error) {
				if assert.NoError(t, err) {
					fExists, err := afero.Exists(fs, "/output-dir/2023-05-16T20-16-35_category.json")

					assert.NoError(t, err)
					assert.True(t, fExists)
				}
			},
		},
		{
			&api.Response{
				Timestamp: "Tue May 16 2023 20:16:35 BST",
			},
			func(t *testing.T, err error) {
				if assert.NoError(t, err) {
					fExists, err := afero.Exists(fs, "/output-dir/2023-05-16T20-16-35_category.json")

					assert.NoError(t, err)
					assert.True(t, fExists)
				}
			},
		},
		{
			&api.Response{
				Timestamp: "",
			},
			func(t *testing.T, err error) {
				if assert.NoError(t, err) {
					fExists, err := afero.Exists(fs, fmt.Sprintf("/output-dir/%s_category.json", viper.GetTime(CfgRunTime).Format(viper.GetString(FlagTimeFormat))))

					assert.NoError(t, err)
					assert.True(t, fExists)
				}
			},
		},
	}

	for _, test := range tests {
		test.assert(t, write(test.data, "category"))
	}
}

func TestLoop(t *testing.T) {
	defer resetVars()

	viper.Set(FlagCategory, []string{"cat1", "cat2"})

	wc := &MockClient{}

	wc.On("Get", mock.Anything, "cat1").Return(&api.Response{Timestamp: "Mon May 16 2023 20:40:05 BST"}, nil)
	wc.On("Get", mock.Anything, mock.Anything).Return(&api.Response{}, api.ErrEventOrCategoryNotFound)

	loop(wc, "event")

	fExists, err := afero.Exists(fs, "/output-dir/2023-05-16T20-40-05_cat1.json")

	assert.NoError(t, err)
	assert.True(t, fExists)
}

func TestTrackCBFRunE(t *testing.T) {
	defer resetVars()

	viper.Set(FlagCategory, []string{})

	assert.NoError(t, trackCBFRunE(nil, []string{"event"}))
}

func TestTrackCBFPreRunE(t *testing.T) {
	defer resetVars()

	tests := []struct {
		logLevel string
		expected zerolog.Level
	}{
		{
			"fatal",
			zerolog.FatalLevel,
		},
		{
			"not-a-real-level",
			zerolog.InfoLevel,
		},
	}

	for _, test := range tests {
		viper.Set(CfgRunTime, nil)
		viper.Set(FlagLogLevel, test.logLevel)

		if assert.NoError(t, trackCBFPreRunE(nil, []string{})) {
			assert.IsType(t, time.Time{}, viper.Get(CfgRunTime))
			assert.Equal(t, log.Logger.GetLevel(), test.expected)
		}
	}
}
