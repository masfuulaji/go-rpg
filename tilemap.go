package main

import (
	"encoding/json"
	"os"
)

type TilemapLayerJSON struct {
	Data   []int `json:"data"`
	Height int   `json:"height"`
	Width  int   `json:"width"`
}
type TilemapJSON struct {
	Layers []TilemapLayerJSON `json:"layers"`
}

func NewTilemapJSON(filepath string) (*TilemapJSON, error) {
	contents, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	var tilemap TilemapJSON
	err = json.Unmarshal(contents, &tilemap)
	if err != nil {
		return nil, err
	}
	return &tilemap, nil
}
