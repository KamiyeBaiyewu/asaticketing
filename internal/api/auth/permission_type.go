package auth

// type PermissionType string

// const (
// 	// Admin - User has 'admin' role
// 	Admin PermissionType = "admin"
// 	// Member - User is loged in (we have userID in principal)
// 	Member PermissionType = "member"
// 	// MemberIsTarget - User logged in and the user id passed to the API is the same
// 	MemberIsTarget PermissionType = "memberIsTarget"

// 	// Anyone - anyobe who is not logged in
// 	Anyone PermissionType = "anyone"
// )

// // Admin
// var adminOnly = func(roles []*model.Role) bool {
// 	for _, role := range roles {
// 		switch role.Role {
// 		case model.RoleAdmin:
// 			return true
// 		}
// 	}
// 	return false
// }

// // Logged in as User
// var member = func(principal model.Principal) bool {
// 	return principal.UserID != ""
// }

// // Logged in User - Target User
// var memberIsTarget = func(userID model.UserID, principal model.Principal) bool {
// 	if userID == "" || principal.UserID == "" {
// 		return false
// 	}
// 	if userID != principal.UserID {
// 		return false
// 	}

// 	return true
// }

// var any = func() bool {
// 	return true
// }
