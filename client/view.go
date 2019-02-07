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

func newHistoryBox(history *tui.Box) *tui.Box {
	historyScroll := tui.NewScrollArea(history)
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

func newInputBox(input *tui.Entry) *tui.Box {
	inputBox := tui.NewHBox(input)
	inputBox.SetBorder(true)
	inputBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	return inputBox
}

func newChatView() *tui.Box {
	history := newHistory()
	input := newInput()
	chat := tui.NewVBox(
		newHistoryBox(history),
		newInputBox(input),
	)
	chat.SetSizePolicy(tui.Expanding, tui.Expanding)

	return chat
}

func format(message *pb.Message) *tui.Box {
	return tui.NewHBox(
		tui.NewLabel(time.Unix(message.GetTime().GetSeconds(), 0).Format(time.RFC3339)),
		tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("<%s>", message.GetSender()))),
		tui.NewLabel(message.GetContent()),
		tui.NewSpacer(),
	)
}

func defaultOnInput(history *tui.Box, input *tui.Entry, msgCh chan<- *pb.Message) func(e *tui.Entry) {
	return func(e *tui.Entry) {
		ts := &timestamp.Timestamp{Seconds: time.Now().UTC().Unix()}
		message := &pb.Message{Time: ts, Sender: "John", Content: e.Text()}
		input.SetText("")
		msgCh <- message
	}
}

func defaultOnNewMessage(history *tui.Box) func(message *pb.Message) {
	return func(message *pb.Message) {
		history.Append(format(message))
	}
}

type View struct {
	tui.UI
	onNewMessage     func(msg *pb.Message)
	incomingMessages <-chan *pb.Message
}

func NewView(incomingMessages <-chan *pb.Message, outgoingMessages chan<- *pb.Message) (*View, error) {
	input := newInput()
	history := newHistory()
	input.OnSubmit(defaultOnInput(history, input, outgoingMessages))
	chat := tui.NewVBox(
		newHistoryBox(history),
		newInputBox(input),
	)
	chat.SetSizePolicy(tui.Expanding, tui.Expanding)
	view, err := tui.New(chat)
	if err != nil {
		log.Fatal(err)
	}

	view.SetKeybinding("Esc", func() { view.Quit() })
	view.SetKeybinding("Ctrl-C", func() { view.Quit() })

	return &View{UI: view, onNewMessage: defaultOnNewMessage(history), incomingMessages: incomingMessages}, err
}

func (v *View) Start() {
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
