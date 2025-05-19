package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/wzyjerry/harmonies/pkg/types"
	"google.golang.org/protobuf/encoding/protojson"
	"gopkg.in/yaml.v3"
)

func yaml2json(b []byte) ([]byte, error) {
	var data map[string]any
	err := yaml.Unmarshal(b, &data)
	if err != nil {
		return nil, err
	}
	return json.Marshal(data)
}

func main() {
	var data types.GameData
	file, err := os.ReadFile("../../data/animals.yaml")
	if err != nil {
		panic(err)
	}
	raw, err := yaml2json(file)
	if err != nil {
		panic(err)
	}
	err = protojson.UnmarshalOptions{DiscardUnknown: true}.Unmarshal(raw, &data)
	if err != nil {
		panic(err)
	}
	for _, animal := range data.Animals {
		fmt.Println(PrintCard(animal))
	}
}

func PrintCard(card *types.Card) string {
	var buf bytes.Buffer
	buf.WriteString("[")
	buf.WriteString(card.Name)
	buf.WriteString("]<")
	buf.WriteString(card.Kind.String())
	buf.WriteString(">\nscores: ")
	buf.WriteString(fmt.Sprintf("%v", card.Scores))
	buf.WriteString("\n")
	for _, pattern := range card.Pattern {
		buf.WriteString("  ")
		if pattern.Animal {
			buf.WriteString("Aanimal at ")
		}
		buf.WriteString("<")
		buf.WriteString(pattern.Poi.String())
		buf.WriteString(">")
		buf.WriteString(fmt.Sprintf(" [%d]", pattern.Height))
		if !pattern.Animal {
			buf.WriteString(fmt.Sprintf(" (%d, %d)\n", pattern.DeltaQ, pattern.DeltaR))
		} else {
			buf.WriteString("\n")
		}
	}
	return buf.String()
}
