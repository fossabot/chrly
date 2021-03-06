package cmd

import (
	"fmt"
	"log"

	"github.com/elyby/chrly/auth"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Creates a new token, which allows to interact with Chrly API",
	Run: func(cmd *cobra.Command, args []string) {
		jwtAuth := &auth.JwtAuth{Key: []byte(viper.GetString("chrly.secret"))}
		token, err := jwtAuth.NewToken(auth.SkinScope)
		if err != nil {
			log.Fatalf("Unable to create new token. The error is %v\n", err)
		}

		fmt.Printf("%s\n", token)
	},
}

func init() {
	RootCmd.AddCommand(tokenCmd)
}
