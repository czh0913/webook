package memory

import (
	"context"
	"fmt"
)

type Service struct {
}

func NewSMSService() *Service {
	return &Service{}
}

func (s Service) Send(ctx context.Context, tpl string, args []string, number ...string) error {
	fmt.Println(args)
	return nil
}
