package valkey

import (
	"errors"
	"fmt"

	errs "github.com/jarrodhroberson/ossgo/errors"
	vk "github.com/valkey-io/valkey-go"
)

func NewKeyFunc(src string, objType string) KeyFunc {
	return func(key string) string {
		return fmt.Sprintf("%s:%s:{%s}",src, objType, key)
	}
}

func NewKeyFuncWith(src string, objType string, propName string) KeyFunc {
	return func(key string) string {
		return fmt.Sprintf("%s:%s:{%s}:%s", src, objType, key, propName)
	}
}

func ValkeyResultErrors(vkr vk.ValkeyResult) error {
	var err error
	if vkr.Error() != nil {
		if vkr.NonValkeyError() != nil {
			err = NonValKeyError.WrapWithNoMessage(vkr.NonValkeyError())
		}
		if errors.Is(vkr.Error(), vk.Nil) {
			err = errs.NotFoundError.WrapWithNoMessage(err)
		}
		err = ValKeyError.WrapWithNoMessage(vkr.Error()).WithUnderlyingErrors(err)
	}
	return err
}
