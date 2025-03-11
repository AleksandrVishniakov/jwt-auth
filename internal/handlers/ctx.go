package handlers

import (
	"context"
	"fmt"
	"strconv"
)

func UserIDFromContext(ctx context.Context) (int32, error) {
	id, err := strconv.Atoi(fmt.Sprintf("%v", ctx.Value(userIDKey)))
	if err != nil {
		return 0, err
	}

	return int32(id), nil
}

func PermissionMaskFromContext(ctx context.Context) (int64, error) {
	mask, err := strconv.Atoi(fmt.Sprintf("%v", ctx.Value(permissionMaskKey)))
	if err != nil {
		return 0, err
	}

	return int64(mask), nil
}