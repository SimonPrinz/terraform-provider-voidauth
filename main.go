package main

import (
	"context"
	"flag"
	"log"

	"github.com/SimonPrinz/terraform-provider-voidauth/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

var (
	version string = "dev"
)

func main() {
	//c, _ := client.NewClient("http://localhost:3000", &http.Client{}, client.WithCredentials("auth_admin", "Voidauth!23"))
	//res, err := c.AdminService.DeleteProxyAuth("5b89aeba-64c4-4db9-b160-c390125b520e")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Println(res)

	var debug bool

	flag.BoolVar(&debug, "debug", false, "debug support")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/simonprinz/voidauth",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.New(version), opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}
