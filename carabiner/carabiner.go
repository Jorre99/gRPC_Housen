package main

//go:generate protoc -I ../server_fllower_house/proto/ ../server_fllower_house/proto/HouseServer.proto --go_out=plugins=grpc:../server_fllower_house/proto

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/Jorre99/gRPC_Housen/carabiner/backoff"
	_ "github.com/Jorre99/gRPC_Housen/carabiner/resolver"
	"github.com/Jorre99/gRPC_Housen/carabiner/ui"
	chatpb "github.com/Jorre99/gRPC_Housen/server_fllower_house/proto"
	termui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"google.golang.org/grpc"
)

func envOr(key, def string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return def
}

var (
	username = os.Getenv("CHATTER_USER")
	friend   = os.Getenv("CHATTER_FRIEND")
	address  = envOr("CHATTER_SERVER", "dns-srv:///chatter|tcp|housen.tech") // This is an SRV dns record. service|proto|name.
)

// Gouser implements the chat client state.
type Gouser struct {
	server chatpb.BroadcastClient
	chatP  *ui.FlexParagraph
	inputP *widgets.Paragraph
}

// NewGouser creates a new Gouser.
func NewGouser(server chatpb.BroadcastClient) (*Gouser, error) {
	if err := termui.Init(); err != nil {
		return nil, err
	}

	chatP := ui.NewFlexParagraph()
	inputP := widgets.NewParagraph()

	chatP.SelectedRowStyle = termui.NewStyle(termui.ColorYellow)

	g := &Gouser{
		server: server,
		chatP:  chatP,
		inputP: inputP,
	}
	g.resize()
	return g, nil
}

// Run starts the event loop and runs the client.
func (g *Gouser) Run() {
	go g.ListenForMessages()
	evs := termui.PollEvents()
	for e := range evs {
		switch e.Type {
		case termui.KeyboardEvent:
			switch e.ID {
			case "<C-c>":
				termui.Close()
				return
			case "<Enter>":
				g.server.BroadcastMessage(context.Background(), &chatpb.Message{Content: g.inputP.Text, Id: username, PeerUser: friend})
				g.chatP.AddLine(fmt.Sprintf("%s | %s: %s", time.Now().Format("15:04:05"), username, g.inputP.Text))
				g.inputP.Text = ""
				termui.Render(g.inputP, g.chatP)
			case "<Space>":
				g.inputP.Text += " "
				termui.Render(g.inputP)
			case "<Up>":
				g.chatP.ScrollUp()
				termui.Render(g.chatP)
			case "<Down>":
				g.chatP.ScrollDown()
				termui.Render(g.chatP)
			case "<PageUp>":
				g.chatP.ScrollHalfPageUp()
				termui.Render(g.chatP)
			case "<PageDown>":
				g.chatP.ScrollHalfPageDown()
				termui.Render(g.chatP)
			case "<Home>":
				g.chatP.ScrollTop()
				termui.Render(g.chatP)
			case "<End>":
				g.chatP.ScrollBottom()
				termui.Render(g.chatP)
			case "<Backspace>":
				if len(g.inputP.Text) > 0 {
					g.inputP.Text = g.inputP.Text[:len(g.inputP.Text)-1]
					termui.Render(g.inputP)
				}
			default:
				g.inputP.Text += e.ID
				termui.Render(g.inputP)
			}
		case termui.ResizeEvent:
			g.resize()
		case termui.MouseEvent:
			switch e.ID {
			case "<MouseWheelUp>":
				g.chatP.ScrollUp()
				termui.Render(g.chatP)
			case "<MouseWheelDown>":
				g.chatP.ScrollDown()
				termui.Render(g.chatP)
			}
		}
	}
}

func (g *Gouser) resize() {
	x, y := termui.TerminalDimensions()
	if y < 6 {
		return
	}
	g.chatP.SetRect(0, 0, x, y-3)
	g.inputP.SetRect(0, y-3, x, y)
	termui.Render(g.chatP, g.inputP)
}

// ListenForMessages perpetually connects to the gateway and writes new messages to the chat pane.
func (g *Gouser) ListenForMessages() {
	backoff := backoff.NewBackoff(time.Second, time.Minute, 2)
Outer:
	for {
		g.chatP.AddLine("Creating connection...")
		termui.Render(g.chatP)
		stream, err := g.server.CreateStream(context.Background(), &chatpb.Connect{Id: username})
		if err != nil {
			g.chatP.AddLine(err.Error())
			g.chatP.AddLinef("sleeping %s before retrying", backoff.Get())
			termui.Render(g.chatP)
			time.Sleep(backoff.Get())
			backoff.Incr()
			continue Outer
		}
		backoff.Reset()
	Inner:
		for {
			resp, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					g.chatP.AddLine("EOF")
					termui.Render(g.chatP)
					break Inner
				}
				g.chatP.AddLine(err.Error())
				termui.Render(g.chatP)
				break Inner
			}
			g.chatP.AddLine(fmt.Sprintf("%s|%s: %s", time.Now().Format("15:04:05"), resp.GetId(), resp.GetContent()))
			termui.Render(g.chatP)
		}
	}
}

func mainWith(ctx context.Context) error {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	client := chatpb.NewBroadcastClient(conn)
	g, err := NewGouser(client)
	if err != nil {
		return err
	}

	g.Run()

	m := &chatpb.Message{Content: "quit", PeerUser: friend}
	if _, err := client.BroadcastMessage(ctx, m); err != nil {
		return err
	}
	return nil
}

func main() {
	ctx := context.Background()

	if err := mainWith(ctx); err != nil {
		log.Fatal(err)
	}
}
