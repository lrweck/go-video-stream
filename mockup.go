package main

import "time"

type Video struct {
	Id       string        `json:"id"`
	Poster   string        `json:"Poster"`
	Duration time.Duration `json:"Duration"`
	Name     string        `json:"Name"`
}

func getVideos() []Video {

	return []Video{
		{
			Id:       "tom-and-jerry",
			Poster:   "https://image.tmdb.org/t/p/w500/fev8UFNFFYsD5q7AcYS8LyTzqwl.jpg",
			Duration: 144 * time.Second,
			Name:     "Tom & Jerry",
		},
		{
			Id:       "soul",
			Poster:   "https://image.tmdb.org/t/p/w500/kf456ZqeC45XTvo6W9pW5clYKfQ.jpg",
			Duration: 139 * time.Second,
			Name:     "Soul",
		},
		{
			Id:       "outside-the-wire",
			Poster:   "https://image.tmdb.org/t/p/w500/lOSdUkGQmbAl5JQ3QoHqBZUbZhC.jpg",
			Duration: 152 * time.Second,
			Name:     "Outside the wire",
		},
	}
}
