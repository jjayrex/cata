import * as PresetUtils from '../../core/preset_utils.js';
import { Consumes, Debuffs, Faction, Flask, Food, Glyphs, Potions, RaidBuffs, TristateEffect } from '../../core/proto/common.js';
import {
	AirTotem,
	EarthTotem,
	EnhancementShaman_Options as EnhancementShamanOptions,
	FireTotem,
	ShamanImbue,
	ShamanPrimeGlyph,
	ShamanMajorGlyph,
	ShamanMinorGlyph,
	ShamanShield,
	ShamanSyncType,
	ShamanTotems,
	WaterTotem,
} from '../../core/proto/shaman.js';
import { SavedTalents } from '../../core/proto/ui.js';
import DefaultApl from './apls/default.apl.json';
import P1Gear from './gear_sets/p1.gear.json';
import P2FtGear from './gear_sets/p2_ft.gear.json';
import P2WfGear from './gear_sets/p2_wf.gear.json';
import P3AllianceGear from './gear_sets/p3_alliance.gear.json';
import P3HordeGear from './gear_sets/p3_horde.gear.json';
import P4FtGear from './gear_sets/p4_ft.gear.json';
import P4WfGear from './gear_sets/p4_wf.gear.json';
import PreraidGear from './gear_sets/preraid.gear.json';

// Preset options for this spec.
// Eventually we will import these values for the raid sim too, so its good to
// keep them in a separate file.

export const PRERAID_PRESET = PresetUtils.makePresetGear('Preraid Preset', PreraidGear);
export const P1_PRESET = PresetUtils.makePresetGear('P1 Preset', P1Gear);
export const P2_PRESET_FT = PresetUtils.makePresetGear('P2 Preset FT', P2FtGear);
export const P2_PRESET_WF = PresetUtils.makePresetGear('P2 Preset WF', P2WfGear);
export const P3_PRESET_ALLIANCE = PresetUtils.makePresetGear('P3 Preset [A]', P3AllianceGear, { faction: Faction.Alliance });
export const P3_PRESET_HORDE = PresetUtils.makePresetGear('P3 Preset [H]', P3HordeGear, { faction: Faction.Horde });
export const P4_PRESET_FT = PresetUtils.makePresetGear('P4 Preset FT', P4FtGear);
export const P4_PRESET_WF = PresetUtils.makePresetGear('P4 Preset WF', P4WfGear);

export const ROTATION_PRESET_DEFAULT = PresetUtils.makePresetAPLRotation('Default', DefaultApl);

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/cata/talent-calc and copy the numbers in the url.
export const StandardTalents = {
	name: 'Standard',
	data: SavedTalents.create({
		talentsString: '3020003-2333310013003012321',
		glyphs: Glyphs.create({
			prime1: ShamanPrimeGlyph.GlyphOfLavaLash,
			prime2: ShamanPrimeGlyph.GlyphOfStormstrike,
			prime3: ShamanPrimeGlyph.GlyphOfWindfuryWeapon,
			major1: ShamanMajorGlyph.GlyphOfLightningShield,
			major2: ShamanMajorGlyph.GlyphOfChainLightning,
			major3: ShamanMajorGlyph.GlyphOfFireNova,
			minor1: ShamanMinorGlyph.GlyphOfWaterWalking,
			minor2: ShamanMinorGlyph.GlyphOfRenewedLife,
			minor3: ShamanMinorGlyph.GlyphOfTheArcticWolf,
		}),
	}),
};

export const DefaultOptions = EnhancementShamanOptions.create({
	classOptions: {
		shield: ShamanShield.LightningShield,
		totems: ShamanTotems.create({
			earth: EarthTotem.StrengthOfEarthTotem,
			fire: FireTotem.MagmaTotem,
			water: WaterTotem.ManaSpringTotem,
			air: AirTotem.WindfuryTotem,
		}),
		imbueMh: ShamanImbue.WindfuryWeapon,
	},
	imbueOh: ShamanImbue.FlametongueWeapon,
	syncType: ShamanSyncType.Auto,
});

export const DefaultConsumes = Consumes.create({
	defaultPotion: Potions.PotionOfSpeed,
	flask: Flask.FlaskOfEndlessRage,
	food: Food.FoodFishFeast,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	arcaneBrilliance: true,
	leaderOfThePack: TristateEffect.TristateEffectImproved,
	totemicWrath: true,
	wrathOfAirTotem: true,
	sanctifiedRetribution: true,
	divineSpirit: true,
	battleShout: true,
	demonicPactSp: 500,
});

export const DefaultDebuffs = Debuffs.create({
	bloodFrenzy: true,
	sunderArmor: true,
	curseOfWeakness: TristateEffect.TristateEffectRegular,
	curseOfElements: true,
	faerieFire: TristateEffect.TristateEffectImproved,
	judgementOfWisdom: true,
	misery: true,
	shadowMastery: true,
});
