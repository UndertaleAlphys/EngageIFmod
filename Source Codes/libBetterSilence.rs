// Currently needed because we use these functionality, they'll be removable when the Rust language stabilizes them
#![feature(lazy_cell, ptr_sub_ptr)]

use std::f32::consts::PI;

use engage::{
    battle::BattleInfo,
    gamedata::{
        item::{ItemData, UnitItem, UnitItemList},
        unit::Unit,
        GodData,
    },
};
use skyline::hooks::InlineCtx;
use unity::prelude::OptionalMethod;

/// This is called a proc(edural) macro. You use this to indicate that a function will be used as a hook.
///
/// Pay attention to the argument, offset.
/// This is the address of the start of the function you would like to hook.
/// This address has to be relative to the .text section of the game.
/// If you do not know what any of this means, take the address in Ghidra and remove the starting ``71`` and the zeroes that follow it.
/// Do not forget the 0x indicator, as it denotates that you are providing a hexadecimal value.

const MAGIC: u32 = 6;
const STAFF: u32 = 7;
const MAGIC_WEAPON_FLAG: i32 = 65536;
const ENGAGE_ITEM_FLAG: i32 = 128;
const SILENCE_TYPE: u32 = 11;
trait SilenceTarget {
    fn is_silence_target(&self) -> bool;
}

impl SilenceTarget for ItemData {
    fn is_silence_target(&self) -> bool {
        if self.kind == MAGIC {
            true
        } else if self.kind == STAFF {
            true
        } else if self.flag.value & MAGIC_WEAPON_FLAG != 0 {
            true
        } else {
            false
        }
    }
}

impl SilenceTarget for UnitItem {
    fn is_silence_target(&self) -> bool {
        self.item.is_silence_target()
    }
}

trait EngageItem {
    fn is_engage_item(&self) -> bool;
}

impl EngageItem for ItemData {
    fn is_engage_item(&self) -> bool {
        self.flag.value & ENGAGE_ITEM_FLAG != 0
    }
}

impl EngageItem for UnitItem {
    fn is_engage_item(&self) -> bool {
        self.item.is_engage_item()
    }
}

#[skyline::hook(offset = 0x01A46970, inline)]
pub fn silence_magic_weapon(ctx: &mut InlineCtx) {
    // Guarantee Existence
    let item = unsafe { &*(*ctx.registers[23].x.as_ref() as *const ItemData) };
    if item.is_silence_target() {
        unsafe { *ctx.registers[9].w.as_mut() = 6 };
    }
}

#[skyline::hook(offset = 0x01A46984, inline)]
pub fn silence_engage_attack(ctx: &mut InlineCtx) {
    let skill_flag = unsafe { *ctx.registers[9].w.as_ref() as u32 };
    let skill_flag = skill_flag & !2u32;
    unsafe { *ctx.registers[9].w.as_mut() = skill_flag };
}

trait BattleInfoTrait {
    fn get_side(&self) -> &BattleInfoSide;
}

impl BattleInfoTrait for BattleInfo {
    fn get_side(&self) -> &BattleInfoSide {
        unsafe { battle_info_get_side(self, 0, None) }
    }
}

#[skyline::from_offset(0x01E7F210)]
pub fn battle_info_get_side(this: &BattleInfo, t: i32, method: OptionalMethod) -> &BattleInfoSide;

struct BattleInfoSide {}

impl BattleInfoSide {
    fn get_unit_item(&self) -> &UnitItem {
        unsafe { battle_info_side_get_unit_item(self, None) }
    }
}

#[skyline::from_offset(0x01E74CD0)]
pub fn battle_info_side_get_unit_item(this: &BattleInfoSide, method: OptionalMethod) -> &UnitItem;

struct AIInterferenceSimulator {
    // 0x10 offenseUnit &Unit
    // 0x20: defenseUnit &Unit
    // 0x28: &BattleInfo
    // 0x34: is_not_suitable BYTE
}

impl AIInterferenceSimulator {
    fn get_suitable(&self) -> u8 {
        unsafe {
            let p_self = self as *const AIInterferenceSimulator;
            let p_not_suitable = p_self.byte_add(0x34) as *const u8;
            *p_not_suitable
        }
    }
    fn set_suitable(&self, not_suitable: u8) {
        unsafe {
            let p_self = self as *const AIInterferenceSimulator;
            let p_not_suitable = p_self.byte_add(0x34) as *const u8 as *mut u8;
            *p_not_suitable = not_suitable
        }
    }
    fn get_offense(&self) -> &Unit {
        unsafe {
            let p_self = self as *const AIInterferenceSimulator;
            let p_unit = p_self.byte_add(0x10) as *const *const Unit;
            &**p_unit
        }
    }
    fn get_defense(&self) -> &Unit {
        unsafe {
            let p_self = self as *const AIInterferenceSimulator;
            let p_unit = p_self.byte_add(0x20) as *const *const Unit;
            &**p_unit
        }
    }
    fn get_battle_info(&self) -> &BattleInfo {
        unsafe {
            let p_self = self as *const AIInterferenceSimulator;
            let p_info = p_self.byte_add(0x28) as *const *const BattleInfo;
            &**p_info
        }
    }
}

#[skyline::from_offset(0x01FB4C90)]
pub fn unit_item_list_get_equipped_item(
    unit_item_list: &UnitItemList,
    method: OptionalMethod,
) -> Option<&UnitItem>;

#[skyline::from_offset(0x01FB47E0)]
pub fn unit_item_list_get_equipped_index(
    unit_item_list: &UnitItemList,
    method: OptionalMethod,
) -> i32;

trait ItemListTrait {
    fn get_equipped_item(&self) -> Option<&UnitItem>;
    fn get_equipped_index(&self) -> Option<i32>;
}

impl ItemListTrait for UnitItemList {
    fn get_equipped_item(&self) -> Option<&UnitItem> {
        unsafe { unit_item_list_get_equipped_item(self, None) }
    }
    fn get_equipped_index(&self) -> Option<i32> {
        let idx = unsafe { unit_item_list_get_equipped_index(self, None) };
        if idx == -1 {
            None
        } else {
            Some(idx)
        }
    }
}

#[skyline::hook(offset = 0x01930740)]
pub fn interference_cal_score(ai: &mut AIInterferenceSimulator, method: OptionalMethod) {
    call_original!(ai, method);
    let rod_equipped = ai.get_battle_info().get_side().get_unit_item();
    if rod_equipped.item.usetype != SILENCE_TYPE {
        return;
    }
    // else {
    //     println!("Silence");
    // }
    let defense = ai.get_defense();
    let item_list = &defense.item_list;
    let mut not_suitable = ai.get_suitable();
    // println!("Weapon Lists: ");
    for item_idx in 0..item_list.unit_items.len() as i32 {
        let unit_item = item_list.get_item(item_idx);
        if let Some(unit_item) = unit_item {
            // let item_name = unit_item.item.name.to_string();
            // println!("{}", item_name);
            if defense.can_equip(item_idx, true, true)
                || (defense.is_engaging() && unit_item.is_engage_item())
            {
                if unit_item.is_silence_target() {
                    // println!("All conditions satisfied!");
                    not_suitable = 0;
                    break;
                }
            }
        }
    }
    ai.set_suitable(not_suitable);
}

/// The internal name of your plugin. This will show up in crash logs. Make it 8 characters long at max.
#[skyline::main(name = "BttrSlnc")]
pub fn main() {
    // Install a panic handler for your plugin, allowing you to customize what to do if there's an issue in your code.
    std::panic::set_hook(Box::new(|info| {
        let location = info.location().unwrap();

        // Some magic thing to turn what was provided to the panic into a string. Don't mind it too much.
        // The message will be stored in the msg variable for you to use.
        let msg = match info.payload().downcast_ref::<&'static str>() {
            Some(s) => *s,
            None => match info.payload().downcast_ref::<String>() {
                Some(s) => &s[..],
                None => "Box<Any>",
            },
        };

        // This creates a new String with a message of your choice, writing the location of the panic and its message inside of it.
        // Note the \0 at the end. This is needed because show_error is a C function and expects a C string.
        // This is actually just a result of bad old code and shouldn't be necessary most of the time.
        let err_msg = format!(
            "Custom plugin has panicked at '{}' with the following message:\n{}\0",
            location, msg
        );

        // We call the native Error dialog of the Nintendo Switch with this convenient method.
        // The error code is set to 69 because we do need a value, while the first message displays in the popup and the second shows up when pressing Details.
        skyline::error::show_error(
            69,
            "Custom plugin has panicked! Please open the details and send a screenshot to the developer, then close the game.\n\0",
            err_msg.as_str(),
        );
    }));

    // This is what you call to install your hook(s).
    // If you do not install your hook(s), they will just not execute and nothing will be done with them.
    // It is common to install then in ``main`` but nothing stops you from only installing a hook if some conditions are fulfilled.
    // Do keep in mind that hooks cannot currently be uninstalled, so proceed accordingly.
    //
    // A ``install_hooks!`` variant exists to let you install multiple hooks at once if separated by a comma.
    skyline::install_hooks!(
        silence_magic_weapon,
        silence_engage_attack,
        interference_cal_score,
    );
}
