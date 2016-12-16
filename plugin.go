package main

import (
	"fmt"
	"log"

	"github.com/vallard/spark"
)

type (
	Repo struct {
		Owner string
		Name  string
	}

	Build struct {
		Tag     string
		Event   string
		Number  int
		Commit  string
		Ref     string
		Branch  string
		Author  string
		Status  string
		Link    string
		Started int64
		Created int64
	}

	Config struct {
		Room     string
		Token    string
		Template string
		Markdown string
	}

	Job struct {
		Started int64
	}

	Plugin struct {
		Repo   Repo
		Build  Build
		Config Config
		Job    Job
	}
)

func (p Plugin) Exec() error {

	if p.Config.Room == "" {
		log.Fatal("SPARK_ROOM is not defined")
	}

	if p.Config.Token == "" {
		log.Fatal("SPARK_TOKEN is not defined")
	}
	s := spark.New(p.Config.Token)

	// get the room ID from the room name
	room, err := s.GetRoomWithName(p.Config.Room)
	// if no room matched, then return error
	if err != nil {
		return err
	}

	// create a new message
	m := spark.Message{
		RoomId:   room.Id,
		Markdown: message(p.Repo, p.Build),
	}

	// template defined trumps default markdown
	if p.Config.Template != "" {
		txt, err := RenderTrim(p.Config.Template, p)
		if err != nil {
			return err
		}
		m.Markdown = ""
		m.Text = txt
	}

	// Markdown trumps text
	if p.Config.Markdown != "" {
		txt, err := RenderTrim(p.Config.Markdown, p)
		if err != nil {
			return err
		}
		m.Markdown = txt
		m.Text = ""
	}

	_, err = s.CreateMessage(m)
	return err
}

func message(repo Repo, build Build) string {
	var c string
	if build.Commit == "" {
		c = ""
	} else {
		c = build.Commit[:8]
	}
	//return fmt.Sprintf("*%s* <%s|%s/%s#%s> (%s) by %s",
	var m string
	if build.Status == "success" {
		m = fmt.Sprintf("##Build for %s/%s is successful\n__Build Author:__ %s\n",
			repo.Owner,
			repo.Name,
			build.Author)

	} else {
		m = fmt.Sprintf("#Build for %s/%s FAILED!!! \n__Drone blames build author:__ %s\n",
			repo.Owner,
			repo.Name,
			build.Author)
	}

	return fmt.Sprintf("%s###Build Details\n * [Build Log](%s/%s/%s#%s) \n* __Event__ %s\n",
		m,
		build.Link,
		repo.Owner,
		repo.Name,
		c,
		build.Event,
	)

}
