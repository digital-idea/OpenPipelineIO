package main

// Tasksetting 자료구조이다
type Tasksetting struct {
	ID        string `json:"id"` // Task ID
	Name      string `json:"name"`
	Type      string `json:"type"` // Type: Asset, Shot
	LinuxDev  string `json:"linuxdev"`
	LinuxPub  string `json:"linuxpub"`
	WindowDev string `json:"windowdev"`
	WindowPub string `json:"windowpub"`
	MacOSDev  string `json:"macosdev"`
	MacOSPub  string `json:"macospub"`
}
