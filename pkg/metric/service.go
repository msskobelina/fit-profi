package metric

type Service struct {
	userCreated CounterVec
	loginFailed CounterVec
}

func NewService(
	metric Metric,
	userCreated MetricWithLabels,
	loginFailed MetricWithLabels,
) *Service {
	return &Service{
		userCreated: metric.CounterVec(userCreated),
		loginFailed: metric.CounterVec(loginFailed),
	}
}

func (s *Service) TrackUserCreated(source string) {
	s.userCreated.IncWithNamedValues(map[string]string{
		"source": source,
	})
}

func (s *Service) TrackLoginFailed(reason string) {
	s.loginFailed.IncWithNamedValues(map[string]string{
		"reason": reason,
	})
}
