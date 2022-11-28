package main

import (
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
		if item.Ddline2d != "" {
			// 리소스를 추가한다. 간트챠트를 그리기 위해서 필요하다.
			resource := FullCalendarResource{}
			resource.ID = "!Deadline2D-" + item.ID
			resource.Title = "Deadline"
			resource.ResourceGroup = item.Name
			resources = append(resources, resource)

			deadline2dEvent := FullCalendarEvent{}
			deadline2dEvent.Title = "Deadline - " + item.Name
			deadline2dEvent.Color = "#0A610C"
			deadline2dEvent.Start = item.Ddline2d
			deadline2dEvent.End = item.Ddline2d
			deadline2dEvent.AllDay = true
			deadline2dEvent.Editable = true
			deadline2dEvent.StartEditable = true
			deadline2dEvent.EndEditable = true
			deadline2dEvent.DurationEditable = true
			deadline2dEvent.ResourceEditable = false
			deadline2dEvent.ExtendedProps.ItemName = item.Name
			deadline2dEvent.ExtendedProps.ItemID = item.ID
			deadline2dEvent.ExtendedProps.Project = item.Project
			deadline2dEvent.ExtendedProps.Tags = item.Tag
			deadline2dEvent.ExtendedProps.Key = "deadline2d"
			deadline2dEvent.ResourceId = "!Deadline2D-" + item.ID

			events = append(events, deadline2dEvent)
		}
		// Task별로 날짜 Event를 만들어야 한다.
		for task, value := range item.Tasks {
			if value.Predate != "" || value.Date != "" {
				// 리소스를 추가한다. 간트챠트를 그리기 위해서 필요하다.
				resource := FullCalendarResource{}
				resource.ID = item.ID + task
				resource.Title = task
				resource.ResourceGroup = item.Name
				resources = append(resources, resource)
			}
			if value.Predate != "" {
				taskEvent := FullCalendarEvent{}
				taskEvent.AllDay = true
				taskEvent.Editable = true
				taskEvent.StartEditable = true
				taskEvent.EndEditable = true
				taskEvent.DurationEditable = true
				taskEvent.ResourceEditable = false
				taskEvent.ID = item.ID
				taskEvent.Color = "#93E57B"
				taskEvent.TextColor = "#000000"
				taskEvent.Title = task + " - " + item.Name
				taskEvent.Start = value.Startdate // 작업 시작일 설정
				if value.Startdate == "" {
					taskEvent.Start = value.Predate // 작업마감일을 시작일로 설정한다.
				}
				// fullcalendar 특성상 end 날짜에 1일을 더해야 간트챠트 드로잉시 그래프 모양이 딱 맞다.
				t, err := time.Parse(time.RFC3339, value.Predate)
				if err != nil {
					continue // 사용자가 날짜에 이상하게 문자를 넣거나 하더라도 에러가 나면 안된다. 그냥 넘긴다.
				}
				t = t.Add(time.Hour * 24)
				taskEvent.End = t.Format(time.RFC3339)
				taskEvent.ExtendedProps.Project = item.Project
				taskEvent.ExtendedProps.ItemID = item.ID
				taskEvent.ExtendedProps.ItemName = item.Name
				taskEvent.ExtendedProps.Task = task
				taskEvent.ExtendedProps.Tags = item.Tag
				taskEvent.ExtendedProps.Pipelinestep = value.Pipelinestep
				taskEvent.ExtendedProps.UserID = value.UserID
				taskEvent.ExtendedProps.DeadlineType = "1st"
				taskEvent.ExtendedProps.Key = "tasks"
				taskEvent.ResourceId = item.ID + task
				events = append(events, taskEvent)
			}
			if value.Date != "" {
				taskEvent := FullCalendarEvent{}
				taskEvent.AllDay = true
				taskEvent.Editable = true
				taskEvent.StartEditable = true
				taskEvent.EndEditable = true
				taskEvent.DurationEditable = true
				taskEvent.ResourceEditable = false
				taskEvent.ID = item.ID
				taskEvent.Color = "#29AF1D"
				taskEvent.TextColor = "#000000"
				taskEvent.Title = task + " - " + item.Name
				taskEvent.Start = value.Startdate2nd // 작업 시작일 설정
				if value.Startdate2nd == "" {
					taskEvent.Start = value.Date // 작업마감일을 시작일로 설정한다.
				}
				// fullcalendar 특성상 end 날짜에 1일을 더해야 간트챠트 드로잉시 그래프 모양이 딱 맞다.
				t, err := time.Parse(time.RFC3339, value.Date)
				if err != nil {
					continue // 사용자가 날짜에 이상하게 문자를 넣거나 하더라도 에러가 나면 안된다. 그냥 넘긴다.
				}
				t = t.Add(time.Hour * 24)
				taskEvent.End = t.Format(time.RFC3339)
				taskEvent.ExtendedProps.Project = item.Project
				taskEvent.ExtendedProps.ItemID = item.ID
				taskEvent.ExtendedProps.ItemName = item.Name
				taskEvent.ExtendedProps.Task = task
				taskEvent.ExtendedProps.Tags = item.Tag
				taskEvent.ExtendedProps.Pipelinestep = value.Pipelinestep
				taskEvent.ExtendedProps.UserID = value.UserID
				taskEvent.ExtendedProps.DeadlineType = "2nd"
				taskEvent.ExtendedProps.Key = "tasks"
				taskEvent.ResourceId = item.ID + task
				events = append(events, taskEvent)
			}
		}

	}
	return events, resources, nil
}
