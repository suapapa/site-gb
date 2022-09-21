package main

import (
	"fmt"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/suapapa/site-gb/msg"
)

var (
	tgBot       *tgbotapi.BotAPI
	tgAPIToke   = os.Getenv("TELEGRAM_APITOKEN")
	tgRoomIDStr = os.Getenv("TELEGRAM_ROOM_ID")
)

func sendMsgToTelegram(m *msg.Message) error {
	var err error
	if tgBot == nil {
		tgBot, err = tgbotapi.NewBotAPI(tgAPIToke)
		if err != nil {
			return errors.Wrap(err, "fail to send msg to telegram")
		}
		// tgBot.Debug = true
	}

	tgRoomID, err := strconv.Atoi(tgRoomIDStr)
	if err != nil {
		return errors.Wrap(err, "fail to send msg to telegram")
	}

	mStr := makeGBStr4Telegram(m)
	c := tgbotapi.NewMessage(int64(tgRoomID), mStr)
	// c.ParseMode = tgbotapi.ModeMarkdown // NOT WORKING :(
	if _, err := tgBot.Send(c); err != nil {
		return errors.Wrap(err, "fail to send msg to telegram")
	}
	return nil
}

func makeGBStr4Telegram(m *msg.Message) string {
	gb, ok := m.Data.(msg.GuestBook)
	if !ok {
		return ""
	}
	outFmt := `## %s ##
%s
- %s -`
	out := fmt.Sprintf(outFmt,
		m.Type,
		gb.Content,
		gb.From,
	)
	return out
}
