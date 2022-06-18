package tasks

type Task struct {
	Name   string
	Action func()
	Cron   string
}

var RegisteredTasks []*Task

func registerTask(t *Task) {
	RegisteredTasks = append(RegisteredTasks, t)
}
