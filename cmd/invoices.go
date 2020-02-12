package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xescugc/invoicer/cmd/model"
)

var (
	invoicesCmd = &cobra.Command{
		Use: "invoices",
	}

	invoicesNewCmd = &cobra.Command{
		Use: "new",
		RunE: func(cmd *cobra.Command, args []string) error {
			billing, err := initializeFilesystemBilling()
			if err != nil {
				return err
			}

			b, err := captureInputFromEditor(model.NewInvoice())
			if err != nil {
				return err
			}

			var i model.Invoice
			err = json.Unmarshal(b, &i)
			if err != nil {
				return err
			}

			in, err := i.ToDomain()
			if err != nil {
				return err
			}

			ctx := context.Background()

			err = billing.CreateInvoice(ctx, in, viper.GetString("customer"))
			if err != nil {
				return err
			}

			return nil
		},
	}

	invoicesGetCmd = &cobra.Command{
		Use: "get",
		RunE: func(cmd *cobra.Command, args []string) error {
			billing, err := initializeFilesystemBilling()
			if err != nil {
				return err
			}

			number := args[0]

			ctx := context.Background()

			i, err := billing.GetInvoice(ctx, number)
			if err != nil {
				return err
			}

			b, err := json.MarshalIndent(i, "", " ")
			if err != nil {
				return err
			}

			fmt.Println(string(b))

			return nil
		},
	}

	invoicesListCmd = &cobra.Command{
		Use: "list",
		RunE: func(cmd *cobra.Command, args []string) error {
			billing, err := initializeFilesystemBilling()
			if err != nil {
				return err
			}

			ctx := context.Background()
			ins, err := billing.GetInvoices(ctx)
			if err != nil {
				return err
			}

			for _, in := range ins {
				fmt.Printf(
					"Total: %s \tNumber: %q\t Date: %q: Customer: %q\n",
					in.Total(), in.Number,
					in.Date.Format(model.DefaultDateFormat), in.Customer.Canonical,
				)
			}

			return nil
		},
	}

	invoicesEditCmd = &cobra.Command{
		Use: "edit",
		RunE: func(cmd *cobra.Command, args []string) error {
			billing, err := initializeFilesystemBilling()
			if err != nil {
				return err
			}

			numb := args[0]

			ctx := context.Background()
			i, err := billing.GetInvoice(ctx, numb)
			if err != nil {
				return err
			}

			b, err := captureInputFromEditor(model.NewInvoiceFromDomain(i))
			if err != nil {
				return err
			}

			var ni model.Invoice
			err = json.Unmarshal(b, &ni)
			if err != nil {
				return err
			}

			i, err = ni.ToDomain()
			if err != nil {
				return err
			}

			err = billing.UpdateInvoice(ctx, numb, i, viper.GetString("customer"))
			if err != nil {
				return err
			}

			return nil
		},
	}

	invoicesDeleteCmd = &cobra.Command{
		Use: "delete",
		RunE: func(cmd *cobra.Command, args []string) error {
			billing, err := initializeFilesystemBilling()
			if err != nil {
				return err
			}

			number := args[0]

			ctx := context.Background()

			err = billing.DeleteInvoice(ctx, number)
			if err != nil {
				return err
			}

			return nil
		},
	}
)

func init() {
	invoicesCmd.AddCommand(
		invoicesNewCmd,
		invoicesGetCmd,
		invoicesListCmd,
		invoicesEditCmd,
		invoicesDeleteCmd,
	)

	invoicesNewCmd.Flags().StringP("customer", "c", "", "The Customer canonical, required")
	invoicesNewCmd.Flags().SetAnnotation("customer", cobra.BashCompOneRequiredFlag, []string{"true"})
	viper.BindPFlag("customer", invoicesNewCmd.Flags().Lookup("customer"))

	invoicesEditCmd.Flags().StringP("customer", "c", "", "The Customer canonical")
	viper.BindPFlag("customer", invoicesEditCmd.Flags().Lookup("customer"))
}
