package cli

import jww "github.com/spf13/jwalterweatherman"

func (cli *Cli) Upload(uploadBundlePath, firstName, lastName, email, company, bundleDescription string) error {
	jww.SetStdoutThreshold(jww.LevelTrace)

	// jww.FEEDBACK.Println("Uploading the provided support bundle")

	// contents, err := os.Open(uploadBundlePath)
	// if err != nil {
	// 	jww.ERROR.Printf("Error encountered when trying to read support bundle: %s", err)
	// 	return err
	// }
	// defer contents.Close()

	// bundleName, err := bundle.Upload(contents, firstName, lastName, email, company, bundleDescription)
	// if err != nil {
	// 	jww.ERROR.Printf("Error encountered when uploading support bundle: %s", err)
	// 	return err
	// }

	// jww.FEEDBACK.Printf("Support bundle located at %s was uploaded. This bundle can be referred to as %s.", uploadBundlePath, bundleName)

	jww.FEEDBACK.Println("Uploading support bundles from the command line is not yet implemented.")

	return nil
}
