package types

var DefaultLifecycle = []*LifecycleTask{
	{
		Message: &MessageOptions{"Starting support bundle collection..."},
	},
	{
		Generate: &GenerateOptions{},
	},
	{
		BooleanPrompt: &BooleanPromptOptions{
			Contents: "Done! Do you want to upload the support bundle for analysis?",
			Default:  true,
		},
	},
	{
		Upload: &UploadOptions{},
	},
}
