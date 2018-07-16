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
		Generate: &GenerateOptions{UseDefaults: true},
	},
	{
		Upload: &UploadOptions{},
	},
	{
		Message: &MessageOptions{"Upload complete! Check the analyzed bundle for more information"},
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
