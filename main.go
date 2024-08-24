package main

import (
	"github.com/joho/godotenv"

	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
)

func loadEnv() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("env 파일을 찾지 못했어요")
    }
}

func daily_notic(dg *discordgo.Session) *cron.Cron {
	c := cron.New()
	_, err := c.AddFunc("0 8 * * *", func() {
		channelID := os.Getenv("CHANNEL_ID")
		_, err := dg.ChannelMessageSend(channelID, "좋은 아침이에요, 오늘 할 일을 공유해주세요!")
		if err != nil {
			fmt.Println("메시지 전송 중 에러가 발생했어요: ", err)
		}
	})

	if err != nil {
		log.Fatal("cron 스케줄 등록 중 에러가 발생했어요: ", err)
	}

	return c
}

func main() {

	loadEnv()

	token := os.Getenv("DISCORD_BOT_TOKEN")

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

	dailyCron := daily_notic(dg)
	dailyCron.Start()

	fmt.Println("봇 활성화 중입니다.")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop

	fmt.Println("봇이 꺼졌습니다.")
	dailyCron.Stop()
}
