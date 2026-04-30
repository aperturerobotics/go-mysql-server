//go:build sql_lite

package sql

import (
	"fmt"
	"strings"

	"github.com/dolthub/go-mysql-server/sql/encodings"
)

// CharacterSet represents the character set of a string.
type CharacterSet struct {
	ID               CharacterSetID
	Name             string
	DefaultCollation CollationID
	BinaryCollation  CollationID
	Description      string
	MaxLength        uint8
	Encoder          encodings.Encoder
}

// CharacterSetsIterator iterates over every lite character set available.
type CharacterSetsIterator struct {
	idx int
}

// CharacterSetID represents a character set.
type CharacterSetID uint16

const (
	CharacterSet_armscii8 CharacterSetID = 32
	CharacterSet_ascii    CharacterSetID = 11
	CharacterSet_big5     CharacterSetID = 1
	CharacterSet_binary   CharacterSetID = 63
	CharacterSet_cp1250   CharacterSetID = 26
	CharacterSet_cp1251   CharacterSetID = 51
	CharacterSet_cp1256   CharacterSetID = 57
	CharacterSet_cp1257   CharacterSetID = 59
	CharacterSet_cp850    CharacterSetID = 4
	CharacterSet_cp852    CharacterSetID = 40
	CharacterSet_cp866    CharacterSetID = 36
	CharacterSet_cp932    CharacterSetID = 95
	CharacterSet_dec8     CharacterSetID = 3
	CharacterSet_eucjpms  CharacterSetID = 97
	CharacterSet_euckr    CharacterSetID = 19
	CharacterSet_gb18030  CharacterSetID = 248
	CharacterSet_gb2312   CharacterSetID = 24
	CharacterSet_gbk      CharacterSetID = 28
	CharacterSet_geostd8  CharacterSetID = 92
	CharacterSet_greek    CharacterSetID = 25
	CharacterSet_hebrew   CharacterSetID = 16
	CharacterSet_hp8      CharacterSetID = 6
	CharacterSet_keybcs2  CharacterSetID = 37
	CharacterSet_koi8r    CharacterSetID = 7
	CharacterSet_koi8u    CharacterSetID = 22
	CharacterSet_latin1   CharacterSetID = 8
	CharacterSet_latin2   CharacterSetID = 9
	CharacterSet_latin5   CharacterSetID = 30
	CharacterSet_latin7   CharacterSetID = 41
	CharacterSet_macce    CharacterSetID = 38
	CharacterSet_macroman CharacterSetID = 39
	CharacterSet_sjis     CharacterSetID = 13
	CharacterSet_swe7     CharacterSetID = 10
	CharacterSet_tis620   CharacterSetID = 18
	CharacterSet_ucs2     CharacterSetID = 35
	CharacterSet_ujis     CharacterSetID = 12
	CharacterSet_utf16    CharacterSetID = 54
	CharacterSet_utf16le  CharacterSetID = 56
	CharacterSet_utf32    CharacterSetID = 60
	CharacterSet_utf8mb3  CharacterSetID = 33
	CharacterSet_utf8mb4  CharacterSetID = 255

	CharacterSet_utf8 = CharacterSet_utf8mb3

	// CharacterSet_Unspecified is used when a character set has not been specified, either explicitly or implicitly.
	// This is usually used as an intermediate character set to be later replaced by an analyzer pass or a plan,
	// although it is valid to use it directly. When used, behaves identically to the character set belonging to the
	// default collation, although it will NOT match the aforementioned character set.
	CharacterSet_Unspecified CharacterSetID = 0
)

func ParseCharacterSet(str string) (CharacterSetID, error) {
	switch strings.ToLower(str) {
	case "", "utf8mb4", "utf8":
		return CharacterSet_utf8mb4, nil
	case "utf8mb3":
		return CharacterSet_utf8mb3, nil
	case "binary":
		return CharacterSet_binary, nil
	default:
		return CharacterSet_utf8mb4, nil
	}
}

func (cs CharacterSetID) Name() string {
	switch cs {
	case CharacterSet_binary:
		return "binary"
	case CharacterSet_utf8mb3:
		return "utf8mb3"
	case CharacterSet_Unspecified:
		return ""
	default:
		return "utf8mb4"
	}
}

func (cs CharacterSetID) DefaultCollation() CollationID {
	switch cs {
	case CharacterSet_binary:
		return Collation_binary
	case CharacterSet_utf8mb3:
		return Collation_utf8mb3_general_ci
	case CharacterSet_Unspecified:
		return Collation_Unspecified
	default:
		return Collation_utf8mb4_0900_ai_ci
	}
}

func (cs CharacterSetID) BinaryCollation() CollationID {
	if cs == CharacterSet_binary {
		return Collation_binary
	}
	return Collation_utf8mb4_bin
}

func (cs CharacterSetID) Description() string { return cs.Name() }
func (cs CharacterSetID) MaxLength() int64 {
	if cs == CharacterSet_binary {
		return 1
	}
	return 4
}
func (cs CharacterSetID) String() string { return cs.Name() }
func (cs CharacterSetID) Encoder() encodings.Encoder {
	if cs == CharacterSet_binary {
		return encodings.Binary
	}
	return encodings.Utf8mb4
}

var liteCharacterSets = []CharacterSetID{CharacterSet_binary, CharacterSet_utf8mb3, CharacterSet_utf8mb4}

// SupportedCharsets contains the character sets supported by the sql_lite profile.
var SupportedCharsets = liteCharacterSets

func NewCharacterSetsIterator() *CharacterSetsIterator { return &CharacterSetsIterator{} }

func (csi *CharacterSetsIterator) Next() (CharacterSet, bool) {
	if csi.idx >= len(liteCharacterSets) {
		return CharacterSet{}, false
	}
	cs := liteCharacterSets[csi.idx]
	csi.idx++
	return CharacterSet{ID: cs, Name: cs.Name(), DefaultCollation: cs.DefaultCollation(), BinaryCollation: cs.BinaryCollation(), Description: cs.Description(), MaxLength: uint8(cs.MaxLength()), Encoder: cs.Encoder()}, true
}

func ConvertCharacterSetID(val any) (string, error) {
	switch v := val.(type) {
	case CharacterSetID:
		return v.Name(), nil
	case string:
		return v, nil
	default:
		return "", fmt.Errorf("invalid character set value %T", val)
	}
}
