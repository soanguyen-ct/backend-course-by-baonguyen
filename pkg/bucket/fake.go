package bucket

import (
	"context"
	"fmt"
	"io"
)

func NewFake() *fake {
	return &fake{}
}

type fake struct {
}

func (f *fake) SaveImage(ctx context.Context, name string, r io.Reader) (string, error) {
	return fmt.Sprintf("/static/images/%s", name), nil
}
