package client

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/qingw1230/plato/common/config"
	"github.com/qingw1230/plato/common/sdk"

	"github.com/gookit/color"
	"github.com/rocket049/gocui"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var (
	buf  string
	chat *sdk.Chat
)

type VOT struct {
	Name string
	Msg  string
	Sep  string
}

func (vot VOT) Show(g *gocui.Gui) error {
	v, err := g.View("out")
	if err != nil {
		return nil
	}
	_, _ = fmt.Fprintf(v, "%v:%v%v\n", color.FgGreen.Text(vot.Name), vot.Sep, color.FgYellow.Text(vot.Msg))
	return nil
}

func viewPrint(g *gocui.Gui, name, msg string, newLine bool) {
	var out VOT
	out.Name = name
	out.Msg = msg
	if newLine {
		out.Sep = "\n"
	} else {
		out.Sep = " "
	}
	g.Update(out.Show)
}

// doReceive 不断接收消息
func doReceive(g *gocui.Gui) {
	receiveChan := chat.Receive()
	for msg := range receiveChan {
		switch msg.Type {
		case sdk.MsgTypeText:
			viewPrint(g, msg.Name, msg.Content, false)
		}
	}
}

// doSay 发送消息
func doSay(g *gocui.Gui, cv *gocui.View) {
	v, err := g.View("out")
	if cv != nil && err == nil {
		p := cv.ReadEditor()
		if p != nil {
			msg := &sdk.Message{
				Type:       sdk.MsgTypeText,
				Name:       "qgw",
				FromUserID: "1234",
				ToUserID:   "4321",
				Content:    string(p),
			}
			// TODO: 保证消息显示一致性
			viewPrint(g, "me", msg.Content, false)
			chat.Send(msg)
		}
		v.Autoscroll = true
	}
}

// setHeadText 设置标题框文字
func setHeadText(g *gocui.Gui, msg string) {
	v, err := g.View("head")
	if err == nil {
		v.Clear()
		_, _ = fmt.Fprint(v, color.FgGreen.Text(msg))
	}
}

// quit 退出 sdk 界面
func quit(g *gocui.Gui, _ *gocui.View) error {
	chat.Close()
	v, _ := g.View("out")
	buf = v.Buffer()
	return gocui.ErrQuit
}

func viewUpdate(g *gocui.Gui, cv *gocui.View) error {
	doSay(g, cv)
	l := len(cv.Buffer())
	cv.MoveCursor(0-l, 0, true)
	cv.Clear()
	return nil
}

// viewHead 创建标题框视图
func viewHead(g *gocui.Gui, x0, y0, x1, y1 int) error {
	if v, err := g.SetView("head", x0, y0, x1, y1); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Wrap = false
		v.Overwrite = true
		msg := "开始聊天吧!"
		setHeadText(g, msg)
	}
	return nil
}

// viewOutput 创建输出框视图
func viewOutput(g *gocui.Gui, x0, y0, x1, y1 int) error {
	v, err := g.SetView("out", x0, y0, x1, y1)
	if err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Wrap = true
		v.Overwrite = false
		v.Autoscroll = true
		v.SelBgColor = gocui.ColorRed
		v.Title = "Messages"
	}
	return nil
}

// viewInput 创建输入框视图
func viewInput(g *gocui.Gui, x0, y0, x1, y1 int) error {
	if v, err := g.SetView("in", x0, y0, x1, y1); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Editable = true
		v.Wrap = true
		v.Overwrite = false
		if _, err := g.SetCurrentView("in"); err != nil {
			return err
		}
	}
	return nil
}

// layout 创建标题框、输入框、输出框视图
func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if err := viewHead(g, 1, 1, maxX-1, 3); err != nil {
		return err
	}
	if err := viewOutput(g, 1, 4, maxX-1, maxY-4); err != nil {
		return err
	}
	if err := viewInput(g, 1, maxY-3, maxX-1, maxY-1); err != nil {
		return err
	}
	return nil
}

var pos int

// viewUp 输出框上翻
func viewUp(g *gocui.Gui, cv *gocui.View) error {
	v, err := g.View("out")
	v.Autoscroll = false
	ox, oy := v.Origin()
	if err == nil {
		_ = v.SetOrigin(ox, oy-1)
	}
	return nil
}

// viewDown 输入框下翻
func viewDown(g *gocui.Gui, cv *gocui.View) error {
	v, err := g.View("out")
	_, y := v.Size()
	ox, oy := v.Origin()
	lenNum := len(v.BufferLines())
	if err == nil {
		if oy > lenNum-y-1 {
			v.Autoscroll = true
		} else {
			_ = v.SetOrigin(ox, oy+1)
		}
	}
	return nil
}

// pasteUp 输入框上翻
func pasteUp(g *gocui.Gui, cv *gocui.View) error {
	v, err := g.View("out")
	if err != nil {
		_, _ = fmt.Fprintf(cv, "error:%s", err)
		return nil
	}
	lines := v.BufferLines()
	bufLen := len(lines)
	if pos < bufLen-1 {
		pos++
	}
	cv.Clear()
	_, _ = fmt.Fprintf(cv, "%s", lines[bufLen-pos-1])
	return nil
}

// pasteDown 输入框下翻
func pasteDown(g *gocui.Gui, cv *gocui.View) error {
	v, err := g.View("out")
	if err != nil {
		_, _ = fmt.Fprintf(cv, "errors:%s", err)
		return nil
	}
	if pos > 0 {
		pos--
	}
	lines := v.BufferLines()
	bufLen := len(lines)
	cv.Clear()
	_, _ = fmt.Fprintf(cv, "%s", lines[bufLen-pos-1])
	return nil
}

// RunMain 运行 sdk 客户端
func RunMain() {
	config.Init("/home/qgw/git/plato/im.yaml")
	_, err := net.ListenTCP("tcp", &net.TCPAddr{
		Port: config.GetGatewayServerPort(),
	})
	if err != nil {
		log.Fatalf("net.ListenTCP err:%s", err.Error())
		panic(err)
	}

	chat = sdk.NewChat(net.ParseIP("0.0.0.0"), 8900, "test-im", "1230", "12301230")
	chat.Receive()
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true
	g.Mouse = false
	g.ASCII = false
	g.SetManagerFunc(layout)

	// 注册输入框按键的回调函数
	if err := g.SetKeybinding("in", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("in", gocui.KeyEnter, gocui.ModNone, viewUpdate); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("in", gocui.KeyPgup, gocui.ModNone, viewUp); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("in", gocui.KeyPgdn, gocui.ModNone, viewDown); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("in", gocui.KeyArrowUp, gocui.ModNone, pasteUp); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("in", gocui.KeyArrowDown, gocui.ModNone, pasteDown); err != nil {
		log.Panicln(err)
	}

	go doReceive(g)
	if err := g.MainLoop(); err != nil {
		log.Println(err)
	}
	_ = ioutil.WriteFile("chat.log", []byte(buf), 0644)
}
