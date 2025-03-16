package repositories

const (
	// Таблицы
	QuestionTable = "survey_questions_temp"
	OptionTable   = "survey_options_temp"

	// Поля внешнего ключа
	QuestionFKField = "survey_id"   // для вопросов — внешний ключ на опрос (из surveys_temp или surveys, по необходимости)
	OptionFKField   = "question_id" // для опций — внешний ключ на вопрос (из survey_questions_temp)

	// Поля порядка
	QuestionOrderField = "question_order"
	OptionOrderField   = "option_order"

	// Поля состояния
	QuestionStateField = "question_state"
	OptionStateField   = "option_state"
)
