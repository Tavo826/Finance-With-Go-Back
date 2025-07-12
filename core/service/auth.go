package service

import (
	"context"
	"mime/multipart"
	"personal-finance/core/domain"
	"personal-finance/core/port"
)

type AuthService struct {
	authRepo        port.AuthRepository
	transactionRepo port.TransactionRepository
	adapter         port.ImageAdapter
}

func NewAuthService(authRepo port.AuthRepository, transactionRepo port.TransactionRepository, adapter port.ImageAdapter) *AuthService {
	return &AuthService{
		authRepo,
		transactionRepo,
		adapter,
	}
}

func (as *AuthService) GetAllUsers(ctx context.Context) ([]domain.User, error) {

	users, err := as.authRepo.GetAllUsers(ctx)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		if err.Error() == domain.ErrNoDocuments.Error() {
			return nil, domain.ErrDataNotFound
		}
		return nil, domain.ErrInternal
	}

	return users, nil
}

func (as *AuthService) GetUserById(ctx context.Context, id string) (*domain.User, error) {

	user, err := as.authRepo.GetUserById(ctx, id)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		if err.Error() == domain.ErrNoDocuments.Error() {
			return nil, domain.ErrDataNotFound
		}
		return nil, domain.ErrInternal
	}

	return user, nil
}

func (as *AuthService) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {

	user, err := as.authRepo.CreateUser(ctx, user)
	if err != nil {
		if err == domain.ErrConflictingData {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return user, nil
}

func (as *AuthService) UpdateUser(ctx context.Context, id string, user *domain.User) (*domain.User, error) {

	_, err := as.authRepo.UpdateUser(ctx, id, user)
	if err != nil {
		if err == domain.ErrConflictingData {
			return nil, err
		}
		if err == domain.ErrDataNotFound {
			return nil, domain.ErrDataNotFound
		}
		return nil, domain.ErrInternal
	}

	return user, nil
}

func (as *AuthService) UpdateUserProfileImage(ctx context.Context, file multipart.File, userId string) (*domain.Image, error) {

	uploadedImage, err := as.adapter.UploadImageFromFile(ctx, file, userId)
	if err != nil {
		return nil, err
	}

	return uploadedImage, nil
}

func (as *AuthService) DeleteUserProfileImage(ctx context.Context, publicId string) error {

	return as.adapter.DeleteImage(ctx, publicId)
}

func (as *AuthService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {

	user, err := as.authRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return user, nil
}

func (as *AuthService) VerifyUserEmail(ctx context.Context, email string) (bool, error) {

	_, err := as.authRepo.GetUserByEmail(ctx, email)

	if err != nil {
		if err == domain.ErrDataNotFound {
			return false, nil
		}
		if err.Error() == domain.ErrNoDocuments.Error() {
			return false, nil
		}
		return true, domain.ErrInternal
	}

	return true, nil
}

func (as *AuthService) DeleteUser(ctx context.Context, id string) error {

	return as.authRepo.DeleteUser(ctx, id)
}

func (as *AuthService) DeleteTransactionsByUserId(ctx context.Context, id string) error {

	return as.transactionRepo.DeleteTransactionsByUserId(ctx, id)
}
