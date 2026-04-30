//go:build sql_lite

package sql

import (
	"fmt"
	"hash/fnv"
	"io"
	"strings"
	"unicode/utf8"
)

// Collation represents the collation of a string.
type Collation struct {
	ID                CollationID
	Name              string
	CharacterSet      CharacterSetID
	IsDefault         bool
	IsCompiled        bool
	IsCaseSensitive   bool
	IsAccentSensitive bool
	SortLength        uint8
	PadAttribute      string
	Sorter            CollationSorter
}

// CollationSorter is a collation's sort function.
type CollationSorter func(r rune) int32

// CollationsIterator iterates over available lite collations.
type CollationsIterator struct {
	idx int
}

var collationStringToID = map[string]CollationID{
	"binary":             Collation_binary,
	"utf8mb4_0900_ai_ci": Collation_utf8mb4_0900_ai_ci,
	"utf8mb4_bin":        Collation_utf8mb4_bin,
	"utf8mb4_general_ci": Collation_utf8mb4_general_ci,
	"utf8mb3_general_ci": Collation_utf8mb3_general_ci,
	"utf8_general_ci":    Collation_utf8_general_ci,
}

// CollationID represents the collation's unique identifier.
type CollationID uint16

const (
	Collation_armscii8_bin                CollationID = 64
	Collation_armscii8_general_ci         CollationID = 32
	Collation_ascii_bin                   CollationID = 65
	Collation_ascii_general_ci            CollationID = 11
	Collation_big5_bin                    CollationID = 84
	Collation_big5_chinese_ci             CollationID = 1
	Collation_binary                      CollationID = 63
	Collation_cp1250_bin                  CollationID = 66
	Collation_cp1250_croatian_ci          CollationID = 44
	Collation_cp1250_czech_cs             CollationID = 34
	Collation_cp1250_general_ci           CollationID = 26
	Collation_cp1250_polish_ci            CollationID = 99
	Collation_cp1251_bin                  CollationID = 50
	Collation_cp1251_bulgarian_ci         CollationID = 14
	Collation_cp1251_general_ci           CollationID = 51
	Collation_cp1251_general_cs           CollationID = 52
	Collation_cp1251_ukrainian_ci         CollationID = 23
	Collation_cp1256_bin                  CollationID = 67
	Collation_cp1256_general_ci           CollationID = 57
	Collation_cp1257_bin                  CollationID = 58
	Collation_cp1257_general_ci           CollationID = 59
	Collation_cp1257_lithuanian_ci        CollationID = 29
	Collation_cp850_bin                   CollationID = 80
	Collation_cp850_general_ci            CollationID = 4
	Collation_cp852_bin                   CollationID = 81
	Collation_cp852_general_ci            CollationID = 40
	Collation_cp866_bin                   CollationID = 68
	Collation_cp866_general_ci            CollationID = 36
	Collation_cp932_bin                   CollationID = 96
	Collation_cp932_japanese_ci           CollationID = 95
	Collation_dec8_bin                    CollationID = 69
	Collation_dec8_swedish_ci             CollationID = 3
	Collation_eucjpms_bin                 CollationID = 98
	Collation_eucjpms_japanese_ci         CollationID = 97
	Collation_euckr_bin                   CollationID = 85
	Collation_euckr_korean_ci             CollationID = 19
	Collation_gb18030_bin                 CollationID = 249
	Collation_gb18030_chinese_ci          CollationID = 248
	Collation_gb18030_unicode_520_ci      CollationID = 250
	Collation_gb2312_bin                  CollationID = 86
	Collation_gb2312_chinese_ci           CollationID = 24
	Collation_gbk_bin                     CollationID = 87
	Collation_gbk_chinese_ci              CollationID = 28
	Collation_geostd8_bin                 CollationID = 93
	Collation_geostd8_general_ci          CollationID = 92
	Collation_greek_bin                   CollationID = 70
	Collation_greek_general_ci            CollationID = 25
	Collation_hebrew_bin                  CollationID = 71
	Collation_hebrew_general_ci           CollationID = 16
	Collation_hp8_bin                     CollationID = 72
	Collation_hp8_english_ci              CollationID = 6
	Collation_keybcs2_bin                 CollationID = 73
	Collation_keybcs2_general_ci          CollationID = 37
	Collation_koi8r_bin                   CollationID = 74
	Collation_koi8r_general_ci            CollationID = 7
	Collation_koi8u_bin                   CollationID = 75
	Collation_koi8u_general_ci            CollationID = 22
	Collation_latin1_bin                  CollationID = 47
	Collation_latin1_danish_ci            CollationID = 15
	Collation_latin1_general_ci           CollationID = 48
	Collation_latin1_general_cs           CollationID = 49
	Collation_latin1_german1_ci           CollationID = 5
	Collation_latin1_german2_ci           CollationID = 31
	Collation_latin1_spanish_ci           CollationID = 94
	Collation_latin1_swedish_ci           CollationID = 8
	Collation_latin2_bin                  CollationID = 77
	Collation_latin2_croatian_ci          CollationID = 27
	Collation_latin2_czech_cs             CollationID = 2
	Collation_latin2_general_ci           CollationID = 9
	Collation_latin2_hungarian_ci         CollationID = 21
	Collation_latin5_bin                  CollationID = 78
	Collation_latin5_turkish_ci           CollationID = 30
	Collation_latin7_bin                  CollationID = 79
	Collation_latin7_estonian_cs          CollationID = 20
	Collation_latin7_general_ci           CollationID = 41
	Collation_latin7_general_cs           CollationID = 42
	Collation_macce_bin                   CollationID = 43
	Collation_macce_general_ci            CollationID = 38
	Collation_macroman_bin                CollationID = 53
	Collation_macroman_general_ci         CollationID = 39
	Collation_sjis_bin                    CollationID = 88
	Collation_sjis_japanese_ci            CollationID = 13
	Collation_swe7_bin                    CollationID = 82
	Collation_swe7_swedish_ci             CollationID = 10
	Collation_tis620_bin                  CollationID = 89
	Collation_tis620_thai_ci              CollationID = 18
	Collation_ucs2_bin                    CollationID = 90
	Collation_ucs2_croatian_ci            CollationID = 149
	Collation_ucs2_czech_ci               CollationID = 138
	Collation_ucs2_danish_ci              CollationID = 139
	Collation_ucs2_esperanto_ci           CollationID = 145
	Collation_ucs2_estonian_ci            CollationID = 134
	Collation_ucs2_general_ci             CollationID = 35
	Collation_ucs2_general_mysql500_ci    CollationID = 159
	Collation_ucs2_german2_ci             CollationID = 148
	Collation_ucs2_hungarian_ci           CollationID = 146
	Collation_ucs2_icelandic_ci           CollationID = 129
	Collation_ucs2_latvian_ci             CollationID = 130
	Collation_ucs2_lithuanian_ci          CollationID = 140
	Collation_ucs2_persian_ci             CollationID = 144
	Collation_ucs2_polish_ci              CollationID = 133
	Collation_ucs2_roman_ci               CollationID = 143
	Collation_ucs2_romanian_ci            CollationID = 131
	Collation_ucs2_sinhala_ci             CollationID = 147
	Collation_ucs2_slovak_ci              CollationID = 141
	Collation_ucs2_slovenian_ci           CollationID = 132
	Collation_ucs2_spanish2_ci            CollationID = 142
	Collation_ucs2_spanish_ci             CollationID = 135
	Collation_ucs2_swedish_ci             CollationID = 136
	Collation_ucs2_turkish_ci             CollationID = 137
	Collation_ucs2_unicode_520_ci         CollationID = 150
	Collation_ucs2_unicode_ci             CollationID = 128
	Collation_ucs2_vietnamese_ci          CollationID = 151
	Collation_ujis_bin                    CollationID = 91
	Collation_ujis_japanese_ci            CollationID = 12
	Collation_utf16_bin                   CollationID = 55
	Collation_utf16_croatian_ci           CollationID = 122
	Collation_utf16_czech_ci              CollationID = 111
	Collation_utf16_danish_ci             CollationID = 112
	Collation_utf16_esperanto_ci          CollationID = 118
	Collation_utf16_estonian_ci           CollationID = 107
	Collation_utf16_general_ci            CollationID = 54
	Collation_utf16_german2_ci            CollationID = 121
	Collation_utf16_hungarian_ci          CollationID = 119
	Collation_utf16_icelandic_ci          CollationID = 102
	Collation_utf16_latvian_ci            CollationID = 103
	Collation_utf16_lithuanian_ci         CollationID = 113
	Collation_utf16_persian_ci            CollationID = 117
	Collation_utf16_polish_ci             CollationID = 106
	Collation_utf16_roman_ci              CollationID = 116
	Collation_utf16_romanian_ci           CollationID = 104
	Collation_utf16_sinhala_ci            CollationID = 120
	Collation_utf16_slovak_ci             CollationID = 114
	Collation_utf16_slovenian_ci          CollationID = 105
	Collation_utf16_spanish2_ci           CollationID = 115
	Collation_utf16_spanish_ci            CollationID = 108
	Collation_utf16_swedish_ci            CollationID = 109
	Collation_utf16_turkish_ci            CollationID = 110
	Collation_utf16_unicode_520_ci        CollationID = 123
	Collation_utf16_unicode_ci            CollationID = 101
	Collation_utf16_vietnamese_ci         CollationID = 124
	Collation_utf16le_bin                 CollationID = 62
	Collation_utf16le_general_ci          CollationID = 56
	Collation_utf32_bin                   CollationID = 61
	Collation_utf32_croatian_ci           CollationID = 181
	Collation_utf32_czech_ci              CollationID = 170
	Collation_utf32_danish_ci             CollationID = 171
	Collation_utf32_esperanto_ci          CollationID = 177
	Collation_utf32_estonian_ci           CollationID = 166
	Collation_utf32_general_ci            CollationID = 60
	Collation_utf32_german2_ci            CollationID = 180
	Collation_utf32_hungarian_ci          CollationID = 178
	Collation_utf32_icelandic_ci          CollationID = 161
	Collation_utf32_latvian_ci            CollationID = 162
	Collation_utf32_lithuanian_ci         CollationID = 172
	Collation_utf32_persian_ci            CollationID = 176
	Collation_utf32_polish_ci             CollationID = 165
	Collation_utf32_roman_ci              CollationID = 175
	Collation_utf32_romanian_ci           CollationID = 163
	Collation_utf32_sinhala_ci            CollationID = 179
	Collation_utf32_slovak_ci             CollationID = 173
	Collation_utf32_slovenian_ci          CollationID = 164
	Collation_utf32_spanish2_ci           CollationID = 174
	Collation_utf32_spanish_ci            CollationID = 167
	Collation_utf32_swedish_ci            CollationID = 168
	Collation_utf32_turkish_ci            CollationID = 169
	Collation_utf32_unicode_520_ci        CollationID = 182
	Collation_utf32_unicode_ci            CollationID = 160
	Collation_utf32_vietnamese_ci         CollationID = 183
	Collation_utf8mb3_bin                 CollationID = 83
	Collation_utf8mb3_croatian_ci         CollationID = 213
	Collation_utf8mb3_czech_ci            CollationID = 202
	Collation_utf8mb3_danish_ci           CollationID = 203
	Collation_utf8mb3_esperanto_ci        CollationID = 209
	Collation_utf8mb3_estonian_ci         CollationID = 198
	Collation_utf8mb3_general_ci          CollationID = 33
	Collation_utf8mb3_general_mysql500_ci CollationID = 223
	Collation_utf8mb3_german2_ci          CollationID = 212
	Collation_utf8mb3_hungarian_ci        CollationID = 210
	Collation_utf8mb3_icelandic_ci        CollationID = 193
	Collation_utf8mb3_latvian_ci          CollationID = 194
	Collation_utf8mb3_lithuanian_ci       CollationID = 204
	Collation_utf8mb3_persian_ci          CollationID = 208
	Collation_utf8mb3_polish_ci           CollationID = 197
	Collation_utf8mb3_roman_ci            CollationID = 207
	Collation_utf8mb3_romanian_ci         CollationID = 195
	Collation_utf8mb3_sinhala_ci          CollationID = 211
	Collation_utf8mb3_slovak_ci           CollationID = 205
	Collation_utf8mb3_slovenian_ci        CollationID = 196
	Collation_utf8mb3_spanish2_ci         CollationID = 206
	Collation_utf8mb3_spanish_ci          CollationID = 199
	Collation_utf8mb3_swedish_ci          CollationID = 200
	Collation_utf8mb3_tolower_ci          CollationID = 76
	Collation_utf8mb3_turkish_ci          CollationID = 201
	Collation_utf8mb3_unicode_520_ci      CollationID = 214
	Collation_utf8mb3_unicode_ci          CollationID = 192
	Collation_utf8mb3_vietnamese_ci       CollationID = 215
	Collation_utf8mb4_0900_ai_ci          CollationID = 255
	Collation_utf8mb4_0900_as_ci          CollationID = 305
	Collation_utf8mb4_0900_as_cs          CollationID = 278
	Collation_utf8mb4_0900_bin            CollationID = 309
	Collation_utf8mb4_bg_0900_ai_ci       CollationID = 318
	Collation_utf8mb4_bg_0900_as_cs       CollationID = 319
	Collation_utf8mb4_bin                 CollationID = 46
	Collation_utf8mb4_bs_0900_ai_ci       CollationID = 316
	Collation_utf8mb4_bs_0900_as_cs       CollationID = 317
	Collation_utf8mb4_croatian_ci         CollationID = 245
	Collation_utf8mb4_cs_0900_ai_ci       CollationID = 266
	Collation_utf8mb4_cs_0900_as_cs       CollationID = 289
	Collation_utf8mb4_czech_ci            CollationID = 234
	Collation_utf8mb4_da_0900_ai_ci       CollationID = 267
	Collation_utf8mb4_da_0900_as_cs       CollationID = 290
	Collation_utf8mb4_danish_ci           CollationID = 235
	Collation_utf8mb4_de_pb_0900_ai_ci    CollationID = 256
	Collation_utf8mb4_de_pb_0900_as_cs    CollationID = 279
	Collation_utf8mb4_eo_0900_ai_ci       CollationID = 273
	Collation_utf8mb4_eo_0900_as_cs       CollationID = 296
	Collation_utf8mb4_es_0900_ai_ci       CollationID = 263
	Collation_utf8mb4_es_0900_as_cs       CollationID = 286
	Collation_utf8mb4_es_trad_0900_ai_ci  CollationID = 270
	Collation_utf8mb4_es_trad_0900_as_cs  CollationID = 293
	Collation_utf8mb4_esperanto_ci        CollationID = 241
	Collation_utf8mb4_estonian_ci         CollationID = 230
	Collation_utf8mb4_et_0900_ai_ci       CollationID = 262
	Collation_utf8mb4_et_0900_as_cs       CollationID = 285
	Collation_utf8mb4_general_ci          CollationID = 45
	Collation_utf8mb4_german2_ci          CollationID = 244
	Collation_utf8mb4_gl_0900_ai_ci       CollationID = 320
	Collation_utf8mb4_gl_0900_as_cs       CollationID = 321
	Collation_utf8mb4_hr_0900_ai_ci       CollationID = 275
	Collation_utf8mb4_hr_0900_as_cs       CollationID = 298
	Collation_utf8mb4_hu_0900_ai_ci       CollationID = 274
	Collation_utf8mb4_hu_0900_as_cs       CollationID = 297
	Collation_utf8mb4_hungarian_ci        CollationID = 242
	Collation_utf8mb4_icelandic_ci        CollationID = 225
	Collation_utf8mb4_is_0900_ai_ci       CollationID = 257
	Collation_utf8mb4_is_0900_as_cs       CollationID = 280
	Collation_utf8mb4_ja_0900_as_cs       CollationID = 303
	Collation_utf8mb4_ja_0900_as_cs_ks    CollationID = 304
	Collation_utf8mb4_la_0900_ai_ci       CollationID = 271
	Collation_utf8mb4_la_0900_as_cs       CollationID = 294
	Collation_utf8mb4_latvian_ci          CollationID = 226
	Collation_utf8mb4_lithuanian_ci       CollationID = 236
	Collation_utf8mb4_lt_0900_ai_ci       CollationID = 268
	Collation_utf8mb4_lt_0900_as_cs       CollationID = 291
	Collation_utf8mb4_lv_0900_ai_ci       CollationID = 258
	Collation_utf8mb4_lv_0900_as_cs       CollationID = 281
	Collation_utf8mb4_mn_cyrl_0900_ai_ci  CollationID = 322
	Collation_utf8mb4_mn_cyrl_0900_as_cs  CollationID = 323
	Collation_utf8mb4_nb_0900_ai_ci       CollationID = 310
	Collation_utf8mb4_nb_0900_as_cs       CollationID = 311
	Collation_utf8mb4_nn_0900_ai_ci       CollationID = 312
	Collation_utf8mb4_nn_0900_as_cs       CollationID = 313
	Collation_utf8mb4_persian_ci          CollationID = 240
	Collation_utf8mb4_pl_0900_ai_ci       CollationID = 261
	Collation_utf8mb4_pl_0900_as_cs       CollationID = 284
	Collation_utf8mb4_polish_ci           CollationID = 229
	Collation_utf8mb4_ro_0900_ai_ci       CollationID = 259
	Collation_utf8mb4_ro_0900_as_cs       CollationID = 282
	Collation_utf8mb4_roman_ci            CollationID = 239
	Collation_utf8mb4_romanian_ci         CollationID = 227
	Collation_utf8mb4_ru_0900_ai_ci       CollationID = 306
	Collation_utf8mb4_ru_0900_as_cs       CollationID = 307
	Collation_utf8mb4_sinhala_ci          CollationID = 243
	Collation_utf8mb4_sk_0900_ai_ci       CollationID = 269
	Collation_utf8mb4_sk_0900_as_cs       CollationID = 292
	Collation_utf8mb4_sl_0900_ai_ci       CollationID = 260
	Collation_utf8mb4_sl_0900_as_cs       CollationID = 283
	Collation_utf8mb4_slovak_ci           CollationID = 237
	Collation_utf8mb4_slovenian_ci        CollationID = 228
	Collation_utf8mb4_spanish2_ci         CollationID = 238
	Collation_utf8mb4_spanish_ci          CollationID = 231
	Collation_utf8mb4_sr_latn_0900_ai_ci  CollationID = 314
	Collation_utf8mb4_sr_latn_0900_as_cs  CollationID = 315
	Collation_utf8mb4_sv_0900_ai_ci       CollationID = 264
	Collation_utf8mb4_sv_0900_as_cs       CollationID = 287
	Collation_utf8mb4_swedish_ci          CollationID = 232
	Collation_utf8mb4_tr_0900_ai_ci       CollationID = 265
	Collation_utf8mb4_tr_0900_as_cs       CollationID = 288
	Collation_utf8mb4_turkish_ci          CollationID = 233
	Collation_utf8mb4_unicode_520_ci      CollationID = 246
	Collation_utf8mb4_unicode_ci          CollationID = 224
	Collation_utf8mb4_vi_0900_ai_ci       CollationID = 277
	Collation_utf8mb4_vi_0900_as_cs       CollationID = 300
	Collation_utf8mb4_vietnamese_ci       CollationID = 247
	Collation_utf8mb4_zh_0900_as_cs       CollationID = 308

	Collation_utf8_general_ci          = Collation_utf8mb3_general_ci
	Collation_utf8_tolower_ci          = Collation_utf8mb3_tolower_ci
	Collation_utf8_bin                 = Collation_utf8mb3_bin
	Collation_utf8_unicode_ci          = Collation_utf8mb3_unicode_ci
	Collation_utf8_icelandic_ci        = Collation_utf8mb3_icelandic_ci
	Collation_utf8_latvian_ci          = Collation_utf8mb3_latvian_ci
	Collation_utf8_romanian_ci         = Collation_utf8mb3_romanian_ci
	Collation_utf8_slovenian_ci        = Collation_utf8mb3_slovenian_ci
	Collation_utf8_polish_ci           = Collation_utf8mb3_polish_ci
	Collation_utf8_estonian_ci         = Collation_utf8mb3_estonian_ci
	Collation_utf8_spanish_ci          = Collation_utf8mb3_spanish_ci
	Collation_utf8_swedish_ci          = Collation_utf8mb3_swedish_ci
	Collation_utf8_turkish_ci          = Collation_utf8mb3_turkish_ci
	Collation_utf8_czech_ci            = Collation_utf8mb3_czech_ci
	Collation_utf8_danish_ci           = Collation_utf8mb3_danish_ci
	Collation_utf8_lithuanian_ci       = Collation_utf8mb3_lithuanian_ci
	Collation_utf8_slovak_ci           = Collation_utf8mb3_slovak_ci
	Collation_utf8_spanish2_ci         = Collation_utf8mb3_spanish2_ci
	Collation_utf8_roman_ci            = Collation_utf8mb3_roman_ci
	Collation_utf8_persian_ci          = Collation_utf8mb3_persian_ci
	Collation_utf8_esperanto_ci        = Collation_utf8mb3_esperanto_ci
	Collation_utf8_hungarian_ci        = Collation_utf8mb3_hungarian_ci
	Collation_utf8_sinhala_ci          = Collation_utf8mb3_sinhala_ci
	Collation_utf8_german2_ci          = Collation_utf8mb3_german2_ci
	Collation_utf8_croatian_ci         = Collation_utf8mb3_croatian_ci
	Collation_utf8_unicode_520_ci      = Collation_utf8mb3_unicode_520_ci
	Collation_utf8_vietnamese_ci       = Collation_utf8mb3_vietnamese_ci
	Collation_utf8_general_mysql500_ci = Collation_utf8mb3_general_mysql500_ci

	Collation_Default                    = Collation_utf8mb4_0900_bin
	Collation_Information_Schema_Default = Collation_utf8mb3_general_ci
	// Collation_Unspecified is used when a collation has not been specified, either explicitly or implicitly. This is
	// usually used as an intermediate collation to be later replaced by an analyzer pass or a plan, although it is
	// valid to use it directly. When used, behaves identically to the default collation, although it will NOT match
	// the default collation.
	Collation_Unspecified CollationID = 0
)

func ParseCollation(characterSetStr string, collationStr string, binary bool) (CollationID, error) {
	if binary {
		if characterSetStr == "binary" {
			return Collation_binary, nil
		}
		return CharacterSet_utf8mb4.BinaryCollation(), nil
	}
	if collationStr != "" {
		if id, ok := collationStringToID[strings.ToLower(collationStr)]; ok {
			return id, nil
		}
		return Collation_Default, fmt.Errorf("unknown collation %q", collationStr)
	}
	if characterSetStr != "" {
		cs, err := ParseCharacterSet(characterSetStr)
		if err != nil {
			return Collation_Default, err
		}
		return cs.DefaultCollation(), nil
	}
	return Collation_Default, nil
}

func (c CollationID) Name() string {
	switch c {
	case Collation_binary:
		return "binary"
	case Collation_utf8mb4_bin:
		return "utf8mb4_bin"
	case Collation_utf8mb4_general_ci:
		return "utf8mb4_general_ci"
	case Collation_utf8mb3_general_ci:
		return "utf8mb3_general_ci"
	case Collation_Unspecified:
		return ""
	default:
		return "utf8mb4_0900_ai_ci"
	}
}

func (c CollationID) CharacterSet() CharacterSetID {
	if c == Collation_binary {
		return CharacterSet_binary
	}
	if c == Collation_Unspecified {
		return CharacterSet_Unspecified
	}
	return CharacterSet_utf8mb4
}

func (c CollationID) WorksWithCharacterSet(cs CharacterSetID) bool {
	if c == Collation_Unspecified || cs == CharacterSet_Unspecified {
		return true
	}
	return c.CharacterSet() == cs
}

func (c CollationID) String() string { return c.Name() }
func (c CollationID) IsDefault() string {
	if c == c.CharacterSet().DefaultCollation() {
		return "Yes"
	}
	return ""
}
func (c CollationID) IsCompiled() string            { return "Yes" }
func (c CollationID) SortLength() uint32            { return 1 }
func (c CollationID) PadAttribute() string          { return "NO PAD" }
func (c CollationID) Equals(other CollationID) bool { return c == other }

func (c CollationID) Collation() Collation {
	return Collation{ID: c, Name: c.Name(), CharacterSet: c.CharacterSet(), IsDefault: c == c.CharacterSet().DefaultCollation(), IsCompiled: true, SortLength: 1, PadAttribute: c.PadAttribute(), Sorter: c.Sorter()}
}

func (c CollationID) WriteWeightString(hash io.Writer, str string) error {
	_, err := io.WriteString(hash, str)
	return err
}

func (c CollationID) HashToUint(str string) (uint64, error) {
	h := fnv.New64a()
	if err := c.WriteWeightString(h, str); err != nil {
		return 0, err
	}
	return h.Sum64(), nil
}

func (c CollationID) HashToBytes(str string) ([]byte, error) {
	h := fnv.New64a()
	if err := c.WriteWeightString(h, str); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

func (c CollationID) Sorter() CollationSorter {
	return func(r rune) int32 { return int32(r) }
}

var liteCollations = []CollationID{Collation_binary, Collation_utf8mb4_0900_ai_ci, Collation_utf8mb4_bin, Collation_utf8mb4_general_ci, Collation_utf8mb3_general_ci}

func NewCollationsIterator() *CollationsIterator { return &CollationsIterator{} }

func (ci *CollationsIterator) Next() (Collation, bool) {
	if ci.idx >= len(liteCollations) {
		return Collation{}, false
	}
	c := liteCollations[ci.idx]
	ci.idx++
	return c.Collation(), true
}

func ConvertCollationID(val any) (string, error) {
	switch v := val.(type) {
	case CollationID:
		return v.Name(), nil
	case string:
		return v, nil
	default:
		return "", fmt.Errorf("invalid collation value %T", val)
	}
}

// TypeWithCollation is implemented on all types that may return a collation.
type TypeWithCollation interface {
	Collation() CollationID
	WithNewCollation(collation CollationID) (Type, error)
	StringWithTableCollation(tableCollation CollationID) string
}

func compareRune(c CollationID, left string, right string) int {
	lr, _ := utf8.DecodeRuneInString(left)
	rr, _ := utf8.DecodeRuneInString(right)
	if lr < rr {
		return -1
	}
	if lr > rr {
		return 1
	}
	return 0
}
