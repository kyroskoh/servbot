package yandexOAuth

import (
	"github.com/khades/servbot/config"
	"github.com/khades/servbot/donationSource"
	"github.com/khades/servbot/httpAPI"
	"goji.io/pat"
)

func Init(httpAPIService *httpAPI.Service, config *config.Config,
	donationSourceService *donationSource.Service) {
	service := Service{config, donationSourceService}
	mux := httpAPIService.NewMux()
	mux.HandleFunc(pat.Post("/yandex/oauth"), httpAPIService.WithMod(service.login))
}
