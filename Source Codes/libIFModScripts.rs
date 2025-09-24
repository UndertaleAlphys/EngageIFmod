// Currently needed because we use these functionality, they'll be removable when the Rust language stabilizes them
#![feature(lazy_cell, ptr_sub_ptr)]

use engage::{
    gamedata::{item::ItemData, Gamedata,unit::Unit},
    script::{DynValue, EventScript, EventScriptCommand, ScriptUtils},
};
use engage::gamedata::skill::{SkillArrayEntity, SkillArrayEntityList};
use engage::gamedata::unit::UnitUtil;
use engage::script::EventResultScriptCommand;
use skyline::hooks::InlineCtx;
use unity::from_offset;
use unity::il2cpp::object::Array;
use unity::prelude::*;
/// This is called a proc(edural) macro. You use this to indicate that a function will be used as a hook.
///
/// Pay attention to the argument, offset.
/// This is the address of the start of the function you would like to hook.
/// This address has to be relative to the .text section of the game.
/// If you do not know what any of this means, take the address in Ghidra and remove the starting ``71`` and the zeroes that follow it.
/// Do not forget the 0x indicator, as it denotates that you are providing a hexadecimal value.

#[repr(C)]
pub struct Vector3 {
    pub fields: Vector3Fields,
}

#[repr(C)]
pub struct Vector3Fields {
    pub x: f32,
    pub y: f32,
    pub z: f32,
}

#[repr(C)]
pub struct Quaternion {
    pub fields: QuaternionFields,
}

#[repr(C)]
pub struct QuaternionFields {
    pub x: f32,
    pub y: f32,
    pub z: f32,
    pub w: f32,
}
#[skyline::hook(offset = 0x024E279C, inline)]
pub fn scripts_regist(ctx: &InlineCtx) {
    let event = unsafe { &*((*ctx.registers[20].x.as_ref()) as *const EventScript) };
    ScriptIF::register(event);
}
pub struct ScriptIF;
impl ScriptIF {
    pub fn register(event: &EventScript) {
        event.register_action("RemoveTransporterItem", remove_transporter_item);
        event.register_function("GetUnitSkillList", get_unit_skill_list_func);
        event.register_action("EffectCreateRotateXYZ", create_map_effect_rotate_xyz);
    }
}

#[skyline::from_offset(0x021A3970)]
pub fn dyn_value_array_get_by_idx(
    array: &Il2CppArray<DynValue>,
    idx: i32,
    method: OptionalMethod,
) -> Option<&DynValue>;
#[skyline::from_offset(0x02E3CD60)]
pub fn dyn_value_get_type(dv: &DynValue, method: OptionalMethod) -> i32;
#[skyline::from_offset(0x022A2260)]
pub fn transporter_delete(idx: i32, method: OptionalMethod);
#[skyline::from_offset(0x022A2570)]
pub fn transporter_delete_item(data: &ItemData, method: OptionalMethod);

#[skyline::from_offset(0x021961F0)]
pub fn try_get_unit(array: &Il2CppArray<DynValue>, idx: i32, warning: bool,method: OptionalMethod) -> &Unit;

#[skyline::from_offset(0x01DBD3E0)]
pub fn app_map_effect_create(
    name: Option<&'static Il2CppString>,
    position: Vector3,
    rotation: Quaternion,
    method: OptionalMethod
);




extern "C" fn remove_transporter_item(args: &Il2CppArray<DynValue>, _method: OptionalMethod) {
    for arg_idx in 0..args.len() {
        const NONE: i32 = i32::MAX;
        const NUMBER: i32 = 3;
        const STRING: i32 = 4;
        match {
            if let Some(dv) = unsafe { dyn_value_array_get_by_idx(args, arg_idx as i32, None) } {
                unsafe { dyn_value_get_type(dv, None) }
            } else {
                NONE
            }
        } {
            NUMBER => {
                let iidx = args.try_get_i32(arg_idx as i32);
                unsafe { transporter_delete(iidx, None) };
            }
            STRING => {
                if let Some(iid) = args.try_get_string(arg_idx as i32) {
                    let idata = ItemData::get(iid);
                    if let Some(idata) = idata {
                        unsafe { transporter_delete_item(idata, None) };
                    }
                }
            }
            _ => {}
        }
    }
}

fn quaternion_multiply(q1: Quaternion, q2: Quaternion) -> Quaternion {
    Quaternion {
        fields: QuaternionFields {
            x: q1.fields.w * q2.fields.x + q1.fields.x * q2.fields.w + q1.fields.y * q2.fields.z - q1.fields.z * q2.fields.y,
            y: q1.fields.w * q2.fields.y - q1.fields.x * q2.fields.z + q1.fields.y * q2.fields.w + q1.fields.z * q2.fields.x,
            z: q1.fields.w * q2.fields.z + q1.fields.x * q2.fields.y - q1.fields.y * q2.fields.x + q1.fields.z * q2.fields.w,
            w: q1.fields.w * q2.fields.w - q1.fields.x * q2.fields.x - q1.fields.y * q2.fields.y - q1.fields.z * q2.fields.z,
        },
    }
}
extern "C" fn create_map_effect_rotate_xyz(args: &Il2CppArray<DynValue>, _method: OptionalMethod) {
    let effectName = args.try_get_string(0);
    let positionX = ((args.try_get_i32(1) as f32 * 5.0)+2.5 );
    let positionY = (args.try_get_i32(3) );
    let positionZ = ((args.try_get_i32(2) as f32 * 5.0)+2.5 );
    let rotationXAngle = args.try_get_i32(4);
    let rotationYAngle = args.try_get_i32(5);
    let rotationZAngle = args.try_get_i32(6);
    let order = args.try_get_i32(7);

    const XYZ: i32 = 0;
    const XZY: i32 = 1;
    const YXZ: i32 = 2;
    const YZX: i32 = 3;
    const ZXY: i32 = 4;
    const ZYX: i32 = 5;

    let position = Vector3 {
        fields: Vector3Fields {
            x: positionX as f32,
            y: positionY as f32,
            z: positionZ as f32,
        },
    };

    let x_angle_rad = (rotationXAngle as f32).to_radians();
    let x_half_angle = x_angle_rad / 2.0;
    let x_sin_half = x_half_angle.sin();
    let x_cos_half = x_half_angle.cos();

    let rotation_x = Quaternion {
        fields: QuaternionFields {
            x: x_sin_half,
            y: 0.0,
            z: 0.0,
            w: x_cos_half,
        },
    };

    let y_angle_rad = (rotationYAngle as f32).to_radians();
    let y_half_angle = -y_angle_rad / 2.0;
    let y_sin_half = y_half_angle.sin();
    let y_cos_half = y_half_angle.cos();

    let rotation_y = Quaternion {
        fields: QuaternionFields {
            x: 0.0,
            y: y_sin_half,
            z: 0.0,
            w: y_cos_half,
        },
    };

    let z_angle_rad = (rotationZAngle as f32).to_radians();
    let z_half_angle = z_angle_rad / 2.0;
    let z_sin_half = z_half_angle.sin();
    let z_cos_half = z_half_angle.cos();

    let rotation_z = Quaternion {
        fields: QuaternionFields {
            x: 0.0,
            y: 0.0,
            z: z_sin_half,
            w: z_cos_half,
        },
    };

    // 根据顺序参数选择不同的旋转顺序
    let rotation = match order {
        XYZ => {
            let rot_xy = quaternion_multiply(rotation_y, rotation_x);
            quaternion_multiply(rotation_z, rot_xy)
        },
        XZY => {
            let rot_xz = quaternion_multiply(rotation_z, rotation_x);
            quaternion_multiply(rotation_y, rot_xz)
        },
        YXZ => {
            let rot_yx = quaternion_multiply(rotation_x, rotation_y);
            quaternion_multiply(rotation_z, rot_yx)
        },
        YZX => {
            let rot_yz = quaternion_multiply(rotation_z, rotation_y);
            quaternion_multiply(rotation_x, rot_yz)
        },
        ZXY => {
            let rot_zx = quaternion_multiply(rotation_x, rotation_z);
            quaternion_multiply(rotation_y, rot_zx)
        },
        ZYX => {
            let rot_zy = quaternion_multiply(rotation_y, rotation_z);
            quaternion_multiply(rotation_x, rot_zy)
        },
        _ => {
            let rot_xy = quaternion_multiply(rotation_y, rotation_x);
            quaternion_multiply(rotation_z, rot_xy)
        }
    };

    match effectName {
        Some(name) => {
            unsafe {
                match std::panic::catch_unwind(|| {
                    app_map_effect_create(Some(name), position, rotation, _method);
                }) {
                    Ok(_) => println!("Successfully created effect: {}", name.to_string()),
                    Err(e) => {
                        println!("Error occurred while creating effect: {}", name.to_string());
                        println!("Panic info: {:?}", e);
                    }
                }
            }
        },
        None => {
        }
    }
}
extern "C" fn get_unit_skill_list_func(args: &Il2CppArray<DynValue>, _method: OptionalMethod) -> &'static DynValue {
    let unit: Option<&Unit> = unsafe {
        let unit_ptr = try_get_unit(args, 0, true, None) as *const Unit;
        if unit_ptr.is_null() {
            None
        } else {
            Some(&*unit_ptr)
        }
    };

    let mut skill_sids = String::new();

    if let Some(unit) = unit {
        if let Some(skill_list) = &unit.fields.mask_skill {
            let skills = &skill_list.fields.list.fields.item;

            for (i, skill_entity) in skills.iter().enumerate() {
                let skill_entity: &SkillArrayEntity = skill_entity;
                let skill_data = skill_entity.get_skill();

                let skill_sid = if let Some(skill) = skill_data {
                    match std::panic::catch_unwind(|| skill.fields.sid.to_string()) {
                        Ok(sid) => sid,
                        Err(_) => "None".to_string()
                    }
                } else {
                    "None".to_string()
                };

                if i > 0 {
                    skill_sids.push(',');
                }
                skill_sids.push_str(&skill_sid);
            }
        }
    }

    let il2cpp_string = Il2CppString::new(&skill_sids);

    DynValue::new_string(&il2cpp_string)
}





/// The internal name of your plugin. This will show up in crash logs. Make it 8 characters long at max.
#[skyline::main(name = "IFScript")]
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
    skyline::install_hooks!(scripts_regist);
}