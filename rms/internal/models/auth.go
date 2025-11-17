package models

import "github.com/golang-jwt/jwt/v5"

type AccessRequest struct {
    Email    string      `json:"email"`
    Password string      `json:"password"`
    Argument interface{} `json:"argument,omitempty"`
}

type LoginResponse struct {
    Token string `json:"token"`
}

type Claims struct {
    Email string `json:"email"`
    Role  Role   `json:"role"`
    jwt.RegisteredClaims
}

type PasswordChangeInput struct {
    OldPassword string `json:"old_password"`
    NewPassword string `json:"new_password"`
}