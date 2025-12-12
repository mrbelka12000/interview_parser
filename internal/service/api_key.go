package service

func (s *Service) GetAPIKey() (string, error) {
	return s.apiKeyRepo.GetOpenAIAPIKeyFromDB()
}

func (s *Service) InsertAPIKey(apiKey string) error {
	return s.apiKeyRepo.InsertOpenAIAPIKey(apiKey)
}

func (s *Service) DeleteAPIKey() error {
	return s.apiKeyRepo.DeleteOpenAIAPIKey()
}
