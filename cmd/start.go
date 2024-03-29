package cmd

import (
	"errors"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/manifoldco/promptui"
	"github.com/nickbazinet/interviewer-app/cmd/config"
	"github.com/nickbazinet/interviewer-app/cmd/record"
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

func getCategoryWorker(done chan []config.Category, p *tea.Program) {
	categories, _ := config.GetCategory("chatgpt")
	done <- categories
	p.Quit()
}

func start(cmd *cobra.Command, args []string) {
	fmt.Println(welcomeMessage())

	p := tea.NewProgram(initialModel())
	ch := make(chan []config.Category, 1)
	go getCategoryWorker(ch, p)
	
	p.Run()

	categories := <-ch
	categoryPromptContent := promptContent{
		"What category should you be evaluated today?",
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
		record.Register_answer(question.Text, promptGetInput(wordPromptContent))
	}
	
}

func getCategoriesName(categories []config.Category) []string {
	var categoriesName []string
	for _, category := range categories {
		categoriesName = append(categoriesName, category.Name)
	}
	return categoriesName
}

func getQuestions(categories []config.Category, categoryName string) ([]config.Question, error) {
	for _, category := range categories {
		if category.Name == categoryName {
			return category.Questions, nil
		}
	}
	return nil, fmt.Errorf("error: no category match the given category name: %s", categoryName)

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
