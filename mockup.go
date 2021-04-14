package main

import "time"

type Video struct {
	Id       string        `json:"id"`
	Poster   string        `json:"Poster"`
	Duration time.Duration `json:"Duration"`
	Name     string        `json:"Name"`
}

// var (
// 	AllVideos :=
// )

func getVideos() []Video {

	return []Video{
		{
			Id:       "tom-and-jerry",
			Poster:   "https://image.tmdb.org/t/p/w500/fev8UFNFFYsD5q7AcYS8LyTzqwl.jpg",
			Duration: 3 * time.Minute,
			Name:     "Tom & Jerry",
		},
		{
			Id:       "soul",
			Poster:   "https://image.tmdb.org/t/p/w500/kf456ZqeC45XTvo6W9pW5clYKfQ.jpg",
			Duration: 4 * time.Minute,
			Name:     "Soul",
		},
		{
			Id:       "outside-the-wire",
			Poster:   "https://image.tmdb.org/t/p/w500/lOSdUkGQmbAl5JQ3QoHqBZUbZhC.jpg",
			Duration: 2 * time.Minute,
			Name:     "Outside the wire",
		},
	}
}
