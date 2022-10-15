package main

import (
	"fmt"
	"time"
)

func ItemsToFCEventsAndFCResource(items []Item) ([]FullCalendarEvent, []FullCalendarResource, error) {
	events := []FullCalendarEvent{}
	resources := []FullCalendarResource{}
	// 프로젝트에 사용될 기본 리소스를 추가한다.
	resource := FullCalendarResource{}
	resource.ID = "" // 빈값이면 상단으로 올라간다.
	resource.Title = "Project Information"
	resources = append(resources, resource)

	for _, item := range items {
		// 리소스를 추가한다. 간트챠트를 그리기 위해서 필요하다.
		resource := FullCalendarResource{}
		resource.ID = item.Name
		resource.Title = item.Name
		resources = append(resources, resource)

		// PM이 설정한 마감일로 Event를 만들어서 날짜 이벤트를 만들어야 한다.
		if item.Ddline2d != "" {
			deadline2dEvent := FullCalendarEvent{}
			deadline2dEvent.Title = "Deadline - " + item.Name
			deadline2dEvent.Color = "#BE2625"
			deadline2dEvent.Start = item.Ddline2d
			deadline2dEvent.End = item.Ddline2d
			deadline2dEvent.AllDay = true
			deadline2dEvent.Editable = true
			deadline2dEvent.StartEditable = true
			deadline2dEvent.EndEditable = true
			deadline2dEvent.DurationEditable = true
			deadline2dEvent.ResourceEditable = false
			deadline2dEvent.ExtendedProps.ItemName = item.Name
			deadline2dEvent.ExtendedProps.Project = item.Project
			deadline2dEvent.ExtendedProps.Tags = item.Tag
			deadline2dEvent.ExtendedProps.Key = "deadline2d"
			deadline2dEvent.ResourceId = item.Name
			events = append(events, deadline2dEvent)
		}

		// Task별로 날짜 Event를 만들어야 한다.
		for task, value := range item.Tasks {
			if value.Date == "" && value.Startdate == "" {
				continue
			}
			taskEvent := FullCalendarEvent{}
			taskEvent.AllDay = true
			taskEvent.Editable = true
			taskEvent.StartEditable = true
			taskEvent.EndEditable = true
			taskEvent.DurationEditable = true
			taskEvent.ResourceEditable = false
			taskEvent.Title = fmt.Sprintf("%s - %s", task, item.Name)
			if value.Startdate == "" {
				taskEvent.Start = value.Date // 작업마감일을 시작일로 설정한다.
			} else {
				taskEvent.Start = value.Startdate // 작업시작일
			}
			// fullcalendar 특성상 end 날짜에 1일을 더해야 간트챠트 드로잉시 그래프 모양이 딱 맞다.
			t, err := time.Parse(time.RFC3339, value.Date)
			if err != nil {
				return events, resources, err
			}
			t = t.Add(time.Hour * 24)
			taskEvent.End = t.Format(time.RFC3339)

			taskEvent.ExtendedProps.ItemName = item.Name
			taskEvent.ExtendedProps.Project = item.Project
			taskEvent.ExtendedProps.Task = task
			taskEvent.ExtendedProps.Tags = item.Tag
			taskEvent.ExtendedProps.Pipelinestep = value.Pipelinestep
			taskEvent.ExtendedProps.UserID = value.UserID
			taskEvent.ExtendedProps.TaskDeadline = value.Date       // 최종마감일
			taskEvent.ExtendedProps.TaskPreDeadline = value.Predate // 1차마감일
			taskEvent.ExtendedProps.TaskStartDate = value.Startdate // 작업시작일
			taskEvent.ExtendedProps.Key = "tasks"
			taskEvent.ResourceId = item.Name
			events = append(events, taskEvent)
		}
	}
	return events, resources, nil
}
