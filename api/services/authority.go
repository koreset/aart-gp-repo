package services

import (
	"api/models"
	"errors"
	"fmt"
	"strings"
)

// Authority matrix actions. Each new state transition that should be gated by
// role + monetary threshold lives here as a constant so the database, services,
// and frontend stay aligned.
const (
	AuthActionSignOffSchedule = "signoff_schedule"
	AuthActionFinanceReview   = "finance_review"
	AuthActionAuthoriseFirst  = "authorise_first"
	AuthActionAuthoriseSecond = "authorise_second"
	AuthActionArchive         = "archive_schedule"
)

// resolveRoleName returns the GPUserRole.RoleName for a user (looked up by
// username then email). Returns empty string when the user has no assigned
// role — callers should treat that as "deny" unless they're in bootstrap mode.
func resolveRoleName(user models.AppUser) (string, error) {
	if user.UserName == "" && user.UserEmail == "" {
		return "", nil
	}

	var orgUser models.OrgUser
	q := DB
	switch {
	case user.UserEmail != "":
		q = q.Where("email = ?", user.UserEmail)
	default:
		q = q.Where("name = ?", user.UserName)
	}
	if err := q.First(&orgUser).Error; err != nil {
		return "", nil
	}
	if orgUser.GPRoleId == 0 {
		return "", nil
	}

	var role models.GPUserRole
	if err := DB.Select("id, role_name").Where("id = ?", orgUser.GPRoleId).First(&role).Error; err != nil {
		return "", err
	}
	return role.RoleName, nil
}

// CanPerform reports whether `user` is permitted to perform `action` on a
// schedule whose NetTotal is `amount`. The check walks AuthorityMatrix rows
// for the user's role and returns true if any active row matches `action` and
// covers `amount` in its [MinAmount, MaxAmount] range. MaxAmount = -1 means
// "no upper bound".
//
// If no matrix rows exist at all (fresh install) the check returns true so
// the matrix can be configured without locking the org out — same shape as
// the permission system's bootstrap behaviour.
func CanPerform(user models.AppUser, action string, amount float64) (bool, error) {
	var totalRows int64
	if err := DB.Model(&models.AuthorityMatrix{}).Count(&totalRows).Error; err != nil {
		return false, err
	}
	if totalRows == 0 {
		return true, nil
	}

	role, err := resolveRoleName(user)
	if err != nil {
		return false, err
	}
	if role == "" {
		return false, nil
	}

	var rows []models.AuthorityMatrix
	if err := DB.Where("action = ? AND role = ? AND is_active = ?", action, role, true).Find(&rows).Error; err != nil {
		return false, err
	}
	for _, r := range rows {
		if amount < r.MinAmount {
			continue
		}
		if r.MaxAmount >= 0 && amount > r.MaxAmount {
			continue
		}
		return true, nil
	}
	return false, nil
}

// RequireAuthority is a thin wrapper that returns a typed error when the
// user doesn't meet the matrix. Keeps controller code readable.
func RequireAuthority(user models.AppUser, action string, amount float64) error {
	ok, err := CanPerform(user, action, amount)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("authority denied: %s is not authorised to %s for an amount of %.2f", strings.TrimSpace(user.UserName+" "+user.UserEmail), action, amount)
	}
	return nil
}

// ErrAuthorityDenied is returned by helpers that need to distinguish a
// missing-authority outcome from a database error.
var ErrAuthorityDenied = errors.New("authority denied")
