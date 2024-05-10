package hunter

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

type PetAbilityType int

// Pet AI doesn't use abilities immediately, so model this with a 1.2s GCD.
const PetGCD = time.Millisecond * 1200

const (
	Unknown PetAbilityType = iota
	AcidSpit
	Bite
	Claw
	DemoralizingScreech
	FireBreath
	Smack
	Stampede
	CorrosiveSpit
	RoarOfCourage
	TailSpin
)

// These IDs are needed for certain talents.
const BiteSpellID = 17253
const ClawSpellID = 16827
const SmackSpellID = 49966

func (hp *HunterPet) NewPetAbility(abilityType PetAbilityType, isPrimary bool) *core.Spell {
	switch abilityType {

	case Bite:
		return hp.newBite()
	case Claw:
		return hp.newClaw()
	case Smack:
		return hp.newSmack()
	case DemoralizingScreech:
		return hp.newDemoralizingScreech()
		//return nil
	case RoarOfCourage: // Agi/Str Buff
		return hp.newRoarOfCourage()
		//return nil
	case FireBreath: // 8% Spell Damage Taken
		return hp.newPetDebuff(PetDebuffSpellConfig{
			SpellID:    24844,
			CD:         time.Second * 30,
			School:     core.SpellSchoolFire,
			DebuffAura: core.FireBreathDebuff,
		})
	case AcidSpit: // 4% Phys Dmg Taken
		return hp.newPetDebuff(PetDebuffSpellConfig{
			SpellID:    55749,
			CD:         time.Second * 10,
			School:     core.SpellSchoolNature,
			DebuffAura: core.AcidSpitAura,
		})
	case CorrosiveSpit: // 10% Armor Reduction
		return hp.newPetDebuff(PetDebuffSpellConfig{
			SpellID:    35387,
			CD:         time.Second * 6,
			School:     core.SpellSchoolNature,
			DebuffAura: core.CorrosiveSpitAura,
		})
	case Stampede: // Bleed Damage 30%
		return hp.newPetDebuff(PetDebuffSpellConfig{
			SpellID:    35290,
			CD:         time.Second * 10,
			School:     core.SpellSchoolPhysical,
			DebuffAura: core.StampedeAura,
		})
	case Unknown:
		return nil
	default:
		panic("Invalid pet ability type")
	}
}

type PetDebuffSpellConfig struct {
	DebuffAura func(*core.Unit) *core.Aura
	SpellID    int32
	School     core.SpellSchool
	GCD        time.Duration
	CD         time.Duration

	OnSpellHitDealt func(*core.Simulation, *core.Spell, *core.SpellResult)
}

func (hp *HunterPet) newPetDebuff(config PetDebuffSpellConfig) *core.Spell {
	auraArray := hp.NewEnemyAuraArray(core.CurseOfElementsAura)
	return hp.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: config.SpellID},
		SpellSchool: config.School, // Adjust the spell school as needed
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagAPL,
		//ClassSpellMask: HunterPetSpellDebuff, // Define or adjust the class spell mask appropriately

		FocusCost: core.FocusCostOptions{
			Cost: 0,
		},
		Cast: core.CastConfig{
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    hp.NewTimer(),
				Duration: hp.hunterOwner.applyLongevity(config.CD),
			},
		},

		ThreatMultiplier: 1,
		FlatThreatBonus:  156,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				auraArray.Get(target).Activate(sim)
			}

			spell.DealOutcome(sim, result)
		},

		RelatedAuras: []core.AuraArray{auraArray},
	})
}

func (hp *HunterPet) newFocusDump(pat PetAbilityType, spellID int32) *core.Spell {
	return hp.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: spellID},
		SpellSchool:    core.SpellSchoolPhysical,
		ProcMask:       core.ProcMaskMeleeMHSpecial,
		ClassSpellMask: HunterPetFocusDump,
		Flags:          core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

		FocusCost: core.FocusCostOptions{
			Cost: 25,
		},
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Duration: time.Millisecond * 3320,
				Timer:    hp.NewTimer(),
			},
			DefaultCast: core.Cast{
				GCD: PetGCD,
			},
			IgnoreHaste: true,
		},
		DamageMultiplierAdditive: 1,
		DamageMultiplier:         1,
		CritMultiplier:           2,
		ThreatMultiplier:         1,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(132, 188) + (spell.MeleeAttackPower() * 0.2)

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
		},
	})
}

func (hp *HunterPet) newBite() *core.Spell {
	return hp.newFocusDump(Bite, BiteSpellID)
}
func (hp *HunterPet) newClaw() *core.Spell {
	return hp.newFocusDump(Claw, ClawSpellID)
}
func (hp *HunterPet) newSmack() *core.Spell {
	return hp.newFocusDump(Smack, SmackSpellID)
}

type PetSpecialAbilityConfig struct {
	Type    PetAbilityType
	SpellID int32
	School  core.SpellSchool
	GCD     time.Duration
	CD      time.Duration

	OnSpellHitDealt func(*core.Simulation, *core.Spell, *core.SpellResult)
}

func (hp *HunterPet) newSpecialAbility(config PetSpecialAbilityConfig) *core.Spell {
	var flags core.SpellFlag
	var applyEffects core.ApplySpellResults
	var procMask core.ProcMask
	onSpellHitDealt := config.OnSpellHitDealt
	if config.School == core.SpellSchoolPhysical {
		flags = core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage
		procMask = core.ProcMaskSpellDamage
		applyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {

			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMeleeSpecialHitAndCrit)
			if onSpellHitDealt != nil {
				onSpellHitDealt(sim, spell, result)
			}

		}
	} else {
		procMask = core.ProcMaskMeleeMHSpecial
		applyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMagicHitAndCrit)
			if onSpellHitDealt != nil {
				onSpellHitDealt(sim, spell, result)
			}
		}
	}

	return hp.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: config.SpellID},
		SpellSchool: config.School,
		ProcMask:    procMask,
		Flags:       flags,

		DamageMultiplier: 1, //* hp.hunterOwner.markedForDeathMultiplier(),
		CritMultiplier:   2,
		ThreatMultiplier: 1,

		FocusCost: core.FocusCostOptions{
			Cost: 0,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: config.GCD,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    hp.NewTimer(),
				Duration: hp.hunterOwner.applyLongevity(config.CD),
			},
		},
		ApplyEffects: applyEffects,
	})
}

func (hp *HunterPet) newRoarOfCourage() *core.Spell {
	actionID := core.ActionID{SpellID: 24604}
	return hp.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    hp.NewTimer(),
				Duration: hp.hunterOwner.applyLongevity(time.Second * 45),
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hp.IsEnabled()
		},
		ApplyEffects: func(sim *core.Simulation, unit *core.Unit, _ *core.Spell) {
			core.RoarOfCourageAura(unit).Activate(sim)
		},
	})
}
func (hp *HunterPet) newDemoralizingScreech() *core.Spell {
	debuffs := hp.NewEnemyAuraArray(core.DemoralizingScreechAura)

	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type: DemoralizingScreech,

		GCD:     PetGCD,
		CD:      time.Second * 10,
		SpellID: 55487,
		School:  core.SpellSchoolPhysical,
		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					debuffs.Get(aoeTarget).Activate(sim)
				}
			}
		},
	})
}

func (hp *HunterPet) newStampede() *core.Spell {
	debuffs := hp.NewEnemyAuraArray(core.StampedeAura)
	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    Stampede,
		CD:      time.Second * 60,
		SpellID: 57386,
		School:  core.SpellSchoolPhysical,
		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() {
				debuffs.Get(result.Target).Activate(sim)
			}
		},
	})
}
