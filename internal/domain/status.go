// Package domain реализует доменные файлы игры: слово, словарь, сессия, статус, результат, а также покрытие тестами
package domain

// Status - состояние игры
type Status int

// Возможные состояния игры
const (
	StatusUnknown    Status = 0
	StatusInProgress Status = 1
	StatusWin        Status = 2
	StatusLose       Status = 3
)

// IsFinished сообщает, завершена ли игры (либо WIN, либо LOSE)
func (s Status) IsFinished() bool {
	return s == StatusWin || s == StatusLose
}
