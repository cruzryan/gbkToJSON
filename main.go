package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type Reference struct {
	Authors string
	Title   string
	Journal string
}

type Plasmid struct {
	Locus      string
	Definition string
	Accession  string
	Version    string
	Keywords   string
	Source     string
	Organism   string
	References []Reference
	Features   []map[string]interface{}
	DNA        string
}

type Location struct {
	Start int
	End   int
}

//Features

type Source struct {
	Location Location
	Organism string
	MolType  string
}

type Gene struct {
	Location Location
	Label    string
	Note     string

	/*
		Kinds: "gene, terminator", "misc"
	*/
	Kind        string
	FeatureType string
}

type Promoter struct {
	Location    Location
	Label       string
	Note        string
	Gene        string
	FeatureType string
}

type PromoterBind struct {
	Complement  Location
	Label       string
	Note        string
	FeatureType string
}

type RepOrigin struct {
	Complement Location
	//Left or Right
	Direction   string
	Label       string
	Note        string
	FeatureType string
}

type CDS struct {
	Location    Location
	Codon_start int
	Gene        string
	Product     string
	Label       string
	Note        string
	Translation string
	FeatureType string
}

var (
	feature_types = [...]string{
		"source", "gene", "promoter", "CDS", "terminator",
		"misc_feature", "rep_origin", "primer_bind",
	}
)

func getFileData(path string) string {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	return string(bytes)
}

func nextLine(file *[]rune, pos int) int {
	for i := pos; i < len(*file); i++ {
		if (*file)[i] == '\n' {
			return i
		}
	}

	return -1
}

func currentWord(file *[]rune, text string, pos int) bool {
	return string((*file)[pos:pos+len(text)]) == text
}

/*
	If the initial word of the line IS the word we want:
		then make the value of the plasmid the rest of the line
*/
func asignRestOfTheLine(file *[]rune, pos *int, word string) string {
	isWord := currentWord(file, word, *pos)
	if isWord {
		eol := nextLine(file, *pos)
		//TO-DO: Maybe converting to "string" is kinda dumb?
		result := string((*file)[*pos+len(word) : eol])
		*pos = eol
		return result
	}
	return ""
}

func isFeature(file *[]rune, pos int) (bool, string) {
	for k := 0; k < len(feature_types); k++ {
		if currentWord(file, feature_types[k], pos) {
			return true, feature_types[k]
		}
	}
	return false, ""
}

//Shoutout to Lex Fridman
func lex() *Plasmid {

	plasmid := new(Plasmid)

	string_data := getFileData("tests_reSources\\addgeneplasmid.gbk")

	file := []rune(string_data)

	for i := 0; i < len(file); i++ {

		// fmt.Print(" ", i, " ", string(file[i]))

		if plasmid.Locus == "" {
			plasmid.Locus = asignRestOfTheLine(&file, &i, "LOCUS")
		}

		if plasmid.Definition == "" {
			plasmid.Definition = asignRestOfTheLine(&file, &i, "DEFINITION")
		}

		if plasmid.Accession == "" {
			plasmid.Accession = asignRestOfTheLine(&file, &i, "ACCESSION")
		}

		if plasmid.Version == "" {
			plasmid.Version = asignRestOfTheLine(&file, &i, "VERSION")
		}

		if plasmid.Keywords == "" {
			plasmid.Keywords = asignRestOfTheLine(&file, &i, "KEYWORDS")
		}

		if plasmid.Source == "" {
			plasmid.Source = asignRestOfTheLine(&file, &i, "SOURCE")
		}

		if plasmid.Organism == "" {
			plasmid.Organism = asignRestOfTheLine(&file, &i, "ORGANISM")
		}

		if currentWord(&file, "REFERENCE", i) {
			eol := nextLine(&file, i)
			newRef := new(Reference)

			//Search for end of reference
			eor := 0
			for k := eol; k < len(file); k++ {
				if currentWord(&file, "REFERENCE", k) || currentWord(&file, "FEATURES", k) {
					eor = k
					break
				}
			}

			//References
			//TO-DO: Clean this up
			for l := eol; l < eor; l++ {
				if currentWord(&file, "AUTHORS", l) {
					weol := nextLine(&file, l)
					newRef.Authors = string((file)[l+len("AUTHORS") : weol])
					l = weol
				}

				if currentWord(&file, "TITLE", l) {
					weol := nextLine(&file, l)
					newRef.Title = string((file)[l+len("TITLE") : weol])
					l = weol
				}

				if currentWord(&file, "JOURNAL", l) {
					weol := nextLine(&file, l)
					newRef.Journal = string((file)[l+len("JOURNAL") : weol])
					l = weol
				}
			}
			plasmid.References = append(plasmid.References, *newRef)
		}

		if currentWord(&file, "FEATURES", i) {

			all_features := []map[string]interface{}{}

			//TO-THINK: Are this many loops necessary? is there a better way to do it?
			for m := i; m < len(file); m++ {
				//If you find the DNA section, stop searching for features
				if currentWord(&file, "ORIGIN", m) {
					break
				}

				feature := map[string]interface{}{}
				var feature_name bytes.Buffer
				var feature_content bytes.Buffer
				//is looping through each feature big brain? or small brain?
				// for k := 0; k < len(feature_types); k++ {

				// 	if currentWord(&file, feature_types[k], m) {
				// 		feature["FeatureType"] = feature_types[k]

				// 		//Before we append add its properties dude lol

				// 		all_features = append(all_features, feature)
				// 		feature = map[string]interface{}{}
				// 	}
				// }

				feature_found, name := isFeature(&file, m)

				if feature_found {
					feature["FeatureType"] = name
					all_features = append(all_features, feature)
					feature = map[string]interface{}{}
				} else {
					//TO-DO: fix this if so it skips /Qualitifier
					if file[m] == '/' && file[m+1] != 'Q' {
						fmt.Println(string(file[m : m+10]))
						name_done := false

						for n := m; n < len(file); n++ {
							if file[n] == '=' {
								name_done = true
								continue
							}

							if !name_done {
								//TO-DO: make it not change to string
								feature_name.WriteString(string(file[n]))
							} else {
								feature_content.WriteString(string(file[n]))
							}
						}
						fmt.Println(feature_name.String(), feature_content.String())

						feature["a"] = "b"
						// feature["feature_name.String()"] = feature_content.String()
					}
				}

			}

			plasmid.Features = all_features

		}

	}

	return plasmid
}

func prettyPrint(data interface{}) {
	var p []byte
	//    var err := error
	p, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s \n", p)
}

func main() {
	plas := lex()
	prettyPrint(&plas)
}
