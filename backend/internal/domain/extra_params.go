package domain

import "encoding/json"

type ExtraParams interface {
	isExtraParams()
}

type SingleChoiceExtraParams struct {
	Required bool `json:"required"`
}

func (s SingleChoiceExtraParams) isExtraParams() {}

type MultiChoiceExtraParams struct {
	Required        bool `json:"required"`
	MinAnswersCount int  `json:"minAnswersCount,omitempty"`
	MaxAnswersCount int  `json:"maxAnswersCount,omitempty"`
}

func (m MultiChoiceExtraParams) isExtraParams() {}

type ConsentExtraParams struct{}

func (c ConsentExtraParams) isExtraParams() {}

type EmailExtraParams struct {
	Required bool `json:"required"`
}

func (e EmailExtraParams) isExtraParams() {}

type RatingExtraParams struct {
	Required   bool `json:"required"`
	StarsCount int  `json:"starsCount,omitempty"`
}

func (r RatingExtraParams) isExtraParams() {}

type DateExtraParams struct {
	Required bool   `json:"required"`
	MinDate  string `json:"minDate,omitempty"`
	MaxDate  string `json:"maxDate,omitempty"`
}

func (d DateExtraParams) isExtraParams() {}

type TextExtraParams struct {
	Required  bool `json:"required"`
	MaxLength int  `json:"maxLength,omitempty"`
}

func (t TextExtraParams) isExtraParams() {}

type NumberExtraParams struct {
	Required  bool    `json:"required"`
	MinNumber float64 `json:"minNumber,omitempty"`
	MaxNumber float64 `json:"maxNumber,omitempty"`
}

func (n NumberExtraParams) isExtraParams() {}

// ParseExtraParams парсит json.RawMessage в структуру ExtraParams в зависимости от типа вопроса
func ParseExtraParams(raw json.RawMessage, questionType QuestionType) (ExtraParams, error) {
	var base map[string]interface{}
	if err := json.Unmarshal(raw, &base); err != nil {
		return nil, err
	}

	switch questionType {
	case SingleChoice:
		var params SingleChoiceExtraParams
		if err := json.Unmarshal(raw, &params); err != nil {
			return nil, err
		}
		return params, nil
	case MultiChoice:
		var params MultiChoiceExtraParams
		if err := json.Unmarshal(raw, &params); err != nil {
			return nil, err
		}
		return params, nil
	case Consent:
		return ConsentExtraParams{}, nil
	case Email:
		var params EmailExtraParams
		if err := json.Unmarshal(raw, &params); err != nil {
			return nil, err
		}
		return params, nil
	case Rating:
		var params RatingExtraParams
		if err := json.Unmarshal(raw, &params); err != nil {
			return nil, err
		}
		return params, nil
	case Date:
		var params DateExtraParams
		if err := json.Unmarshal(raw, &params); err != nil {
			return nil, err
		}
		return params, nil
	case ShortText, LongText:
		var params TextExtraParams
		if err := json.Unmarshal(raw, &params); err != nil {
			return nil, err
		}
		return params, nil
	case Number:
		var params NumberExtraParams
		if err := json.Unmarshal(raw, &params); err != nil {
			return nil, err
		}
		return params, nil
	default:
		return nil, ErrInvalidQuestionType
	}
}
