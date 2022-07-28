package generators

import "time"

type Directory struct {
	Id        string    `tr:"primary_key"`
	CreatedAt time.Time `tr:"DisplayName:Created At;DisplayOrder:1"`
	FileName  string    `tr:"DisplayName:File Name;DisplayOrder:2"`
	Owner     string    `tr:"DisplayName:Owner;DisplayOrder:3"`
	IsFolder  bool      `tr:"DisplayName:Is Folder;DisplayOrder:4"`
}
