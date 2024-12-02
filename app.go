package main

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/ebfe/scard"
	"github.com/hikarisakira/nhi-pcsc-reader/models"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	traditionalChinese "golang.org/x/text/encoding/traditionalchinese"
)

var (
	selectAPDU      = []byte{0x00, 0xA4, 0x04, 0x00, 0x10, 0xD1, 0x58, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x11, 0x00}
	readProfileAPDU = []byte{0x00, 0xca, 0x11, 0x00, 0x02, 0x00, 0x00}
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	go a.startCardReader()
}

func (a *App) startCardReader() {
	for {
		if info, err := readCard(); err == nil {
			runtime.EventsEmit(a.ctx, "card-status", info)
		} else {
			runtime.EventsEmit(a.ctx, "card-status", models.NhicFormat{IsCardExist: false})
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func readCard() (models.NhicFormat, error) {
	ctx, card, err := connectCard()
	if err != nil {
		return models.NhicFormat{}, err
	}
	defer func() {
		card.Disconnect(scard.LeaveCard)
		ctx.Release()
	}()

	rsp, err := readCardData(card)
	if err != nil {
		return models.NhicFormat{}, err
	}

	return parseCardData(rsp), nil
}

func connectCard() (*scard.Context, *scard.Card, error) {
	ctx, err := scard.EstablishContext()
	if err != nil {
		return nil, nil, err
	}

	readers, err := ctx.ListReaders()
	if err != nil || len(readers) == 0 {
		ctx.Release()
		return nil, nil, fmt.Errorf("無法取得讀卡機")
	}

	card, err := ctx.Connect(readers[0], scard.ShareShared, scard.ProtocolAny)
	if err != nil {
		ctx.Release()
		return nil, nil, err
	}

	return ctx, card, nil
}

func readCardData(card *scard.Card) ([]byte, error) {
	if _, err := card.Transmit(selectAPDU); err != nil {
		return nil, err
	}
	return card.Transmit(readProfileAPDU)
}

func parseCardData(rsp []byte) models.NhicFormat {
	big5 := traditionalChinese.Big5.NewDecoder()
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
	}
}
