package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/xescugc/invoicer/user"

	"github.com/spf13/cobra"
)

var (
	userCmd = &cobra.Command{
		Use: "user",
	}

	userNewCmd = &cobra.Command{
		Use: "new",
		RunE: func(cmd *cobra.Command, args []string) error {
			billing, err := initializeFilesystemBilling()
			if err != nil {
				return err
			}
			b, err := captureInputFromEditor(user.User{})
			if err != nil {
				return err
			}

			var u user.User
			err = json.Unmarshal(b, &u)
			if err != nil {
				return err
			}

			ctx := context.Background()
			err = billing.CreateUser(ctx, &u)
			if err != nil {
				return err
			}

			return nil
		},
	}

	userGetCmd = &cobra.Command{
		Use: "get",
		RunE: func(cmd *cobra.Command, args []string) error {
			billing, err := initializeFilesystemBilling()
			if err != nil {
				return err
			}

			ctx := context.Background()
			u, err := billing.GetUser(ctx)
			if err != nil {
				return err
			}

			b, err := json.MarshalIndent(u, "", " ")
			if err != nil {
				return err
			}

			fmt.Println(string(b))

			return nil
		},
	}

	userEditCmd = &cobra.Command{
		Use: "edit",
		RunE: func(cmd *cobra.Command, args []string) error {
			billing, err := initializeFilesystemBilling()
			if err != nil {
				return err
			}

			ctx := context.Background()
			oldUser, err := billing.GetUser(ctx)
			if err != nil {
				return err
			}

			b, err := captureInputFromEditor(oldUser)
			if err != nil {
				return err
			}

			var u user.User
			err = json.Unmarshal(b, &u)
			if err != nil {
				return err
			}

			err = billing.UpdateUser(ctx, &u)
			if err != nil {
				return err
			}

			return nil
		},
	}
)

func init() {
	userCmd.AddCommand(
		userNewCmd,
		userGetCmd,
		userEditCmd,
	)
}
