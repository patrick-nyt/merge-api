package main

type service struct {
	config Config
	remote Fetcher
	local  Fetcher
}

func NewService(config Config) *service {
	remote := NewRemoteClient(config.RemoteSource)
	local := NewLocalClient(config.LocalSource)
	return &service{config, remote, local}
}
