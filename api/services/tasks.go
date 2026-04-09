package services

import (
	"api/models"
	"fmt"
)

func GenerateTask(task models.Task) {
	if task.Creator == "" {
		//We assume this was a task that's a follow up from a creator generated task... Find it
		var previousTask models.Task
		DB.Where("product_id =?", task.ProductID).First(&previousTask)
		// So we are retiring the previous task and starting a new oen, right?
		previousTask.Status = TASK_COMPLETE
		DB.Save(&previousTask)
		task.Creator = previousTask.Creator
		task.Assignee = previousTask.Creator
	}

	DB.Save(&task)
}

func DeleteTask(id int) {
	err := DB.Table("tasks").Where("id =?", id).Delete(&models.Task{}).Error
	if err != nil {
		fmt.Println(err)
	}
}

func GetActiveTasks(user models.AppUser) []models.Task {
	var tasks []models.Task
	DB.Where("status = ? and assignee = ?", TASK_ACTIVE, user.UserEmail).Find(&tasks)
	return tasks
}
