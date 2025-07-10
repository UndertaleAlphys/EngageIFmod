// Currently needed because we use these functionality, they'll be removable when the Rust language stabilizes them
#![feature(lazy_cell, ptr_sub_ptr)]

use std::default;

use engage::{
    gamedata::{
        item::ItemData,
        skill::{SkillArray, SkillData},
        unit::Unit,
    },
    unitpool::UnitPool,
};
use unity::prelude::*;
/// This is called a proc(edural) macro. You use this to indicate that a function will be used as a hook.
///
/// Pay attention to the argument, offset.
/// This is the address of the start of the function you would like to hook.
/// This address has to be relative to the .text section of the game.
/// If you do not know what any of this means, take the address in Ghidra and remove the starting ``71`` and the zeroes that follow it.
/// Do not forget the 0x indicator, as it denotates that you are providing a hexadecimal value.

#[unity::from_offset("App", "Unit", "get_EngageCount")]
pub fn unit_get_engage_count(this: &Unit, method_info: OptionalMethod) -> i32;

#[unity::from_offset("App", "Unit", "GetEngageCountLimit")]
pub fn unit_get_engage_count_limit(this: &Unit, method_info: OptionalMethod) -> i32;

#[unity::from_offset("App", "Unit", "GetMaxHp")]
pub fn unit_get_max_hp(this: &Unit, method_info: OptionalMethod) -> i32;

#[unity::from_offset("App", "Unit", "get_Hp")]
pub fn unit_get_hp(this: &Unit, method_info: OptionalMethod) -> i32;

trait UnitGetter {
    fn get_engage_count(&self) -> i32;
    fn get_engage_count_limit(&self) -> i32;
    fn get_max_hp(&self) -> i32;
    fn get_hp(&self) -> i32;
}

impl UnitGetter for Unit {
    fn get_engage_count(&self) -> i32 {
        unsafe { unit_get_engage_count(self, None) }
    }
    fn get_engage_count_limit(&self) -> i32 {
        unsafe { unit_get_engage_count_limit(self, None) }
    }
    fn get_max_hp(&self) -> i32 {
        unsafe { unit_get_max_hp(self, None) }
    }
    fn get_hp(&self) -> i32 {
        unsafe { unit_get_hp(self, None) }
    }
}

trait UnitCalculator {
    fn can_charge_engage_count(&self) -> bool;
}

impl UnitCalculator for Unit {
    fn can_charge_engage_count(&self) -> bool {
        self.get_engage_count() < self.get_engage_count_limit() && !self.is_engaging()
    }
}

#[skyline::hook(offset = 0x01DEA700)]
pub fn can_item_use(
    unit: Option<&Unit>,
    item: Option<&ItemData>,
    target: Option<&Unit>,
    use_type: i32,
    skills: Option<&SkillArray>,
    method_info: OptionalMethod,
) -> bool {
    let o_result = call_original!(unit, item, target, use_type, skills, method_info);
    if unit.is_none() || item.is_none() || target.is_none() || use_type == 0 {
        o_result
    } else {
        let unit = unit.unwrap();
        let item = item.unwrap();
        let target = target.unwrap();
        match item.usetype {
            2 => {
                let skills = item.get_equip_skills();
                if skills.find_sid(Il2CppString::new("SID_祝福")).is_some()
                    && target.can_charge_engage_count()
                {
                    true
                } else {
                    o_result
                }
            }
            3 => {
                let skills = item.get_equip_skills();
                if skills.find_sid(Il2CppString::new("SID_祝福")).is_some()
                    && target.can_charge_engage_count()
                {
                    true
                } else {
                    o_result
                }
            }
            _ => o_result,
        }
    }
}

/// The internal name of your plugin. This will show up in crash logs. Make it 8 characters long at max.
#[skyline::main(name = "CItemUse")]
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
    skyline::install_hook!(can_item_use);
}
