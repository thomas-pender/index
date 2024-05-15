package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	no_searches   = 10
	index         = "./index.txt"
	error_message = "ERROR -- error usage\n\t./index [enter | search] [--args]\n"
	index_size    = 5000
)

// helpers /////////////////////////////////////////////////////////////////////

func get_types(args []string) (map[string]int, map[string]int) {
	author, subject := make(map[string]int), make(map[string]int)
	for _, arg := range args {
		s := strings.Split(arg, "=")

		if len(s) <= 1 || len(s) > 2 || s[0][:2] != "--" {
			continue
		}

		t := strings.Split(s[1], " ")
		switch s[0] {
		case "--author":
			for _, t1 := range t {
				author[t1]++
			}
		case "--subject":
			for _, t1 := range t {
				subject[t1]++
			}
		}
	}
	return author, subject
}

func print_entry(line string) {
	entry := strings.Split(line, " | ")
	fmt.Printf("\tINDEX   : %s\n\tAUTHORS : %s\n\tTITLE   : %s\n\tSUBJECT : %s\n\n", entry[0], entry[3], entry[1], entry[2])
}

// enter index entry ///////////////////////////////////////////////////////////

func get_index_no() int {
	f, err := os.Open(index)
	if err != nil {
		log.Fatal(err)
	}

	fi, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}
	if fi.Size() == 0 {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
		return 1
	}

	count := 0
	indices := make([]int, 0, index_size)

	input := bufio.NewScanner(f)
	for input.Scan() {
		line := strings.Split(input.Text(), " ")
		n, err := strconv.Atoi(line[0])
		if err != nil {
			log.Fatal(err)
		}
		indices = append(indices, n)
		count++
	}
	sort.Ints(indices)
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

	for i, j := range indices {
		if i+1 != j {
			return i + 1
		}
	}
	return count + 1
}

func get_title(args []string) string {
	for _, arg := range args {
		t := strings.Split(arg, "=")
		if t[0] == "--title" {
			return t[1]
		}
	}
	return ""
}

func cat_type(t map[string]int) string {
	var s []string
	for name, num := range t {
		for i := 1; i <= num; i++ {
			s = append(s, name)
		}
	}
	sort.Strings(s)
	return strings.Join(s, " ")
}

func enter(args []string) {
	f, err := os.OpenFile(index, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}

	title := get_title(args)
	author, subject := get_types(args)
	if len(author) == 0 || len(subject) == 0 || title == "" {
		log.Fatal(error_message)
	}

	line := strconv.Itoa(get_index_no()) + " | " + cat_type(subject) + " | " + title + " | " + cat_type(author) + "\n"
	if _, err := f.WriteString(line); err != nil {
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Entered article:")
	print_entry(line)
}

// search index entry //////////////////////////////////////////////////////////

func check_line(author map[string]int, subject map[string]int, line string) (bool, bool) {
	var count int
	var flag_author, flag_subject bool

	entry := strings.Split(line, " | ")
	line_authors := strings.Split(entry[3], " ")
	line_subject := strings.Split(entry[1], " ")

	count = 0
	for _, name := range line_authors {
		if author[name] > 0 {
			count++
		}
	}
	if count == len(author) {
		flag_author = true
	}

	count = 0
	for _, sub := range line_subject {
		if subject[sub] > 0 {
			count++
		}
	}
	if count == len(subject) {
		flag_subject = true
	}

	return flag_author, flag_subject
}

func line_search(author, subject map[string]int, f *os.File) bool {
	var line string
	input := bufio.NewScanner(f)
	var flag_author, flag_subject, flag_found bool

	flag_found = true
	for input.Scan() {
		line = input.Text()
		flag_author, flag_subject = check_line(author, subject, line)

		if len(author) > 0 && len(subject) > 0 && flag_author && flag_subject {
			flag_found = false
			print_entry(line)
		} else if len(author) > 0 && len(subject) == 0 && flag_author {
			flag_found = false
			print_entry(line)
		} else if len(author) == 0 && len(subject) > 0 && flag_subject {
			flag_found = false
			print_entry(line)
		}
	}

	return flag_found
}

func search(args []string) {
	f, err := os.Open(index)
	if err != nil {
		log.Fatal(err)
	}

	fi, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}
	if fi.Size() == 0 {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
		log.Fatal("ERROR -- index is empty")
	}

	author, subject := get_types(args)
	if len(author) == 0 && len(subject) == 0 {
		log.Fatal(error_message)
	}

	if line_search(author, subject, f) {
		fmt.Println("ERROR -- no such entry")
	}
}

// remove entry ////////////////////////////////////////////////////////////////

func get_indices(args []string) []string {
	indices := make([]string, 0, index_size)
	for _, arg := range args {
		s := strings.Split(arg, "=")

		if s[0] != "--index" {
			continue
		}

		t := strings.Split(s[1], " ")
		for _, num := range t {
			indices = append(indices, num)
		}
	}
	if len(indices) == 0 {
		log.Fatal(error_message)
	}

	return indices
}

func remove_line(args []string, f *os.File) []string {
	var flag bool
	indices := get_indices(args)
	var line []string
	lines := make([]string, 0, index_size)
	input := bufio.NewScanner(f)
	for input.Scan() {
		flag = false
		line = strings.Split(input.Text(), " | ")
		for _, num := range indices {
			if line[0] == num {
				flag = true
				break
			}
		}
		if flag {
			continue
		}
		lines = append(lines, strings.Join(line, " | ")+"\n")
	}
	return lines
}

func entry_remove(args []string) {
	f, err := os.OpenFile(index, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}

	fi, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}
	if fi.Size() == 0 {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
		log.Fatal("ERROR -- index is empty")
	}

	lines := remove_line(args, f)

	if err := f.Truncate(0); err != nil {
		log.Fatal(err)
	}
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		log.Fatal(err)
	}

	for _, s := range lines {
		if _, err := f.WriteString(s); err != nil {
			log.Fatal(err)
		}
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

// remove duplicates ///////////////////////////////////////////////////////////

func get_duplicates(f *os.File) map[string][]string {
	var line []string
	var meta string
	entries := make(map[string][]string)
	input := bufio.NewScanner(f)
	for input.Scan() {
		line = strings.Split(input.Text(), " | ")
		meta = strings.Join(line[1:], " | ")
		entries[meta] = append(entries[meta], line[0])
	}
	return entries
}

func remove_duplicates() {
	f, err := os.OpenFile(index, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}

	fi, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}
	if fi.Size() == 0 {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
		log.Fatal("ERROR -- index is empty")
	}

	entries := get_duplicates(f)

	if err := f.Truncate(0); err != nil {
		log.Fatal(err)
	}
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		log.Fatal(err)
	}

	for info, num := range entries {
		sort.Strings(num)
		if _, err := f.WriteString(num[0] + " | " + info + "\n"); err != nil {
			log.Fatal(err)
		}
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

// exhaustive searching ////////////////////////////////////////////////////////

func list_search(arg []string, f *os.File) bool {
	flag := true
	var line []string
	input := bufio.NewScanner(f)
	for input.Scan() {
	}
}

func list_subjects(args []string) {
	f, err := os.Open(index)
	if err != nil {
		log.Fatal(err)
	}

	fi, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}
	if fi.Size() == 0 {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
		log.Fatal("ERROR -- index is empty")
	}

	arg := strings.Split(args[1], "=")
	if len(args) <= 1 || (arg[0] != "--author" && arg[0] != "--subject") {
		log.Fatal(error_message)
	}
}

// driver //////////////////////////////////////////////////////////////////////

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		log.Fatal(error_message)
	}

	switch args[0] {
	case "enter":
		if len(args) <= 1 {
			log.Fatal(error_message)
		}
		enter(args[1:])
	case "search":
		if len(args) <= 1 {
			log.Fatal(error_message)
		}
		search(args[1:])
	case "remove-index":
		if len(args) <= 1 {
			log.Fatal(error_message)
		}
		entry_remove(args[1:])
	case "remove-duplicates":
		remove_duplicates()
	}
}
