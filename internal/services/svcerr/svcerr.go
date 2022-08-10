package svcerr

// type ServiceError interface {
// 	Error() string
// 	ErrorValue() error
// 	Category() string
// }

func New(err error, category ServiceErrorCategory) *ServiceError {
	return &ServiceError{
		err,
		category,
	}
}

type ServiceError struct {
	error
	category ServiceErrorCategory
}

func (se *ServiceError) ErrorValue() error {
	return se.error
}

func (se *ServiceError) Category() ServiceErrorCategory {
	return se.category
}

type ServiceErrorCategory string
