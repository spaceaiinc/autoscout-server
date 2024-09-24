package responses

import "github.com/spaceaiinc/autoscout-server/domain/entity"

type TaskGroup struct {
	TaskGroup *entity.TaskGroup `json:"task_group"`
}

func NewTaskGroup(taskGroup *entity.TaskGroup) TaskGroup {
	return TaskGroup{
		TaskGroup: taskGroup,
	}
}

type TaskGroupList struct {
	TaskGroupList []*entity.TaskGroup `json:"task_group_list"`
}

func NewTaskGroupList(taskGroups []*entity.TaskGroup) TaskGroupList {
	return TaskGroupList{
		TaskGroupList: taskGroups,
	}
}

type TaskGroupListAndMaxPage struct {
	MaxPageNumber uint                `json:"max_page_number"`
	TaskGroupList []*entity.TaskGroup `json:"task_group_list"`
}

func NewTaskGroupListAndMaxPage(enterpriseProfiles []*entity.TaskGroup, maxPageNumber uint) TaskGroupListAndMaxPage {
	return TaskGroupListAndMaxPage{
		MaxPageNumber: maxPageNumber,
		TaskGroupList: enterpriseProfiles,
	}
}

type SoundOutGroupList struct {
	SoundOutGroupList       []*entity.SoundOutGroup `json:"sound_out_group_list"`
	AlreadyCreatedGroupList []*entity.SoundOutGroup `json:"already_created_group_list"`
}

func NewSoundOutGroupList(soundOutGroupList, alreadyCreatedGroupList []*entity.SoundOutGroup) SoundOutGroupList {
	return SoundOutGroupList{
		SoundOutGroupList:       soundOutGroupList,
		AlreadyCreatedGroupList: alreadyCreatedGroupList,
	}
}
