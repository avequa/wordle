package cli

// stats - статистика игрока за один запуск программы (только завершенные партии)
type stats struct {
	played     int
	wins       int
	streak     int // текущая серия побед подряд
	bestStreak int // лучшая серия побед за запуск
}

// add - учитывает результат завершенной партии
func (s *stats) add(win bool) {
	s.played++

	if !win {
		s.streak = 0

		return
	}

	s.wins++
	s.streak++
	s.bestStreak = max(s.bestStreak, s.streak)
}
