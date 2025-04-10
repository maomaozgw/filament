package main

import "github.com/maomaozgw/filament/cmd"

// @title           Filament Management API
// @version         1.0
// @description     This is a web server for 3d printer filament management.

// @contact.name   MaomaoZGW
// @contact.url    https://github.com/maomaozgw
// @contact.email  maomaozgw@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	cmd.Execute()
}
