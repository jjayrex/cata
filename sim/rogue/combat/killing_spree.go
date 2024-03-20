package combat

import (
	"math"
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (comRogue *CombatRogue) registerKillingSpreeSpell() {
	mhWeaponSwing := comRogue.GetOrRegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 51690, Tag: 1}, // actual spellID is 57841
		SpellSchool:      core.SpellSchoolPhysical,
		ProcMask:         core.ProcMaskMeleeMHSpecial,
		Flags:            core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,
		DamageMultiplier: 1,
		CritMultiplier:   comRogue.MeleeCritMultiplier(false),
		ThreatMultiplier: 1,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 0 +
				spell.Unit.MHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) +
				spell.BonusWeaponDamage()

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialNoBlockDodgeParry)
		},
	})
	ohWeaponSwing := comRogue.GetOrRegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 51690, Tag: 2}, // actual spellID is 57842
		SpellSchool:      core.SpellSchoolPhysical,
		ProcMask:         core.ProcMaskMeleeOHSpecial,
		Flags:            core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,
		DamageMultiplier: 1 * comRogue.DWSMultiplier(),
		CritMultiplier:   comRogue.MeleeCritMultiplier(false),
		ThreatMultiplier: 1,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 0 +
				spell.Unit.OHNormalizedWeaponDamage(sim, spell.MeleeAttackPower()) +
				spell.BonusWeaponDamage()

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialNoBlockDodgeParry)
		},
	})
	comRogue.KillingSpreeAura = comRogue.RegisterAura(core.Aura{
		Label:    "Killing Spree",
		ActionID: core.ActionID{SpellID: 51690},
		Duration: time.Second*2 + 1,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			comRogue.SetGCDTimer(sim, core.NeverExpires)
			comRogue.PseudoStats.DamageDealtMultiplier *= core.TernaryFloat64(comRogue.HasPrimeGlyph(proto.RoguePrimeGlyph_GlyphOfKillingSpree), 1.3, 1.2)
			core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				Period:          time.Millisecond * 500,
				NumTicks:        5,
				TickImmediately: true,
				OnAction: func(s *core.Simulation) {
					targetCount := sim.GetNumTargets()
					target := comRogue.CurrentTarget
					if targetCount > 1 {
						newUnitIndex := int32(math.Ceil(float64(targetCount)*sim.RandomFloat("Killing Spree"))) - 1
						target = sim.GetTargetUnit(newUnitIndex)
					}
					mhWeaponSwing.Cast(sim, target)
					ohWeaponSwing.Cast(sim, target)
				},
			})
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			comRogue.SetGCDTimer(sim, sim.CurrentTime)
			comRogue.PseudoStats.DamageDealtMultiplier /= core.TernaryFloat64(comRogue.HasPrimeGlyph(proto.RoguePrimeGlyph_GlyphOfKillingSpree), 1.3, 1.2)
		},
	})
	comRogue.KillingSpree = comRogue.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 51690},
		Flags:    core.SpellFlagAPL,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    comRogue.NewTimer(),
				Duration: time.Minute * 2,
			},
		},

		ApplyEffects: func(sim *core.Simulation, u *core.Unit, s2 *core.Spell) {
			comRogue.BreakStealth(sim)
			comRogue.KillingSpreeAura.Activate(sim)
		},
	})

	comRogue.AddMajorCooldown(core.MajorCooldown{
		Spell:    comRogue.KillingSpree,
		Type:     core.CooldownTypeDPS,
		Priority: core.CooldownPriorityDefault,
		ShouldActivate: func(sim *core.Simulation, c *core.Character) bool {
			if comRogue.CurrentEnergy() > 40 || comRogue.AdrenalineRushAura.IsActive() {
				return false
			}
			return true
		},
	})
}

func (comRogue *CombatRogue) registerKillingSpreeCD() {
	if !comRogue.Talents.KillingSpree {
		return
	}
	comRogue.registerKillingSpreeSpell()
}
