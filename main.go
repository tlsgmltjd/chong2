package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("env 파일을 찾지 못했어요")
	}
}

func addDailyMorningReminder(c *cron.Cron, dg *discordgo.Session) {
	_, err := c.AddFunc("0 8 * * *", func() {
		channelID := os.Getenv("DAILY_CHANNEL_ID")
		_, err := dg.ChannelMessageSend(channelID, "좋은 아침이에요, 오늘 할 일을 공유해주세요!")
		if err != nil {
			fmt.Println("메시지 전송 중 에러가 발생했어요: ", err)
		}
	})
	if err != nil {
		log.Fatal("cron 스케줄 등록 중 에러가 발생했어요: ", err)
	}
}

func addDailyEveningReminder(c *cron.Cron, dg *discordgo.Session) {
	_, err := c.AddFunc("0 23 * * *", func() {
		channelID := os.Getenv("DAILY_CHANNEL_ID")
		_, err := dg.ChannelMessageSend(channelID, "오늘 하루도 수고했어요, 푹 쉬고 내일도 정진합시다!")
		if err != nil {
			fmt.Println("메시지 전송 중 에러가 발생했어요: ", err)
		}
	})
	if err != nil {
		log.Fatal("cron 스케줄 등록 중 에러가 발생했어요: ", err)
	}
}

func addMondayStudyReminder(c *cron.Cron, dg *discordgo.Session) {
	_, err := c.AddFunc("0 13 * * 1", func() {
		channelID := os.Getenv("NOTICE_CHANNEL_ID")
		_, err := dg.ChannelMessageSend(channelID, "@everyone 오늘 10교시에 스터디 일정이 진행되요! 작성한 블로그를 점검해주세요 ☺️")
		if err != nil {
			fmt.Println("메시지 전송 중 에러가 발생했어요: ", err)
		}
	})
	if err != nil {
		log.Fatal("cron 스케줄 등록 중 에러가 발생했어요: ", err)
	}
}

func addWeekendReminder(c *cron.Cron, dg *discordgo.Session) {
	_, err := c.AddFunc("0 15 * * 5,6,0", func() {
		channelID := os.Getenv("DAILY_CHANNEL_ID")
		_, err := dg.ChannelMessageSend(channelID, "주말에도 정진합시다! ☺️")
		if err != nil {
			fmt.Println("메시지 전송 중 에러가 발생했어요: ", err)
		}
	})
	if err != nil {
		log.Fatal("cron 스케줄 등록 중 에러가 발생했어요: ", err)
	}
}

func addCorn(dg *discordgo.Session) *cron.Cron {
	c := cron.New()

	// 매일 오전 8시 데일리 알림
	addDailyMorningReminder(c, dg)

	// 매일 저녁 11시 하루 마무리 알림
	addDailyEveningReminder(c, dg)

	// 매주 월요일 오후 1시 스터디 알림
	addMondayStudyReminder(c, dg)

	// 매주 금, 토, 일 오후 3시 주말 알림
	addWeekendReminder(c, dg)

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

	dailyCron := addCorn(dg)
	dailyCron.Start()

	fmt.Println("봇 활성화 중입니다.")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop

	fmt.Println("봇이 꺼졌습니다.")
	dailyCron.Stop()
}
