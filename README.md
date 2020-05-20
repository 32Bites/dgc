# dgc
A DiscordGo command router with tons of useful features

If you find any bugs or if you have a feature request, please tell me by creating an issue.


## Usage
Here's a commented example which should help you using this library:
```go
func main() {
    // Open a simple Discord session
    token := "Your token"
    session, err := discordgo.New("Bot " + token)
    if err != nil {
        panic(err)
    }
    err = session.Open()
    if err != nil {
        panic(err)
    }

    // Create a dgc router
    // NOTE: The dgc.Create function makes sure all internal and external maps of the struct get initialized, so you should use it in every case!
    router := dgc.Create(&dgc.Router{
        // We will allow '!', '$' and the bot mention as command prefixes
        // NOTE: The first prefix (in our case '!') will be used as the prefix in the default help messages
        Prefixes: []string{
            "!",
            "$",
            "<@!" + session.State.User.ID + ">",
        },

        // Whether or not the parser should ignore the case of our prefixes (this would be redundant in our case)
        IgnorePrefixCase: false,

        // Whether or not bots should be allowed to execute our commands
        BotsAllowed: false,

        // We can define commands in here, but in this example we will use the provided method (later)
        Commands: []*dgc.Command{
            // ...
        },

        // We can define middlewares in here, but in this example we will use the provided method (later)
        Middlewares: map[string][]dgc.Middleware{
            // ...
        },

        // The ping handler will be executed if the message only contains the bot's mention (no arguments)
        PingHandler: func(ctx *dgc.Ctx) {
            _, err := ctx.Session.ChannelMessageSend(ctx.Event.ChannelID, "Pong!")
            if err != nil {
                // Error handling
            }
        },
    })

    // Add a simple command
    router.RegisterCmd(&dgc.Command{
        // The general name of the command
        Name: "hey",

        // The aliases of the command
        Aliases: []string{
            "hi",
            "hello",
        },

        // A brief description of the commands functionality
        Description: "Greets you",

        // The correct usage of the command
        Usage: "hey",

        // An example how to use the command
        Example: "hey",

        // Commands may have flags. They will be used for middleware selection and can also be used for grouping
        Flags: []string{
            "greeting",
        },

        // Whether or not the parser should ignore the case of our command
        IgnoreCase: true,

        // A list of sub commands
        SubCommands: []*dgc.Command{
            &dgc.Command{
                Name: "world",
                Description: "Greets the world",
                Usage: "hey world",
                Example: "greet world",
                Flags: []string{
                    "greeting",
                },
                IgnoreCase: true,
                Handler: func(ctx *dgc.Ctx) {
                    _, err := ctx.Session.ChannelMessageSend(ctx.Event.ChannelID, "Hello, world.")
                    if err != nil {
                        // Error handling
                    }
                },
            },
        },

        // dgc supports rate limiting. You can define a rate limiter here.
        RateLimiter: dgc.NewRateLimiter(5*time.Second, 2*time.Second, func (ctx *Ctx) {
            _, err := ctx.Sesion.ChannelMessageSend(ctx.Event.ChannelID, "You are being rate limited!")
            if err != nil {
                // Error handling
            }

            // HINT: You can get the timestamp when the next execution is allowed like this:
            nextExecution := ctx.CustomObjects["dgc_nextExecution"].(time.Time)
        }),

        // The handler of the command
        Handler: func(ctx *dgc.Ctx) {
            _, err := ctx.Session.ChannelMessageSend(ctx.Event.ChannelID, "Hello.")
            if err != nil {
                // Error handling
            }
        },
    })

    // Add a simple middleware that injects a custom object into the context
    // This middleware will be executed on every command, because we wildcard it using the '*'
    // NOTE: You have to return true or false. If you return false, the command will not be executed
    router.AddMiddleware("*", func(ctx *dgc.Ctx) bool {
        // Inject a custom object into the context
        ctx.CustomObjects.Set("myObjectName", "Hello, world")
        return true
    })

    // This middleware will only be executed for commands that implement the 'greeting' flag
    router.AddMiddleware("greeting", func(ctx *dgc.Ctx) bool {
        // Inject a custom object into the context
        ctx.CustomObjects.Set("foo", "bar")
        return true
    })

    // Enable the default help command
    router.RegisterDefaultHelpCommand(session, nil)

    // Initialize the router to make it functional
    router.Initialize(session)
}
```

## Arguments
This library provides an enhanced argument parsing system. The following example will show you how to use it:
```go
func SimpleCommandHandler(ctx *dgc.Ctx) {
    // Retrieve the arguments out of the context
    arguments := ctx.Arguments

    // Print all the arguments with spaces between each other to the console
    fmt.Println(arguments.Raw())

    // Print the amount of given arguments to the console
    fmt.Println(arguments.Amount())


    // Get the first argument
    argument := arguments.Get(0) // will be an argument with an emoty raw string if there is no argument at this index

    // Print it's raw value to the console
    fmt.Println(argument.Raw())

    // Parse the argument into a boolean
    parsedBool, err := argument.AsBool()
    if err != nil {
        // argument is no boolean
    }

    // Parse the argument into an int32
    parsedInt, err := argument.AsInt()
    if err != nil {
        // argument is no int32
    }
    // NOTE: You can use similiar methods (ex. AsInt8, AsInt16...) to parse the argument into these values

    // Parse the argument into a user ID if it is a user mention
    userID := argument.AsUserMentionID()
    if userID == "" {
        // argument is no user mention
    }
    // NOTE: You can also use the methods AsRoleMentionID and AsChannelMentionID
}
```
