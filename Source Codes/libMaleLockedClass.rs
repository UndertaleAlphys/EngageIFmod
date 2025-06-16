// Currently needed because we use these functionality, they'll be removable when the Rust language stabilizes them
#![feature(lazy_cell, ptr_sub_ptr)]

use engage::gamedata::skill::SkillArrayEntityList;
use engage::gamedata::unit::Unit;
use engage::gamedata::unit::UnitEdit;
use engage::gamedata::JobData;
use skyline::hooks::InlineCtx;
// / This is called a proc(edural) macro. You use this to indicate that a function will be used as a hook.
// /
// / Pay attention to the argument, offset.
// / This is the address of the start of the function you would like to hook.
// / This address has to be relative to the .text section of the game.
// / If you do not know what any of this means, take the address in Ghidra and remove the starting ``71`` and the zeroes that follow it.
// / Do not forget the 0x indicator, as it denotates that you are providing a hexadecimal value.

/*
Analysis
ClassChangeCheck
x20 : *Unit - at least in the scope that we want to edit.
*x19.byte_add(0x10) : *JobData

got the registers we want in ClassChangeCheck().
*/

const MALE_LOCK_FLAG: i32 = 16;
const MALE_GENDER: i32 = 1;

fn class_change_check_get_unit(ctx: &InlineCtx) -> &Unit {
    unsafe { &*(*ctx.registers[20].x.as_ref() as *const Unit) }
}

fn class_change_check_get_job_data(ctx: &InlineCtx) -> &JobData {
    unsafe {
        let x19 = *ctx.registers[19].x.as_ref() as *const *const JobData;
        let job_data_ptr = *x19.byte_add(0x10);
        &*job_data_ptr
    }
}

#[skyline::hook(offset = 0x19C6E40, inline)]
pub fn check_male_gender(ctx: &mut InlineCtx) {
    let unit = class_change_check_get_unit(ctx);
    let job_data = class_change_check_get_job_data(ctx);
    // Class is male only
    let w24_set0 = if (job_data.flag.value & MALE_LOCK_FLAG) != 0 {
        if (unit.edit.is_enabled()) {
            unit.edit.gender != MALE_GENDER
        } else {
            unit.get_person().gender != MALE_GENDER
        }
    } else {
        false
    };
    if w24_set0 {
        unsafe {
            *ctx.registers[24].w.as_mut() = 0;
        }
    }
}

/// The internal name of your plugin. This will show up in crash logs. Make it 8 characters long at max.
#[skyline::main(name = "MLckClz")]
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
    skyline::install_hooks!(check_male_gender,);
}
