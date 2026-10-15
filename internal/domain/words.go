// Package domain реализует доменные файлы игры: слово, словарь, сессия, статус, результат, а также покрытие тестами
package domain

import (
	"fiveletters/internal/apperr"
	"fmt"
	"slices"
)

// Уровни сложности: у каждого уровня свой набор загадываемых слов
const (
	LevelEasy   = "easy"
	LevelMedium = "medium"
	LevelHard   = "hard"
)

// Категории слов
const (
	CategoryAnimals = "animals"
	CategoryNature  = "nature"
	CategoryFood    = "food"
	CategoryHome    = "home"
	CategoryCity    = "city"
	CategoryStudy   = "study"
	CategorySport   = "sport"
)

// WordInfo - слово встроенного словаря с уровнем сложности, категорией и подсказкой
type WordInfo struct {
	Word     string
	Level    string
	Category string
	Hint     string
}

// Levels - все уровни сложности
func Levels() []string {
	return []string{LevelEasy, LevelMedium, LevelHard}
}

// Categories - все категории слов
func Categories() []string {
	return []string{
		CategoryAnimals,
		CategoryNature,
		CategoryFood,
		CategoryHome,
		CategoryCity,
		CategoryStudy,
		CategorySport,
	}
}

// Words - встроенный словарь: 85 слов с уровнем, категорией и подсказкой
func Words() []WordInfo { //nolint:funlen
	return []WordInfo{
		{"арбуз", LevelEasy, CategoryFood, "Большая полосатая ягода с красной сладкой мякотью"},
		{"баран", LevelEasy, CategoryAnimals, "Рогатое животное с густой шерстью, упрямо смотрит на новые ворота"},
		{"белка", LevelEasy, CategoryAnimals, "Рыжий зверёк с пушистым хвостом, запасает орехи на зиму"},
		{"берег", LevelMedium, CategoryNature, "Край суши у воды"},
		{"билет", LevelEasy, CategoryCity, "Его покупают, чтобы проехать в поезде или попасть в кино"},
		{"буква", LevelEasy, CategoryStudy, "Знак алфавита - в этой игре их пять"},
		{"вагон", LevelMedium, CategoryCity, "Часть поезда, в которой едут пассажиры"},
		{"весна", LevelEasy, CategoryNature, "Время года, когда тает снег и прилетают птицы"},
		{"ветер", LevelEasy, CategoryNature, "Движение воздуха, которое качает деревья"},
		{"вечер", LevelMedium, CategoryNature, "Время суток между днём и ночью"},
		{"вилка", LevelEasy, CategoryHome, "Столовый прибор с зубцами"},
		{"волна", LevelEasy, CategoryNature, "Бежит по морю и разбивается о берег"},
		{"ворон", LevelMedium, CategoryAnimals, "Крупная чёрная птица, в сказках живёт триста лет"},
		{"время", LevelMedium, CategoryStudy, "Его измеряют часами, и его всегда не хватает"},
		{"глава", LevelMedium, CategoryStudy, "Часть книги"},
		{"голос", LevelMedium, CategoryStudy, "Им поют и говорят"},
		{"город", LevelEasy, CategoryCity, "Большой населённый пункт с улицами и домами"},
		{"гость", LevelMedium, CategoryHome, "Тот, кого встречают у порога и угощают чаем"},
		{"груша", LevelEasy, CategoryFood, "Фрукт, по форме похожий на лампочку"},
		{"дверь", LevelEasy, CategoryHome, "Её открывают, чтобы войти в комнату"},
		{"дождь", LevelEasy, CategoryNature, "Капли воды, падающие из туч"},
		{"доска", LevelEasy, CategoryStudy, "На ней пишут мелом в классе"},
		{"дупло", LevelHard, CategoryNature, "Дыра в стволе дерева, домик для белки или дятла"},
		{"жираф", LevelEasy, CategoryAnimals, "Самое высокое животное на планете, с очень длинной шеей"},
		{"забор", LevelMedium, CategoryHome, "Ограда вокруг участка"},
		{"завод", LevelMedium, CategoryCity, "Большое предприятие, где делают машины или сталь"},
		{"закат", LevelMedium, CategoryNature, "Солнце опускается за горизонт"},
		{"замок", LevelHard, CategoryHome, "Бывает на двери, а бывает рыцарский - зависит от ударения"},
		{"зебра", LevelEasy, CategoryAnimals, "Полосатая африканская лошадка, её именем назван пешеходный переход"},
		{"земля", LevelEasy, CategoryNature, "Наша планета, а ещё почва под ногами"},
		{"зерно", LevelMedium, CategoryFood, "Семечко пшеницы или ржи, из него мелют муку"},
		{"игрок", LevelEasy, CategorySport, "Участник игры"},
		{"икона", LevelHard, CategoryCity, "Изображение святого в храме; так же называют значок программы"},
		{"кабан", LevelHard, CategoryAnimals, "Дикий родственник свиньи с клыками"},
		{"камин", LevelMedium, CategoryHome, "Очаг в стене комнаты, у которого греются зимой"},
		{"канат", LevelMedium, CategorySport, "Толстая верёвка, её перетягивают командами"},
		{"карта", LevelMedium, CategoryStudy, "Бывает географической, игральной и банковской"},
		{"кисть", LevelMedium, CategoryStudy, "Инструмент художника, а ещё часть руки"},
		{"класс", LevelEasy, CategoryStudy, "Комната для уроков и ребята, которые учатся вместе"},
		{"книга", LevelEasy, CategoryStudy, "Её читают, перелистывая страницы"},
		{"колос", LevelHard, CategoryNature, "Верхушка стебля пшеницы, в которой созревают зёрна"},
		{"комар", LevelMedium, CategoryAnimals, "Пищит над ухом летней ночью и кусается"},
		{"конус", LevelHard, CategoryStudy, "Геометрическое тело, похожее на рожок мороженого"},
		{"корка", LevelHard, CategoryFood, "Твёрдая верхняя часть хлеба или кожура апельсина"},
		{"кость", LevelMedium, CategoryAnimals, "Любимое лакомство собаки, из них состоит скелет"},
		{"крыло", LevelHard, CategoryAnimals, "Есть у птицы, у самолёта и у мельницы"},
		{"кубик", LevelEasy, CategorySport, "Его бросают в настольных играх, у него шесть граней"},
		{"кулак", LevelMedium, CategorySport, "Сжатая ладонь; боксёр бьёт им в перчатке"},
		{"лампа", LevelEasy, CategoryHome, "Светит на столе по вечерам"},
		{"лапша", LevelMedium, CategoryFood, "Длинные полоски теста; иногда её вешают на уши"},
		{"ларек", LevelHard, CategoryCity, "Маленький магазинчик-будка на улице"},
		{"лента", LevelMedium, CategoryHome, "Узкая полоска ткани, её вплетают в косу"},
		{"лимон", LevelEasy, CategoryFood, "Кислый жёлтый цитрус, его кладут в чай"},
		{"лодка", LevelEasy, CategoryCity, "Небольшое судно с вёслами"},
		{"магия", LevelHard, CategoryStudy, "Волшебство и колдовство"},
		{"масло", LevelMedium, CategoryFood, "Его намазывают на хлеб и льют на сковородку"},
		{"место", LevelMedium, CategoryStudy, "Точка в пространстве, а ещё кресло в зрительном зале"},
		{"метро", LevelEasy, CategoryCity, "Подземный городской транспорт"},
		{"мираж", LevelHard, CategoryNature, "Обманчивое видение в пустыне"},
		{"молот", LevelHard, CategoryHome, "Тяжёлый инструмент кузнеца"},
		{"мороз", LevelMedium, CategoryNature, "Сильный холод, рисует узоры на окнах"},
		{"мысль", LevelMedium, CategoryStudy, "То, что приходит в голову"},
		{"набор", LevelHard, CategoryStudy, "Комплект предметов, например цветных карандашей"},
		{"напор", LevelHard, CategoryHome, "Сила, с которой вода бьёт из крана"},
		{"нитка", LevelMedium, CategoryHome, "Её вдевают в иголку"},
		{"ножка", LevelHard, CategoryHome, "Есть у стола, у стула и у гриба"},
		{"обувь", LevelMedium, CategoryHome, "Общее название для ботинок, туфель и сапог"},
		{"огонь", LevelMedium, CategoryNature, "Горит в костре и греет"},
		{"озеро", LevelEasy, CategoryNature, "Водоём, окружённый сушей, например Байкал"},
		{"около", LevelHard, CategoryStudy, "Наречие со значением «рядом, поблизости»"},
		{"отбор", LevelHard, CategorySport, "Этап соревнований, после которого остаются лучшие"},
		{"парус", LevelHard, CategoryCity, "Полотно на мачте, которое ловит ветер"},
		{"ручка", LevelEasy, CategoryStudy, "Ею пишут в тетради"},
		{"слово", LevelEasy, CategoryStudy, "Состоит из букв; в этой игре его нужно угадать"},
		{"тучка", LevelMedium, CategoryNature, "Маленькое облако, из которого может пойти дождь"},
		{"улица", LevelEasy, CategoryCity, "Дорога в городе между рядами домов"},
		{"финал", LevelMedium, CategorySport, "Решающая игра турнира"},
		{"халат", LevelMedium, CategoryHome, "Домашняя одежда с поясом, а ещё рабочая одежда врача"},
		{"цапля", LevelHard, CategoryAnimals, "Длинноногая птица, подолгу стоит на одной ноге в болоте"},
		{"чайка", LevelMedium, CategoryAnimals, "Морская птица, громко кричит над волнами"},
		{"шкала", LevelHard, CategoryStudy, "Ряд делений на градуснике или линейке"},
		{"школа", LevelEasy, CategoryStudy, "Туда дети ходят учиться с первого по одиннадцатый класс"},
		{"щенок", LevelEasy, CategoryAnimals, "Детёныш собаки"},
		{"экран", LevelMedium, CategoryCity, "На него смотрят в кинотеатре"},
		{"ягода", LevelEasy, CategoryFood, "Малина, клубника или черника"},
	}
}

// AnswerPool - слова для загадывания с нужным уровнем и категорией; пустая строка - любой уровень или категория
func AnswerPool(level, category string) ([]string, error) {
	if level != "" && !slices.Contains(Levels(), level) {
		return nil, fmt.Errorf("%q: %w", level, apperr.ErrUnknownLevel)
	}

	if category != "" && !slices.Contains(Categories(), category) {
		return nil, fmt.Errorf("%q: %w", category, apperr.ErrUnknownCategory)
	}

	var pool []string

	for _, info := range Words() {
		if (level == "" || info.Level == level) && (category == "" || info.Category == category) {
			pool = append(pool, info.Word)
		}
	}

	if len(pool) == 0 {
		return nil, fmt.Errorf("уровень %q, категория %q: %w", level, category, apperr.ErrEmptyDict)
	}

	return pool, nil
}

// HintFor - подсказка к слову встроенного словаря; для неизвестного слова - пустая строка
func HintFor(word string) string {
	for _, info := range Words() {
		if info.Word == word {
			return info.Hint
		}
	}

	return ""
}
