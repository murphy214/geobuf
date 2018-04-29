package main

import (
  "os"
  g "github.com/murphy214/geobuf"
  "github.com/urfave/cli"
)

func main() {
  app := cli.NewApp()

  app.Action = func(c *cli.Context) error {
    infilename := c.Args().Get(0)
    g.ReadGeobufCSV(infilename)

    //g.Convert_Geobuf(infilename,outfilename)

    //fmt.Printf("Hello %q", c.Args().Get(0))
    return nil
  }

  app.Run(os.Args)
}