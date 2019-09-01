package mutation

import (
	"context"
	"fmt"
	"net/mail"
	"strconv"

	"github.com/rafaelsq/boiler/pkg/graphql/internal/entity"
	"github.com/rafaelsq/boiler/pkg/iface"
	"github.com/rafaelsq/boiler/pkg/log"
	"github.com/rafaelsq/errors"
	"github.com/vektah/gqlparser/gqlerror"
)

func NewMutation(service iface.Service) *Mutation {
	return &Mutation{
		service: service,
	}
}

type Mutation struct {
	service iface.Service
}

func (m *Mutation) AddUser(ctx context.Context, input entity.AddUserInput) (*entity.UserResponse, error) {
	userID, err := m.service.AddUser(ctx, input.Name)
	if err != nil {
		log.Log(errors.New("fail to add user").SetParent(err))
		return nil, fmt.Errorf("service failed")
	}

	return &entity.UserResponse{User: &entity.User{ID: strconv.Itoa(userID)}}, nil
}

func (m *Mutation) AddEmail(ctx context.Context, input entity.AddEmailInput) (*entity.EmailResponse, error) {
	userID, err := strconv.Atoi(input.UserID)
	if err != nil || userID == 0 {
		return nil, &gqlerror.Error{
			Message: "invalid userID",
			Extensions: map[string]interface{}{
				"code": "-1",
			},
		}
	}

	address, err := mail.ParseAddress(input.Address)
	if err != nil {
		return nil, &gqlerror.Error{
			Message: "invalid email address",
			Extensions: map[string]interface{}{
				"code": "-2",
			},
		}
	}

	emailID, err := m.service.AddEmail(ctx, userID, address.Address)
	if err != nil {
		if er := errors.Cause(err); er == iface.ErrAlreadyExists {
			return nil, &gqlerror.Error{
				Message: er.Error(),
				Extensions: map[string]interface{}{
					"code": er.(*errors.Error).Args["code"].(string),
				},
			}
		}

		log.Log(errors.New("fail to add email").SetParent(err))
		return nil, fmt.Errorf("service failed")
	}

	return &entity.EmailResponse{Email: &entity.Email{ID: strconv.Itoa(emailID)}}, nil
}
