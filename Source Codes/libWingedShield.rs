// Currently needed because we use these functionality, they'll be removable when the Rust language stabilizes them
#![feature(lazy_cell, ptr_sub_ptr)]

use engage::gamedata::{unit::Unit, terrain::TerrainData};
use skyline::hooks::InlineCtx;
use unity::prelude::*;
/// This is called a proc(edural) macro. You use this to indicate that a function will be used as a hook.
///
/// Pay attention to the argument, offset.
/// This is the address of the start of the function you would like to hook.
/// This address has to be relative to the .text section of the game.
/// If you do not know what any of this means, take the address in Ghidra and remove the starting ``71`` and the zeroes that follow it.
/// Do not forget the 0x indicator, as it denotates that you are providing a hexadecimal value.
trait HasWingedShield {
    fn has_winged_shield(&self) -> bool;
}

impl HasWingedShield for Unit {
    fn has_winged_shield(&self) -> bool {
        self.has_sid(Il2CppString::new("SID_翼盾"))
    }
}

#[skyline::hook(offset = 0x01A34C90)]
pub fn is_terrain_invalid(unit: &Unit, terrain: &TerrainData, method: OptionalMethod) -> bool {
    if unit.has_winged_shield() {
        true
    } else {
        call_original!(unit, terrain, method)
    }
}

#[skyline::hook(offset = 0x01E3BA04, inline)]
pub fn phase_start_damage(ctx: &mut InlineCtx) {
    let unit = unsafe { &*(*ctx.registers[20].x.as_ref() as *const Unit) };
    if unit.has_winged_shield() {
        unsafe { *ctx.registers[0].w.as_mut() = 1 }
    }
}

#[skyline::hook(offset = 0x01E3BBA0, inline)]
pub fn phase_start_heal(ctx: &mut InlineCtx) {
    let unit = unsafe { &*(*ctx.registers[20].x.as_ref() as *const Unit) };
    if unit.has_winged_shield() {
        unsafe { *ctx.registers[0].w.as_mut() = 1 }
    }
}

#[skyline::hook(offset = 0x01E44DDC, inline)]
pub fn map_terrain_single(ctx: &mut InlineCtx) {
    let unit = unsafe { &*(*ctx.registers[20].x.as_ref() as *const Unit) };
    if unit.has_winged_shield() {
        unsafe { *ctx.registers[0].w.as_mut() = 1 }
    }
}

#[skyline::hook(offset = 0x01E7A938, inline)]
pub fn battle_set_terrain_1(ctx: &mut InlineCtx) {
    let unit = unsafe { &*(*ctx.registers[23].x.as_ref() as *const Unit) };
    if unit.has_winged_shield() {
        unsafe { *ctx.registers[0].w.as_mut() = 1 }
    }
}

#[skyline::hook(offset = 0x01E7A99C, inline)]
pub fn battle_set_terrain_2(ctx: &mut InlineCtx) {
    let unit = unsafe { &*(*ctx.registers[19].x.as_ref() as *const Unit) };
    if unit.has_winged_shield() {
        unsafe { *ctx.registers[0].w.as_mut() = 1 }
    }
}

#[skyline::hook(offset = 0x02470958, inline)]
pub fn battle_cal_interrupt(ctx: &mut InlineCtx) {
    let unit = unsafe { &*(*ctx.registers[21].x.as_ref() as *const Unit) };
    if unit.has_winged_shield() {
        unsafe { *ctx.registers[0].w.as_mut() = 1 }
    }
}

/// The internal name of your plugin. This will show up in crash logs. Make it 8 characters long at max.
#[skyline::main(name = "WngShield")]
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
        is_terrain_invalid,
        phase_start_damage,
        phase_start_heal,
        map_terrain_single,
        battle_set_terrain_1,
        battle_set_terrain_2,
        battle_cal_interrupt,
    );
}