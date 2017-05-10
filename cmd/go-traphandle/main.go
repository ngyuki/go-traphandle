package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/ngyuki/go-traphandle/action"
	"github.com/ngyuki/go-traphandle/config"
	"github.com/ngyuki/go-traphandle/format"
	"github.com/ngyuki/go-traphandle/match"
	"github.com/ngyuki/go-traphandle/parser"
	"github.com/ngyuki/go-traphandle/types"
)

type options struct {
	server *string
	config *string
}

type Route struct {
	Match     *match.Match
	Templates map[string]*format.Template
	Actions   []action.Acter
	Config    *config.Config
	*config.MatchConfig
}

func main() {

	log.SetFlags(0)

	opt := options{}
	opt.server = flag.String("server", "", "Server bind address (default pipe mode)")
	opt.config = flag.String("config", "config.yml", "Configuration filename")
	flag.Parse()

	routes := processConfigure(*opt.config)

	if len(*opt.server) > 0 {
		startServer(*opt.server, func(input []byte) {
			processTrapHandler(input, routes)
		})
	} else {
		input, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}
		processTrapHandler(input, routes)
	}
}

func processConfigure(filename string) []*Route {

	absname, err := filepath.Abs(filename)
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := config.Load(absname)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("load config %s\n%v", filename, dump(cfg))

	routes := make([]*Route, 0, len(cfg.Matches))

	for _, mcfg := range cfg.Matches {

		mat, err := match.NewMatch(&mcfg)
		if err != nil {
			log.Fatal(err)
		}

		tmpls := make(map[string]*format.Template)
		for name, fcfg := range mcfg.Formats {
			template, err := format.NewTemplate(name, fcfg)
			if err != nil {
				log.Fatal(err)
			}
			tmpls[name] = template
		}

		acts, err := action.NewActions(&mcfg.Actions)
		if err != nil {
			log.Fatal(err)
		}

		routes = append(routes, &Route{
			Match:       mat,
			Actions:     acts,
			Templates:   tmpls,
			Config:      cfg,
			MatchConfig: &mcfg,
		})
	}

	return routes
}

func processTrapHandler(input []byte, routes []*Route) {

	trap := parser.Parse(input)
	log.Printf("trap receive\n%s", dump(trap))

	for _, route := range routes {

		values, ok := route.Match.Match(trap)
		if ok == false {
			continue
		}

		values = prepareValues(route, trap, values)

		for _, act := range route.Actions {
			log.Printf("%T %v", act, act)
			err := act.Act(values)
			if err != nil {
				log.Println(err)
			}
		}

		if route.Fallthrough == false {
			break
		}
	}
}

func prepareValues(route *Route, trap *types.Trap, values map[string]string) map[string]string {

	for name, val := range route.Config.Defaults {
		values[name] = val
	}

	values["date"] = time.Now().String()
	values["ipaddr"] = trap.Ipaddr

	for name, tmpl := range route.Templates {
		if out, err := tmpl.Apply(values); err != nil {
			log.Printf("format %v ... %v", name, err)
		} else {
			values[name] = out
		}
	}

	return values
}
