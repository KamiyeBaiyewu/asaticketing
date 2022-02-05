package auth

// type Permissions interface {
// 	Wrap(next http.HandlerFunc, permissionTypes ...PermissionType) http.HandlerFunc
// 	Check(r *http.Request, permissionTypes ...PermissionType) bool
// }

// type permissions struct {
// 	DB    database.Database
// 	cache gcache.Cache
// }

// // NewPermission crates a new permission structure
// func NewPermission(db database.Database) Permissions {

// 	p := &permissions{
// 		DB: db,
// 	}
// 	p.cache = gcache.New(200).LRU().LoaderExpireFunc(func(key interface{}) (interface{}, *time.Duration, error) {
// 		userID := key.(model.UserID)
// 		roles, err := p.DB.GetRoleByUser(context.Background(), userID)
// 		if err != nil {
// 			return nil, nil, err
// 		}
// 		expire := 1 * time.Minute
// 		return roles, &expire, nil
// 	}).Build()

// 	return p

// }

// func (p *permissions) getRoles(userID model.UserID) ([]*model.Role, error) {
// 	roles, err := p.cache.Get(userID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return roles.([]*model.Role), nil
// }

// func (p *permissions) withRoles(principal model.Principal, roleFunc func([]*model.Role) bool) (bool, error) {
// 	if principal.UserID == model.NilUserID {
// 		return false, nil
// 	}

// 	roles, err := p.getRoles(principal.UserID)
// 	if err != nil {
// 		return false, err
// 	}
// 	return roleFunc(roles), nil
// }

// // Wrap - reserves certain endpoint to a role
// func (p *permissions) Wrap(next http.HandlerFunc, permissionTypes ...PermissionType) http.HandlerFunc {

// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 		if allowed := p.Check(r, permissionTypes...); !allowed {
// 			utils.WriteError(w, http.StatusUnauthorized, "Permission Denied", nil)
// 			return
// 		}

// 		next.ServeHTTP(w, r)
// 	})
// }

// func (p *permissions) Check(r *http.Request, permissionTypes ...PermissionType) bool {

// 	principal := GetPrincipal(r)

// 	for _, permissionType := range permissionTypes {
// 		switch permissionType {
// 		case Admin:
// 			if allowed, _ := p.withRoles(principal, adminOnly); allowed {
// 				return true
// 			}
// 		case MemberIsTarget:
// 			targetID := model.UserID(mux.Vars(r)["userID"])
// 			if allowed := memberIsTarget(targetID, principal); allowed {
// 				return true
// 			}
// 		case Member:
// 			if allowed := member(principal); allowed {
// 				return true
// 			}
// 		case Anyone:
// 			if allowed := any(); allowed {
// 				return true
// 			}
// 		}
// 	}
// 	return false

// }
