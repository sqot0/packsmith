package sources

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"

	"github.com/pelletier/go-toml/v2"
)

type FabricMod struct {
	ID string `json:"id"`
}

type QuiltMod struct {
	QuiltLoader struct {
		ID string `json:"id"`
	} `json:"quilt_loader"`
}

type ForgeModsToml struct {
	Mods []struct {
		ModID string `toml:"modId"`
	} `toml:"mods"`
}

func GetModID(path string) (string, error) {
	reader, err := zip.OpenReader(path)
	if err != nil {
		return "", err
	}
	defer reader.Close()

	for _, file := range reader.File {
		switch file.Name {

		case "fabric.mod.json":
			return parseFabric(file)

		case "quilt.mod.json":
			return parseQuilt(file)

		case "META-INF/mods.toml":
			return parseForge(file)

		case "META-INF/neoforge.mods.toml":
			return parseForge(file)
		}
	}

	return "", fmt.Errorf("mod metadata not found")
}

func parseFabric(file *zip.File) (string, error) {
	rc, err := file.Open()
	if err != nil {
		return "", err
	}
	defer rc.Close()

	var mod FabricMod

	if err := json.NewDecoder(rc).Decode(&mod); err != nil {
		return "", err
	}

	return mod.ID, nil
}

func parseQuilt(file *zip.File) (string, error) {
	rc, err := file.Open()
	if err != nil {
		return "", err
	}
	defer rc.Close()

	var mod QuiltMod

	if err := json.NewDecoder(rc).Decode(&mod); err != nil {
		return "", err
	}

	return mod.QuiltLoader.ID, nil
}

func parseForge(file *zip.File) (string, error) {
	rc, err := file.Open()
	if err != nil {
		return "", err
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		return "", err
	}

	var mods ForgeModsToml

	if err := toml.Unmarshal(data, &mods); err != nil {
		return "", err
	}

	id := mods.Mods[0].ModID

	return id, nil
}
