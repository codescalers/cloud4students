/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/codescalers/cloud4students/cmd"
	"github.com/codescalers/cloud4students/validators"
	"github.com/rs/zerolog/log"
	"gopkg.in/validator.v2"
)

func init() {
	// validations
	err := validator.SetValidationFunc("ssh", validators.ValidateSSHKey)
	if err != nil {
		log.Error().Err(err).Send()
	}

	err = validator.SetValidationFunc("password", validators.ValidatePassword)
	if err != nil {
		log.Error().Err(err).Send()
	}

	err = validator.SetValidationFunc("mail", validators.ValidateMail)
	if err != nil {
		log.Error().Err(err).Send()
	}
}

func main() {
	cmd.Execute()
}
