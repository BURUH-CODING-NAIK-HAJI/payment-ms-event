package eventrepository

type EventRepositoryInterface interface {
	Welcome() string
}

type EventRepository struct {
}

func New() EventRepositoryInterface {
	return &EventRepository{}
}

func (eventrepository *EventRepository) Welcome() string {
	return "Hai From Repository"
}
