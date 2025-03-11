package usecases

type UserModel struct {
	ID int32
	Login string
	PasswordHash string
	Role string
	PermissionMask int64
}