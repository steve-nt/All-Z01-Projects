package services

import (
	"context"
	"errors"
	"fmt"
	"forum/internal/backend/models"
	"forum/internal/backend/repository"
	"forum/internal/utils"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	users repository.UserRepository
}

func NewUserService(ur repository.UserRepository, post repository.PostRepository) *UserService {
	return &UserService{users: ur}
}

func (s *UserService) GetProfile(ctx context.Context, r *http.Request) (resp models.User, err error) {

	user := ctx.Value("user").(models.User)
	if user.UUID == "" {
		return resp, errors.New("User Action from Guest")
	}
	err = s.users.GetProfile(ctx, &user)
	if err != nil {
		return resp, err
	}
	return user, nil
}

func (s *UserService) CountUnseenNotifications(profresp *models.ProfileResponse) {

	counter := 0

	for _, notification := range profresp.Notifications {
		if !notification.Seen {
			counter++
		}
	}

	profresp.UnseenNotifications = counter
}

func (s *UserService) SeeNotification(ctx context.Context, r *http.Request) error {
	user := ctx.Value("user").(models.User)
	if user.UUID == "" {
		return errors.New("User Action from Guest")
	}

	userId := user.UUID
	notId := strings.TrimPrefix(r.URL.Path, "/see-notification/")
	if err := s.users.SeeNotification(ctx, notId, userId); err != nil {
		return err
	}
	return nil
}

// ValidationError wraps multiple validation errors
type ValidationError struct {
	Errors []string
}

func (e ValidationError) Error() string {
	if len(e.Errors) == 0 {
		return "validation error"
	}
	return e.Errors[0] // Return first error for error interface
}

func (s *UserService) Register(ctx context.Context, mail, username, password, repeat_password, role string) error {
	clean, validationErrors := utils.SanitizeAndValidateNewUser(mail, username, password, repeat_password, role)
	if len(validationErrors) > 0 {
		return ValidationError{Errors: validationErrors}
	}

	// Έλεγχοι ύπαρξης
	existingUser, _ := s.users.FindByUsername(ctx, clean.Username)
	if existingUser.Username != "" {
		return ValidationError{Errors: []string{"Username already exists"}}
	}
	existingMail, _ := s.users.FindByEmail(ctx, clean.Mail)
	if existingMail.Mail != "" {
		return ValidationError{Errors: []string{"Email already exists"}}
	}

	uid, _ := uuid.NewV4()
	hash, err := bcrypt.GenerateFromPassword([]byte(clean.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u := models.User{
		UUID:         uid.String(),
		Mail:         clean.Mail,
		Username:     clean.Username,
		Password:     string(hash),
		Role:         clean.Role,
		CreationDate: time.Now(),
		Verified:     false,
	}

	// Αποθήκευση χρήστη
	if err := s.users.Create(ctx, u); err != nil {
		return err
	}

	// Στείλε email επιβεβαίωσης (ασύγχρονα)
	go func() {
		link := fmt.Sprintf("http://localhost:8080/verify?uuid=%s", uid.String())
		if err := utils.SendVerificationEmail(u.Mail, link); err != nil {
			log.Println("Error sending verification email:", err)
		}
	}()

	// Spawn a goroutine that will delete the user if still unverified after 5 minutes.
	go func(userUUID string) {
		// wait threshold
		time.Sleep(5 * time.Minute)

		// check current state
		checkCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		existing, err := s.users.FindByUUID(checkCtx, userUUID)
		if err != nil {
			// no-op: log if you want; repository may return error when not found
			return
		}
		if existing.UUID != "" && !existing.Verified {
			// Attempt deletion (add DeleteByUUID to your UserRepository if missing)
			_ = s.users.DeleteByUUID(checkCtx, userUUID)
		}
	}(uid.String())

	return nil
}

func (s *UserService) ValidateCredentials(ctx context.Context, username, password string) (models.User, error) {
	u, err := s.users.FindByUsername(ctx, username)
	if err != nil {
		return models.User{}, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return models.User{}, errors.New("invalid credentials")
	}
	return u, nil
}

func (s *UserService) VerifyEmail(ctx context.Context, uuid string) error {
	// 1. Βρες τον χρήστη
	user, err := s.users.FindByUUID(ctx, uuid)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// 2. Έλεγξε αν είναι ήδη verified
	if user.Verified {
		return fmt.Errorf("email already verified")
	}

	// 3. Ενημέρωσε τη βάση ότι είναι verified
	err = s.users.UpdateVerification(ctx, uuid, true)
	if err != nil {
		return fmt.Errorf("failed to update verification: %w", err)
	}

	return nil
}

func (s *UserService) ResendVerificationEmail(ctx context.Context, email string) error {
	// Έλεγχος αν υπάρχει χρήστης
	user, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to fetch user: %w", err)
	}
	if user.Mail == "" {
		return fmt.Errorf("user not found")
	}
	if user.Verified {
		return fmt.Errorf("user already verified")
	}

	// Δημιουργία verification link (όπως στο Register)
	link := fmt.Sprintf("http://localhost:8080/verify?uuid=%s", user.UUID)

	// Αποστολή email (προαιρετικά ασύγχρονα)
	go func() {
		if err := utils.SendVerificationEmail(user.Mail, link); err != nil {
			log.Printf("Error resending verification email to %s: %v", user.Mail, err)
		}
	}()

	return nil
}

func (s *UserService) GetByUsername(ctx context.Context, username string) (models.User, error) {
	return s.users.FindByUsername(ctx, username)
}
