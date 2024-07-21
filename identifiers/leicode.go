package identifiers

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"unicode"

	"github.com/joomcode/errorx"
)

const INVALID_VALUE = "\u0000"

var InvalidLeiCodeTrait = errorx.RegisterTrait("invalid_lei_code")
var InvalidLeiCode = errorx.IllegalArgument.NewSubtype("invalid_lei_code", InvalidLeiCodeTrait)
var InvalidFormat = errorx.IllegalFormat.NewSubtype("invalid_format", InvalidLeiCodeTrait)
var InvalidCheckSum = errorx.IllegalState.NewSubtype("InvalidCheckSum", InvalidLeiCodeTrait)
var InvalidLength = errorx.IllegalFormat.NewSubtype("invalid_length", InvalidLeiCodeTrait)

type LeiCode string

func (l LeiCode) LocalOperatingUnit() string {
	return string(l[0:4])
}

func (l LeiCode) EntityIdentifier() string {
	return string(l[6:18])
}

func (l LeiCode) String() string {
	return string(l)
}

func NewLeiCode(lou string, identifier string) (LeiCode, error) {
	if len(lou) != 4 {
		return INVALID_VALUE, errors.Join(InvalidLeiCode.NewWithNoMessage(), InvalidLength.New("lou must be 4 characters long; len(%s) = %d", lou, len(lou)))
	}
	if err := isUppercaseAlphaNumeric(lou); err != nil {
		return INVALID_VALUE, errors.Join(InvalidLeiCode.NewWithNoMessage(), err)
	}
	if len(identifier) != 12 {
		return INVALID_VALUE, errors.Join(InvalidLeiCode.NewWithNoMessage(), InvalidLength.New("entity identifier must be 12 characters long; len(%s) = %d", identifier, len(identifier)))
	}
	if err := isUppercaseAlphaNumeric(identifier); err != nil {
		return INVALID_VALUE, errors.Join(InvalidLeiCode.NewWithNoMessage(), err)
	}
	reservedCharPadding := "00"
	checksum := calculateChecksum(fmt.Sprintf("%s%s%s", lou, reservedCharPadding, identifier))
	return LeiCode(fmt.Sprintf("%s%s%s%s", lou, reservedCharPadding, identifier, checksum)), nil
}

func LeiCodeFrom(code string) (LeiCode, error) {
	leiCode := LeiCode(code)
	if err := isValidLeiCode(leiCode); err != nil {
		return INVALID_VALUE, err
	} else {
		return leiCode, nil
	}
}

func calculateChecksum(lei string) string {
	mods := ""
	for _, char := range lei {
		if char > 64 && char < 91 {
			char -= 55
			mods += strconv.Itoa(int(char))
		} else {
			mods += string(char)
		}
	}
	mods += "00"
	bigint, ok := new(big.Int).SetString(mods, 10)
	if !ok {
		panic(errorx.InternalError.New("error converting %s to bigint", mods))
	}
	ninetySeven := new(big.Int).SetInt64(97)
	var mod = big.NewInt(0)
	_, mod = new(big.Int).DivMod(bigint, ninetySeven, mod)
	ninetyEight := big.NewInt(98)
	return ninetyEight.Sub(ninetyEight, mod).String()
}

// IsValid checks if the format of the LEI code is valid
func isValidLeiCode(leiCode LeiCode) error {
	if len(leiCode) != 20 {
		return errors.Join(InvalidLeiCode.NewWithNoMessage(), InvalidLength.New("leiCode must be exactly 20 characters; %d", len(leiCode)))
	}

	if err := isFormatValid(leiCode); err != nil {
		return errors.Join(InvalidLeiCode.NewWithNoMessage(), err)
	}

	return isChecksumValid(leiCode)
}

// isFormatValid checks if a string contains only uppercase letters or digits
func isFormatValid(leiCode LeiCode) error {
	if err := isUppercaseAlphaNumeric(string(leiCode[:18])); err != nil {
		return err
	} else {
		return nil
	}
}

func isUppercaseAlphaNumeric(s string) error {
	for idx, c := range s {
		if !((c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')) {
			return InvalidFormat.New("%c at index %d must be Uppercase AlphaNumeric", c, idx)
		}
	}
	return nil
}

// isChecksumValid checks if the LEI code's checksum is valid
func isChecksumValid(leiCode LeiCode) error {
	if !unicode.IsDigit(rune(leiCode[18])) || !unicode.IsDigit(rune(leiCode[19])) {
		return errors.Join(InvalidLeiCode.NewWithNoMessage(), InvalidFormat.New("leiCode checksum must be digits; %s", leiCode[18:19]))
	}

	var m int64 = 0
	for _, c := range leiCode {
		if unicode.IsDigit(c) {
			m = (m*10 + int64(c) - 48) % 97
		} else {
			m = (m*100 + int64(c) - 55) % 97
		}
	}

	if m == 1 {
		return nil
	} else {
		return InvalidCheckSum.New("invalid checksum; %s", leiCode[18:19])
	}
}
