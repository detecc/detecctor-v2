package reply

import . "github.com/detecc/detecctor-v2/internal/model/reply"

type (
	Builder struct {
		buildActions []replyHandler
	}

	replyHandler func(r *Reply)
)

//NewReplyBuilder - constructor
func NewReplyBuilder() *Builder {
	return &Builder{
		buildActions: []replyHandler{},
	}
}

//TypeMessage sets the type to TypeMessage
func (b *Builder) TypeMessage() *Builder {
	b.buildActions = append(b.buildActions, func(r *Reply) {
		r.ReplyType = TypeMessage
	})
	return b
}

//TranslatableMessage sets the type to TranslatableMessage
func (b *Builder) TranslatableMessage() *Builder {
	b.buildActions = append(b.buildActions, func(r *Reply) {
		r.ReplyType = TranslatableMessage
	})
	return b
}

//TypePhoto sets the type to TypePhoto
func (b *Builder) TypePhoto() *Builder {
	b.buildActions = append(b.buildActions, func(r *Reply) {
		r.ReplyType = TypePhoto
	})
	return b
}

//TypeAudio sets the type to TypeAudio
func (b *Builder) TypeAudio() *Builder {
	b.buildActions = append(b.buildActions, func(r *Reply) {
		r.ReplyType = TypeAudio
	})
	return b
}

//WithContent sets the content of the Reply message
func (b *Builder) WithContent(content interface{}) *Builder {
	b.buildActions = append(b.buildActions, func(r *Reply) {
		r.Content = content
	})
	return b
}

//ForChat sets the chatId
func (b *Builder) ForChat(chatId string) *Builder {
	b.buildActions = append(b.buildActions, func(r *Reply) {
		r.ChatId = chatId
	})
	return b
}

//Build builds the Reply object
func (b *Builder) Build() Reply {
	emp := Reply{
		Content:   nil,
		ReplyType: -1,
		ChatId:    "0",
	}

	for _, a := range b.buildActions {
		a(&emp)
	}

	defer func() {
		// Empty builder
		b.buildActions = []replyHandler{}
	}()

	return emp
}
