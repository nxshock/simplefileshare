package main

import "fmt"

func sizeToApproxHuman(s int64) string {
	t := []struct {
		Name string
		Val  int64
	}{
		{"EiB", 1 << 60},
		{"PiB", 1 << 50},
		{"TiB", 1 << 40},
		{"GiB", 1 << 30},
		{"MiB", 1 << 20},
		{"KiB", 1 << 10}}

	var v float64
	for i := 0; i < len(t); i++ {
		v = float64(s) / float64(t[i].Val)
		if v < 1.0 {
			continue
		}

		return fmt.Sprintf("%.1f %s", v, t[i].Name)
	}

	return fmt.Sprintf("%.1f KiB", v)
}

func nvl(a1, a2 string) string {
	if a1 != "" {
		return a1
	}

	return a2
}
