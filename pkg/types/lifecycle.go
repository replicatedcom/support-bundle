package types

var GenerateOnlyLifecycleTasks = []*LifecycleTask{
	{
		Message: &MessageOptions{"Starting support bundle collection..."},
	},
	{
		Generate: &GenerateOptions{},
	},
	{
		Message: &MessageOptions{"Generation complete!"},
	},
}

var DefaultLifecycleTasks = []*LifecycleTask{
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
	{
		Message: &MessageOptions{"Upload complete! Check the analyzed bundle for more information"},
	},
}
