package cli

import "github.com/jessevdk/go-flags"

type GeneratorFlags struct {
	DisableForm bool `short:"d" long:"disable-form" description:"Disable form"`
	DisableLogs bool `short:"l" long:"disable-logs" description:"Disable logs, only return result if no output is specified or errors"`

	Host     string `short:"h" long:"host" description:"Pocketbase host"`
	Email    string `short:"e" long:"email" description:"Pocketbase email"`
	Password string `short:"p" long:"password" description:"Pocketbase password"`

	EncryptionPassword string `short:"c" long:"encryption-password" description:"credentials.enc.env password"`

	AllCollections     bool     `short:"a" long:"collections-all" description:"Select all collections include system collections"`
	CollectionsInclude []string `short:"i" long:"collections-include" description:"Collections to include (Overrides default selection or all collections)"`
	CollectionsExclude []string `short:"x" long:"collections-exclude" description:"Collections to exclude"`

	Output string `short:"o" long:"output" description:"Output target, if not specified, generated contents will be printed in console"`
}

func ParseArgs() (*GeneratorFlags, []string, error) {
	generatorFlags := &GeneratorFlags{}

	args, err := flags.Parse(generatorFlags)
	if err != nil {
		return nil, nil, err
	}

	return generatorFlags, args, nil
}
