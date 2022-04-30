package db

import "lproxy/internal/auth/providers"

type builder struct {
	userRepo    *userRepository
	sessionRepo *sessionRepository
	roleRepo    *roleRepository
}

func (b *builder) userEntityByUserId(userId string) (providers.User, error) {
	user, userErr := b.userRepo.findById(userId)
	if userErr != nil {
		return providers.User{}, userErr
	}

	return b.userEntityFromDbEntity(*user)
}

func (b *builder) userEntityFromDbEntity(user User) (providers.User, error) {
	roles, rolesErr := b.roleRepo.findForUser(user.Id)
	if rolesErr != nil {
		return providers.User{}, rolesErr
	}
	var sessionRoles []providers.Role
	for _, v := range *roles {
		sessionRoles = append(sessionRoles, providers.Role{
			Id:     v.Id,
			Name:   v.Name,
			Code:   v.Code,
			Source: v.Source,
		})
	}

	return providers.User{
		Id:       user.Id,
		Username: user.Username,
		Email:    user.Email,
		Phone:    user.Phone,
		Provider: ProviderType,
		Source:   ProviderType,
		Roles:    sessionRoles,
	}, nil
}
