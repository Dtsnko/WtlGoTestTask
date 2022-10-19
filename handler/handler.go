package handler

type Handler struct {
	RecordHand RecordHandler
	TaskHand   TaskHandler
}

func New() *Handler {

	return &Handler{RecordHand: RecordHandler{}, TaskHand: TaskHandler{}}
}
