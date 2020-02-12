package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xescugc/invoicer/customer"
)

var (
	customersCmd = &cobra.Command{
		Use: "customers",
	}

	customersNewCmd = &cobra.Command{
		Use: "new",
		RunE: func(cmd *cobra.Command, args []string) error {
			billing, err := initializeFilesystemBilling()
			if err != nil {
				return err
			}

			b, err := captureInputFromEditor(customer.Customer{})
			if err != nil {
				return err
			}

			var c customer.Customer
			err = json.Unmarshal(b, &c)
			if err != nil {
				return err
			}

			ctx := context.Background()
			err = billing.CreateCustomer(ctx, &c)
			if err != nil {
				return err
			}

			return nil
		},
	}

	customersGetCmd = &cobra.Command{
		Use:  "get",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			billing, err := initializeFilesystemBilling()
			if err != nil {
				return err
			}

			can := args[0]

			ctx := context.Background()
			c, err := billing.GetCustomer(ctx, can)
			if err != nil {
				return err
			}

			b, err := json.MarshalIndent(c, "", " ")
			if err != nil {
				return err
			}

			fmt.Println(string(b))

			return nil
		},
	}

	customersListCmd = &cobra.Command{
		Use: "list",
		RunE: func(cmd *cobra.Command, args []string) error {
			billing, err := initializeFilesystemBilling()
			if err != nil {
				return err
			}

			ctx := context.Background()
			cs, err := billing.GetCustomers(ctx)
			if err != nil {
				return err
			}

			for _, c := range cs {
				fmt.Printf("Name: %s \tCanonical: %s\n", c.Name, c.Canonical)
			}

			return nil
		},
	}

	customersEditCmd = &cobra.Command{
		Use:  "edit",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			billing, err := initializeFilesystemBilling()
			if err != nil {
				return err
			}

			can := args[0]

			ctx := context.Background()
			c, err := billing.GetCustomer(ctx, can)
			if err != nil {
				return err
			}

			b, err := captureInputFromEditor(c)
			if err != nil {
				return err
			}

			var nc customer.Customer
			err = json.Unmarshal(b, &nc)
			if err != nil {
				return err
			}

			err = billing.UpdateCustomer(ctx, can, &nc)
			if err != nil {
				return err
			}

			return nil
		},
	}

	customersDeleteCmd = &cobra.Command{
		Use:  "delete",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			billing, err := initializeFilesystemBilling()
			if err != nil {
				return err
			}

			can := args[0]

			ctx := context.Background()
			err = billing.DeleteCustomer(ctx, can)
			if err != nil {
				return err
			}

			return nil
		},
	}
)

func init() {
	customersCmd.AddCommand(
		customersNewCmd,
		customersGetCmd,
		customersListCmd,
		customersEditCmd,
		customersDeleteCmd,
	)
}
