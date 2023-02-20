package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func main() {
	defs, err := loadItemDefs()
	if err != nil {
		panic(err)
	}

	if false {
		// lifetime guaranteed rares
		items := generateItems(defs, TaggedBundleDefs{
			{
				// Random Drop Pool Guaranteed Rare
				// (drops after playing 100 hours, max once per 90 days, max 5 times lifetime)
				Item:     7021,
				Quantity: 5,
			},
		})

		sortItems(items)

		printItems(defs, items)
	}

	if true {
		// potential daily drops for simulated players
		const (
			daysForSimulation = 30
			dailyPlayerCount  = 5000

			playtime15MinutesWeight  = 10  // 0.25 hours
			playtime30MinutesWeight  = 100 // 0.5 hours
			playtime45MinutesWeight  = 150 // 0.75 hours
			playtime60MinutesWeight  = 150 // 1 hour
			playtime75MinutesWeight  = 10  // 1.25 hours
			playtime90MinutesWeight  = 5   // 1.5 hours
			playtime180MinutesWeight = 1   // 3 hours
			playtime270MinutesWeight = 1   // 4.5 hours
			playtime360MinutesWeight = 1   // 6 hours
			playtime450MinutesWeight = 1   // 7.5 hours

			officerMarineWeight        = 5
			specialWeaponsMarineWeight = 5
			medicMarineWeight          = 6
			techMarineWeight           = 8

			missionFallbackWeight                   = 0
			missionWorkshopCompetitionWeight        = 1
			missionWorkshopCampaignAWeight          = 1
			missionWorkshopCampaignBWeight          = 1
			missionWorkshopBonusAWeight             = 1
			missionWorkshopBonusBWeight             = 1
			missionStandaloneOfficialMissionsWeight = 10
			missionEndlessWeight                    = 1
			missionDeathmatchWeight                 = 1
			missionJacobsRestWeight                 = 100
			missionArea9800Weight                   = 10
			missionOperationCleansweepWeight        = 10
			missionResearch7Weight                  = 10
			missionTearsForTarnorWeight             = 10
			missionTilarus5Weight                   = 10
			missionLanasEscapeWeight                = 10
			missionParanoiaWeight                   = 10
			missionNamHumanumWeight                 = 10
			missionBioGenCorporationWeight          = 10
			missionAccident32Weight                 = 40
			missionAdanaxisWeight                   = 40

			playtimeTotalWeight = playtime15MinutesWeight + playtime30MinutesWeight + playtime45MinutesWeight + playtime60MinutesWeight + playtime75MinutesWeight + playtime90MinutesWeight + playtime180MinutesWeight + playtime270MinutesWeight + playtime360MinutesWeight + playtime450MinutesWeight
			totalMarineWeight   = officerMarineWeight + specialWeaponsMarineWeight + medicMarineWeight + techMarineWeight
			totalMissionWeight  = missionFallbackWeight + missionWorkshopCompetitionWeight + missionWorkshopCampaignAWeight + missionWorkshopCampaignBWeight + missionWorkshopBonusAWeight + missionWorkshopBonusBWeight + missionStandaloneOfficialMissionsWeight + missionEndlessWeight + missionDeathmatchWeight + missionJacobsRestWeight + missionArea9800Weight + missionOperationCleansweepWeight + missionResearch7Weight + missionTearsForTarnorWeight + missionTilarus5Weight + missionLanasEscapeWeight + missionParanoiaWeight + missionNamHumanumWeight + missionBioGenCorporationWeight + missionAccident32Weight + missionAdanaxisWeight

			totalNormalDrops       = daysForSimulation * dailyPlayerCount * (5*playtimeTotalWeight - 4*playtime15MinutesWeight + 3*playtime30MinutesWeight - 2*playtime45MinutesWeight - playtime60MinutesWeight) / playtimeTotalWeight
			totalMissionDrops      = totalNormalDrops / 2
			totalMarineDrops       = totalNormalDrops - totalMissionDrops
			totalExtendedFarmDrops = daysForSimulation * dailyPlayerCount * (playtime90MinutesWeight + 2*playtime180MinutesWeight + 3*playtime270MinutesWeight + 4*playtime360MinutesWeight + 5*playtime450MinutesWeight) / playtimeTotalWeight

			totalMarineDropsMissed  = totalMarineDrops - (totalMarineDrops * officerMarineWeight / totalMarineWeight) - (totalMarineDrops * specialWeaponsMarineWeight / totalMarineWeight) - (totalMarineDrops * medicMarineWeight / totalMarineWeight) - (totalMarineDrops * techMarineWeight / totalMarineWeight)
			totalMissionDropsMissed = totalMissionDrops - (totalMissionDrops * missionFallbackWeight / totalMissionWeight) - (totalMissionDrops * missionWorkshopCompetitionWeight / totalMissionWeight) - (totalMissionDrops * missionWorkshopCampaignAWeight / totalMissionWeight) - (totalMissionDrops * missionWorkshopCampaignBWeight / totalMissionWeight) - (totalMissionDrops * missionWorkshopBonusAWeight / totalMissionWeight) - (totalMissionDrops * missionWorkshopBonusBWeight / totalMissionWeight) - (totalMissionDrops * missionStandaloneOfficialMissionsWeight / totalMissionWeight) - (totalMissionDrops * missionEndlessWeight / totalMissionWeight) - (totalMissionDrops * missionDeathmatchWeight / totalMissionWeight) - (totalMissionDrops * missionJacobsRestWeight / totalMissionWeight) - (totalMissionDrops * missionArea9800Weight / totalMissionWeight) - (totalMissionDrops * missionOperationCleansweepWeight / totalMissionWeight) - (totalMissionDrops * missionResearch7Weight / totalMissionWeight) - (totalMissionDrops * missionTearsForTarnorWeight / totalMissionWeight) - (totalMissionDrops * missionTilarus5Weight / totalMissionWeight) - (totalMissionDrops * missionLanasEscapeWeight / totalMissionWeight) - (totalMissionDrops * missionParanoiaWeight / totalMissionWeight) - (totalMissionDrops * missionNamHumanumWeight / totalMissionWeight) - (totalMissionDrops * missionBioGenCorporationWeight / totalMissionWeight) - (totalMissionDrops * missionAccident32Weight / totalMissionWeight) - (totalMissionDrops * missionAdanaxisWeight / totalMissionWeight)
		)

		fmt.Printf("Simulating total drops for %d players playing for %d days...\n\n", dailyPlayerCount, daysForSimulation)

		items := generateItems(defs, TaggedBundleDefs{
			// Unless specified:
			// Drop tables are chosen at the end of a mission.
			// 50% chance for marine class; 50% chance for one of the others.
			// Drops can happen every 15 minutes, up to 5 times per day.
			{
				// Random Drop Pool Fallback
				Item:     7000,
				Quantity: totalMissionDrops*missionFallbackWeight/totalMissionWeight + totalMissionDropsMissed,
			},
			{
				// Random Drop Pool Workshop Competition
				Item:     7001,
				Quantity: totalMissionDrops * missionWorkshopCompetitionWeight / totalMissionWeight,
			},
			{
				// Random Drop Pool Workshop Campaign A
				Item:     7002,
				Quantity: totalMissionDrops * missionWorkshopCampaignAWeight / totalMissionWeight,
			},
			{
				// Random Drop Pool Workshop Campaign B
				Item:     7003,
				Quantity: totalMissionDrops * missionWorkshopCampaignBWeight / totalMissionWeight,
			},
			{
				// Random Drop Pool Workshop Bonus A
				Item:     7004,
				Quantity: totalMissionDrops * missionWorkshopBonusAWeight / totalMissionWeight,
			},
			{
				// Random Drop Pool Workshop Bonus B
				Item:     7005,
				Quantity: totalMissionDrops * missionWorkshopBonusBWeight / totalMissionWeight,
			},
			{
				// Random Drop Pool Standalone Official Missions
				Item:     7006,
				Quantity: totalMissionDrops * missionStandaloneOfficialMissionsWeight / totalMissionWeight,
			},
			{
				// Random Drop Pool Endless
				Item:     7007,
				Quantity: totalMissionDrops * missionEndlessWeight / totalMissionWeight,
			},
			{
				// Random Drop Pool Deathmatch
				Item:     7008,
				Quantity: totalMissionDrops * missionDeathmatchWeight / totalMissionWeight,
			},
			{
				// Random Drop Pool Jacob's Rest
				Item:     7009,
				Quantity: totalMissionDrops * missionJacobsRestWeight / totalMissionWeight,
			},
			{
				// Random Drop Pool Area 9800
				Item:     7010,
				Quantity: totalMissionDrops * missionArea9800Weight / totalMissionWeight,
			},
			{
				// Random Drop Pool Operation Cleansweep
				Item:     7011,
				Quantity: totalMissionDrops * missionOperationCleansweepWeight / totalMissionWeight,
			},
			{
				// Random Drop Pool Research 7
				Item:     7012,
				Quantity: totalMissionDrops * missionResearch7Weight / totalMissionWeight,
			},
			{
				// Random Drop Pool Tears for Tarnor
				Item:     7013,
				Quantity: totalMissionDrops * missionTearsForTarnorWeight / totalMissionWeight,
			},
			{
				// Random Drop Pool Tilarus-5
				Item:     7014,
				Quantity: totalMissionDrops * missionTilarus5Weight / totalMissionWeight,
			},
			{
				// Random Drop Pool Lana's Escape
				Item:     7015,
				Quantity: totalMissionDrops * missionLanasEscapeWeight / totalMissionWeight,
			},
			{
				// Random Drop Pool Paranoia
				Item:     7016,
				Quantity: totalMissionDrops * missionParanoiaWeight / totalMissionWeight,
			},
			{
				// Random Drop Pool Nam Humanum
				Item:     7017,
				Quantity: totalMissionDrops * missionNamHumanumWeight / totalMissionWeight,
			},
			{
				// Random Drop Pool BioGen Corporation
				Item:     7018,
				Quantity: totalMissionDrops * missionBioGenCorporationWeight / totalMissionWeight,
			},
			{
				// Random Drop Pool Accident 32
				Item:     7019,
				Quantity: totalMissionDrops * missionAccident32Weight / totalMissionWeight,
			},
			{
				// Random Drop Pool Adanaxis
				Item:     7020,
				Quantity: totalMissionDrops * missionAdanaxisWeight / totalMissionWeight,
			},
			{
				// Random Drop Pool Marine Class Officer
				Item:     7025,
				Quantity: totalMarineDrops * officerMarineWeight / totalMarineWeight,
			},
			{
				// Random Drop Pool Marine Class Special Weapons
				Item:     7026,
				Quantity: totalMarineDrops*specialWeaponsMarineWeight/totalMarineWeight + totalMarineDropsMissed,
			},
			{
				// Random Drop Pool Marine Class Medic
				Item:     7027,
				Quantity: totalMarineDrops * medicMarineWeight / totalMarineWeight,
			},
			{
				// Random Drop Pool Marine Class Tech
				Item:     7028,
				Quantity: totalMarineDrops * techMarineWeight / totalMarineWeight,
			},
			{
				// Random Drop Pool Extended Farm
				// (every 90 minutes, no daily limit)
				Item:     7029,
				Quantity: totalExtendedFarmDrops,
			},
		})

		sortItems(items)

		printItems(defs, items)
	}
}

func sortItems(items TaggedBundleDefs) {
	sort.Slice(items, func(i, j int) bool {
		if items[i].Quantity != items[j].Quantity {
			return items[i].Quantity > items[j].Quantity
		}

		if items[i].Item != items[j].Item {
			return items[i].Item < items[j].Item
		}

		return fmt.Sprint(items[i].Tags) < fmt.Sprint(items[j].Tags)
	})
}

func printItems(defs map[int32]*ItemDef, items TaggedBundleDefs) {
	for _, item := range items {
		def := defs[item.Item]
		name := def.NameEnglish
		if name == "" {
			name = def.Name
		}
		if name == "" {
			name = fmt.Sprintf("UNNAMED ITEM #%d", item.Item)
		}

		displayType := def.DisplayTypeEnglish
		if displayType == "" {
			displayType = def.DisplayType
		}
		if displayType == "" {
			displayType = "<no display type>"
		}

		allTags := append(append(KeyValuePairs(nil), def.Tags...), item.Tags...)

		tagStrings := make([]string, len(allTags))
		for i, kv := range allTags {
			value := kv.Value
			if def.AccessoryTag == kv.Key {
				id, err := strconv.ParseInt(kv.Value, 10, 32)
				if err == nil {
					value = defs[int32(id)].NameEnglish
					if value == "" {
						value = defs[int32(id)].Name
					}
					if value == "" {
						value = kv.Value
					}
				}
			}

			tagStrings[i] = kv.Key + ":" + value
		}
		tags := strings.Join(tagStrings, ";")

		uniqueStar := ""
		if len(item.Tags) != 0 {
			uniqueStar = "*"
		}

		fmt.Printf("%dx\t\t#%d%s %s (%s)\t\t%s\n", item.Quantity, item.Item, uniqueStar, name, displayType, tags)
	}
}
