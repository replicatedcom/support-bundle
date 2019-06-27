package types

var GenerateOnlyLifecycleTasks = []LifecycleTask{
	{
		Message: &MessageOptions{"Starting support bundle collection..."},
	},
	{
		Generate: &GenerateOptions{UseDefaults: true},
	},
	{
		Message: &MessageOptions{"Generation complete!"},
	},
}

var DefaultLifecycleTasks = []LifecycleTask{
	{
		Message: &MessageOptions{"Starting support bundle collection..."},
	},
	{
		Notes: &NotesOptions{
			Prompt: "Enter a note: ",
		},
	},
	{
		Generate: &GenerateOptions{UseDefaults: true},
	},
	{
		Upload: &UploadOptions{},
	},
	{
		Message: &MessageOptions{"Upload complete! Check the analyzed bundle for more information"},
	},
}

var DefaultWatchLifecycleTasks = []LifecycleTask{
	{
		Message: &MessageOptions{"Starting support bundle collection..."},
	},
	{
		Generate: &GenerateOptions{UseDefaults: true},
	},
	{
		Upload: &UploadOptions{},
	},
	{
		Message: &MessageOptions{"Upload complete! Visit the Troubleshoot page in Ship to view the results"},
	},
}

func GetUseDefaults(tasks []LifecycleTask) bool {
	for _, task := range tasks {
		if task.Generate != nil && task.Generate.UseDefaults {
			return true
		}
	}
	return false
}
