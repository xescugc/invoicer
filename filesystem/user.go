package filesystem

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/xescugc/invoicer/billing"
	"github.com/xescugc/invoicer/user"
)

type UserRepository struct {
	baseDir string
}

func NewUserRepository(baseDir string) *UserRepository {
	return &UserRepository{
		baseDir: baseDir,
	}
}

func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	b, err := json.MarshalIndent(u, "", " ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath.Join(r.baseDir, "user.json"), b, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) Find(ctx context.Context) (*user.User, error) {
	b, err := ioutil.ReadFile(filepath.Join(r.baseDir, "user.json"))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, billing.ErrNotFoundUser
		}
		return nil, err
	}

	var u user.User
	err = json.Unmarshal(b, &u)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *UserRepository) Update(ctx context.Context, u *user.User) error {
	b, err := json.MarshalIndent(u, "", " ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath.Join(r.baseDir, "user.json"), b, 0644)
	if err != nil {
		return err
	}

	return nil
}
