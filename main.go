package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echolog "github.com/labstack/gommon/log"
	"github.com/syndtr/goleveldb/leveldb"
)

var port = flag.Int("port", 8080, "the port")
var dbPath = flag.String("db_path", "memo_testdata/db", "the path to the db")

var wordKey = []byte("word.word")
var noteKey = []byte("note.note")

type Note struct {
	ID       string `json:"id" xml:"id" form:"id" query:"id"`
	Content  string `json:"content" xml:"content" form:"content" query:"content"`
	Modified int64  `json:"modified" xml:"modified" form:"modified" query:"modified"`
}

type NoteList struct {
	Notes []Note `json:"notes" xml:"notes" form:"notes" query:"notes"`
}

type Word struct {
	Name string `json:"name" xml:"name" form:"name" query:"name"`
}

type WordList struct {
	Words []Word `json:"words" xml:"words" form:"words" query:"words"`
}

func (l *WordList) AddWordIfNotExist(w Word) {
	existed := false
	for _, i := range l.Words {
		if i.Name == w.Name {
			existed = true
			break
		}
	}
	if !existed {
		l.Words = append(l.Words, w)
	}
}

func (l *WordList) RemoveWordIfNotExist(w Word) {
	var result []Word
	for _, i := range l.Words {
		if i.Name != w.Name {
			result = append(result, i)
		}
	}
	l.Words = result
}

func main() {
	println("hello world")
	db, err := leveldb.OpenFile(*dbPath, nil)
	if err != nil {
		log.Fatalf("failed open db:%v", err)
	}
	e := echo.New()
	e.Logger.SetLevel(echolog.DEBUG)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/", "static")
	e.GET("/", func(c echo.Context) error {
		return c.File("static/index.html")
	})
	e.POST("/api/word", func(c echo.Context) error {
		w := new(Word)
		if err := c.Bind(w); err != nil {
			return fmt.Errorf("read Word from request payload failed %v", err)
		}
		rawList, err := db.Get(wordKey, nil)
		if err != nil && !errors.Is(err, leveldb.ErrNotFound) {
			return fmt.Errorf("read WordList from db failed %v", err)
		}
		wordList := new(WordList)
		if err == nil {
			if err := json.Unmarshal(rawList, wordList); err != nil {
				return fmt.Errorf("unmarshal WordList (%s) failed %v", rawList, err)
			}
		}
		wordList.AddWordIfNotExist(*w)
		updateRawWord, err := json.Marshal(wordList)
		if err != nil {
			return fmt.Errorf("marshal WordList (%v) failed %v", wordList, err)
		}
		if err := db.Put(wordKey, updateRawWord, nil); err != nil {
			return fmt.Errorf("write wordList to db failed %v", err)
		}
		return c.JSON(http.StatusCreated, wordList)
	})
	e.GET("/api/word", func(c echo.Context) error {
		rawList, err := db.Get(wordKey, nil)
		if err != nil && !errors.Is(err, leveldb.ErrNotFound) {
			return fmt.Errorf("read WordList from db failed %v", err)
		}
		wordList := new(WordList)
		if err == nil {
			if err := json.Unmarshal(rawList, wordList); err != nil {
				return fmt.Errorf("unmarshal WordList (%s) failed %v", rawList, err)
			}
		}
		return c.JSON(http.StatusCreated, wordList)
	})
	e.DELETE("/api/word", func(c echo.Context) error {
		w := new(Word)
		if err := c.Bind(w); err != nil {
			return fmt.Errorf("read Word from request payload failed %v", err)
		}
		rawList, err := db.Get(wordKey, nil)
		if err != nil && !errors.Is(err, leveldb.ErrNotFound) {
			return fmt.Errorf("read WordList from db failed %v", err)
		}
		wordList := new(WordList)
		if err == nil {
			if err := json.Unmarshal(rawList, wordList); err != nil {
				return fmt.Errorf("unmarshal WordList (%s) failed %v", rawList, err)
			}
		}
		wordList.RemoveWordIfNotExist(*w)
		updateRawWord, err := json.Marshal(wordList)
		if err != nil {
			return fmt.Errorf("marshal WordList (%v) failed %v", wordList, err)
		}
		if err := db.Put(wordKey, updateRawWord, nil); err != nil {
			return fmt.Errorf("write wordList to db failed %v", err)
		}
		return c.JSON(http.StatusCreated, wordList)
	})

	e.GET("/api/note", func(c echo.Context) error {
		rawList, err := db.Get(noteKey, nil)
		if err != nil && !errors.Is(err, leveldb.ErrNotFound) {
			return fmt.Errorf("read NoteList from db failed %v", err)
		}
		noteList := new(NoteList)
		if err == nil {
			if err := json.Unmarshal(rawList, noteList); err != nil {
				return fmt.Errorf("unmarshal noteList (%s) failed %v", rawList, err)
			}
		}
		return c.JSON(http.StatusCreated, noteList)
	})
	e.PUT("/api/note", func(c echo.Context) error {
		bodyBytes, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			return fmt.Errorf("read body error %v", err)
		}
		bodyString := string(bodyBytes)
		notes := strings.Split(bodyString, "\n\n")
		noteList := new(NoteList)
		for _, note := range notes {
			note = strings.Trim(note, " \n\t")
			if len(note) == 0 {
				continue
			}
			n := new(Note)
			n.Content = note
			noteList.Notes = append(noteList.Notes, *n)
		}
		updateRawNote, err := json.Marshal(noteList)
		if err != nil {
			return fmt.Errorf("marshal noteList (%v) failed %v", noteList, err)
		}
		if err := db.Put(noteKey, updateRawNote, nil); err != nil {
			return fmt.Errorf("write wordList to db failed %v", err)
		}
		return c.JSON(http.StatusCreated, noteList)
	})
	// e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", *port)))
	e.Logger.Fatal(e.StartTLS(fmt.Sprintf(":%d", *port), "certs/RSA-cert.pem", "certs/RSA-privkey.pem"))
}
