package main

import (
	"fmt"
	"log"
	"time"

	"github.com/AkihikoITOH/protobuf-over-zeromq-concept/client/pb"
	"github.com/golang/protobuf/ptypes/timestamp"

	"github.com/marcusolsson/tui-go"
)

func newHistory() *tui.Box {
	return tui.NewVBox()
}

func (v *View) newHistoryBox() *tui.Box {
	historyScroll := tui.NewScrollArea(v.history)
	historyScroll.SetAutoscrollToBottom(true)

	historyBox := tui.NewVBox(historyScroll)
	historyBox.SetBorder(true)

	return historyBox
}

func newInput() *tui.Entry {
	input := tui.NewEntry()
	input.SetFocused(true)
	input.SetSizePolicy(tui.Expanding, tui.Maximum)

	return input
}

func (v *View) newInputBox() *tui.Box {
	inputBox := tui.NewHBox(v.input)
	inputBox.SetBorder(true)
	inputBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	return inputBox
}

func (v *View) newChatView() *tui.Box {
	chat := tui.NewVBox(
		v.newHistoryBox(),
		v.newInputBox(),
	)
	chat.SetSizePolicy(tui.Expanding, tui.Expanding)

	return chat
}

func format(message *pb.Message) *tui.Box {
	return tui.NewHBox(
		tui.NewLabel(time.Unix(message.GetTime().GetSeconds(), 0).Format(time.RFC3339)),
		tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("<%s (%s)>", message.GetSender().GetName(), message.GetSender().GetUuid()))),
		tui.NewLabel(message.GetContent()),
		tui.NewSpacer(),
	)
}

func (v *View) defaultOnInput() func(e *tui.Entry) {
	return func(e *tui.Entry) {
		ts := &timestamp.Timestamp{Seconds: time.Now().UTC().Unix()}
		message := &pb.Message{Time: ts, Sender: v.user, Content: e.Text()}
		v.input.SetText("")
		v.outgoingMessages <- message
	}
}

func (v *View) defaultOnNewMessage() func(message *pb.Message) {
	return func(message *pb.Message) {
		v.history.Append(format(message))
		v.Repaint()
	}
}

type View struct {
	tui.UI
	user             *pb.User
	input            *tui.Entry
	history          *tui.Box
	onNewMessage     func(msg *pb.Message)
	outgoingMessages chan<- *pb.Message
	incomingMessages <-chan *pb.Message
}

func NewView(user *pb.User, incomingMessages <-chan *pb.Message, outgoingMessages chan<- *pb.Message) *View {
	input := newInput()
	history := newHistory()
	view := &View{UI: nil, input: input, history: history, user: user, incomingMessages: incomingMessages, outgoingMessages: outgoingMessages}

	return view
}

func (v *View) BuildUI() error {
	v.input.OnSubmit(v.defaultOnInput())
	v.onNewMessage = v.defaultOnNewMessage()
	chat := v.newChatView()
	view, err := tui.New(chat)
	if err != nil {
		log.Fatal(err)
		return err
	}
	v.UI = view

	v.SetKeybinding("Esc", func() { v.Quit() })
	v.SetKeybinding("Ctrl-C", func() { v.Quit() })

	return nil
}

func (v *View) Start() {
	if v.UI == nil {
		return
	}
	go func() {
		for {
			msg, ok := <-v.incomingMessages
			if !ok {
				break
			}
			v.onNewMessage(msg)
		}
	}()

	if err := v.Run(); err != nil {
		log.Fatal(err)
	}
}
