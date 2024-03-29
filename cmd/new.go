/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

type SurveyData struct {
	Title       string
	Description string
}

type Post struct {
	Title       string
	Description string
	Image       string
	Slug        string
	CreatedAt   string
}

var (
	unsplashProviderURL = "https://source.unsplash.com/1600x900"
	postsDir            = "www/gabriel-blog/src/_posts/"
	staticAssetsDir     = "www/gabriel-blog/public/static/"
)

var qs = []*survey.Question{
	{
		Name:     "Title",
		Prompt:   &survey.Input{Message: "Post Title"},
		Validate: survey.Required,
	},
	{
		Name:     "Description",
		Prompt:   &survey.Input{Message: "Post Description"},
		Validate: survey.Required,
	},
}

/*
	Get survey data
	Generate timestamp
	Generate slug
	Generate image

	Assemble the file
*/

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Creates a new post using user input",
	Run: func(cmd *cobra.Command, args []string) {
		var answers SurveyData
		err := survey.Ask(qs, &answers)
		must(err)

		homedir, err := homedir.Dir()
		must(err)

		postsDir = fmt.Sprintf("%s/%s", homedir, postsDir)
		staticAssetsDir = fmt.Sprintf("%s/%s", homedir, staticAssetsDir)

		today := time.Now()
		postDate := fmt.Sprintf("%d-%d-%d", today.Month(), today.Day(), today.Year())
		postSlug := generateSlug(answers.Title)
		postImage := fmt.Sprintf("%s%s.jpg", staticAssetsDir, postSlug)

		must(downloadImage(unsplashProviderURL, postImage))

		file, err := os.Create(fmt.Sprintf("%s%s.md", postsDir, postSlug))
		must(err)

		post := Post{
			Title:       answers.Title,
			Slug:        postSlug,
			Description: answers.Description,
			CreatedAt:   postDate,
			Image:       fmt.Sprintf("/static/%s.jpg", postSlug),
		}
		must(writeToFile(file, post))
	},
}

func generateSlug(title string) string {
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")

	return slug
}

func downloadImage(fileUrl, filePath string) error {
	res, err := http.Get(fileUrl)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	out, err := os.Create(filePath)
	if err != nil {
		return err
	}

	_, err = io.Copy(out, res.Body)
	if err != nil {
		return err
	}

	return nil
}

func writeToFile(file *os.File, data Post) error {
	w := bufio.NewWriter(file)
	content := []string{
		"---",
		fmt.Sprintf("title: %s", data.Title),
		fmt.Sprintf("description: %s", data.Description),
		fmt.Sprintf("createdAt: %s", data.CreatedAt),
		fmt.Sprintf("image: %s", data.Image),
		fmt.Sprintf("slug: %s", data.Slug),
		"---",
		"# This post was autogenerated by my [blog-automator](https://github.com/Gabriel2233/blog-automator) project",
		"<WarningBox />",
	}

	for _, c := range content {
		_, err := w.WriteString(c + "\n")
		if err != nil {
			return err
		}
	}

	err := w.Flush()
	if err != nil {
		return err
	}

	return nil
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func init() {
	rootCmd.AddCommand(newCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
