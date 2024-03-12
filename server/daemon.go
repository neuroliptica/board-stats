package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	TimeLimit = time.Hour
	Timeout   = 10 * time.Second
)

type Daemon struct {
	f     func() (*Record, error)
	table string
	last  time.Time
}

func (d *Daemon) Run() {
	for ; true; time.Sleep(Timeout) {
		if time.Since(d.last) < TimeLimit {
			continue
		}
		r, err := d.f()
		if err != nil {
			logger.Error().Msg(err.Error())
			continue
		}
		err = db.InsertRecord(d.table, r)
		if err != nil {
			logger.Error().Msg(err.Error())
			continue
		}
		d.last = time.Now()
		logger.Info().Fields(map[string]interface{}{
			"board":  d.table,
			"record": *r,
		}).Msg("inserted")
	}
}

func FetchSosach(board string) (*Record, error) {
	res, err := http.Get(fmt.Sprintf(
		"https://2ch.hk/%s/index.json",
		board,
	))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	c, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	type Res struct {
		BoardSpeed uint `json:"board_speed"`
	}
	rj := Res{}
	err = json.Unmarshal(c, &rj)
	if err != nil {
		return nil, err
	}

	return &Record{
		Timestamp: uint64(time.Now().UnixMilli()),
		Posts:     uint64(rj.BoardSpeed),
	}, nil
}
