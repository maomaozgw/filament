package serve

import (
	"context"

	"github.com/maomaozgw/filament/pkg/server"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "serve",
	Short: "serve server",
	Run: func(cmd *cobra.Command, args []string) {
		s, err := server.NewServer(server.Option{
			Addr: ":8080",
			Orm: server.OrmOpt{
				Type: "sqlite",
				Sqlite: &server.SqliteOpt{
					Path: "filament.db",
				},
			},
		})
		if err != nil {
			panic(err)
		}
		s.Run(context.Background())
	},
}
