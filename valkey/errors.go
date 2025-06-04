package valkey

import (
	"github.com/joomcode/errorx"
)

var ValKeyNamespace = errorx.NewNamespace("valkey")
var ValKeyTrait = errorx.RegisterTrait("valkey")
var ValKeyError = errorx.NewType(ValKeyNamespace, "valkey_error", ValKeyTrait)
var NonValKeyError = errorx.NewType(ValKeyNamespace, "non_valkey_error")
