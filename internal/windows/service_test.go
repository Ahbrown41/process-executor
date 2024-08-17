package windows

import (
	"context"
	"github.com/stretchr/testify/suite"
	"testing"
	"tool-commerce-store-data-sender/internal/config"
)

type WindowsServiceSuite struct {
	storeConfig config.Store
	appConfig   config.Config
	ctx         context.Context
	suite.Suite
}

func TestWindowsServiceSuite(t *testing.T) {
	suite.Run(t, new(WindowsServiceSuite))
}

func (suite *WindowsServiceSuite) SetupSuite() {
	suite.ctx = context.Background()
	storeConfig := config.Store{
		StoreNum:               "a",
		StoreProductServiceURL: "b",
		StoreDataCache:         "c",
		Interval:               "5s",
	}
	suite.storeConfig = storeConfig
	appConfig := config.Config{}
	appConfig.Store = storeConfig
	appConfig.Oauth.ClientID = "CLIENTID"
	appConfig.Oauth.ClientSecret = "SECRET"
	appConfig.Oauth.Audience = "AUDIENCE"
	suite.appConfig = appConfig
}

func (suite *WindowsServiceSuite) TestNew() {
	svc := New(suite.storeConfig, suite.ctx)
	suite.Assert().NotNil(svc)
}
