package main

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/ebfe/scard"
	"github.com/hikarisakira/nhi-pcsc-reader/models"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.org/x/text/encoding/traditionalchinese"
)

var (
	// 定義 APDU 指令
	selectAPDU      = []byte{0x00, 0xA4, 0x04, 0x00, 0x10, 0xD1, 0x58, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x11, 0x00}
	readProfileAPDU = []byte{0x00, 0xca, 0x11, 0x00, 0x02, 0x00, 0x00}
)

// App struct
type App struct {
	ctx      context.Context
	cardChan chan models.NhicFormat
}

// NewApp creates a new App application struct
func NewApp() *App {

	return &App{
		cardChan: make(chan models.NhicFormat),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	go a.startCardReader()
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) startCardReader() {
	for {
		info, err := readCard()
		if err != nil {
			runtime.EventsEmit(a.ctx, "card-status", models.NhicFormat{IsCardExist: false})
			time.Sleep(1 * time.Second)
			continue
		}
		runtime.EventsEmit(a.ctx, "card-status", info)
		time.Sleep(500 * time.Millisecond)
	}
}

func readCard() (models.NhicFormat, error) {
	context, err := scard.EstablishContext()
	if err != nil {
		return models.NhicFormat{}, err
	}
	defer context.Release()

	readers, err := context.ListReaders()
	if err != nil || len(readers) == 0 {
		return models.NhicFormat{}, fmt.Errorf("無法取得讀卡機")
	}

	card, err := context.Connect(readers[0], scard.ShareShared, scard.ProtocolAny)
	if err != nil {
		return models.NhicFormat{}, err
	}
	defer card.Disconnect(scard.LeaveCard)

	// 發送 select APDU
	rsp, err := card.Transmit(selectAPDU)
	if err != nil {
		return models.NhicFormat{}, err
	}

	// 讀取卡片資料
	rsp, err = card.Transmit(readProfileAPDU)
	if err != nil {
		return models.NhicFormat{}, err
	}

	// 解析資料
	big5 := traditionalchinese.Big5.NewDecoder()
	nameBig5 := bytes.TrimRight(rsp[12:18], "\x00")
	nameUtf8, _ := big5.Bytes(nameBig5)

	return models.NhicFormat{
		CardNumber:  string(rsp[0:12]),
		Name:        string(nameUtf8),
		IdNumber:    string(rsp[32:42]),
		Birthday:    string(rsp[43:49]),
		Sex:         string(rsp[49:50]),
		CardDate:    string(rsp[51:57]),
		IsCardExist: true,
	}, nil
}
