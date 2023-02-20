package main

import (
	"math/rand"
	"sort"
)

type TaggedBundleDef struct {
	Item     int32
	Quantity int32
	Tags     KeyValuePairs
}

type TaggedBundleDefs []TaggedBundleDef

var rng = rand.New(rand.NewSource(0))

func generateItems(defs map[int32]*ItemDef, items TaggedBundleDefs) TaggedBundleDefs {
	var items1, items2 TaggedBundleDefs
	items1 = append(items1, items...)

	any := true
	for any {
		any = false

		for _, item := range items1 {
			def := defs[item.Item]
			switch def.Type {
			case "item", "tag_tool":
				items2 = addMergeItem(items2, item.Item, item.Quantity, item.Tags)
			case "playtimegenerator", "generator":
				any = true

				totalWeight := int64(0)
				for _, option := range def.Bundle {
					totalWeight += int64(option.Quantity)
				}

				for i := int32(0); i < item.Quantity; i++ {
					tags := append(KeyValuePairs(nil), item.Tags...)

					for _, tgid := range def.TagGenerators {
						tgdef := defs[tgid]
						totalTagWeight := int64(0)

						for _, option := range tgdef.TagGeneratorValues {
							totalTagWeight += int64(option.Weight)
						}

						tagWeight := rng.Int63n(totalTagWeight)

						for _, option := range tgdef.TagGeneratorValues {
							tagWeight -= int64(option.Weight)
							if tagWeight < 0 {
								tags = append(tags, KeyValuePair{
									Key:   tgdef.TagGeneratorName,
									Value: option.Value,
								})
								break
							}
						}
					}

					weight := rng.Int63n(totalWeight)
					for _, option := range def.Bundle {
						weight -= int64(option.Quantity)
						if weight < 0 {
							items2 = addMergeItem(items2, option.Item, 1, tags)
							break
						}
					}
				}
			case "bundle":
				any = true

				for _, b := range def.Bundle {
					items2 = addMergeItem(items2, b.Item, b.Quantity*item.Quantity, append(KeyValuePairs(nil), item.Tags...))
				}
			default:
				panic("unhandled item type: " + def.Type)
			}
		}

		items1, items2 = items2, items1[:0]
	}

	return items1
}

func addMergeItem(items TaggedBundleDefs, id, quantity int32, tags KeyValuePairs) TaggedBundleDefs {
	for i := range items {
		if items[i].Item == id && sameTags(items[i].Tags, tags) {
			items[i].Quantity += quantity

			return items
		}
	}

	return append(items, TaggedBundleDef{
		Item:     id,
		Quantity: quantity,
		Tags:     tags,
	})
}

func sameTags(a, b KeyValuePairs) bool {
	if len(a) != len(b) {
		return false
	}

	// take the easy (but inefficient) way out
	a1 := append(KeyValuePairs(nil), a...)
	b1 := append(KeyValuePairs(nil), b...)
	sort.Slice(a1, func(i, j int) bool {
		if a1[i].Key == a1[j].Key {
			return a1[i].Value < a1[j].Value
		}

		return a1[i].Key < a1[j].Key
	})
	sort.Slice(b1, func(i, j int) bool {
		if b1[i].Key == b1[j].Key {
			return b1[i].Value < b1[j].Value
		}

		return b1[i].Key < b1[j].Key
	})

	for i, kv := range a1 {
		if b1[i] != kv {
			return false
		}
	}

	return true
}
