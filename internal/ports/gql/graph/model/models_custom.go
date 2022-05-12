package model

import "FICSIT-Ordis/internal/domain"

type Command domain.Command

func (c Command) Domain() *domain.Command {
	cmd := domain.Command(c)
	return &cmd
}

func ModelCommand(c domain.Command) *Command {
	cmd := Command(c)
	return &cmd
}

type Response domain.Response
type ResponseInput domain.Response
