package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (deathKnight *DeathKnight) registerEmpowerRuneWeaponSpell() {
	actionID := core.ActionID{SpellID: 47568}
	cdTimer := deathKnight.NewTimer()
	cd := time.Minute * 5

	deathKnight.EmpowerRuneWeapon = deathKnight.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: cd,
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			deathKnight.RegenAllRunes(sim)

			amountOfRunicPower := 25.0
			deathKnight.AddRunicPower(sim, amountOfRunicPower, deathKnight.EmpowerRuneWeapon.RunicPowerMetrics())
		},
	})

	// Temp stuff for testing
	if deathKnight.Talents.SummonGargoyle {
		deathKnight.AddMajorCooldown(core.MajorCooldown{
			Spell:    deathKnight.EmpowerRuneWeapon,
			Priority: core.CooldownPriorityDefault,
			Type:     core.CooldownTypeDPS,
			CanActivate: func(sim *core.Simulation, character *core.Character) bool {
				if deathKnight.CurrentBloodRunes() > 0 {
					return false
				}
				if deathKnight.CurrentFrostRunes() > 0 {
					return false
				}
				if deathKnight.CurrentUnholyRunes() > 0 {
					return false
				}
				return deathKnight.CanEmpowerRuneWeapon(sim)
			},
		})
	}
}

func (deathKnight *DeathKnight) CanEmpowerRuneWeapon(sim *core.Simulation) bool {
	return deathKnight.EmpowerRuneWeapon.IsReady(sim)
}

func (deathKnight *DeathKnight) CastEmpowerRuneWeapon(sim *core.Simulation, target *core.Unit) bool {
	if deathKnight.CanEmpowerRuneWeapon(sim) {
		deathKnight.CastEmpowerRuneWeapon(sim, target)
		return true
	}
	return false
}
