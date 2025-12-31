package discord

import (
	"fmt"
	"time"

	"github.com/hugolgst/rich-go/client"
	"github.com/sqot0/packsmith/backend/internal/config"
	"github.com/sqot0/packsmith/backend/internal/logger"
)

var startTime time.Time

func Init() {
	err := client.Login("1455868971067637763")
	if err != nil {
		panic(err)
	}

	startTime = time.Now()

	err = client.SetActivity(client.Activity{
		State:      "Choosing a modpack",
		Details:    "Smithing mod packs",
		LargeImage: "https://i.imgur.com/EhxmU2E.png",
		Timestamps: &client.Timestamps{
			Start: &startTime,
		},
	})

	if err != nil {
		logger.Log.Println("Discord RPC initialization error:", err)
	}
}

func OpenProject(cfg *config.Config) {
	var smallImage string

	if cfg.Loader == "forge" {
		smallImage = "https://i.imgur.com/O9acTGw.png"
	} else if cfg.Loader == "fabric" {
		smallImage = "https://i.imgur.com/lLTttOy.png"
	} else {
		smallImage = ""
	}

	err := client.SetActivity(client.Activity{
		State:      fmt.Sprintf("Modpack: %s %s", cfg.Name, cfg.Minecraft),
		Details:    "Smithing mod packs",
		LargeImage: "https://i.imgur.com/EhxmU2E.png",
		SmallImage: smallImage,
		Timestamps: &client.Timestamps{
			Start: &startTime,
		},
	})

	if err != nil {
		logger.Log.Println("Discord RPC update error:", err)
	}
}
