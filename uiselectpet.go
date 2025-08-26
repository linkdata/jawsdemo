package main

import (
	"github.com/linkdata/jaws"
)

func newUiSelectPet() *jaws.NamedBoolArray {
	nba := jaws.NewNamedBoolArray()
	nba.Add("", "--Please choose an option--")
	nba.Add("dog", "Dog")
	nba.Add("cat", "Cat")
	nba.Add("hamster", "Hamster")
	nba.Add("parrot", "Parrot")
	nba.Add("spider", "Spider")
	return nba
}
