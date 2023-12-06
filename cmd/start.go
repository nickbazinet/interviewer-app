/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the Interviewer Application",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: start,
}

func start(cmd *cobra.Command, args []string) {
	fmt.Println(welcomeMessage())

	categories, _ := getCategory("chatgpt")
	fmt.Println(categories)
	categoryPromptContent := promptContent{
		fmt.Sprintf("What category should you be evaluated today?"),
		"Please select a category.",
		getCategoriesName(categories),
	}
	selectedCategory := promptGetSelect(categoryPromptContent)
	fmt.Printf("Let's get started with the question related to %s.\n", selectedCategory)

	questions, err := getQuestions(categories, selectedCategory)
	if err != nil {
		fmt.Println(err)
	}
	for _, question := range questions { 
		wordPromptContent := simplePromptContent{
			question.Text,
			"Please answer a question.",
		}
		 promptGetInput(wordPromptContent)
	}
	
}

func getCategoriesName(categories []Category) []string {
	var categoriesName []string
	for _, category := range categories {
		categoriesName = append(categoriesName, category.Name)
	}
	return categoriesName
}

func getQuestions(categories []Category, categoryName string) ([]Question, error) {
	for _, category := range categories {
		if category.Name == categoryName {
			return category.Questions, nil
		}
	}
	return nil, errors.New(fmt.Sprintf(" Error: No Category match the given category name."))

}

type simplePromptContent struct {
	label    string
	errorMsg string
}

type promptContent struct {
	label	 string
	errorMsg string
	options  []string
}

func promptGetInput(pc simplePromptContent) string {
    validate := func(input string) error {
        if len(input) <= 0 {
            return errors.New(pc.errorMsg)
        }
        return nil
    }

    templates := &promptui.PromptTemplates{
        Prompt:  "{{ . }} ",
        Valid:   "{{ . | green }} ",
        Invalid: "{{ . | red }} ",
        Success: "{{ . | bold }} ",
    }

    prompt := promptui.Prompt{
        Label:     pc.label,
        Templates: templates,
        Validate:  validate,
    }

    result, err := prompt.Run()
    if err != nil {
        fmt.Printf("Prompt failed %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Input: %s\n", result)

    return result
}


func promptGetSelect(pc promptContent) string {
    items := pc.options 
    index := -1
    var result string
    var err error

    for index < 0 {
        prompt := promptui.Select{
            Label:    pc.label,
            Items:    items,
        }

        index, result, err = prompt.Run()

        if index == -1 {
            items = append(items, result)
        }
    }

    if err != nil {
        fmt.Printf("Prompt failed %v\n", err)
        os.Exit(1)
    }

    return result
}


func init() {
	rootCmd.AddCommand(startCmd)
}
