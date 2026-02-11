package ServiceErrors

type ServiceError struct {
	EndUserMessage   string `json:"errorMessage"`
	Err              error
	TechnicalMessage string `json:"technicalMessage"`
}

func (s *ServiceError) Error() string {
	return s.EndUserMessage
}
