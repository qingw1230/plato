package client

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/gookit/color"
	"github.com/qingw1230/plato/common/sdk"
	"github.com/rocket049/gocui"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var (
	buf  string
	chat *sdk.Chat
	pos  int // 用于输入框上、下翻
)

// setHeadText 设置标题框文字
func setHeadText(g *gocui.Gui, msg string) {
	v, err := g.View("head")
	if err == nil {
		v.Clear()
		fmt.Fprint(v, color.FgGreen.Text(msg))
	}
}

// headView 创建标题框 view
func headView(g *gocui.Gui, x0, y0, x1, y1 int) error {
	v, err := g.SetView("head", x0, y0, x1, y1)
	if err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Overwrite = true
		v.Wrap = false
		msg := "start chat!"
		setHeadText(g, msg)
	}
	return nil
}

// outputView 创建输出框 view
func outputView(g *gocui.Gui, x0, y0, x1, y1 int) error {
	v, err := g.SetView("out", x0, y0, x1, y1)
	if err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.SelBgColor = gocui.ColorRed
		v.Overwrite = false
		v.Frame = true
		v.Wrap = true
		v.Autoscroll = true
		v.Title = "Messages"
	}
	return nil
}

// inputView 创建输入框 view
func inputView(g *gocui.Gui, x0, y0, x1, y1 int) error {
	v, err := g.SetView("in", x0, y0, x1, y1)
	if err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Editable = true
		v.Overwrite = false
		v.Wrap = true
		if _, err := g.SetCurrentView("in"); err != nil {
			return err
		}
	}
	return nil
}

// layout 创建标题框、输入框、输出框 view
func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if err := headView(g, 1, 1, maxX-1, 3); err != nil {
		return err
	}
	if err := outputView(g, 1, 4, maxX-1, maxY-4); err != nil {
		return err
	}
	if err := inputView(g, 1, maxY-3, maxX-1, maxY-1); err != nil {
		return err
	}
	return nil
}

func printView(g *gocui.Gui, name, msg string, newLine bool) {
	var sep string
	if newLine {
		sep = "\n"
	} else {
		sep = " "
	}

	v, err := g.View("out")
	if err != nil {
		log.Panicln(err)
	}
	fmt.Fprintf(v, "%v:%v%v\n", color.FgGreen.Text(name), sep, color.FgYellow.Text(msg))
}

func receiveMessage(g *gocui.Gui) {
	receiveChan := chat.Receive()
	for msg := range receiveChan {
		switch msg.Type {
		case sdk.MsgTypeText:
			printView(g, msg.Name, msg.Content, false)
		}
	}
}

func sendMessage(g *gocui.Gui, iv *gocui.View) {
	ov, err := g.View("out")
	if iv != nil && err == nil {
		b := iv.ReadEditor()
		if b != nil {
			msg := &sdk.Message{
				Type:       sdk.MsgTypeText,
				Name:       chat.Nick,
				FromUserID: chat.UserID,
				ToUserID:   "1231",
				Content:    string(b),
			}
			printView(g, "me", msg.Content, false)
			chat.Send(msg)
		}
	}
	ov.Autoscroll = true
}

// updateView 更新 view
func updateView(g *gocui.Gui, iv *gocui.View) error {
	sendMessage(g, iv)
	l := len(iv.Buffer())
	iv.MoveCursor(0-l, 0, true)
	iv.Clear()
	return nil
}

// quit 退出 sdk 界面
func quit(g *gocui.Gui, _ *gocui.View) error {
	chat.Close()
	v, _ := g.View("out")
	buf = v.Buffer()
	return gocui.ErrQuit
}

// pasteUp 输入框上翻
func pasteUp(g *gocui.Gui, iv *gocui.View) error {
	ov, err := g.View("out")
	if err != nil {
		fmt.Fprintf(iv, "errors:%s", err)
		return nil
	}
	lines := ov.BufferLines()
	bufLen := len(lines)
	if bufLen == 0 {
		return nil
	}
	if pos < bufLen-1 {
		pos++
	}
	iv.Clear()
	s := lines[bufLen-1-pos]
	idx := strings.Index(s, ":")
	fmt.Fprintf(iv, "%s", s[idx+2:])
	return nil
}

// pasteDown 输入框下翻
func pasteDown(g *gocui.Gui, iv *gocui.View) error {
	ov, err := g.View("out")
	if err != nil {
		fmt.Fprintf(iv, "errors:%s", err)
		return nil
	}
	if pos > 0 {
		pos--
	}
	lines := ov.BufferLines()
	bufLen := len(lines)
	if bufLen == 0 {
		return nil
	}
	iv.Clear()
	s := lines[bufLen-1-pos]
	idx := strings.Index(s, ":")
	fmt.Fprintf(iv, "%s", s[idx+2:])
	return nil
}

func RunMain() {
	g, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true
	g.Mouse = false
	g.ASCII = false
	g.SetManagerFunc(layout)

	chat = sdk.NewChat("127.0.0.1:8080", "test-im", "1230", "123456")

	// 注册输入框按键的回调函数
	if err := g.SetKeybinding("in", gocui.KeyEnter, gocui.ModNone, updateView); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("in", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("in", gocui.KeyArrowUp, gocui.ModNone, pasteUp); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("in", gocui.KeyArrowDown, gocui.ModNone, pasteDown); err != nil {
		log.Panicln(err)
	}

	go receiveMessage(g)
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Println(err)
	}

	os.WriteFile("logs/chat.log", []byte(buf), 0644)
}
