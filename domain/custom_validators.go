package domain

import (
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/dlclark/regexp2"

	"github.com/go-playground/validator/v10"
)

const (
	StrongPasswordTag = "strongpassword"
	ValidateImagesTag = "validateImages"
	General           = "general"
	MaxImageSize      = 5 * 1024 * 1024
)

var AllowedImagesExtensions = map[string]bool{
	".png":  true,
	".jpg":  true,
	".jpeg": true,
}

func SetupCustomValidations(validator *validator.Validate) error {
	if err := validator.RegisterValidation(StrongPasswordTag, strongPasswordValidator); err != nil {
		return err
	}

	if err := validator.RegisterValidation(ValidateImagesTag, imageFileValidator); err != nil {
		return err
	}

	return nil
}

func strongPasswordValidator(fl validator.FieldLevel) bool {
	pattern := `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[!@#$&*])[A-Za-z\d!@#$&*]{8,}$`

	re := regexp2.MustCompile(pattern, 0)

	match, _ := re.MatchString(fl.Field().String())
	return match
}

func imageFileValidator(fl validator.FieldLevel) bool {
	param := fl.Param()
	maxImagesAllowed, err := strconv.Atoi(param)
	if err != nil {
		return false
	}

	images := fl.Field().Interface().([]*multipart.FileHeader)

	if len(images) == 0 {
		return true
	}

	if len(images) > maxImagesAllowed {
		return false
	}

	for _, file := range images {
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if !AllowedImagesExtensions[ext] {
			return false
		}

		if file.Size > MaxImageSize {
			return false
		}
	}

	return true
}
