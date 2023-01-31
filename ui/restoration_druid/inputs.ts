import { RaidTarget } from '../core/proto/common.js';
import { Spec } from '../core/proto/common.js';
import { NO_TARGET } from '../core/proto_utils/utils.js';
import { ActionId } from '../core/proto_utils/action_id.js';
import { Player } from '../core/player.js';
import { EventID, TypedEvent } from '../core/typed_event.js';

import * as InputHelpers from '../core/components/input_helpers.js';

import {
	RestorationDruid_Options as DruidOptions,
} from '../core/proto/druid.js';


// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const SelfInnervate = InputHelpers.makeSpecOptionsBooleanIconInput<Spec.SpecRestorationDruid>({
	fieldName: 'innervateTarget',
	id: ActionId.fromSpellId(29166),
	extraCssClasses: [
		'within-raid-sim-hide',
	],
	getValue: (player: Player<Spec.SpecRestorationDruid>) => player.getSpecOptions().innervateTarget?.targetIndex != NO_TARGET,
	setValue: (eventID: EventID, player: Player<Spec.SpecRestorationDruid>, newValue: boolean) => {
		const newOptions = player.getSpecOptions();
		newOptions.innervateTarget = RaidTarget.create({
			targetIndex: newValue ? 0 : NO_TARGET,
		});
		player.setSpecOptions(eventID, newOptions);
	},
});

export const RestorationDruidRotationConfig = {
	inputs: [
	],
};
