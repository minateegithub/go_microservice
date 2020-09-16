package controllers

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/minateegithub/go_microservice/models"
	"github.com/nyaruka/phonenumbers"

	"github.com/go-playground/validator/v10"
)

//ValidationError struct
type ValidationError struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

const phoneNumberRegexExp = "^+?([0-9]{2})?[-. ]?([0-9]{4})[-. ]?([0-9]{4})$"

var phoneNumberRegex = regexp.MustCompile(phoneNumberRegexExp)

var validate *validator.Validate

//InitValidator - Initialize the validators
func InitValidator() {
	validate = validator.New()
	validate.RegisterValidation("notblank", NotBlank)
}

//Descriptive - Returns errors in more descriptive way
func Descriptive(verr validator.ValidationErrors) []ValidationError {
	errs := []ValidationError{}

	for _, f := range verr {
		err := f.ActualTag()
		if f.Param() != "" {
			err = fmt.Sprintf("%s=%s", err, f.Param())
		}
		if strings.Contains(err, "datetime") {
			err = "invalid"
		}
		errs = append(errs, ValidationError{Field: f.Namespace(), Reason: err})
	}

	return errs
}

//NotBlank - Custom validator
func NotBlank(fl validator.FieldLevel) bool {
	return strings.TrimSpace(fl.Field().String()) != ""
}

//validatePhoneNumber - Validates phone number
func validatePhoneNumber(phoneNumber *string) ValidationError {

	if phoneNumber != nil && len(strings.TrimSpace(*phoneNumber)) > 0 {
		regionCode := os.Getenv("PHONE_NUMBER_REGION")
		num, err := phonenumbers.Parse(*phoneNumber, regionCode)

		if err != nil || !phonenumbers.IsPossibleNumber(num) {
			return ValidationError{Field: "Enrollee.PhoneNumber", Reason: "invalid"}
		}
	}
	return ValidationError{}
}

//validateNewEnroleeData - Validates new Enrollee data
func validateNewEnroleeData(enrollee models.Enrollee) []ValidationError {
	errs := []ValidationError{}

	validationErr := validate.Struct(enrollee)
	if validationErr != nil {
		var verr validator.ValidationErrors
		if errors.As(validationErr, &verr) {
			errs = Descriptive(verr)
		}
	}
	phoneNumberErr := validatePhoneNumber(enrollee.PhoneNumber)
	if phoneNumberErr != (ValidationError{}) {
		errs = append(errs, phoneNumberErr)
	}

	return errs
}

//validateEditedEnroleeData - Validates Enrollee data in case of edit
func validateEditedEnroleeData(enrollee models.Enrollee) []ValidationError {
	errs := []ValidationError{}

	validationErr := validate.Struct(enrollee)
	if validationErr != nil {
		var verr validator.ValidationErrors
		if errors.As(validationErr, &verr) {
			for _, f := range verr {

				if (f.Kind() != reflect.Ptr) || !isNil(f.Value()) {
					err := f.ActualTag()
					if f.Param() != "" {
						err = fmt.Sprintf("%s=%s", err, f.Param())
					}
					if strings.Contains(err, "datetime") {
						err = "invalid"
					}
					errs = append(errs, ValidationError{Field: f.Namespace(), Reason: err})
				}
			}
		}
	}
	phoneNumberErr := validatePhoneNumber(enrollee.PhoneNumber)
	if phoneNumberErr != (ValidationError{}) {
		errs = append(errs, phoneNumberErr)
	}

	return errs
}

//validateNewDependentData - Validates new dependent data
func validateNewDependentData(dependent models.Dependent) []ValidationError {
	errs := []ValidationError{}

	validationErr := validate.Struct(dependent)
	if validationErr != nil {
		var verr validator.ValidationErrors
		if errors.As(validationErr, &verr) {
			errs = Descriptive(verr)
		}
	}
	return errs
}

//validateEditedDependentData - Validates dependent data in case of edit
func validateEditedDependentData(dependent models.Dependent) []ValidationError {
	errs := []ValidationError{}

	validationErr := validate.Struct(dependent)
	if validationErr != nil {
		var verr validator.ValidationErrors
		if errors.As(validationErr, &verr) {
			for _, f := range verr {

				if (f.Kind() != reflect.Ptr) || !isNil(f.Value()) {
					err := f.ActualTag()
					if f.Param() != "" {
						err = fmt.Sprintf("%s=%s", err, f.Param())
					}
					if strings.Contains(err, "datetime") {
						err = "invalid"
					}
					errs = append(errs, ValidationError{Field: f.Namespace(), Reason: err})
				}
			}
		}
	}

	return errs
}

func isNil(i interface{}) bool {
	return i == nil || reflect.ValueOf(i).IsNil()
}
