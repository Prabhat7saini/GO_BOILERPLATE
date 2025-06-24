package utils

import (
	common "boiler-platecode/src/common/const"
	"errors"
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
)

func CreateUserValidationErrors(err error) (exception *common.Exception, fieldErrors map[string]string) {
	var ve validator.ValidationErrors
	fieldErrors = make(map[string]string)

	if errors.As(err, &ve) {
		for _, fe := range ve {
			field := fe.Field()
			tag := fe.Tag()

			switch field {
			case "PASSWORD":
				if tag == "strongpassword" {
					fieldErrors["password"] = "Password must be 8-12 characters and include uppercase, lowercase, digit, and special character."
				} else {
					fieldErrors["password"] = "Invalid password."
				}
			case "Email":
				if tag == "email" {
					fieldErrors["email"] = "Please enter a valid email address."
				}
			case "Name":
				if tag == "min" {
					fieldErrors["name"] = "Name must be at least 3 characters long."
				}
			default:
				fieldErrors[field] = fmt.Sprintf("Invalid %s", field)
			}
		}


			// Log field errors one by one
		log.Println("CreateUserValidationErrors:")
		for field, msg := range fieldErrors {
			log.Printf("  %s: %s", field, msg)
		}
	
		// Return a pointer
		return &common.Exception{
			Code:           "VAL001",
			Message:        "Validation Failed",
			HttpStatusCode: 400,
		}, nil
	}

	return &common.Exception{
		Code:           "VAL000",
		Message:        "Invalid request payload",
		HttpStatusCode: 400,
	}, nil
}

