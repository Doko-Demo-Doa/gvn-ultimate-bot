package slack

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/shomali11/slacker"
)

func Bootstrap() {
	slack_bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	definition := slacker.CommandDefinition{
		Description: "ping",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			response.Reply("custom")
			response.ReportError(errors.New("oops"))
		},
	}

	slack_bot.Command("ping", &definition)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := slack_bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
