package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"path"
	"time"

	"github.com/arylatt/go-cbf/api"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	DefaultFileTimeFormat = "2006-01-02T15-04-05"

	FlagCategory   = "category"
	FlagTimeFormat = "time-format"
	FlagOutputDir  = "output-dir"
	FlagLogLevel   = "log-level"
	FlagInterval   = "interval"

	CfgRunTime = "run-time"
)

var (
	trackCBFCmd = &cobra.Command{
		Use:          fmt.Sprintf("track-cbf event [-c %s]... [-t %s] [-o %s] [-l %s] [-i %s]", FlagCategory, FlagTimeFormat, FlagOutputDir, FlagLogLevel, FlagInterval),
		Short:        "Get the JSON data for the provided category and Cambridge Beer Festival event.",
		Example:      "track-cbf cbf2023",
		Args:         cobra.ExactArgs(1),
		PreRunE:      trackCBFPreRunE,
		RunE:         trackCBFRunE,
		SilenceUsage: true,
	}

	fs = afero.NewOsFs()
)

func init() {
	trackCBFCmd.PersistentFlags().StringArrayP(FlagCategory, "c", api.KnownCategories, "Search for specific categories.")
	trackCBFCmd.PersistentFlags().StringP(FlagTimeFormat, "t", DefaultFileTimeFormat, "Date/time format to use for file name prefixes.")
	trackCBFCmd.PersistentFlags().StringP(FlagOutputDir, "o", "", "Directory to output files to.")
	trackCBFCmd.PersistentFlags().StringP(FlagLogLevel, "l", zerolog.LevelInfoValue, "Log level.")
	trackCBFCmd.PersistentFlags().DurationP(FlagInterval, "i", time.Duration(0), "How often to poll data.")

	viper.BindPFlags(trackCBFCmd.PersistentFlags())
}

func filename(timestamp time.Time, category string) (filename string) {
	return fmt.Sprintf("%s_%s.json", timestamp.Format(viper.GetString(FlagTimeFormat)), category)
}

func write(data *api.Response, category string) (err error) {
	dataTime, err := data.Updated()
	if err != nil {
		dataTime = viper.GetTime(CfgRunTime)

		log.Warn().Err(err).Msgf("Failed to get time from data, using command invocation time: %s.", dataTime.Format(viper.GetString(FlagTimeFormat)))

		err = nil
	}

	outputFile := path.Join(viper.GetString(FlagOutputDir), filename(dataTime, category))
	log := log.With().Str("output-file", outputFile).Logger()

	if _, err = fs.Stat(outputFile); err == nil {
		log.Debug().Msg("Not overwriting existing file.")
		return nil
	}

	file, err := fs.Create(outputFile)
	if err != nil {
		log.Err(err).Msg("Failed to create/truncate file.")
		return
	}

	defer file.Close()

	err = json.NewEncoder(file).Encode(data)
	if err != nil {
		log.Err(err).Msg("Failed to write data to file.")
	} else {
		log.Info().Msgf("Fetched data!")
	}

	return
}

func loop(wc api.Client, event string) {
	evLog := log.With().Str("event", event).Logger()

	for _, category := range viper.GetStringSlice(FlagCategory) {
		log := evLog.With().Str("category", category).Logger()

		resp, err := wc.Get(event, category)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to get data from API.")
			continue
		}

		write(resp, category)
	}
}

func trackCBFRunE(cmd *cobra.Command, args []string) (err error) {
	wc, err := api.NewWebClient()
	if err != nil {
		log.Error().Msgf("Failed to create web client: %s.", err.Error())
		return
	}

	ctx, cancel := context.WithCancel(cmd.Context())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	defer signal.Stop(sigCh)

	go func() {
		<-sigCh
		log.Info().Msg("Interrupt received, stopping...")
		cancel()
	}()

	timedLoop(ctx, func() { loop(wc, args[0]) })

	return
}

func timedLoop(ctx context.Context, looper func()) {
	loopInterval := viper.GetDuration(FlagInterval)

	looper()

	if loopInterval == time.Duration(0) {
		return
	}

	for {
		log.Info().Msgf("Waiting for %s", loopInterval)

		select {
		case <-time.After(loopInterval):
			looper()
		case <-ctx.Done():
			return
		}
	}
}

func trackCBFPreRunE(cmd *cobra.Command, args []string) (err error) {
	log.Logger = log.Output(zerolog.NewConsoleWriter())

	logLvlStr := viper.GetString(FlagLogLevel)
	logLvl, err := zerolog.ParseLevel(logLvlStr)
	if err != nil {
		log.Warn().Err(err).Str(FlagLogLevel, logLvlStr).Msg("Failed to parse log level, defaulting to INFO.")

		logLvl = zerolog.InfoLevel
		err = nil
	}

	log.Logger = log.Level(logLvl)

	viper.Set(CfgRunTime, time.Now())

	return
}

func main() {
	if err := trackCBFCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
