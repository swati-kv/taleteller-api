package store

import "context"

type StoryStorer interface {
	Create(ctx context.Context)
}
