// Currently needed because we use these functionality, they'll be removable when the Rust language stabilizes them
#![feature(lazy_cell, ptr_sub_ptr)]

use engage::{
    battle::BattleInfoSide,
    calculator::*,
    gamedata::{
        terrain::{self, TerrainData},
        unit::Unit,
    },
    map::{
        image::{MapImage, MapImageTerrain},
        overlap,
    },
    util::get_instance,
};
use skyline::nn::friends::Profile_IsValid;
use unity::prelude::*;
/// This is called a proc(edural) macro. You use this to indicate that a function will be used as a hook.
///
/// Pay attention to the argument, offset.
/// This is the address of the start of the function you would like to hook.
/// This address has to be relative to the .text section of the game.
/// If you do not know what any of this means, take the address in Ghidra and remove the starting ``71`` and the zeroes that follow it.
/// Do not forget the 0x indicator, as it denotates that you are providing a hexadecimal value.
#[skyline::from_offset(0x01E8B320)]
pub fn battle_info_side_get_terrain(this: &BattleInfoSide, method: OptionalMethod) -> &TerrainData;

#[skyline::from_offset(0x02064ED0)]
pub fn map_image_terrain_get_data(
    this: &MapImageTerrain,
    x: i32,
    z: i32,
    method: OptionalMethod,
) -> &TerrainData;

#[skyline::from_offset(0x01DFDA90)]
pub fn map_overlap_get_data(x: i32, z: i32, method: OptionalMethod) -> &'static TerrainData;

#[skyline::from_offset(0x021E33A0)]
pub fn terrain_is_fly_enable(terrain: &TerrainData, method: OptionalMethod) -> bool;

#[skyline::from_offset(0x01A54D00)]
pub fn unit_get_move_type(this: &Unit, method: OptionalMethod) -> i32;

#[skyline::from_offset(0x01E8B250)]
pub fn battle_info_side_get_unit(this: &BattleInfoSide, method: OptionalMethod) -> Option<&Unit>;

#[skyline::from_offset(0x021E3380)]
pub fn terrain_data_is_immune_break(this: &TerrainData, method: OptionalMethod) -> bool;

trait UnitTerrainTrait {
    fn get_map_terrain(&self) -> &TerrainData;
    fn get_map_overlap(&self) -> &TerrainData;
    fn is_terrain_valid(&self, terrain: &TerrainData) -> bool;
    fn get_total_terrain_avo(&self) -> f32;
    fn get_total_terrain_def(&self) -> f32;
    fn get_total_terrain_heal(&self) -> f32;
    fn get_total_terrain_mov(&self) -> f32;
    fn get_total_terrain_immune_break(&self) -> f32;
}

impl UnitTerrainTrait for Unit {
    fn get_map_terrain(&self) -> &TerrainData {
        let im = get_instance::<MapImage>();
        let imt = im.terrain;
        unsafe { map_image_terrain_get_data(imt, self.x.into(), self.z.into(), None) }
    }
    fn get_map_overlap(&self) -> &TerrainData {
        unsafe { map_overlap_get_data(self.x.into(), self.z.into(), None) }
    }
    fn is_terrain_valid(&self, terrain: &TerrainData) -> bool {
        const FLY: i32 = 3;
        if unsafe { unit_get_move_type(self, None) == FLY } {
            if self.has_sid("SID_翼盾".into()) {
                true
            } else {
                false
            }
        } else {
            true
        }
    }
    fn get_total_terrain_avo(&self) -> f32 {
        let terrain = self.get_map_terrain();
        let overlap = self.get_map_overlap();
        let mut result = 0;
        if self.is_terrain_valid(terrain) {
            result += terrain.avoid;
            if self.is_player() {
                result += terrain.player_avoid;
            } else if self.is_enemy() {
                result += terrain.enemy_avoid;
            }
        }
        if self.is_terrain_valid(overlap) {
            result += overlap.avoid;
            if self.is_player() {
                result += overlap.player_avoid;
            } else if self.is_enemy() {
                result += overlap.enemy_avoid;
            }
        }
        result as f32
    }
    fn get_total_terrain_def(&self) -> f32 {
        let terrain = self.get_map_terrain();
        let overlap = self.get_map_overlap();
        let mut result = 0;
        if self.is_terrain_valid(terrain) {
            result += terrain.defense;
            if self.is_player() {
                result += terrain.player_defense;
            } else if self.is_enemy() {
                result += terrain.enemy_defense;
            }
        }
        if self.is_terrain_valid(overlap) {
            result += overlap.defense;
            if self.is_player() {
                result += overlap.player_defense;
            } else if self.is_enemy() {
                result += overlap.enemy_defense;
            }
        }
        result as f32
    }
    fn get_total_terrain_heal(&self) -> f32 {
        let terrain = self.get_map_terrain();
        let overlap = self.get_map_overlap();
        let result = if self.is_terrain_valid(terrain) {
            terrain.heal
        } else {
            0
        } + if self.is_terrain_valid(overlap) {
            overlap.heal
        } else {
            0
        };
        result as f32
    }

    fn get_total_terrain_mov(&self) -> f32 {
        let terrain = self.get_map_terrain();
        let overlap = self.get_map_overlap();
        let result = if self.is_terrain_valid(terrain) {
            terrain.move_first
        } else {
            0
        } + if self.is_terrain_valid(overlap) {
            overlap.move_first
        } else {
            0
        };
        result as f32
    }
    fn get_total_terrain_immune_break(&self) -> f32 {
        let terrain = self.get_map_terrain();
        let overlap = self.get_map_overlap();
        let terrain_immune = self.is_terrain_valid(terrain)
            && unsafe { terrain_data_is_immune_break(terrain, None) };
        let overlap_immune = self.is_terrain_valid(overlap)
            && unsafe { terrain_data_is_immune_break(overlap, None) };
        if terrain_immune || overlap_immune {
            1f32
        } else {
            0f32
        }
    }
}

trait UnitForceTrait {
    fn is_player(&self) -> bool;
    fn is_enemy(&self) -> bool;
}

impl UnitForceTrait for Unit {
    fn is_player(&self) -> bool {
        if let Some(force) = self.force {
            force.force_type == 0 || force.force_type >= 2 && force.force_type <= 5
        } else {
            false
        }
    }
    fn is_enemy(&self) -> bool {
        if let Some(force) = self.force {
            force.force_type == 1
        } else {
            false
        }
    }
}

#[unity::hook("App", "UnitCalculator", "AddCommand")]
pub fn add_command_hook(calculator: &mut CalculatorManager, method_info: OptionalMethod) {
    call_original!(calculator, method_info);
    let terrain_avo_c: &mut CalculatorCommand = calculator.find_command("地形回避");
    let terrain_def_c: &mut CalculatorCommand = calculator.find_command("地形防御");
    let terrain_heal_c =
        il2cpp::instantiate_class::<GameCalculatorCommand>(terrain_avo_c.get_class().clone())
            .unwrap();
    let terrain_mov_c =
        il2cpp::instantiate_class::<GameCalculatorCommand>(terrain_avo_c.get_class().clone())
            .unwrap();
    let terrain_immune_break_c =
        il2cpp::instantiate_class::<GameCalculatorCommand>(terrain_avo_c.get_class().clone())
            .unwrap();
    terrain_avo_c
        .get_class_mut()
        .get_virtual_method_mut("GetImpl")
        .map(|method| method.method_ptr = get_terrain_avo_command_unit as _);
    terrain_def_c
        .get_class_mut()
        .get_virtual_method_mut("GetImpl")
        .map(|method| method.method_ptr = get_terrain_def_command_unit as _);
    // Heal Command
    terrain_heal_c
        .get_class_mut()
        .get_virtual_method_mut("get_Name")
        .map(|method| method.method_ptr = get_terrain_heal_command_name as _);
    terrain_heal_c
        .get_class_mut()
        .get_virtual_method_mut("GetImpl")
        .map(|method| method.method_ptr = get_terrain_heal_command_unit as _);
    terrain_heal_c.get_class_mut().get_vtable_mut()[31].method_ptr =
        get_terrain_heal_command_battle_info as *mut u8;
    calculator.add_command(terrain_heal_c);

    // Foe Heal Command
    let terrain_heal_c2 =
        il2cpp::instantiate_class::<GameCalculatorCommand>(terrain_heal_c.get_class().clone())
            .unwrap();
    let foe_heal_c = terrain_heal_c2.reverse();
    calculator.add_command(foe_heal_c);

    // Mov Command
    terrain_mov_c
        .get_class_mut()
        .get_virtual_method_mut("get_Name")
        .map(|method| method.method_ptr = get_terrain_mov_command_name as _);
    terrain_mov_c
        .get_class_mut()
        .get_virtual_method_mut("GetImpl")
        .map(|method| method.method_ptr = get_terrain_mov_command_unit as _);
    terrain_mov_c.get_class_mut().get_vtable_mut()[31].method_ptr =
        get_terrain_mov_command_battle_info as *mut u8;
    calculator.add_command(terrain_mov_c);

    // Foe Mov Command
    let terrain_mov_c2 =
        il2cpp::instantiate_class::<GameCalculatorCommand>(terrain_mov_c.get_class().clone())
            .unwrap();
    let foe_mov_c = terrain_mov_c2.reverse();
    calculator.add_command(foe_mov_c);

    // Immune Break Command
    terrain_immune_break_c
        .get_class_mut()
        .get_virtual_method_mut("get_Name")
        .map(|method| method.method_ptr = get_terrain_immune_break_command_name as _);
    terrain_immune_break_c
        .get_class_mut()
        .get_virtual_method_mut("GetImpl")
        .map(|method| method.method_ptr = get_terrain_immune_break_command_unit as _);
    terrain_immune_break_c.get_class_mut().get_vtable_mut()[31].method_ptr =
        get_terrain_immune_break_command_battle_info as *mut u8;
    calculator.add_command(terrain_immune_break_c);

    // Foe Immune Break Command
    let terrain_immune_break_c2 = il2cpp::instantiate_class::<GameCalculatorCommand>(
        terrain_immune_break_c.get_class().clone(),
    )
    .unwrap();
    let foe_immune_break_c = terrain_immune_break_c2.reverse();
    calculator.add_command(foe_immune_break_c);
}

fn get_terrain_heal_command_name(
    _this: &GameCalculatorCommand,
    method: OptionalMethod,
) -> &'static Il2CppString {
    "TerrainHeal".into()
}
fn get_terrain_mov_command_name(
    _this: &GameCalculatorCommand,
    method: OptionalMethod,
) -> &'static Il2CppString {
    "TerrainMov".into()
}

fn get_terrain_immune_break_command_name(
    _this: &GameCalculatorCommand,
    method: OptionalMethod,
) -> &'static Il2CppString {
    "TerrainImmuneBreak".into()
}

fn get_terrain_avo_command_unit(
    _this: &GameCalculatorCommand,
    unit: &Unit,
    method_info: OptionalMethod,
) -> f32 {
    let avo = unit.get_total_terrain_avo();
    avo
}

fn get_terrain_def_command_unit(
    _this: &GameCalculatorCommand,
    unit: &Unit,
    method_info: OptionalMethod,
) -> f32 {
    let def = unit.get_total_terrain_def();
    def
}

fn get_terrain_heal_command_unit(
    _this: &GameCalculatorCommand,
    unit: &Unit,
    method_info: OptionalMethod,
) -> f32 {
    let heal = unit.get_total_terrain_heal();
    heal
}

fn get_terrain_mov_command_unit(
    _this: &GameCalculatorCommand,
    unit: &Unit,
    method_info: OptionalMethod,
) -> f32 {
    let mov = unit.get_total_terrain_mov();
    mov
}

fn get_terrain_immune_break_command_unit(
    _this: &GameCalculatorCommand,
    unit: &Unit,
    method_info: OptionalMethod,
) -> f32 {
    let immune_break = unit.get_total_terrain_immune_break();
    immune_break
}

fn get_terrain_heal_command_battle_info(
    _this: &GameCalculatorCommand,
    side: &BattleInfoSide,
    method_info: OptionalMethod,
) -> f32 {
    let unit = unsafe { battle_info_side_get_unit(side, None) };
    if let Some(unit) = unit {
        get_terrain_heal_command_unit(_this, unit, None)
    } else {
        0f32
    }
}

fn get_terrain_mov_command_battle_info(
    _this: &GameCalculatorCommand,
    side: &BattleInfoSide,
    method_info: OptionalMethod,
) -> f32 {
    let unit = unsafe { battle_info_side_get_unit(side, None) };
    if let Some(unit) = unit {
        get_terrain_mov_command_unit(_this, unit, None)
    } else {
        0f32
    }
}

fn get_terrain_immune_break_command_battle_info(
    _this: &GameCalculatorCommand,
    side: &BattleInfoSide,
    method_info: OptionalMethod,
) -> f32 {
    let unit = unsafe { battle_info_side_get_unit(side, None) };
    if let Some(unit) = unit {
        get_terrain_immune_break_command_unit(_this, unit, None)
    } else {
        0f32
    }
}

/// The internal name of your plugin. This will show up in crash logs. Make it 8 characters long at max.
#[skyline::main(name = "IfCmd")]
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
    skyline::install_hook!(add_command_hook);
}
