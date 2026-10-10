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
		if logger.Log != nil {
			logger.Log.Printf("Discord RPC login failed (Discord may not be running): %v", err)
		}
		return
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
		if logger.Log != nil {
			logger.Log.Println("Discord RPC initialization error:", err)
		}
	}
}

func OpenProject(cfg *config.Config) {
	var smallImage string

	switch cfg.Loader {
	case "forge":
		smallImage = "https://i.imgur.com/O9acTGw.png"
	case "fabric":
		smallImage = "https://i.imgur.com/lLTttOy.png"
	case "neoforge":
		smallImage = "https://i.imgur.com/ooUR90U.png"
	case "quilt":
		smallImage = "https://i.imgur.com/5lkHEvL.png"
	default:
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
		if logger.Log != nil {
			logger.Log.Println("Discord RPC update error:", err)
		}
	}
}
