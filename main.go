package main

import (
	"context"
	"embed"
	"log"

	"github.com/vvitem/item_all/internal/storage/sqlite"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	storageRuntime := bootstrapStorage(context.Background(), defaultStorageOpener{}, sqlite.DefaultConfig(""))
	app := NewApp(storageRuntime.err)

	err := wails.Run(&options.App{
		Title:     "ItemAll",
		Width:     1024,
		Height:    700,
		MinWidth:  800,
		MinHeight: 560,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 2, G: 6, B: 23, A: 1},
		Bind: []interface{}{
			app,
		},
		OnShutdown: func(context.Context) {
			if err := storageRuntime.Close(); err != nil {
				log.Printf("storage shutdown failed: %v", err)
			}
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
