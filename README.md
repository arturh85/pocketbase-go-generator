# pocketbase-go-generator

Application / Library to generate go types from pocketbase collections

## Installation & Usage

### Standalone executable

pocketbase-go-generator can be run as a standalone application accessing an existing pocketbase-server with credentials.

1. Download the latest version from [Releases](https://github.com/arturh85/pocketbase-go-generator/releases).
2. Extract the archive and run the pocketbase-go-generator executable

By default the generator will prompt you for your pocketbase credentials, the collections to export and the output path.

```bash
$ pocketbase-go-generator

---


  Hostname
  >

  Email address
  >

  Password
  >

enter next

---

  Select collections to generate types from
  > ✓ users (10 fields)
    ✓ everything (17 fields)
    • _mfas (System, 6 fields)
    • _otps (System, 7 fields)
    • _externalAuths (System, 7 fields)
    • _authOrigins (System, 6 fields)
    • _superusers (System, 8 fields)

x toggle • ↑ up • ↓ down • / filter • enter submit • ctrl+a select all
```

After submitting the credentials, you can save them in a credentials.env file. You have the choice to save them plain or encrypted with a custom passphrase. So when you run the pocketbase-go-generator again, you can skip the credentials and just enter the encryption passphrase if you chose an encrypted credentials file.

If you don't want to use the built-in prompts, you can use flags to enter the required information:

```
-a, --collections-all               Select all collections include system collections
-x, --collections-exclude strings   Collections to exclude
-i, --collections-include strings   Collections to include (Overrides default selection or all collections)
-d, --disable-form                  Disable form
-l, --disable-logs                  Disable logs, only return result if no output is specified or errors
-e, --email string                  Pocketbase email
-c, --encryption-password string    credentials.enc.env password
-h, --help                          help for generate-go
-u, --host-url string               Pocketbase host url (e. g. http://127.0.0.1:8090)
    --non-required-optional         Make non required fields optional properties (with question mark)
-o, --output string                 Output file path
-p, --password string               Pocketbase password
```

To export all collections that are not marked as system collections (e.g., _superusers), you can type the following command

```bash
$ pocketbase-go-generator -d -u 127.0.0.1:8090 -e [SUPERUSER_EMAIL] -p [SUPERUSER_PASSWORD] -o [OUTPUT_FILE_PATH]
```

Executing this command will cause the generator to connect to the specified PocketBase server, retrieve all collections and save the go definitions to the specified file.

Alternatively, you can print the definitions directly to the console with the `-l` flag and without the `-o`.

```bash
$ pocketbase-go-generator -d -u 127.0.0.1:8090 -e [SUPERUSER_EMAIL] -p [SUPERUSER_PASSWORD] -l
```

### Implement in Go

You can use the pocketbase-go-generator implemented in your pocketbase project either as a command or as a hook. With a hook you can automatically generate a new go file whenever a collection is updated, created or deleted.

You can add the library to your `go.mod` using this command:

```bash
$ go get -u github.com/arturh85/pocketbase-go-generator
```

Examples for both cases are available in the `./example` directory.

#### Implement as command

If you register the command you can use it just like the standalone executable without entering the credentials.

```go
package main

import (
	"github.com/arturh85/pocketbase-go-generator/pkg/pocketbase-go-generator"
	"github.com/pocketbase/pocketbase"
	"github.com/rs/zerolog/log"
)

func main() {
	app := pocketbase.New()

	pocketbase_go_generator.RegisterCommand(app)

	if err := app.Start(); err != nil {
		log.Fatal().Err(err)
	}
}
```

You can run the generate command by typing:

```bash
$ go run ./path/to/main.go generate-go
```

The following options are available:

```
  -a, --collections-all               Select all collections include system collections
  -x, --collections-exclude strings   Collections to exclude
  -i, --collections-include strings   Collections to include (Overrides default selection or all collections)
  -h, --help                          help for generate-go
      --non-required-optional         Make non required fields optional properties (with question mark)
  -o, --output string                 Output file path
```

#### Implement as a hook

If you want to automatically generate new go definitions whenever a collection is updated, created, or deleted, you can use the following example:

```go
package main

import (
	"github.com/arturh85/pocketbase-go-generator/pkg/pocketbase-go-generator"
	"github.com/pocketbase/pocketbase"
	"github.com/rs/zerolog/log"
)

func main() {
	app := pocketbase.New()

	pocketbase_go_generator.RegisterHook(app, &pocketbase_go_generator.GeneratorOptions{
		Output: "test.go",
	})

	if err := app.Start(); err != nil {
		log.Fatal().Err(err)
	}
}
```

When running the pocketbase-server with `go run ./path/to/main.go serve` and performing a collection change, the go definitions are saved in `test.go`.


### Inspiration and Thanks

This project was forked from the excellent [pocketbase-ts-generator](https://github.com/Vogeslu/pocketbase-ts-generator) and changed to output go instead of typescript.