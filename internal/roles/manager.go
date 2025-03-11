package roles

import (
	"context"
	"fmt"
	"log/slog"
)

var permKeys = map[string]Permission{
	"crud_personal_issues":      CanCRUDPersonalIssues,
	"comment_personal_issues":   CanCommentPersonalIssues,
	"crud_personal_comments":    CanCRUDPersonalComments,
	"close_personal_issues":     CanClosePersonalIssues,
	"update_user_role":          CanUpdateUserRole,
	"comment_external_issues":   CanCommentExternalIssues,
	"close_external_issues":     CanCloseExternalIssues,
	"see_issues_list":           CanSeeIssuesList,
	"collect_issues_statistics": CanCollectIssuesStatistics,
}

type RoleStorage interface {
	UpsertRole(
		ctx context.Context,
		alias string,
		mask int64,
		isDefault bool,
		isSuper bool,
	) (id int32, err error)
}

type RoleConfig struct {
	ID    int32
	Alias string
}

type RolesManager struct {
	log     *slog.Logger
	storage RoleStorage
}

func NewManager(
	log *slog.Logger,
	storage RoleStorage,
) *RolesManager {
	return &RolesManager{
		log:     log,
		storage: storage,
	}
}

func (r *RolesManager) CreateRole(
	ctx context.Context,
	alias string,
	permissions []string,
	isDefault bool,
	isSuper bool,
) (err error) {
	const src = "RolesManager.CreateRole"
	log := r.log.With(slog.String("src", src))

	mask, countPerms := r.maskFromPermArray(permissions)

	// ignoring role id
	_, err = r.storage.UpsertRole(ctx, alias, mask, isDefault, isSuper)
	if err != nil {
		return fmt.Errorf("%s: failed to save %s role: %w", src, alias, err)
	}

	log.Info("role indexed",
		slog.String("alias", alias),
		slog.Int("permissions_granted", countPerms),
		slog.Int("total_permissions", len(permissions)),
	)

	return nil
}

func (r *RolesManager) maskFromPermArray(permissions []string) (mask int64, cnt int) {
	mask = 0
	cnt = 0

	for _, key := range permissions {
		if perm, exists := permKeys[key]; exists {
			mask = AddPermission(mask, perm)
			cnt += 1
		}
	}

	return mask, cnt
}
