// Currently needed because we use these functionality, they'll be removable when the Rust language stabilizes them
#![feature(lazy_cell, ptr_sub_ptr)]

use engage::{
    gamedata::{ai::MoveLimitRange, skill::{SkillArray, SkillData}, terrain::TerrainData, unit::Unit, Gamedata},
    map::{
        image::{MapImage, MapImageTerrain},
        terrain::MapTerrain,
    },
    sequence::mapsequence::human::MapSequenceHuman,
    util::get_instance,
};
use skyline::{hooks::InlineCtx, libc::PTHREAD_MUTEX_NORMAL};
use unity::{prelude::*, system::SystemByte};
/// This is called a proc(edural) macro. You use this to indicate that a function will be used as a hook.
///
/// Pay attention to the argument, offset.
/// This is the address of the start of the function you would like to hook.
/// This address has to be relative to the .text section of the game.
/// If you do not know what any of this means, take the address in Ghidra and remove the starting ``71`` and the zeroes that follow it.
/// Do not forget the 0x indicator, as it denotates that you are providing a hexadecimal value.

fn read_move_image(template: i64, x: i32, z: i32) -> i32 {
    let result = unsafe {
        let temp_o = template as *const u64;
        let temp_f = temp_o.byte_add(0x10);
        // Read MoveImage
        let image_o = *temp_f.byte_add(0x38) as *const u64;
        let image_f = image_o.byte_add(0x10);
        let sb_array = *image_f as *const i8;
        let m_items = sb_array.byte_add(0x18);
        *(m_items.add((z * 32 + x + 8) as usize))
    };
    result as i32
}

#[skyline::hook(offset = 0x02C33BF0)]
pub fn get_distance(template: i64, x: i32, z: i32, method: OptionalMethod) -> i32 {
    read_move_image(template, x, z)
}

#[skyline::hook(offset = 0x01A1E8B0)]
pub fn clear_move_distance(unit: Option<&mut Unit>, method: u64) {
    if let Some(unit) = unit {
        unit.move_distance = 0;
    }
}

#[skyline::hook(offset = 0x0268111C, inline)]
pub fn dont_clear_move_distance(ctx: &mut InlineCtx) {
    unsafe { *ctx.registers[0].x.as_mut() = 0 }
}

#[skyline::hook(offset = 0x02680F80)]
pub fn map_sequence_mind_fixed(map_sequence_mind: i64, method: OptionalMethod) {
    call_original!(map_sequence_mind, method);
    let unit = unsafe {
        let map_sequence_mind_o = map_sequence_mind as *const *const Unit;
        let p_unit = *map_sequence_mind_o.byte_add(0x70);
        if p_unit == std::ptr::null() {
            None
        } else {
            let p_unit = p_unit as *mut Unit;
            Some(&mut *p_unit)
        }
    };
    clear_move_distance(unit, 0);
    // if let Some(unit) = unit{
    //     unit.private_skill.remove_sid(Il2CppString::new("SID_Canto_Flag"));
    // }
}

// #[unity::from_offset("App", "Unit", "RemovePrivateSkill")]
// pub fn remove_skill(unit: &Unit, sid: &Il2CppString, method: OptionalMethod) -> bool;
#[skyline::from_offset(0x01A5D430)]
pub fn add_skill(unit: &Unit, skill_data: &SkillData, method: OptionalMethod) -> bool;

static mut LAST_MOVE_POWER: i32 = 0;
#[skyline::hook(offset = 0x02C1E7F0)]
pub fn unit_move_xz(
    template: i64,
    unit: &Unit,
    x: i32,
    z: i32,
    move_power: i32,
    flag: i64,
    weapon_flag: i64,
    method_info: OptionalMethod,
) {
    let real_move_power = if unit.status.value & 0x40000 != 0 {
        if unit.has_sid(Il2CppString::new("SID_再移動＋＋")) {
            if unit.has_sid(Il2CppString::new("SID_Canto_Flag")) {
                unsafe { LAST_MOVE_POWER }
            } else {
                let unit_mov = unit.get_capability(0xA, true);
                let move_power = (unit_mov - unit.move_distance).max(0);
                unsafe { add_skill(unit, SkillData::get(Il2CppString::new("SID_Canto_Flag")).unwrap(), None) };
                unsafe { LAST_MOVE_POWER = move_power };
                move_power
            }
        } else {
            move_power
        }
    } else {
        move_power
    };
    call_original!(
        template,
        unit,
        x,
        z,
        real_move_power,
        flag,
        weapon_flag,
        method_info
    )
}

/// The internal name of your plugin. This will show up in crash logs. Make it 8 characters long at max.
#[skyline::main(name = "CantoSkl")]
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
        get_distance,
        clear_move_distance,
        dont_clear_move_distance,
        map_sequence_mind_fixed,
        unit_move_xz,
    );
}
