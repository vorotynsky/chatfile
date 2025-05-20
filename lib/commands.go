package chatfile

type CommandName string

// Command represents an entry in a chatfile.
// Each command modifies the [Context] state in a specific way when applied.
// Commands are the building blocks that define how a conversation is structured.
// The system parses the chatfile, converting each entry into a Command instance, then executes them
// sequentially to build up the final Context state.
type Command interface {
	Name() CommandName

	// Apply executes the command's action via modifying [Context].
	Apply(*Context)
}

type FromCommand struct {
	ModelName string
}

func (c *FromCommand) Name() CommandName {
	return "FROM"
}

func (c *FromCommand) Apply(ctx *Context) {
	ctx.CurrentModel = ModelName(c.ModelName)
}

type PromptCommand struct {
	Role    Role
	Message string
}

func (c *PromptCommand) Name() CommandName {
	return "PROMPT"
}

func (c *PromptCommand) Apply(ctx *Context) {
	ctx.History.Append(c.Role, c.Message)
}
