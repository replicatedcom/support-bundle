package commands

import (
	"github.com/replicatedcom/support-bundle/pkg/cli"
	"github.com/spf13/cobra"
)

type uploadOptions struct {
	uploadBundlePath  string
	firstName         string
	lastName          string
	email             string
	company           string
	bundleDescription string
}

func NewUploadCommand(cli *cli.Cli) *cobra.Command {
	opts := uploadOptions{}

	cmd := &cobra.Command{
		Use:   "upload",
		Short: "Upload a support bundle to share",
		Long: `Upload an existing support bundle. This will be secret, and you'll receive a name
when uploaded that can be shared with support staff to access this bundle.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cli.Upload(opts.uploadBundlePath, opts.firstName, opts.lastName, opts.email, opts.company, opts.bundleDescription)
		},
	}

	cmd.Flags().StringVarP(&opts.uploadBundlePath, "path", "p", "supportbundle.tar.gz", "Path to the bundle that should be uploaded")
	cmd.Flags().StringVar(&opts.firstName, "firstname", "", "Your first name")
	cmd.Flags().StringVar(&opts.lastName, "lastname", "", "Your last name")
	cmd.Flags().StringVar(&opts.email, "email", "", "Your email")
	cmd.Flags().StringVar(&opts.company, "company", "", "The name of your company")
	cmd.Flags().StringVar(&opts.bundleDescription, "description", "No description provided", "A description for the issue being experienced")

	return cmd
}
