package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/xanzy/go-gitlab"
)

type WebHookHandler struct {
	Consumer IEventConsumer
}

func (h WebHookHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	logInfo := fmt.Sprintf("[%s@%s]", request.RemoteAddr, request.RequestURI)
	log.Printf("%s webhook called", logInfo)
	defer request.Body.Close()

	payload, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Printf("%s read request failed: %s", logInfo, err)
		return
	}

	if request.Method != http.MethodPost {
		log.Printf("%s invalid method: %s", logInfo, request.Method)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	eventCode := request.Header.Get("X-Gitlab-Event")
	if strings.TrimSpace(eventCode) == "" {
		log.Printf("%s request is not a GitLab Event", logInfo)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	eventType := gitlab.EventType(eventCode)

	event, err := gitlab.ParseWebhook(eventType, payload)
	if err != nil {
		log.Printf("%s cannot parse event: %+s", logInfo, err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("%s got event: %+v", logInfo, event)

	if h.Consumer != nil {
		h.Consumer.HandleEvent(event)
	}
}
