package main

import (
	"atlas2sentinel/config"
	"atlas2sentinel/pkg/atlas"
	msSentinel "atlas2sentinel/pkg/sentinel"
	"context"
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()

	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	confFile := flag.String("config", "config.yml", "The YAML configuration file.")
	flag.Parse()

	conf := config.Config{}
	if err := conf.Load(*confFile); err != nil {
		logger.WithError(err).WithField("config", *confFile).Fatal("failed to load configuration")
	}

	if err := conf.Validate(); err != nil {
		logger.WithError(err).WithField("config", *confFile).Fatal("invalid configuration")
	}

	logrusLevel, err := logrus.ParseLevel(conf.Log.Level)
	if err != nil {
		logger.WithError(err).Error("invalid log level provided")
		logrusLevel = logrus.InfoLevel
	}
	logger.SetLevel(logrusLevel)

	//

	atlasClient, err := atlas.New(logger, conf.Atlas.URL, conf.Atlas.PublicKey, conf.Atlas.PrivateKey)
	if err != nil {
		logger.WithError(err).Fatal("could not create atlas client")
	}

	sentinel, err := msSentinel.New(logger, msSentinel.Credentials{
		TenantID:       conf.Microsoft.TenantID,
		ClientID:       conf.Microsoft.AppID,
		ClientSecret:   conf.Microsoft.SecretKey,
		SubscriptionID: conf.Microsoft.SubscriptionID,
		ResourceGroup:  conf.Microsoft.ResourceGroup,
		WorkspaceName:  conf.Microsoft.WorkspaceName,
	})
	if err != nil {
		logger.WithError(err).Fatal("could not create MS Sentinel client")
	}

	//

	if conf.Microsoft.UpdateTable {
		if err := sentinel.CreateTable(ctx, logger, conf.Microsoft.RetentionDays); err != nil {
			logger.WithError(err).Fatal("failed to create MS Sentinel table")
		}
	}

	//

	logEvents, err := atlasClient.GetLogs(ctx, conf.Atlas.LookBackDays)
	if err != nil {
		logger.WithError(err).Fatal("could not fetch atlas host logs")
	}

	// TODO remove debug
	for _, log := range logEvents {
		fmt.Println(string(log))
	}

	logs, err := atlas.ConvertLogtToMap(logger, logEvents)
	if err != nil {
		logger.WithError(err).Errorf("could not parse audit events")
	}

	//

	logger.WithField("total", len(logs)).Info("collected all Atlas logs")

	//

	if err := sentinel.SendLogs(ctx, logger,
		conf.Microsoft.DataCollection.Endpoint,
		conf.Microsoft.DataCollection.RuleID,
		conf.Microsoft.DataCollection.StreamName,
		logs); err != nil {
		logger.WithError(err).Fatal("could not ship logs to sentinel")
	}

	//

	logger.WithField("total", len(logs)).Info("successfully sent logs to sentinel")
}
