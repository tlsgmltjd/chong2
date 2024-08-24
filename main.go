package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
)

func main() {
	token := ""

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("discord 세션 생성 중 에러가 발생했어요: ", err)
		return
	}

	err = dg.Open()
	if err != nil {
		fmt.Println("discord 봇과 연결 중 에러가 발생했어요: ", err)
		return
	}
	defer dg.Close()

	c := cron.New()
	_, err = c.AddFunc("0 8 * * *", func() {
		channelID := ""
		_, err := dg.ChannelMessageSend(channelID, "좋은 아침이에요, 오늘 할 일을 공유해주세요!")
		if err != nil {
			fmt.Println("메시지 전송 중 에러가 발생했어요: ", err)
		}
	})

	if err != nil {
		log.Fatal("cron 스케쥴 등록중 에러가 발생했어요: ", err)
	}

	c.Start()

	fmt.Println("봇 활성화 중입니다.")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop

	fmt.Println("봇이 꺼졌습니다.")
	c.Stop()
}
