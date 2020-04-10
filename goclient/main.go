package main

import (
	"context"
	"io"
	"log"
	"time"

	_ "github.com/Jorre99/gRPC_Housen/goclient/resolver"
	chatpb "github.com/Jorre99/gRPC_Housen/server_fllower_house/proto"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"google.golang.org/grpc"
)

var (
	address = "dns-srv:///grpclb|tcp|chatter.housen.tech" // This is an SRV dns record. service|proto|name.
)

// Gouser implements the chat client state.
type Gouser struct {
	server chatpb.BroadcastClient
	chatP  *widgets.Paragraph
	inputP *widgets.Paragraph
}

// NewGouser creates a new Gouser.
func NewGouser(server chatpb.BroadcastClient) (*Gouser, error) {
	if err := ui.Init(); err != nil {
		return nil, err
	}

	x, y := ui.TerminalDimensions()

	chatP := widgets.NewParagraph()
	chatP.SetRect(0, 0, x, y-4)
	inputP := widgets.NewParagraph()
	inputP.SetRect(0, y-4, x, y)

	g := &Gouser{
		server: server,
		chatP:  chatP,
		inputP: inputP,
	}
	ui.Render(inputP, chatP)
	return g, nil
}

// Run starts the event loop and runs the client.
func (g *Gouser) Run() {
	go g.ListenForMessages()
	evs := ui.PollEvents()
	for e := range evs {
		switch e.Type {
		case ui.KeyboardEvent:
			switch e.ID {
			case "<C-c>":
				ui.Close()
				return
			case "<Enter>":
				g.server.BroadcastMessage(context.Background(), &chatpb.Message{Content: g.inputP.Text})
				g.inputP.Text = ""
				ui.Render(g.inputP)
			case "<Space>":
				g.inputP.Text += " "
				ui.Render(g.inputP)
			case "<Backspace>":
				if len(g.inputP.Text) > 0 {
					g.inputP.Text = g.inputP.Text[:len(g.inputP.Text)-1]
					ui.Render(g.inputP)
				}
			default:
				g.inputP.Text += e.ID
				ui.Render(g.inputP)
			}
		}
	}
}

// ListenForMessages perpetually connects to the gateway and writes new messages to the chat pane.
func (g *Gouser) ListenForMessages() {
Outer:
	for {
		g.chatP.Text += "Creating connection...\n"
		ui.Render(g.chatP)
		stream, err := g.server.CreateStream(context.Background(), &chatpb.Connect{})
		if err != nil {
			g.chatP.Text += err.Error()
			g.chatP.Text += "\n"
			ui.Render(g.chatP)
			continue Outer
		}
	Inner:
		for {
			resp, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					g.chatP.Text += "EOF\n"
					ui.Render(g.chatP)
					break Inner
				}
				g.chatP.Text += err.Error()
				g.chatP.Text += "\n"
				ui.Render(g.chatP)
				break Inner
			}
			g.chatP.Text += resp.GetContent()
			g.chatP.Text += "\n"
			ui.Render(g.chatP)
		}
	}
}

func mainWith(ctx context.Context) error {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(50*time.Second))
	if err != nil {
		return err
	}
	defer conn.Close()
	log.Printf("connected to: %s", conn.Target())

	client := chatpb.NewBroadcastClient(conn)
	g, err := NewGouser(client)
	if err != nil {
		return err
	}

	g.Run()

	m := &chatpb.Message{Content: "quit"}
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
