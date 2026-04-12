package instrsbx

import (
	//"fmt"
)

//var InstrMap map[byte]*Instruction

//func init() {
//	InstrMap = make(map[byte]*Instruction)
//	for _, i := range Instructions {
//		InstrMap[i.Opcode] = i
//	}
//}

var Instructions map[byte]*Instruction = map[byte]*Instruction{
	0x80: &Instruction{0x80, 0, 0, false, false, Inline_None, "beep"},
	0x81: &Instruction{0x81, 0, 0, false, false, Inline_None, "halt_81"},
	0x82: &Instruction{0x82, 0, 0, false, false, Inline_None, "tape_nmi_shenanigans"},
	0x83: &Instruction{0x83, 0, 0, false, false, Inline_None, "tape_wait"},
	0x84: &Instruction{0x84, 0, 0, false, false, Inline_Word, "jump_abs"},
	0x85: &Instruction{0x85, 0, 0, false, false, Inline_Word, "call_abs"},
	0x86: &Instruction{0x86, 0, 0, false, false, Inline_None, "return"},
	0x87: &Instruction{0x87, 0, 0, false, false, Inline_None, "loop"},
	0x88: &Instruction{0x88, 0, 1, false, false, Inline_None, "play_sound"},
	0x89: &Instruction{0x89, 3, 1, false, false, Inline_None, "draw_string"},
	0x8A: &Instruction{0x8A, 0, 1, false, false, Inline_Word, "pop_string_to_addr"},
	0x8B: &Instruction{0x8B, 1, 0, false, false, Inline_None, "unknown_8B"},
	0x8C: &Instruction{0x8C, 0, 1,  true, false, Inline_None, "string_length"},
	0x8D: &Instruction{0x8D, 0, 1,  true,  true, Inline_None, "string_to_int"},
	0x8E: &Instruction{0x8E, 0, 2, false, false, Inline_None, "string_concat"},
	0x8F: &Instruction{0x8F, 0, 2,  true, false, Inline_None, "strings_equal"},

	0x90: &Instruction{0x90, 0, 2,  true, false, Inline_None, "strings_not_equal"},
	0x91: &Instruction{0x91, 0, 2,  true, false, Inline_None, "string_less_than"},
	0x92: &Instruction{0x92, 0, 2,  true, false, Inline_None, "string_less_than_equal"},
	0x93: &Instruction{0x93, 0, 2,  true, false, Inline_None, "strings_greater_than_equal"},
	0x94: &Instruction{0x94, 0, 2,  true, false, Inline_None, "strings_greater_than"},
	0x95: &Instruction{0x95, 1, 0, false, false, Inline_None, "tape_nmi_shenanigans_set"},
	0x96: &Instruction{0x96, 0, 0, false, false, Inline_Word, "set_word_4E"},
	0x97: &Instruction{0x97, 2, 0, false, false, Inline_None, "load_two_screens"},
	0x98: &Instruction{0x98, 1, 0, false, false, Inline_None, "unknown_98"},
	0x99: &Instruction{0x99, 1, 0, false, false, Inline_None, "enable_audio"},
	0x9A: &Instruction{0x9A, 0, 0, false, false, Inline_None, "disable_audio"},
	0x9B: &Instruction{0x9B, 0, 0, false, false, Inline_None, "halt_9B"},
	0x9C: &Instruction{0x9C, 0, 0, false, false, Inline_None, "toggle_44FE"},
	0x9D: &Instruction{0x9D, 2, 0, false, false, Inline_None, "something_tape"},
	0x9E: &Instruction{0x9E, 2, 0, false, false, Inline_None, "draw_and_show_screen"},
	0x9F: &Instruction{0x9F, 6, 0, false, false, Inline_None, "copy_tiles"},

	0xA0: &Instruction{0xA0, 2, 0,  true,  true, Inline_None, "unknown_A0"},
	0xA1: &Instruction{0xA1, 1, 0, false, false, Inline_None, "load_rom_screen"},
	0xA2: &Instruction{0xA2, 3, 0, false, false, Inline_None, "buffer_palette"},
	0xA3: &Instruction{0xA3, 1, 0, false, false, Inline_None, "sprite_setup"},
	0xA4: &Instruction{0xA4, 3, 0, false, false, Inline_None, "no_sp_wait"},
	0xA5: &Instruction{0xA5, 1, 0, false, false, Inline_None, "set_470A"},
	0xA6: &Instruction{0xA6, 1, 0, false, false, Inline_None, "set_470B"},
	0xA7: &Instruction{0xA7, 0, 0, false, false, Inline_None, "call_asm"},
	0xA8: &Instruction{0xA8, 5, 0, false, false, Inline_None, "play_noise"},
	0xA9: &Instruction{0xA9, 1, 0, false, false, Inline_None, "restore_tiles"},
	0xAA: &Instruction{0xAA, 1, 0, false, false, Inline_None, "long_jump"},
	0xAB: &Instruction{0xAB, 1, 0, false, false, Inline_None, "long_call"},
	0xAC: &Instruction{0xAC, 0, 0, false, false, Inline_None, "long_return"},
	0xAD: &Instruction{0xAD, 1, 0,  true,  true, Inline_None, "absolute"},
	0xAE: &Instruction{0xAE, 1, 0,  true,  true, Inline_None, "compare"},
	0xAF: &Instruction{0xAF, 0, 1,  true,  true, Inline_None, "push_char"},

	0xB0: &Instruction{0xB0, 1, 0, false, false, Inline_None, "pop_char"},
	0xB1: &Instruction{0xB1, 1, 0, false, false, Inline_None, "to_hex_string"},
	0xB2: &Instruction{0xB2, 0, 0,  true,  true, Inline_None, "read_mic"},
	0xB3: &Instruction{0xB3, 7, 0, false, false, Inline_None, "unknown_B3"},
	0xB4: &Instruction{0xB4, 0, 0, false, false, Inline_None, "indirect_copy_471A_4E"},
	0xB5: &Instruction{0xB5, 0, 0, false, false, Inline_None, "string_copy"},
	0xB6: &Instruction{0xB6, 0, 0, false, false, Inline_None, "word4E_to_word471A"},
	0xB7: &Instruction{0xB7, 0, 0, false, false, Inline_Word, "push_var"},
	0xB8: &Instruction{0xB8, 0, 0, false, false, Inline_Word, "push_word"},
	0xB9: &Instruction{0xB9, 0, 0, false, false, Inline_Word, "push_var_indexed"},
	0xBA: &Instruction{0xBA, 0, 0, false, false, Inline_Word, "push_data_indirect"},
	0xBB: &Instruction{0xBB, 0, 0, false, false, Inline_NullTerm, "push_data"},
	0xBC: &Instruction{0xBC, 0, 0, false, false, Inline_Word, "push_string_from_table"},
	0xBD: &Instruction{0xBD, 0, 0, false, false, Inline_Word, "pop_into"},
	0xBE: &Instruction{0xBE, 0, 0, false, false, Inline_Word, "write_to_table"},
	0xBF: &Instruction{0xBF, 1, 0, false, false, Inline_Word, "jump_not_zero"},

	0xC0: &Instruction{0xC0, 1, 0, false, false, Inline_Word, "jump_zero"},
	0xC1: &Instruction{0xC1, 1, 0, false, false, Inline_CountDefault, "jump_switch"},
	0xC2: &Instruction{0xC2, 1, 0,  true, false, Inline_None, "equals_zero"},
	0xC3: &Instruction{0xC3, 2, 0,  true, false, Inline_None, "and_a_b"},
	0xC4: &Instruction{0xC4, 2, 0,  true, false, Inline_None, "or_a_b"},
	0xC5: &Instruction{0xC5, 2, 0,  true, false, Inline_None, "equal"},
	0xC6: &Instruction{0xC6, 2, 0,  true, false, Inline_None, "not_equal"},
	0xC7: &Instruction{0xC7, 2, 0,  true, false, Inline_None, "less_than"},
	0xC8: &Instruction{0xC8, 2, 0,  true, false, Inline_None, "less_than_equal"},
	0xC9: &Instruction{0xC9, 2, 0,  true, false, Inline_None, "greater_than"},
	0xCA: &Instruction{0xCA, 2, 0,  true, false, Inline_None, "greater_than_equal"},
	0xCB: &Instruction{0xCB, 2, 0,  true,  true, Inline_None, "add"},
	0xCC: &Instruction{0xCC, 2, 0,  true,  true, Inline_None, "subtract"},
	0xCD: &Instruction{0xCD, 2, 0,  true,  true, Inline_None, "multiply"},
	0xCE: &Instruction{0xCE, 2, 0,  true,  true, Inline_None, "divide"},
	0xCF: &Instruction{0xCF, 1, 0,  true,  true, Inline_None, "negate"},

	0xD0: &Instruction{0xD0, 1, 0,  true, false, Inline_None, "modulo"},
	0xD1: &Instruction{0xD1, 2, 0,  true, false, Inline_None, "read_controller"},
	0xD2: &Instruction{0xD2, 2, 0,  true,  true, Inline_None, "unknown_D2"},
	0xD3: &Instruction{0xD3, 2, 0, false, false, Inline_None, "unknown_D3"},
	0xD4: &Instruction{0xD4, 3, 0, false, false, Inline_None, "set_cursor_location"},
	0xD5: &Instruction{0xD5, 1, 0, false, false, Inline_None, "wait_for_tape"},
	0xD6: &Instruction{0xD6, 1, 1, false, false, Inline_None, "truncate_string"},
	0xD7: &Instruction{0xD7, 1, 1, false, false, Inline_None, "trim_string_end"},
	0xD8: &Instruction{0xD8, 1, 1, false, false, Inline_None, "trim_string_start"},
	0xD9: &Instruction{0xD9, 2, 1, false, false, Inline_None, "substring"},
	0xDA: &Instruction{0xDA, 1, 0, false, false, Inline_None, "int_to_string"},
	0xDB: &Instruction{0xDB, 3, 0, false, false, Inline_None, "no_bg_wait"},
	0xDC: &Instruction{0xDC, 5, 0, false, false, Inline_None, "set_attr"},
	0xDD: &Instruction{0xDD, 5, 0, false, false, Inline_None, "fill_box"},
	0xDE: &Instruction{0xDE, 3, 0, false, false, Inline_None, "put_pixel"},
	0xDF: &Instruction{0xDF, 3, 0, false, false, Inline_None, "draw_image"},

	0xE0: &Instruction{0xE0, 2, 0,  true,  true, Inline_None, "modulo_B8"}, // BUT DIFFERENT FROM 0xD0 BECAUSE REASONS
	0xE1: &Instruction{0xE1, 4, 0, false, false, Inline_None, "put_pixel_alt"},
	0xE2: &Instruction{0xE2, 7, 0, false, false, Inline_None, "setup_sprite"},
	0xE3: &Instruction{0xE3, 1, 0,  true,  true, Inline_None, "deref_ptr_stack"},
	0xE4: &Instruction{0xE4, 2, 0, false, false, Inline_None, "swap_ram_bank"},
	0xE5: &Instruction{0xE5, 1, 0, false, false, Inline_None, "disable_sprite"},
	0xE6: &Instruction{0xE6, 1, 0, false, false, Inline_None, "tape_nmi_setup"},
	0xE7: &Instruction{0xE7, 7, 0, false, false, Inline_None, "draw_metasprite"},
	0xE8: &Instruction{0xE8, 1, 0, false, false, Inline_None, "setup_tape_nmi"}, // what is the difference with 0xE6???
	0xE9: &Instruction{0xE9, 0, 0, false, false, Inline_Word, "setup_loop"},
	0xEA: &Instruction{0xEA, 0, 0, false, false, Inline_Word, "string_write_to_table"},
	0xEB: &Instruction{0xEB, 4, 0, false, false, Inline_None, "draw_overlay"},
	0xEC: &Instruction{0xEC, 2, 0, false, false, Inline_None, "scroll"},
	0xED: &Instruction{0xED, 1, 0, false, false, Inline_None, "disable_sprites"},
	0xEE: &Instruction{0xEE, 1, 0, false, false, Inline_CountNoDefault, "call_switch"},
	0xEF: &Instruction{0xEF, 6, 1, false, false, Inline_None, "draw_debug_sprites"},

	0xF0: &Instruction{0xF0, 0, 1, false, false, Inline_None, "disable_some_sprites"},
	0xF1: &Instruction{0xF1, 4, 0, false, false, Inline_None, "unknown_F1"},
	0xF2: &Instruction{0xF2, 0, 0, false, false, Inline_None, "halt_F2"},
	0xF3: &Instruction{0xF3, 0, 1, false, false, Inline_None, "halt_F3"},
	0xF4: &Instruction{0xF4, 0, 0, false, false, Inline_None, "halt_F4"},
	0xF5: &Instruction{0xF5, 0, 0,  true,  true, Inline_None, "halt_F5"},
	0xF6: &Instruction{0xF6, 1, 0, false, false, Inline_None, "halt_F6"},
	0xF7: &Instruction{0xF7, 0, 1, false, false, Inline_None, "halt_F7"},
	0xF8: &Instruction{0xF8, 2, 0, false, false, Inline_None, "halt_F8"},
	0xF9: &Instruction{0xF9, 0, 0,  true,  true, Inline_None, "get_741"},
	0xFA: &Instruction{0xFA, 0, 0,  true,  true, Inline_None, "get_742"},
	0xFB: &Instruction{0xFB, 1, 0, false, false, Inline_None, "jump_arg_a"},
	0xFC: &Instruction{0xFC, 2, 0,  true,  true, Inline_None, "get_palette_color"},
	0xFD: &Instruction{0xFD, 0, 0, false, false, Inline_None, "halt_FD"},
	0xFE: &Instruction{0xFE, 4, 0, false, false, Inline_None, "draw_rom_char"},
	0xFF: &Instruction{0xFF, 7, 0,  true, false, Inline_None, "break_engine"},
}

//var Instructions map[byte]*Instruction = map[byte]*Instruction{
//	// WordArgs, StringArgs, ReturnWord, ReturnString, InlineType, Name
//	0x80: &Instruction{0, 0, false, false, Inline_None, "play_beep"},
//	0x81: &Instruction{0, 0, false, false, Inline_None, "halt"},
//	0x82: &Instruction{0, 0, false, false, Inline_None, "tape_nmi_shenanigans"},
//	0x83: &Instruction{0, 0, false, false, Inline_None, "tape_wait"},
//	0x84: &Instruction{0, 0, false, false, Inline_Word, "jump_abs"},
//	0x85: &Instruction{0, 0, false, false, Inline_Word, "call_abs"},
//	0x86: &Instruction{0, 0, false, false, Inline_None, "return"},
//	0x87: &Instruction{0, 0, false, false, Inline_None, "loop"},
//	0x88: &Instruction{0, 0, 0,  false, "play_sound"},
//	0x89: &Instruction{3, 0, 0,  false, "draw_string"},
//	0x8A: &Instruction{0, 2, 0,  false, "pop_string_to_addr"},
//	0x8B: &Instruction{1, 0, 0,  false, ""},
//	0x8C: &Instruction{0, 0, 1,  false, "string_length"},
//	0x8D: &Instruction{0, 0, 1,  false, "string_to_int"},
//	0x8E: &Instruction{0, 0, 16, false, "string_concat"},
//	0x8F: &Instruction{0, 0, 1,  false, "strings_equal"},

//	0x90: &Instruction{0, 0, 1,  false,  "strings_not_equal"},
//	0x91: &Instruction{0, 0, 1,  false,  "string_less_than"},
//	0x92: &Instruction{0, 0, 1,  false,  "string_less_than_equal"},
//	0x93: &Instruction{0, 0, 1,  false,  "string_greater_than_equal"},
//	0x94: &Instruction{0, 0, 1,  false,  "string_greater_than"},
//	0x95: &Instruction{1, 0, 0,  false,  "tape_nmi_shenigans_set"},
//	0x96: &Instruction{0, 2, 0,  true,   "set_word_4E"},
//	0x97: &Instruction{2, 0, 0,  false,  "load_two_screens"},
//	0x98: &Instruction{1, 0, 0,  false,  ""},
//	0x99: &Instruction{1, 0, 0,  false,  "enable_audio"},
//	0x9A: &Instruction{0, 0, 0,  false,  "disable_audio"},
//	0x9B: &Instruction{0, 0, 0,  false,  "halt"},
//	0x9C: &Instruction{0, 0, 0,  false,  "toggle_44FE"},
//	0x9D: &Instruction{2, 0, 0,  false,  "something_tape"},
//	0x9E: &Instruction{2, 0, 0,  false,  "draw_and_show_screen"},
//	0x9F: &Instruction{6, 0, 0,  false,  "copy_tiles"},

//	0xA0: &Instruction{2, 0, 1,  false,  ""},
//	0xA1: &Instruction{1, 0, 0,  false,  "load_rom_screen"},
//	0xA2: &Instruction{1, 0, 0,  false,  "buffer_palette"},
//	0xA3: &Instruction{1, 0, 0,  false,  "sprite_setup"},
//	0xA4: &Instruction{3, 0, 0,  false,  "no_sp_wait"},
//	0xA5: &Instruction{1, 0, 0,  false,  "set_470A"},
//	0xA6: &Instruction{1, 0, 0,  false,  "set_470B"},
//	0xA7: &Instruction{0, 0, 0,  false,   "call_asm"},
//	0xA8: &Instruction{5, 0, 0,  false,  "play_noise"},
//	0xA9: &Instruction{1, 0, 0,  false,  "restore_tiles"},
//	0xAA: &Instruction{1, 0, 0,  false,  "long_jump"},
//	0xAB: &Instruction{1, 0, 0,  false,  "long_call"},
//	0xAC: &Instruction{0, 0, 0,  false,  "long_return"},
//	0xAD: &Instruction{1, 0, 1,  false,  "absolute"},
//	0xAE: &Instruction{1, 0, 1,  false,  "compare"},
//	0xAF: &Instruction{0, 0, 1,  false,  "push_char"},

//	0xB0: &Instruction{1, 0, 16, false, "pop_char"},
//	0xB1: &Instruction{1, 0, 16, false, "to_hex_string"},
//	0xB2: &Instruction{0, 0, 1,  false,  "read_mic"},
//	0xB3: &Instruction{7, 0, 0,  false,  ""},
//	0xB4: &Instruction{0, 0, 0,  false,  "indirect_copy_471A_4E"},
//	0xB5: &Instruction{0, 0, 0,  false,  "string_copy"},
//	0xB6: &Instruction{0, 0, 0,  false,  "word4E_to_word471A"},
//	0xB7: &Instruction{0, 2, 0,  false,  "push_var"},
//	0xB8: &Instruction{0, 2, 0,  true,  "push_word"},
//	0xB9: &Instruction{0, 2, 0,  false, "push_var_indexed"},
//	0xBA: &Instruction{0, 2, 0,  false, "push_data_indirect"},
//	0xBB: &Instruction{0, -1, 0, false, "push_data"},
//	0xBC: &Instruction{0, 2, 0,  false, "push_string_from_table"},
//	0xBD: &Instruction{0, 2, 0,  false,  "pop_into"},
//	0xBE: &Instruction{0, 2, 0,  false,  "write_to_table"},
//	0xBF: &Instruction{0, 2, 0,  false,  "jump_not_zero"},

//	0xC0: &Instruction{1, 2, 0,  false,  "jump_zero"},
//	0xC1: &Instruction{1, -2, 0, false,  "jump_switch"},
//	0xC2: &Instruction{1, 0, 1,  false,  "equals_zero"},
//	0xC3: &Instruction{2, 0, 1,  false,  "and_a_b"},
//	0xC4: &Instruction{2, 0, 1,  false,  "or_a_b"},
//	0xC5: &Instruction{2, 0, 1,  false,  "equal"},
//	0xC6: &Instruction{2, 0, 1,  false,  "not_equal"},
//	0xC7: &Instruction{2, 0, 1,  false,  "less_than"},
//	0xC8: &Instruction{2, 0, 1,  false,  "less_than_equal"},
//	0xC9: &Instruction{2, 0, 1,  false,  "greater_than"},
//	0xCA: &Instruction{2, 0, 1,  false,  "greater_than_equal"},
//	0xCB: &Instruction{2, 0, 1,  false,  "sum"},
//	0xCC: &Instruction{2, 0, 1,  false,  "subtract"},
//	0xCD: &Instruction{2, 0, 1,  false,  "multiply"},
//	0xCE: &Instruction{2, 0, 1,  false,  "signed_divide"},
//	0xCF: &Instruction{1, 0, 1,  false,  "negate"},

//	0xD0: &Instruction{1, 0, 1,  false,  "modulus"},
//	0xD1: &Instruction{2, 0, 1,  false,  "expansion_controller"},
//	0xD2: &Instruction{2, 0, 1,  false,  ""},
//	0xD3: &Instruction{2, 0, 16, false,  ""},
//	0xD4: &Instruction{3, 0, 0,  false,  "set_cursor_location"},
//	0xD5: &Instruction{1, 0, 0,  false,  "wait_for_tape"},
//	0xD6: &Instruction{1, 0, 16, false,  "truncate_string"},
//	0xD7: &Instruction{1, 0, 16, false,  "trim_string"},
//	0xD8: &Instruction{1, 0, 16, false,  "trim_string_start"},
//	0xD9: &Instruction{2, 0, 16, false,  "substring"},
//	0xDA: &Instruction{1, 0, 16, false,  "int_to_string"},
//	0xDB: &Instruction{3, 0, 0,  false,  "no_bg_wait"},
//	0xDC: &Instruction{5, 0, 0,  false,  "set_attr"}, // fucks with attribute data
//	0xDD: &Instruction{5, 0, 0,  false,  "fill_box"},
//	0xDE: &Instruction{3, 0, 0,  false,  "put_pixel"},
//	0xDF: &Instruction{3, 0, 0,  false,  "draw_image"},

//	0xE0: &Instruction{2, 0, 1,  false,  "modulo"},
//	0xE1: &Instruction{4, 0, 0,  false,  "put_pixel_alt"},
//	0xE2: &Instruction{7, 0, 0,  false,  "setup_sprite"},
//	0xE3: &Instruction{1, 0, 1,  false,  "deref_ptr_stack"},
//	0xE4: &Instruction{2, 0, 0,  false,  "swap_ram_bank"},
//	0xE5: &Instruction{1, 0, 0,  false,  "disable_sprite"},
//	0xE6: &Instruction{1, 0, 0,  false,  "tape_nmi_setup"},
//	0xE7: &Instruction{7, 0, 0,  false,  "draw_metasprite"},
//	0xE8: &Instruction{1, 0, 0,  false,  "setup_tape_nmi"},
//	0xE9: &Instruction{0, 2, 0,  false,  "setup_loop"},
//	0xEA: &Instruction{0, 2, 0,  false,  "string_write_to_table"},
//	0xEB: &Instruction{4, 0, 0,  false,  "draw_overlay"},
//	0xEC: &Instruction{2, 0, 0,  false,  "scroll"},
//	0xED: &Instruction{1, 0, 0,  false,  "disable_sprites"},
//	0xEE: &Instruction{1, -3, 0, false,  "call_switch"},
//	0xEF: &Instruction{6, 0, 0,  false,  "draw_debug_sprites"},

//	0xF0: &Instruction{0, 0, 0,  false,  "disable_sprites"},
//	0xF1: &Instruction{4, 0, 0,  false,  ""},
//	0xF2: &Instruction{0, 0, 0,  false,  "halt"},
//	0xF3: &Instruction{0, 0, 0,  false,  "halt"},
//	0xF4: &Instruction{0, 0, 16, false,  "halt"},
//	0xF5: &Instruction{1, 0, 1,  false,  "halt"},
//	0xF6: &Instruction{1, 0, 0,  false,  "halt"},
//	0xF7: &Instruction{0, 0, 0,  false,  "halt"},
//	0xF8: &Instruction{2, 0, 0,  false,  "halt"},
//	0xF9: &Instruction{0, 0, 1,  false,  ""},
//	0xFA: &Instruction{0, 0, 1,  false,  ""},
//	0xFB: &Instruction{1, 0, 0,  false,  "jump_arg_a"},
//	0xFC: &Instruction{2, 0, 1,  false,  "get_palette_color"},
//	0xFD: &Instruction{0, 0, 16, false,  "halt"},
//	0xFE: &Instruction{4, 2, 0,  false,  "draw_rom_char"},
//	0xFF: &Instruction{0, 0, 0,  false,  "break_engine"},
//}

type InlineType int

const (
	Inline_None InlineType = iota
	Inline_Word
	Inline_NullTerm
	Inline_CountDefault
	Inline_CountNoDefault
)

type Instruction struct {
	Opcode byte

	WordArgs   int
	StringArgs int

	ReturnWord   bool
	ReturnString bool

	Inline InlineType

	Name string

	//ArgCount  int  // stack arguments
	//OpCount   int  // inline operands.  length in bytes.
	//			   // -1: nul-terminated
	//			   // -2: first byte is count, followed by that number of words
	//			   // -3: like -2, but with no default.  code continues after list on OOB
	//RetCount  int  // return count
	//InlineImmediate bool // don't turn the inline value into a variable
	//Name      string
}

func (i Instruction) String() string {
	if i.Name != "" {
		//return fmt.Sprintf("$%02X_%s", i.Opcode, i.Name)
		return i.Name
	}

	//return fmt.Sprintf("unknown_0x%02X", i.Opcode)
	return "unknown"
}

