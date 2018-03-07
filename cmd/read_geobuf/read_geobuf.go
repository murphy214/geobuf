package main

import (
  "os"
  g "github.com/murphy214/geobuf"
  "github.com/urfave/cli"
  "github.com/murphy214/geobuf/read_stdout"
)

func main() {
  app := cli.NewApp()

  app.Action = func(c *cli.Context) error {
    infilename := c.Args().Get(0)
    geobuf := g.ReaderFile(infilename)
    geobuf_stdout.ReadGeobuf(geobuf)

    //g.Convert_Geobuf(infilename,outfilename)

    //fmt.Printf("Hello %q", c.Args().Get(0))
    return nil
  }

  app.Run(os.Args)
}