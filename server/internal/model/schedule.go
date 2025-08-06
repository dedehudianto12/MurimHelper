package model

type Schedule struct {
	ID				string `json:"id"`
	StartTime		string `json:"startTime"`
	EndTime 		string `json:"endTime"`
	Task 			string `json:"task"`
	Description		string `json:"description"`
}