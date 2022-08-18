package utils

type IEventConsumer interface {
	HandleEvent(event interface{})
}
