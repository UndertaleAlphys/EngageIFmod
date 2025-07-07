// Currently needed because we use these functionality, they'll be removable when the Rust language stabilizes them
#![feature(lazy_cell, ptr_sub_ptr)]

use skyline::hooks::InlineCtx;

/*
Analysis:
0: MaxHP
1: Str
2: Dex
3: Spd
4: Lck
5: Def
6: Mag
7: Res
8: Bld
9: Sight
A: Mov
*/

#[skyline::hook(offset = 0x1A1D7D8, inline)]
pub fn bound_max_hp_emblem(ctx: &mut InlineCtx) {
    unsafe {
        *ctx.registers[1].w.as_mut() = i8::MIN as u32;
        *ctx.registers[2].w.as_mut() = i8::MAX as u32;
    }
}

#[skyline::hook(offset = 0x1A1D808, inline)]
pub fn bound_str_emblem(ctx: &mut InlineCtx) {
    unsafe {
        *ctx.registers[1].w.as_mut() = i8::MIN as u32;
        *ctx.registers[2].w.as_mut() = i8::MAX as u32;
    }
}

#[skyline::hook(offset = 0x1A1D83C, inline)]
pub fn bound_dex_emblem(ctx: &mut InlineCtx) {
    unsafe {
        *ctx.registers[1].w.as_mut() = i8::MIN as u32;
        *ctx.registers[2].w.as_mut() = i8::MAX as u32;
    }
}

#[skyline::hook(offset = 0x1A1D870, inline)]
pub fn bound_spd_emblem(ctx: &mut InlineCtx) {
    unsafe {
        *ctx.registers[1].w.as_mut() = i8::MIN as u32;
        *ctx.registers[2].w.as_mut() = i8::MAX as u32;
    }
}

#[skyline::hook(offset = 0x1A1D8A4, inline)]
pub fn bound_lck_emblem(ctx: &mut InlineCtx) {
    unsafe {
        *ctx.registers[1].w.as_mut() = i8::MIN as u32;
        *ctx.registers[2].w.as_mut() = i8::MAX as u32;
    }
}

#[skyline::hook(offset = 0x1A1D8DC, inline)]
pub fn bound_def_emblem(ctx: &mut InlineCtx) {
    unsafe {
        *ctx.registers[1].w.as_mut() = i8::MIN as u32;
        *ctx.registers[2].w.as_mut() = i8::MAX as u32;
    }
}

#[skyline::hook(offset = 0x1A1D914, inline)]
pub fn bound_mag_emblem(ctx: &mut InlineCtx) {
    unsafe {
        *ctx.registers[1].w.as_mut() = i8::MIN as u32;
        *ctx.registers[2].w.as_mut() = i8::MAX as u32;
    }
}

#[skyline::hook(offset = 0x1A1D948, inline)]
pub fn bound_res_emblem(ctx: &mut InlineCtx) {
    unsafe {
        *ctx.registers[1].w.as_mut() = i8::MIN as u32;
        *ctx.registers[2].w.as_mut() = i8::MAX as u32;
    }
}

#[skyline::hook(offset = 0x1A1D97C, inline)]
pub fn bound_bld_emblem(ctx: &mut InlineCtx) {
    unsafe {
        *ctx.registers[1].w.as_mut() = i8::MIN as u32;
        *ctx.registers[2].w.as_mut() = i8::MAX as u32;
    }
}

#[skyline::hook(offset = 0x1A1D9B4, inline)]
pub fn bound_sight_emblem(ctx: &mut InlineCtx) {
    unsafe {
        *ctx.registers[1].w.as_mut() = i8::MIN as u32;
        *ctx.registers[2].w.as_mut() = i8::MAX as u32;
    }
}

#[skyline::hook(offset = 0x1A1D9EC, inline)]
pub fn bound_mov_emblem(ctx: &mut InlineCtx) {
    unsafe {
        *ctx.registers[1].w.as_mut() = i8::MIN as u32;
        *ctx.registers[2].w.as_mut() = i8::MAX as u32;
    }
}

#[skyline::hook(offset = 0x1A1C934, inline)]
pub fn bound_max_hp(ctx: &mut InlineCtx) {
    unsafe {
        *ctx.registers[1].w.as_mut() = i8::MIN as u32;
        *ctx.registers[2].w.as_mut() = i8::MAX as u32;
    }
}

#[skyline::hook(offset = 0x1A1C964, inline)]
pub fn bound_str(ctx: &mut InlineCtx) {
    unsafe {
        *ctx.registers[1].w.as_mut() = i8::MIN as u32;
        *ctx.registers[2].w.as_mut() = i8::MAX as u32;
    }
}

#[skyline::hook(offset = 0x1A1C998, inline)]
pub fn bound_dex(ctx: &mut InlineCtx) {
    unsafe {
        *ctx.registers[1].w.as_mut() = i8::MIN as u32;
        *ctx.registers[2].w.as_mut() = i8::MAX as u32;
    }
}

#[skyline::hook(offset = 0x1A1C9CC, inline)]
pub fn bound_spd(ctx: &mut InlineCtx) {
    unsafe {
        *ctx.registers[1].w.as_mut() = i8::MIN as u32;
        *ctx.registers[2].w.as_mut() = i8::MAX as u32;
    }
}

#[skyline::hook(offset = 0x1A1CA00, inline)]
pub fn bound_lck(ctx: &mut InlineCtx) {
    unsafe {
        *ctx.registers[1].w.as_mut() = i8::MIN as u32;
        *ctx.registers[2].w.as_mut() = i8::MAX as u32;
    }
}

#[skyline::hook(offset = 0x1A1CA38, inline)]
pub fn bound_def(ctx: &mut InlineCtx) {
    unsafe {
        *ctx.registers[1].w.as_mut() = i8::MIN as u32;
        *ctx.registers[2].w.as_mut() = i8::MAX as u32;
    }
}

#[skyline::hook(offset = 0x1A1CA70, inline)]
pub fn bound_mag(ctx: &mut InlineCtx) {
    unsafe {
        *ctx.registers[1].w.as_mut() = i8::MIN as u32;
        *ctx.registers[2].w.as_mut() = i8::MAX as u32;
    }
}

#[skyline::hook(offset = 0x1A1CAA4, inline)]
pub fn bound_res(ctx: &mut InlineCtx) {
    unsafe {
        *ctx.registers[1].w.as_mut() = i8::MIN as u32;
        *ctx.registers[2].w.as_mut() = i8::MAX as u32;
    }
}

#[skyline::hook(offset = 0x1A1CAD8, inline)]
pub fn bound_bld(ctx: &mut InlineCtx) {
    unsafe {
        *ctx.registers[1].w.as_mut() = i8::MIN as u32;
        *ctx.registers[2].w.as_mut() = i8::MAX as u32;
    }
}

#[skyline::hook(offset = 0x1A1CB10, inline)]
pub fn bound_sight(ctx: &mut InlineCtx) {
    unsafe {
        *ctx.registers[1].w.as_mut() = i8::MIN as u32;
        *ctx.registers[2].w.as_mut() = i8::MAX as u32;
    }
}

#[skyline::hook(offset = 0x1A1CB48, inline)]
pub fn bound_mov(ctx: &mut InlineCtx) {
    unsafe {
        *ctx.registers[1].w.as_mut() = i8::MIN as u32;
        *ctx.registers[2].w.as_mut() = i8::MAX as u32;
    }
}

/// The internal name of your plugin. This will show up in crash logs. Make it 8 characters long at max.
#[skyline::main(name = "EnhnsLmt")]
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
        bound_max_hp_emblem,
        bound_str_emblem,
        bound_dex_emblem,
        bound_spd_emblem,
        bound_lck_emblem,
        bound_def_emblem,
        bound_mag_emblem,
        bound_res_emblem,
        bound_bld_emblem,
        bound_sight_emblem,
        bound_mov_emblem,
        bound_max_hp,
        bound_str,
        bound_dex,
        bound_spd,
        bound_lck,
        bound_def,
        bound_mag,
        bound_res,
        bound_bld,
        bound_sight,
        bound_mov,
    );
}
