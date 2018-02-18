package healthcheckers

var (
	_ Checker = new(CheckService)
)

type Checker interface {
	Message() string
	Check() error
}

type CheckService struct {
	ServiceName string
	Action func() error
}

func(s *CheckService) Message() string {
	return s.ServiceName
}

func(s *CheckService) Check() error {
	return s.Action()
}