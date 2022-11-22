package main

import (
	"flag"
	"fmt"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	_ "github.com/wowsims/wotlk/sim/encounters" // Needed for preset encounters.
	"github.com/wowsims/wotlk/tools/database"
)

// go run ./tools/database/gen_db

var outDir = flag.String("outDir", "assets", "Path to output directory for writing generated .go files.")

func main() {
	flag.Parse()
	if *outDir == "" {
		panic("outDir flag is required!")
	}

	dbDir := fmt.Sprintf("%s/database", *outDir)
	inputsDir := fmt.Sprintf("%s/db_inputs", *outDir)

	itemTooltips := database.NewWowheadItemTooltipManager(fmt.Sprintf("%s/wowhead_item_tooltips.csv", inputsDir)).Read()
	spellTooltips := database.NewWowheadSpellTooltipManager(fmt.Sprintf("%s/wowhead_spell_tooltips.csv", inputsDir)).Read()

	db := database.NewWowDatabase()
	db.Encounters = core.PresetEncounters

	for _, response := range itemTooltips {
		if response.IsEquippable() {
			db.MergeItem(response.ToItemProto())
		} else if response.IsGem() {
			db.MergeGem(response.ToGemProto())
		}
	}

	db.MergeItems(database.ItemOverrides)
	db.MergeGems(database.GemOverrides)
	db.MergeEnchants(database.EnchantOverrides)

	for _, enchant := range db.Enchants {
		if enchant.ItemId != 0 {
			db.AddItemIcon(enchant.ItemId, itemTooltips)
		}
		if enchant.SpellId != 0 {
			db.AddSpellIcon(enchant.SpellId, spellTooltips)
		}
	}

	for _, itemID := range database.ExtraItemIcons {
		db.AddItemIcon(itemID, itemTooltips)
	}

	ApplyGlobalFilters(db)
	ApplySimmableFilters(db)
	db.WriteBinaryAndJson(fmt.Sprintf("%s/db.bin", dbDir), fmt.Sprintf("%s/db.json", dbDir))
}

// Filters out entities which shouldn't be included anywhere.
func ApplyGlobalFilters(db *database.WowDatabase) {
	db.Items = core.FilterMap(db.Items, func(_ int32, item *proto.UIItem) bool {
		if _, ok := database.ItemDenyList[item.Id]; ok {
			return false
		}

		for _, pattern := range database.DenyListNameRegexes {
			if pattern.MatchString(item.Name) {
				return false
			}
		}
		return true
	})

	db.Gems = core.FilterMap(db.Gems, func(_ int32, gem *proto.UIGem) bool {
		if _, ok := database.GemDenyList[gem.Id]; ok {
			return false
		}

		for _, pattern := range database.DenyListNameRegexes {
			if pattern.MatchString(gem.Name) {
				return false
			}
		}
		return true
	})

	db.ItemIcons = core.FilterMap(db.ItemIcons, func(_ int32, icon *proto.IconData) bool {
		return icon.Name != "" && icon.Icon != ""
	})
	db.SpellIcons = core.FilterMap(db.SpellIcons, func(_ int32, icon *proto.IconData) bool {
		return icon.Name != "" && icon.Icon != ""
	})
}

// Filters out entities which shouldn't be included in the sim.
func ApplySimmableFilters(db *database.WowDatabase) {
	db.Items = core.FilterMap(db.Items, func(_ int32, item *proto.UIItem) bool {
		if _, ok := database.ItemAllowList[item.Id]; ok {
			return true
		}

		if item.Quality < proto.ItemQuality_ItemQualityUncommon {
			return false
		} else if item.Quality > proto.ItemQuality_ItemQualityLegendary {
			return false
		} else if item.Quality < proto.ItemQuality_ItemQualityEpic {
			if item.Ilvl < 145 {
				return false
			}
			if item.Ilvl < 149 && item.SetName == "" {
				return false
			}
		} else {
			// Epic and legendary items might come from classic, so use a lower ilvl threshold.
			if item.Ilvl < 140 {
				return false
			}
		}
		if item.Ilvl == 0 {
			fmt.Printf("Missing ilvl: %s\n", item.Name)
		}

		return true
	})

	db.Gems = core.FilterMap(db.Gems, func(id int32, gem *proto.UIGem) bool {
		if gem.Quality < proto.ItemQuality_ItemQualityUncommon {
			return false
		}
		return true
	})
}
