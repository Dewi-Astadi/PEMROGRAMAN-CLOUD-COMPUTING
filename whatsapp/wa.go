package wa

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

var client *whatsmeow.Client

func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		fmt.Println("Received a message!", v.Message.GetConversation())
		fmt.Println("=> dari saya =", v.Info.IsFromMe)
		fmt.Println("=> pesan =", v.Message.GetConversation())
		fmt.Println("=> Server =", v.Info.MessageSource.Chat.Server)
		fmt.Println("=> apakah group = ", v.Info.IsGroup)
		fmt.Println("=> apakah broadcast =", v.Info.IsIncomingBroadcast())

		if !v.Info.IsFromMe &&
			!v.Info.IsGroup &&
			!v.Info.IsIncomingBroadcast() {

			fmt.Println("PENGIRIM =", v.Info.Sender.User)
			pesan := v.Message.GetConversation()
			fmt.Println("PESAN ASLI =[" + pesan + "]")

			// Normalisasi: lowercase dan trim spasi
			pesanNormalized := strings.ToLower(strings.TrimSpace(pesan))
			fmt.Println("PESAN NORMALIZED =[" + pesanNormalized + "]")

			var id_wa []string
			id_wa = append(id_wa, v.Info.ID)

			client.MarkRead(context.Background(), id_wa, time.Now(), v.Info.Chat, v.Info.Sender)

			client.SubscribePresence(context.Background(), v.Info.Sender)

			client.SendPresence(context.Background(), types.PresenceAvailable)

			time.Sleep(2 * time.Second)

			client.SendChatPresence(context.Background(), v.Info.Sender, types.ChatPresenceComposing, types.ChatPresenceMediaText)

			time.Sleep(3 * time.Second)

			client.SendChatPresence(context.Background(), v.Info.Sender, types.ChatPresencePaused, types.ChatPresenceMediaText)

			// STRICT CHECK: Hanya test atau tes
			fmt.Println("=== CEK KONDISI ===")
			fmt.Println("pesan == 'test': ", pesanNormalized == "test")
			fmt.Println("pesan == 'tes': ", pesanNormalized == "tes")

			if pesanNormalized == "test" || pesanNormalized == "tes" {
				fmt.Println(">>> MEMBALAS PESAN <<<")
				kirimPesan(v.Info.Sender)
			} else {
				fmt.Println(">>> TIDAK MEMBALAS (bukan test/tes) <<<")
			}
		}
	}
}

func KonekWa() {
	dbLog := waLog.Stdout("Database", "DEBUG", true)
	ctx := context.Background()
	container, err := sqlstore.New(ctx, "sqlite3", "file:examplestore.db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}
	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		panic(err)
	}
	clientLog := waLog.Stdout("Client", "DEBUG", true)
	client = whatsmeow.NewClient(deviceStore, clientLog)
	client.AddEventHandler(eventHandler)

	if client.Store.ID == nil {
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			panic(err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		err = client.Connect()
		if err != nil {
			panic(err)
		}
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()
}

func kirimPesan(JIDPenerima types.JID) {
	JIDPenerima.Device = 0
	_, err := client.SendMessage(
		context.Background(),
		JIDPenerima,
		&waE2E.Message{
			Conversation: proto.String("UJI COBA\n PESAN OTOMATIS"),
		},
	)
	if err != nil {
		fmt.Println("ERROR kirim pesan:", err)
	} else {
		fmt.Println("Pesan berhasil terkirim!")
	}
}
