package bot

import (
	"fmt"
	"log"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/xanzy/go-gitlab"
)

const (
	maxEvents     = 100
	updateTimeout = 60
)

type Bot struct {
	wg           sync.WaitGroup
	api          *tgbotapi.BotAPI
	eventChannel chan interface{}
	chatID       int64
}

func New(token string, chatID int64) (*Bot, error) {
	var err error
	bot := &Bot{
		eventChannel: make(chan interface{}, maxEvents),
		chatID:       chatID,
	}

	bot.api, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("connect to Telegram API failed: %w", err)
	}

	bot.wg.Add(1)
	go func() {
		defer bot.wg.Done()
		bot.loop()
	}()

	return bot, nil
}

func (bot *Bot) loop() {

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = updateTimeout

	updatesChannel := bot.api.GetUpdatesChan(updateConfig)

	for {
		select {
		case <-updatesChannel:
			// someday implement commands
		case event := <-bot.eventChannel:
			switch concreteEvent := event.(type) {
			case *gitlab.PipelineEvent:
				log.Println("[bot] got pipeline event, notify")
				if concreteEvent.ObjectAttributes.Status == "failed" {
					bot.sendReport(concreteEvent)
				}
			default:
				log.Printf("[bot] got unknown event: %T", concreteEvent)
			}
		}
	}

}

func (bot *Bot) sendReport(event *gitlab.PipelineEvent) {
	pipelineUrl := fmt.Sprintf("%s/-/pipelines/%d", event.Project.WebURL, event.ObjectAttributes.ID)
	report := fmt.Sprintf("<b>Build failed</b>\n<a href=\"%s\">#%d</a>\n\nProject: %s\nBranch: %s\nCommit: <code>%s</code>\nAuthor: %s\n",
		pipelineUrl, event.ObjectAttributes.ID, event.Project.PathWithNamespace, event.ObjectAttributes.Ref, event.Commit.ID, event.User.Name)
	for _, build := range event.Builds {
		if build.ArtifactsFile.Filename != "" {
			downloadUrl := fmt.Sprintf("%s/-/jobs/%d/artifacts/download?file_type=archive", event.Project.WebURL, build.ID)
			report += fmt.Sprintf("\n<a href=\"%s\">Download %s</a>", downloadUrl, build.ArtifactsFile.Filename)
		}
	}
	msg := tgbotapi.NewMessage(bot.chatID, report)
	msg.ParseMode = "HTML"
	bot.api.Send(msg)
}

func (bot *Bot) HandleEvent(event interface{}) {
	bot.eventChannel <- event
}
