package generator

import (
	"os"
	"path/filepath"
)

func GenerateStubFiles(cwd string) {
	files := map[string]string{
		"internal/api/authutil/authutil.go": `package authutil

func Authenticate(username, password string) (bool, error) {
    // TODO: implement
    return false, nil
}`,
		"internal/api/v1/auth.go": `package v1

func AuthHandler() {
    // TODO: implement
}`,
		"internal/api/api.go": `package api

func InitAPI() {
    // TODO: implement
}`,
		"internal/api/middleware.go": `package api

func Middleware() {
    // TODO: implement
}`,
		"internal/api/system.go": `package api

func SystemHandler() {
    // TODO: implement
}`,
		"internal/common/bcrypt.go": `package common

func HashPassword(pw string) (string, error) {
    // TODO: implement
    return "", nil
}`,
		"internal/common/httphelpers.go": `package common

func JSONResponse() {
    // TODO: implement
}`,
		"internal/common/sqlhelpers.go": `package common

func SQLHelper() {
    // TODO: implement
}`,
		"internal/config/config.go": `package config

func Load() {
    // TODO: implement
}`,
		"internal/logging/loggedmodule/loggedmodule.go": `package loggedmodule

func LogModule() {
    // TODO: implement
}`,
		"internal/logging/logging.go": `package logging

func InitLogger() {
    // TODO: implement
}`,
		"internal/services/auth.go": `package services

func AuthService() {
    // TODO: implement
}`,
		"internal/validation/validation.go": `package validation

func Validate() {
    // TODO: implement
}`,
		"internal/version/version.go": `package version

const Version = "0.1.0"`,
		"pkg/problemdetail/problemdetail.go": `package problemdetail

func ProblemDetail() {
    // TODO: implement
}`,
	}

	for path, content := range files {
		full := filepath.Join(cwd, path)
		os.WriteFile(full, []byte(content), 0644)
	}
}
