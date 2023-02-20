package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

func loadItemDefs() (map[int32]*ItemDef, error) {
	names, err := filepath.Glob("item-schema-*.json")
	if err != nil {
		return nil, err
	}

	defs := make(map[int32]*ItemDef)

	for _, name := range names {
		err = loadItemDefsFromFile(defs, name)
		if err != nil {
			return nil, err
		}
	}

	return defs, nil
}

func loadItemDefsFromFile(defs map[int32]*ItemDef, name string) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()

	var data struct {
		AppID int32      `json:"appid"`
		Items []*ItemDef `json:"items"`

		// game-specific fields; these will vary per game
		TranslatorNote string `json:"translator_note"`
	}

	dec := json.NewDecoder(f)
	dec.UseNumber()
	dec.DisallowUnknownFields()

	err = dec.Decode(&data)
	if err != nil {
		return err
	}

	for _, item := range data.Items {
		if _, ok := defs[item.ID]; ok {
			return fmt.Errorf("duplicate item id %d", item.ID)
		}

		defs[item.ID] = item
	}

	return nil
}

type ItemDef struct {
	ID   int32  `json:"itemdefid"`
	Type string `json:"type"`

	// not doing multi-language localized versions of these for this toy implementation
	Name               string `json:"name"`
	NameEnglish        string `json:"name_english"`
	Description        string `json:"description"`
	DescriptionEnglish string `json:"description_english"`
	DisplayType        string `json:"display_type"`
	DisplayTypeEnglish string `json:"display_type_english"`

	IconURL    string `json:"icon_url"`
	Tradable   *bool  `json:"tradable"`
	Marketable *bool  `json:"marketable"`
	AutoStack  *bool  `json:"auto_stack"`

	DropInterval  int32 `json:"drop_interval"`
	UseDropWindow *bool `json:"use_drop_window"`
	DropWindow    int32 `json:"drop_window"`
	UseDropLimit  *bool `json:"use_drop_limit"`
	DropLimit     int32 `json:"drop_limit"`

	Bundle               BundleDefs       `json:"bundle"`
	Tags                 KeyValuePairs    `json:"tags"`
	AllowedTagsFromTools KeyValuePairs    `json:"allowed_tags_from_tools"`
	AccessoryTag         string           `json:"accessory_tag"`
	Exchange             string           `json:"exchange"` // TODO
	TagGenerators        IDList           `json:"tag_generators"`
	TagGeneratorName     string           `json:"tag_generator_name"`
	TagGeneratorValues   ValueWeightPairs `json:"tag_generator_values"`

	// game-specific fields; these will vary per game
	TranslatorNote              string `json:"translator_note"`
	ItemSlot                    string `json:"item_slot"`
	AccessoryDescription        string `json:"accessory_description"`
	AccessoryDescriptionEnglish string `json:"accessory_description_english"`
}

type KeyValuePair struct {
	Key   string
	Value string
}

func (p *KeyValuePair) UnmarshalText(b []byte) error {
	i := bytes.IndexByte(b, ':')
	if i == len(b)-1 {
		return fmt.Errorf("expected string after %q", b)
	}

	if i == -1 {
		p.Key, p.Value = string(b), ""
	} else {
		p.Key, p.Value = string(b[:i]), string(b[i+1:])
	}

	return nil
}

type KeyValuePairs []KeyValuePair

func (p *KeyValuePairs) UnmarshalText(b []byte) error {
	pairs := bytes.Split(b, []byte{';'})

	*p = make(KeyValuePairs, len(pairs))

	for i, kv := range pairs {
		err := (*p)[i].UnmarshalText(kv)
		if err != nil {
			return err
		}
	}

	return nil
}

type ValueWeightPair struct {
	Value  string
	Weight int32
}

func (p *ValueWeightPair) UnmarshalText(b []byte) error {
	i := bytes.IndexByte(b, ':')
	if i == len(b)-1 {
		return fmt.Errorf("expected number after %q", b)
	}

	if i == -1 {
		p.Weight = 1
	} else {
		x, err := strconv.ParseInt(string(b[i+1:]), 10, 32)
		if err != nil {
			return err
		}

		if x <= 0 {
			return fmt.Errorf("invalid weight: %d", x)
		}

		p.Weight = int32(x)

		b = b[:i]
	}

	p.Value = string(b)

	return nil
}

type ValueWeightPairs []ValueWeightPair

func (p *ValueWeightPairs) UnmarshalText(b []byte) error {
	pairs := bytes.Split(b, []byte{';'})

	*p = make(ValueWeightPairs, len(pairs))

	for i, vw := range pairs {
		err := (*p)[i].UnmarshalText(vw)
		if err != nil {
			return err
		}
	}

	return nil
}

type BundleDef struct {
	Item     int32
	Quantity int32
}

func (d *BundleDef) UnmarshalText(b []byte) error {
	i := bytes.IndexByte(b, 'x')
	if i == -1 {
		d.Quantity = 1
	} else {
		x, err := strconv.ParseInt(string(b[i+1:]), 10, 32)
		if err != nil {
			return err
		}

		if x <= 0 {
			return fmt.Errorf("invalid quantity: %d", x)
		}

		d.Quantity = int32(x)

		b = b[:i]
	}

	x, err := strconv.ParseInt(string(b), 10, 32)
	if err != nil {
		return err
	}

	if x <= 0 || x >= 1000000000 {
		return fmt.Errorf("invalid item id: %d", x)
	}

	d.Item = int32(x)

	return nil
}

type BundleDefs []BundleDef

func (d *BundleDefs) UnmarshalText(b []byte) error {
	defs := bytes.Split(b, []byte{';'})

	*d = make(BundleDefs, len(defs))

	for i, def := range defs {
		err := (*d)[i].UnmarshalText(def)
		if err != nil {
			return err
		}
	}

	return nil
}

type IDList []int32

func (l *IDList) UnmarshalText(b []byte) error {
	ids := bytes.Split(b, []byte{';'})

	*l = make(IDList, len(ids))

	for i, id := range ids {
		x, err := strconv.ParseInt(string(id), 10, 32)
		if err != nil {
			return err
		}

		if x <= 0 || x >= 1000000000 {
			return fmt.Errorf("invalid item id: %d", x)
		}

		(*l)[i] = int32(x)
	}

	return nil
}
