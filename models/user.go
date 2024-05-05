package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRole string

const (
	UserRoleUser       UserRole = "user"
	UserRoleAdmin      UserRole = "admin"
	UserRoleSuperAdmin UserRole = "super_admin"
)

// SignUpInput 👈 SignUpInput struct
type SignUpInput struct {
	Username        string `json:"username" bson:"username" binding:"required"`
	Email           string `json:"email" bson:"email" binding:"required"`
	Password        string `json:"password" bson:"password" binding:"required,min=8"`
	PasswordConfirm string `json:"passwordConfirm" bson:"passwordConfirm,omitempty" binding:"required"`
	// RawPassword        string           `json:"raw_password"`
	Role               UserRole         `json:"role" bson:"role"`
	VerificationCode   string           `json:"verificationCode,omitempty" bson:"verificationCode,omitempty"`
	ResetPasswordToken string           `json:"resetPasswordToken,omitempty" bson:"resetPasswordToken,omitempty"`
	ResetPasswordAt    time.Time        `json:"resetPasswordAt,omitempty" bson:"resetPasswordAt,omitempty"`
	Verified           bool             `json:"verified" bson:"verified"`
	IsFirstLogin       bool             `json:"is_first_login" bson:"is_first_login"`

	CreatedAt          time.Time        `json:"created_at" bson:"created_at"`
	UpdatedAt          time.Time        `json:"updated_at" bson:"updated_at"`
}

// SignInInput 👈 SignInInput struct
type SignInInput struct {
	Email    string `json:"email" bson:"email" binding:"required"`
	Password string `json:"password" bson:"password" binding:"required"`
}

// DBResponse 👈 DBResponse struct
type DBResponse struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id"`
	Username           string             `json:"username" bson:"username"`
	Email              string             `json:"email" bson:"email"`
	Password           string             `json:"password" bson:"password"`
	Role               UserRole           `json:"role" bson:"role"`
	VerificationCode   string             `json:"verificationCode,omitempty" bson:"verificationCode"`
	ResetPasswordToken string             `json:"resetPasswordToken,omitempty" bson:"resetPasswordToken,omitempty"`
	ResetPasswordAt    time.Time          `json:"resetPasswordAt,omitempty" bson:"resetPasswordAt,omitempty"`
	Verified           bool               `json:"verified" bson:"verified"`
	IsFirstLogin       bool               `json:"is_first_login" bson:"is_first_login"`
	Followers          int                `json:"followers" bson:"followers"`
	Following          int                `json:"following" bson:"following"`
	Images             *[]ImageResponse `json:"images" bson:"images"`
	Library            *[]DBLibrary     `json:"libraries" bson:"libraries"`
	CreatedAt          time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt          time.Time          `json:"updated_at" bson:"updated_at"`
}

type UpdateInput struct {
	Username            string    `json:"username,omitempty" bson:"username,omitempty"`
	Email           string    `json:"email,omitempty" bson:"email,omitempty"`
	Password        string    `json:"password,omitempty" bson:"password,omitempty"`
	Role            UserRole  `json:"role,omitempty" bson:"role,omitempty"`
	ProfileImage    string    `json:"profile_image,omitempty" bson:"profile_image,omitempty"`
	UpdatedAt       time.Time `json:"updated_at" bson:"updated_at"`
}

// 👈 UserResponse struct
type UserResponse struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Username  string             `json:"username,omitempty" bson:"username,omitempty"`
	Email     string             `json:"email,omitempty" bson:"email,omitempty"`
	Role      string             `json:"role,omitempty" bson:"role,omitempty"`
	Followers          int              `json:"Followers" bson:"followers"`
	Following          int              `json:"following" bson:"following"`
	Images             *[]ImageResponse `json:"images" bson:"images"`
	Library            *[]DBLibrary     `json:"libraries" bson:"libraries"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// ForgotPasswordInput 👈 ForgotPasswordInput struct
type ForgotPasswordInput struct {
	Email string `json:"email" bson:"email" binding:"required"`
}

// ResetPasswordInput 👈 ResetPasswordInput struct
type ResetPasswordInput struct {
	Password        string `json:"password" bson:"password"`
	PasswordConfirm string `json:"passwordConfirm,omitempty" bson:"passwordConfirm,omitempty"`
}

func FilteredResponse(user *DBResponse) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Followers: user.Followers,
		Following: user.Following,
		Images:    user.Images,
		Library:   user.Library,
	}
}
