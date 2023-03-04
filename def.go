package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

	// not doing any special multi-language handling for this toy implementation
	Name                  string `json:"name"`
	NameBrazilian         string `json:"name_brazilian"`
	NameCzech             string `json:"name_czech"`
	NameDanish            string `json:"name_danish"`
	NameDutch             string `json:"name_dutch"`
	NameEnglish           string `json:"name_english"`
	NameFinnish           string `json:"name_finnish"`
	NameFrench            string `json:"name_french"`
	NameGerman            string `json:"name_german"`
	NameHungarian         string `json:"name_hungarian"`
	NameItalian           string `json:"name_italian"`
	NameJapanese          string `json:"name_japanese"`
	NameKoreanA           string `json:"name_koreana"`
	NameNorwegian         string `json:"name_norwegian"`
	NamePolish            string `json:"name_polish"`
	NamePortuguese        string `json:"name_portuguese"`
	NameRomanian          string `json:"name_romanian"`
	NameRussian           string `json:"name_russian"`
	NameSChinese          string `json:"name_schinese"`
	NameSpanish           string `json:"name_spanish"`
	NameSwedish           string `json:"name_swedish"`
	NameTChinese          string `json:"name_tchinese"`
	NameThai              string `json:"name_thai"`
	NameTurkish           string `json:"name_turkish"`
	NameUkrainian         string `json:"name_ukrainian"`
	Description           string `json:"description"`
	DescriptionBrazilian  string `json:"description_brazilian"`
	DescriptionCzech      string `json:"description_czech"`
	DescriptionDanish     string `json:"description_danish"`
	DescriptionDutch      string `json:"description_dutch"`
	DescriptionEnglish    string `json:"description_english"`
	DescriptionFinnish    string `json:"description_finnish"`
	DescriptionFrench     string `json:"description_french"`
	DescriptionGerman     string `json:"description_german"`
	DescriptionHungarian  string `json:"description_hungarian"`
	DescriptionItalian    string `json:"description_italian"`
	DescriptionJapanese   string `json:"description_japanese"`
	DescriptionKoreanA    string `json:"description_koreana"`
	DescriptionNorwegian  string `json:"description_norwegian"`
	DescriptionPolish     string `json:"description_polish"`
	DescriptionPortuguese string `json:"description_portuguese"`
	DescriptionRomanian   string `json:"description_romanian"`
	DescriptionRussian    string `json:"description_russian"`
	DescriptionSChinese   string `json:"description_schinese"`
	DescriptionSpanish    string `json:"description_spanish"`
	DescriptionSwedish    string `json:"description_swedish"`
	DescriptionTChinese   string `json:"description_tchinese"`
	DescriptionThai       string `json:"description_thai"`
	DescriptionTurkish    string `json:"description_turkish"`
	DescriptionUkrainian  string `json:"description_ukrainian"`
	DisplayType           string `json:"display_type"`
	DisplayTypeEnglish    string `json:"display_type_english"`
	DisplayTypeGerman     string `json:"display_type_german"`
	DisplayTypeItalian    string `json:"display_type_italian"`
	DisplayTypeJapanese   string `json:"display_type_japanese"`
	DisplayTypeRussian    string `json:"display_type_russian"`

	IconURL         string    `json:"icon_url"`
	NameColor       *HexColor `json:"name_color"`
	BackgroundColor *HexColor `json:"background_color"`
	Tradable        bool      `json:"tradable"`
	Marketable      bool      `json:"marketable"`
	AutoStack       bool      `json:"auto_stack"`

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
	TranslatorNote                 string     `json:"translator_note"`
	ItemSlot                       string     `json:"item_slot"`
	CompressedDynamicProps         StringList `json:"compressed_dynamic_props"`
	AfterDescription               string     `json:"after_description"`
	AccessoryDescription           string     `json:"accessory_description"`
	AccessoryDescriptionBrazilian  string     `json:"accessory_description_brazilian"`
	AccessoryDescriptionCzech      string     `json:"accessory_description_czech"`
	AccessoryDescriptionDanish     string     `json:"accessory_description_danish"`
	AccessoryDescriptionDutch      string     `json:"accessory_description_dutch"`
	AccessoryDescriptionEnglish    string     `json:"accessory_description_english"`
	AccessoryDescriptionFinnish    string     `json:"accessory_description_finnish"`
	AccessoryDescriptionFrench     string     `json:"accessory_description_french"`
	AccessoryDescriptionGerman     string     `json:"accessory_description_german"`
	AccessoryDescriptionHungarian  string     `json:"accessory_description_hungarian"`
	AccessoryDescriptionItalian    string     `json:"accessory_description_italian"`
	AccessoryDescriptionJapanese   string     `json:"accessory_description_japanese"`
	AccessoryDescriptionKoreanA    string     `json:"accessory_description_koreana"`
	AccessoryDescriptionNorwegian  string     `json:"accessory_description_norwegian"`
	AccessoryDescriptionPolish     string     `json:"accessory_description_polish"`
	AccessoryDescriptionPortuguese string     `json:"accessory_description_portuguese"`
	AccessoryDescriptionRomanian   string     `json:"accessory_description_romanian"`
	AccessoryDescriptionRussian    string     `json:"accessory_description_russian"`
	AccessoryDescriptionSChinese   string     `json:"accessory_description_schinese"`
	AccessoryDescriptionSpanish    string     `json:"accessory_description_spanish"`
	AccessoryDescriptionSwedish    string     `json:"accessory_description_swedish"`
	AccessoryDescriptionTChinese   string     `json:"accessory_description_tchinese"`
	AccessoryDescriptionThai       string     `json:"accessory_description_thai"`
	AccessoryDescriptionTurkish    string     `json:"accessory_description_turkish"`
	AccessoryDescriptionUkrainian  string     `json:"accessory_description_ukrainian"`
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

type StringList []string

func (l *StringList) UnmarshalText(b []byte) error {
	if len(b) == 0 {
		*l = nil

		return nil
	}

	*l = strings.Split(string(b), ";")

	return nil
}

type HexColor struct {
	R, G, B uint8
}

func (c *HexColor) UnmarshalText(b []byte) error {
	n, err := fmt.Sscanf(string(b), "%02x%02x%02x", &c.R, &c.G, &c.B)
	if err == nil && n != 3 {
		return fmt.Errorf("invalid hex color: %q", b)
	}

	return err
}
