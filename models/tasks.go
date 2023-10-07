package models

import "sync"

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

var (
	Tasks   = []*Task{}
	TasksMu sync.RWMutex
	NextID  = 1
)
