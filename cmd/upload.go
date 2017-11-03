// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload a support bundle to share",
	Long: `Upload an existing support bundle. This will be secret, and you'll receive a name
when uploaded that can be shared with support staff to access this bundle.`,
	RunE: upload,
}

var uploadBundlePath string
var firstName string
var lastName string
var email string
var company string
var bundleDescription string

func init() {
	RootCmd.AddCommand(uploadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uploadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// uploadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	uploadCmd.Flags().StringVarP(&uploadBundlePath, "path", "p", "supportbundle.tar.gz", "Path to the bundle that should be uploaded")
	uploadCmd.Flags().StringVar(&firstName, "firstname", "", "Your first name")
	uploadCmd.Flags().StringVar(&lastName, "lastname", "", "Your last name")
	uploadCmd.Flags().StringVar(&email, "email", "", "Your email")
	uploadCmd.Flags().StringVar(&company, "company", "", "The name of your company")
	uploadCmd.Flags().StringVar(&bundleDescription, "description", "No description provided", "A description for the issue being experienced")
}

func upload(cmd *cobra.Command, args []string) error {
	jww.SetStdoutThreshold(jww.LevelTrace)

	// jww.FEEDBACK.Println("Uploading the provided support bundle")

	// contents, err := os.Open(uploadBundlePath)
	// if err != nil {
	// 	jww.ERROR.Printf("Error encountered when trying to read support bundle: %s\n", err)
	// 	return err
	// }
	// defer contents.Close()

	// bundleName, err := bundle.Upload(contents, firstName, lastName, email, company, bundleDescription)
	// if err != nil {
	// 	jww.ERROR.Printf("Error encountered when uploading support bundle: %s\n", err)
	// 	return err
	// }

	// jww.FEEDBACK.Printf("Support bundle located at %s was uploaded. This bundle can be referred to as %s.\n", uploadBundlePath, bundleName)

	jww.FEEDBACK.Println("Uploading support bundles from the command line is not yet implemented.")

	return nil
}
