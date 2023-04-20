package main

import (
	"context"
	"log"
	"os"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"strings"
	"strconv"
	"math/rand"
	"fmt"
)

type Context struct {
	C context.Context
	O events.MessageNewObject
}

const (
	CmdPrefix = "mjs"
)

var (
	vk *api.VK
	cmds = map[string] func(*Context, []string) string {
		"rand" : cmdRand,
		"chance" : cmdChance,
	}
)

func (c *Context)sendGenMessage(msg string) {
	b := params.NewMessagesSendBuilder()
	b.Message(msg)
	b.RandomID(0)
	b.PeerID(c.O.Message.PeerID)

	_, err := vk.MessagesSend(b.Params)
	if err != nil {
		log.Fatal(err)
	}
}

func cmdChance(c *Context, args []string) string {
	ev := strings.Join(args, " ")

	msg := fmt.Sprintf("Шанс на то, что %s = %d%%", ev, rand.Intn(102))
	c.sendGenMessage(msg)

	return ""
}

func cmdRand(c *Context, args []string) string {
	var (
		n1, n2 int
		err error
	)

	n1 = 0
	n2 = 100

	if len(args) == 1 {
		numStr := args[0]
		n2, err = strconv.Atoi(numStr)
		if err != nil {
			return "ЧИСЛО, ОСЁЛ, БЛЯДЬ"
		}
	} else if len(args) == 2 {
		numStr := args[0]
		n1, err = strconv.Atoi(numStr)
		if err != nil {
			return "ЧИСЛО, СУКА, БЛЯДЬ"
		}


		numStr = args[1]
		n2, err = strconv.Atoi(numStr)
		if err != nil {
			return "ЧИСЛО, СУКА, БЛЯДЬ"
		}
	}

	if n1 >= n2 {
		return fmt.Sprintf("НЕПРАВИЛЬНЫЕ ЧИСЛА, ИДИОТ (%d, %d)", n1, n2)
	}

	n := rand.Intn(n2-n1) + n1
	msg := fmt.Sprintf("Твоё тупое число между %d и %d это %d", n1, n2, n)
	c.sendGenMessage(msg)


	return ""
}

func handleMessage(c context.Context, obj events.MessageNewObject) {
	log.Printf("%d: %s", obj.Message.PeerID, obj.Message.Text)
	ctx := &Context{c, obj}
	txt := obj.Message.Text
	if !strings.HasPrefix(txt, CmdPrefix) {
		return
	}

	fields := strings.Fields(txt)[1:]
	if len(fields) < 1{
		ctx.sendGenMessage("Чё те надо?")
		return
	}

	cmdName := fields[0]
	fields = fields[1:]
	f, ok := cmds[cmdName]
	if !ok {
		ctx.sendGenMessage("Нет такой команды, еблан")
		return
	}
	err := f(ctx, fields)
	if err != "" {
		ctx.sendGenMessage(err)
	}

}

func main() {
	token := os.Getenv("TOKEN")
	vk = api.NewVK(token)

	// get information about the group
	group, err := vk.GroupsGetByID(nil)
	if err != nil {
		log.Fatal(err)
	}

	// Initializing Long Poll
	lp, err := longpoll.NewLongPoll(vk, group[0].ID)
	if err != nil {
		log.Fatal(err)
	}

	// New message event
	lp.MessageNew(handleMessage)

	// Run Bots Long Poll
	log.Println("Start Long Poll")
	if err := lp.Run(); err != nil {
		log.Fatal(err)
	}
}

