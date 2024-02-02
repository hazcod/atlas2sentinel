package atlas

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/atlas-sdk/v20231115005/admin"
)

const (
	maxFetch = 100
)

type Atlas struct {
	Logger *logrus.Logger
	client *admin.APIClient
}

func New(l *logrus.Logger, url, pubKey, privKey string) (*Atlas, error) {
	if l == nil {
		return nil, errors.New("logger is nil")
	}

	atlasOpts := []admin.ClientModifier{
		admin.UseDigestAuth(pubKey, privKey),
		admin.UseBaseURL(url),
		admin.UseDebug(l.IsLevelEnabled(logrus.DebugLevel)),
	}

	sdk, err := admin.NewClient(atlasOpts...)
	if err != nil {
		return nil, fmt.Errorf("invalid pub/priv key: %v", err)
	}

	atlas := Atlas{
		Logger: l,
		client: sdk,
	}

	return &atlas, nil
}
