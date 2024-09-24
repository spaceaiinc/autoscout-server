package responses

import "github.com/spaceaiinc/autoscout-server/domain/entity"

type MessageTemplate struct {
	MessageTemplate *entity.MessageTemplate `json:"message_template"`
}

func NewMessageTemplate(messageTemplate *entity.MessageTemplate) MessageTemplate {
	return MessageTemplate{
		MessageTemplate: messageTemplate,
	}
}

type MessageTemplateList struct {
	MessageTemplateList []*entity.MessageTemplate `json:"message_template_list"`
}

func NewMessageTemplateList(messageTemplates []*entity.MessageTemplate) MessageTemplateList {
	return MessageTemplateList{
		MessageTemplateList: messageTemplates,
	}
}
