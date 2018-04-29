package main

import (
  "fmt"
  "os"
  g "github.com/murphy214/geobuf"
  "github.com/urfave/cli"
)

func main() {
  app := cli.NewApp()

  app.Action = func(c *cli.Context) error {
    infilename := c.Args().Get(0)
    outfilename := c.Args().Get(1)
    fmt.Println("Converting: ",infilename,"to geojson filename:", outfilename)

    g.ConvertGeobuf(infilename,outfilename)

    //fmt.Printf("Hello %q", c.Args().Get(0))
    return nil
  }

  app.Run(os.Args)
}