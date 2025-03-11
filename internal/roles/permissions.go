package roles

type Permission int64

const (
	CanCRUDPersonalIssues = 1 << iota
	CanCommentPersonalIssues
	CanCRUDPersonalComments
	CanClosePersonalIssues

	CanUpdateUserRole

	CanCommentExternalIssues
	CanCloseExternalIssues

	CanSeeIssuesList
	CanCollectIssuesStatistics
)

func HasPermission(mask int64, permission Permission) bool {
	return mask&int64(permission) != 0
}

func AddPermission(mask int64, permission Permission) int64 {
	return mask | int64(permission)
}
