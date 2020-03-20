package main

// あ

import (
	"database/sql"
	"fmt"
	"github.com/mattn/go-gimei"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"math/rand"
	"strconv"
	"time"
)

func I2S(i int) string {
	m := i % 60
	h := (i - m) / 60
	return fmt.Sprintf("%d:%02d", h, m)
}

type Record struct {
	Date   time.Time
	Name   string
	TypeA  string
	ValueA string
	ValueB string
}

func NewRecord(Date time.Time, Name, TypeA string, A, B int) *Record {
	p := new(Record)
	p.Date = Date
	p.Name = Name
	p.TypeA = TypeA
	p.ValueA = I2S(A)
	p.ValueB = I2S(B)
	return p
}

// メンバー
type Member struct {
	Team       string
	Name       string
	MemberType string
	TypeAs     []string
}

func NewMember(Team, Name string, TypeAs []string) *Member {
	p := new(Member)
	p.Team = Team
	p.Name = Name
	p.MemberType = "normal"
	p.TypeAs = TypeAs
	return p
}

func (p *Member) GenerateRecords(Date time.Time) []*Record {
	recs := []*Record{}
	for _, TypeA := range p.TypeAs {
		a := rand.Intn(60)*10 - 300
		if a < 0 {
			a = 0
		}
		b := rand.Intn(60)*10 - 300
		if b < 0 {
			b = 0
		}
		if !(a == 0 && b == 0) {
			recs = append(recs, NewRecord(Date, p.Name, TypeA, a, b))
		}
	}
	return recs
}

// チーム
type Team struct {
	Name    string
	Members []*Member
}

//
func NewTeam(Members int, Codes []string) *Team {
	p := new(Team)
	p.Members = []*Member{}
	for i := 0; i < Members; i++ {
		name := gimei.NewName()
		if len(p.Members) == 0 {
			p.Name = name.Last.Kanji() + `G`
		}
		m := NewMember(p.Name, name.Last.Kanji()+`　`+name.First.Kanji(), Codes)
		p.Members = append(p.Members, m)
	}
	return p
}

//
func main() {

	rand.Seed(time.Now().Unix())

	teams := []*Team{}
	tmp := NewTeam(9, []string{"AAA", "AAB", "AAK", "123"})
	tmp.Members[0].MemberType = "manager"
	teams = append(teams, tmp)
	tmp = NewTeam(10, []string{"AAC", "AAD", "AAK", "AAL", "123", "232"})
	tmp.Members[0].MemberType = "manager"
	teams = append(teams, tmp)
	tmp = NewTeam(18, []string{"AAE", "AAF", "AAK", "123", "111", "145"})
	tmp.Members[0].MemberType = "manager"
	teams = append(teams, tmp)
	teams = append(teams, NewTeam(17, []string{"AAG", "AAH", "AAK", "234"}))
	teams = append(teams, NewTeam(12, []string{"AAI", "AAJ", "AAK", "AAL", "432"}))

	finyear := 2019

	format := "2006-01-02"
	first, _ := time.Parse(format, strconv.Itoa(finyear)+"-04-01")
	last, _ := time.Parse(format, strconv.Itoa(finyear+1)+"-04-01")

	recs := []*Record{}
	for dt := first; dt.Before(last); dt = dt.AddDate(0, 0, 1) {
		for _, team := range teams {
			for _, member := range team.Members {
				recs = append(recs, member.GenerateRecords(dt)...)
			}
		}
	}

	// DBに挿入
	Conn, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatal(err)
	}
	// その1
	Tx, _ := Conn.Begin()
	s := `DROP TABLE IF EXISTS t1;`
	log.Println(s)
	_, err = Tx.Exec(s)
	if err != nil {
		log.Fatal(err)
	}
	s = `CREATE TABLE t1(Date TEXT, Name TEXT, TypeA TEXT, TypeB TEXT, ValueA TEXT, ValueB TEXT);`
	log.Println(s)
	_, err = Tx.Exec(s)
	if err != nil {
		log.Fatal(err)
	}
	for _, rec := range recs {
		log.Println(rec)
		_, err = Tx.Exec(`INSERT INTO t1 VALUES(?, ?, ?, ?, ?, ?);`,
			rec.Date.Format(format), rec.Name, rec.TypeA, "XXXX", rec.ValueA, rec.ValueB)
		if err != nil {
			log.Fatal(err)
		}
	}
	Tx.Commit()
	// その2
	Tx, _ = Conn.Begin()
	s = `DROP TABLE IF EXISTS t2;`
	log.Println(s)
	_, err = Conn.Exec(s)
	if err != nil {
		log.Fatal(err)
	}
	s = `CREATE TABLE t2(Team TEXT, Name TEXT, Type TEXT);`
	log.Println(s)
	_, err = Conn.Exec(s)
	if err != nil {
		log.Fatal(err)
	}

	for _, team := range teams {
		for _, member := range team.Members {
			log.Println(member)
			_, err = Tx.Exec(`INSERT INTO t2 VALUES(?, ?, ?);`,
				member.Team, member.Name, member.MemberType)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	Tx.Commit()
	Conn.Close()
}
