package story

import "context"

type Service interface {
	Create(ctx context.Context, createRequest CreateStoryRequest) (err error)
}
