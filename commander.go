package main

import (
	"fmt"
	"log"
	"time"

	. "github.com/bearyinnovative/lili/commands"
	. "github.com/bearyinnovative/lili/commands/house"
	. "github.com/bearyinnovative/lili/model"
	. "github.com/bearyinnovative/lili/util"

	"github.com/dustin/go-humanize"
)

func RunCommander() {
	cmds := []CommandType{
		NewHackerNewsSlack(),
		NewHackerNewsAll(),
		NewMatsumotooooooInstagram(),
		NewDabieCatInstagram(),
		NewArkDomeInstagram(),
		NewArkDomeInstagram2(),
		NewRoCryInstagram(),
		NewArkDomeV2(),

		NewIMessageZhihu(),
		NewTelegramZhihu(),
		NewWhatsAppZhihu(),
		NewBearyChatZhihu(),
		NewDingDingZhihu(),
		NewSlackZhihu(),

		NewBearyChatV2EX(),
		NewDingDingV2EX(),
		NewSlackV2EX(),
	}

	cmds = append(cmds, DealCommands...)

	for i := 0; i < len(cmds); i++ {
		// fmt.Printf("%+v\n", cmds[i])
		start(cmds[i])
	}

	// FIXME:
	// wait forever
	select {}
}

func start(c CommandType) {
	// trigger once
	go func() {
		fetchAndNotify(c)
	}()

	ticker := time.NewTicker(c.Interval())
	// schedule run
	go func() {
		for _ = range ticker.C {
			fetchAndNotify(c)
		}
	}()
}

func fetchAndNotify(c CommandType) {
	items, err := c.Fetch()
	if err != nil {
		return
	}

	notifiedCount := 0

	for _, item := range items {
		created, err := DBContext.UpsertItem(item)
		if LogIfErr(err) {
			continue
		}

		if !created {
			continue
		}

		if !item.InDays(31) {
			log.Println("too old to notify:", item.Desc)
			continue
		}

		notifiedCount += 1

		// notify
		text := fmt.Sprintf("%s (%s)", item.Desc, humanize.Time(item.Created))
		for _, n := range c.Notifiers() {
			err = n.Notify(text, item.Images)
			LogIfErr(err)
			if err == nil {
				err = DBContext.MarkNotified(item)
				LogIfErr(err)
			}
		}
	}

	log.Printf("[%s] fetched %d items, notified %d", c.Name(), len(items), notifiedCount)
}
