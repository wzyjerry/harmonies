package main

import (
	"encoding/json"

	"github.com/wzyjerry/harmonies/data"
	"github.com/wzyjerry/harmonies/internal/biz/finder"
	"github.com/wzyjerry/harmonies/pkg/types"
	"google.golang.org/protobuf/encoding/protojson"
	"gopkg.in/yaml.v3"
)

var game types.GameData

func yaml2json(b []byte) ([]byte, error) {
	var data map[string]any
	err := yaml.Unmarshal(b, &data)
	if err != nil {
		return nil, err
	}
	return json.Marshal(data)
}

func init() {
	raw, err := yaml2json(data.Data)
	if err != nil {
		panic(err)
	}
	err = protojson.UnmarshalOptions{DiscardUnknown: true}.Unmarshal(raw, &game)
	if err != nil {
		panic(err)
	}
}

func main() {
	f := finder.New([]*types.Card{game.Animals[0]})
	f.Search()
}
