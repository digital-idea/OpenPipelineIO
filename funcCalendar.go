package main

import "fmt"

func ItemsToFCEvents(items []Item) []FullCalendarEvent {
	results := []FullCalendarEvent{}
	for _, item := range items {
		// PM이 설정한 마감일로 Event를 만들어서 날짜 이벤트를 만들어야 한다.
		if item.Ddline2d != "" {
			deadline2dEvent := FullCalendarEvent{}
			deadline2dEvent.Title = item.Name
			deadline2dEvent.Start = item.Ddline2d
			deadline2dEvent.AllDay = true
			deadline2dEvent.Editable = true
			deadline2dEvent.StartEditable = true
			deadline2dEvent.EndEditable = true
			deadline2dEvent.DurationEditable = true
			deadline2dEvent.ResourceEditable = true
			deadline2dEvent.ExtendedProps.ItemID = item.ID
			deadline2dEvent.ExtendedProps.Project = item.Project
			deadline2dEvent.ExtendedProps.Tags = item.Tag
			results = append(results, deadline2dEvent)
		}

		// Task별로 날짜 Event를 만들어야 한다.
		for task, value := range item.Tasks {
			if value.Date == "" && value.Startdate == "" {
				continue
			}
			taskEvent := FullCalendarEvent{}
			taskEvent.Start = value.Startdate // 작업시작일
			taskEvent.Title = fmt.Sprintf("%s-%s", item.Name, task)
			taskEvent.End = value.Date // 최종마감일
			taskEvent.ExtendedProps.ItemID = item.ID
			taskEvent.ExtendedProps.Project = item.Project
			taskEvent.ExtendedProps.Task = task
			taskEvent.ExtendedProps.Tags = item.Tag
			taskEvent.ExtendedProps.Pipelinestep = value.Pipelinestep
			taskEvent.ExtendedProps.UserID = value.UserID
			taskEvent.ExtendedProps.TaskDeadline = value.Date       // 최종마감일
			taskEvent.ExtendedProps.TaskPreDeadline = value.Predate // 1차마감일
			taskEvent.ExtendedProps.TaskStartDate = value.Startdate // 작업시작일
			results = append(results, taskEvent)
		}
	}
	return results
}
