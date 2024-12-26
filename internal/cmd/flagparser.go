package cmd

import (
	"github.com/spf13/cobra"
)

type GeneratorFlags struct {
	DisableForm bool
	DisableLogs bool

	Host     string
	Email    string
	Password string

	EncryptionPassword string

	AllCollections     bool
	CollectionsInclude []string
	CollectionsExclude []string

	Output string
}

func GetGenerateTsCommand(fromPocketBase bool, callback func(cmd *cobra.Command, args []string, generatorFlags *GeneratorFlags)) *cobra.Command {
	generatorFlags := &GeneratorFlags{}

	rootCmd := &cobra.Command{
		Use:   "generate-ts",
		Short: "Generate typescript interfaces from pocketbase_api",
		Long:  "Generate typescript interfaces based on pocketbase_api collection definitions",
		Run: func(cmd *cobra.Command, args []string) {
			callback(cmd, args, generatorFlags)
		},
	}

	if !fromPocketBase {
		rootCmd.PersistentFlags().BoolVarP(&generatorFlags.DisableForm, "disable-form", "d", false, "Disable form")
		rootCmd.PersistentFlags().BoolVarP(&generatorFlags.DisableLogs, "disable-logs", "l", false, "Disable logs, only return result if no output is specified or errors")

		rootCmd.PersistentFlags().StringVarP(&generatorFlags.Host, "host url", "u", "", "Pocketbase host url (e. g. http://127.0.0.1:8090)")
		rootCmd.PersistentFlags().StringVarP(&generatorFlags.Host, "email", "e", "", "Pocketbase email")
		rootCmd.PersistentFlags().StringVarP(&generatorFlags.Host, "password", "p", "", "Pocketbase password")

		rootCmd.PersistentFlags().StringVarP(&generatorFlags.EncryptionPassword, "encryption-password", "c", "", "credentials.enc.env password")
	}

	rootCmd.PersistentFlags().BoolVarP(&generatorFlags.DisableForm, "collections-all", "a", false, "Select all collections include system collections")
	rootCmd.PersistentFlags().StringSliceVarP(&generatorFlags.CollectionsInclude, "collections-include", "i", []string{}, "Collections to include (Overrides default selection or all collections)")
	rootCmd.PersistentFlags().StringSliceVarP(&generatorFlags.CollectionsExclude, "collections-exclude", "x", []string{}, "Collections to exclude")

	rootCmd.PersistentFlags().StringVarP(&generatorFlags.Output, "output", "o", "", "Output file path")

	return rootCmd
}
