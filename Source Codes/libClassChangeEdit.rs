// Currently needed because we use these functionality, they'll be removable when the Rust language stabilizes them
#![feature(lazy_cell, ptr_sub_ptr)]
use std::env::consts;

use engage::gamedata::item::ItemData;
use engage::gamedata::skill::SkillData;
use engage::gamedata::unit::Unit;
use engage::gamedata::{job, Gamedata, JobData};
use skyline::hooks::InlineCtx;
use unity::il2cpp;
use unity::prelude::OptionalMethod;
use unity::system::Il2CppString;

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
LevelReset
x19 : *Unit
x20 : *JobData

got the registers we want in ClassChangeCheck().
*/

trait ClassRank {
    fn is_low_class(&self) -> bool;
    fn is_high_class(&self) -> bool;
    fn is_special_class(&self) -> bool;
}
impl ClassRank for JobData {
    fn is_high_class(&self) -> bool {
        return self.is_high();
    }
    fn is_low_class(&self) -> bool {
        return self.is_low() && self.max_level == 20;
    }
    fn is_special_class(&self) -> bool {
        return self.is_low() && self.max_level != 20;
    }
}
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

fn level_reset_get_unit(ctx: &InlineCtx) -> &Unit {
    unsafe { &*(*ctx.registers[19].x.as_ref() as *const Unit) }
}

fn level_reset_get_job_data(ctx: &InlineCtx) -> &JobData {
    unsafe { &*(*ctx.registers[20].x.as_ref() as *const JobData) }
}

fn disallow_high_to_low_chck(ctx: &InlineCtx) -> bool {
    let job_data = class_change_check_get_job_data(ctx);
    let unit = class_change_check_get_unit(ctx);
    let unit_class = unit.get_job();
    if job_data.is_low_class() {
        if unit_class.is_high_class() {
            false
        } else if unit_class.is_special_class() {
            unit.level < 20 && unit.level > 0
        } else {
            //unit_class is low class
            unit.level > 0
        }
    } else {
        unit.level > 0
    }
}

fn get_class_learning_skill(job: &JobData) -> String {
    unsafe {
        let p_job = job as *const JobData as *const *const Il2CppString;
        let p_skill_name = *p_job.byte_add(0x110);
        if p_skill_name == 0 as *const Il2CppString || (*p_skill_name).to_string() == "" {
            return String::new();
        } else {
            (*p_skill_name).to_string()
        }
    }
}

#[skyline::hook(offset = 0x19C6C6C, inline)]
pub fn disallow_high_to_low_impl(ctx: &mut InlineCtx) {
    let result = disallow_high_to_low_chck(ctx);
    unsafe { *ctx.registers[8].w.as_mut() = result as u32 };
}

#[skyline::hook(offset = 0x19C6C34, inline)]
pub fn disallow_high_to_low_disp(ctx: &mut InlineCtx) {
    let result = disallow_high_to_low_chck(ctx);
    let disp_lv = if result { 1 } else { 99 };
    unsafe { *ctx.registers[0].w.as_mut() = disp_lv };
}

#[skyline::hook(offset = 0x19C6AD8, inline)]
pub fn prevent_same_class_change(ctx: &mut InlineCtx) {
    unsafe { *ctx.registers[8].w.as_mut() = 0 }
}

#[skyline::hook(offset = 0x19C69C0, inline)]
pub fn prevent_same_class_change_normal_disp(ctx: &mut InlineCtx) {
    unsafe { *ctx.registers[0].w.as_mut() = 99 }
}

#[skyline::hook(offset = 0x19C6A68, inline)]
pub fn prevent_same_class_change_special_disp(ctx: &mut InlineCtx) {
    unsafe { *ctx.registers[0].w.as_mut() = 99 }
}

#[skyline::hook(offset = 0x1A088E4, inline)]
pub fn disable_level_addition_on_high_class(ctx: &mut InlineCtx) {
    let w19 = unsafe { *ctx.registers[19].w.as_ref() };
    unsafe { *ctx.registers[8].w.as_mut() = w19 };
}

#[skyline::hook(offset = 0x1A3C848, inline)]
pub fn level_reset(ctx: &mut InlineCtx) {
    let unit = level_reset_get_unit(ctx);
    let unit_class = unit.get_job();
    let job_data = level_reset_get_job_data(ctx);
    let reset_level = if job_data.is_high_class() {
        if unit_class.is_high_class() {
            unit.level
        } else if unit_class.is_special_class() && unit.level > 20 {
            unit.level
        } else {
            20
        }
    } else if job_data.is_low_class() {
        if unit_class.is_high_class() {
            1
        } else if unit_class.is_special_class() {
            if unit.level <= 20 {
                unit.level
            } else {
                1
            }
        } else {
            // low class
            unit.level
        }
    } else {
        // Special class
        unit.level
    };
    unsafe { *ctx.registers[25].w.as_mut() = reset_level as u32 };
}

#[skyline::hook(offset = 0x1A3C7B0)]
pub fn class_change(
    this: &Unit,
    job_data: &JobData,
    item_data: &ItemData,
    method_info: OptionalMethod,
) {
    call_original!(this, job_data, item_data, method_info);
    let current_class = this.get_job();
    if this.level >= current_class.max_level {
        if get_class_learning_skill(job_data) != "" {
            this.learn_job_skill(job_data);
        }
    }
}

/// The internal name of your plugin. This will show up in crash logs. Make it 8 characters long at max.
#[skyline::main(name = "ClzCgEd")]
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
        disallow_high_to_low_impl,
        disallow_high_to_low_disp,
        disable_level_addition_on_high_class,
        level_reset,
        prevent_same_class_change,
        prevent_same_class_change_normal_disp,
        prevent_same_class_change_special_disp,
        class_change,
    );
}
