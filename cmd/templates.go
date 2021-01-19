package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
	"github.com/xescugc/invoicer/cmd/model"
	"github.com/xescugc/invoicer/template"
)

var (
	templatesCmd = &cobra.Command{
		Use: "templates",
	}

	templatesNewCmd = &cobra.Command{
		Use:  "new",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			billing, err := initializeFilesystemBilling()
			if err != nil {
				return err
			}

			b, err := captureInputFromEditor(model.Template{})
			if err != nil {
				return err
			}

			var t template.Template
			err = json.Unmarshal(b, &t)
			if err != nil {
				return err
			}

			body, err := ioutil.ReadFile(args[0])
			if err != nil {
				return err
			}

			t.Template = body

			ctx := context.Background()
			err = billing.CreateTemplate(ctx, &t)
			if err != nil {
				return err
			}

			return nil
		},
	}

	templatesGetCmd = &cobra.Command{
		Use:  "get",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			billing, err := initializeFilesystemBilling()
			if err != nil {
				return err
			}

			can := args[0]

			ctx := context.Background()
			t, err := billing.GetTemplate(ctx, can)
			if err != nil {
				return err
			}

			b, err := json.MarshalIndent(t, "", " ")
			if err != nil {
				return err
			}

			fmt.Println(string(b))

			return nil
		},
	}

	templatesViewCmd = &cobra.Command{
		Use:  "view",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			billing, err := initializeFilesystemBilling()
			if err != nil {
				return err
			}

			can := args[0]

			ctx := context.Background()
			t, err := billing.GetTemplate(ctx, can)
			if err != nil {
				return err
			}

			fmt.Println(string(t.Template))

			return nil
		},
	}

	templatesListCmd = &cobra.Command{
		Use: "list",
		RunE: func(cmd *cobra.Command, args []string) error {
			billing, err := initializeFilesystemBilling()
			if err != nil {
				return err
			}

			ctx := context.Background()
			ts, err := billing.GetTemplates(ctx)
			if err != nil {
				return err
			}

			for _, t := range ts {
				fmt.Printf("Name: %s \tCanonical: %s\n", t.Name, t.Canonical)
			}

			return nil
		},
	}

	templatesEditCmd = &cobra.Command{
		Use:  "edit",
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			billing, err := initializeFilesystemBilling()
			if err != nil {
				return err
			}

			can := args[0]

			ctx := context.Background()
			t, err := billing.GetTemplate(ctx, can)
			if err != nil {
				return err
			}

			b, err := captureInputFromEditor(model.NewTemplateFromDomain(t))
			if err != nil {
				return err
			}

			var nt template.Template
			err = json.Unmarshal(b, &nt)
			if err != nil {
				return err
			}

			var body []byte
			if len(args) == 2 {
				body, err = ioutil.ReadFile(args[1])
				if err != nil {
					return err
				}
			} else {
				body = t.Template
			}

			nt.Template = body

			err = billing.UpdateTemplate(ctx, can, &nt)
			if err != nil {
				return err
			}

			return nil
		},
	}

	templatesDeleteCmd = &cobra.Command{
		Use:  "delete",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			billing, err := initializeFilesystemBilling()
			if err != nil {
				return err
			}

			can := args[0]

			ctx := context.Background()
			err = billing.DeleteTemplate(ctx, can)
			if err != nil {
				return err
			}

			return nil
		},
	}
)

func init() {
	templatesCmd.AddCommand(
		templatesNewCmd,
		templatesGetCmd,
		templatesViewCmd,
		templatesListCmd,
		templatesEditCmd,
		templatesDeleteCmd,
	)
}
