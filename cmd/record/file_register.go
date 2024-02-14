package record

import (
    "io"
    "os"
)

// AppendStringToFile appends a given string to the file specified by filename.
// It adds the string to a new line if the file already contains data.
func Register_answer(related_question, answer string) error  {
    file, err := os.OpenFile("./responses", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer file.Close()

    stat, err := file.Stat()
    if err != nil {
        return err
    }

    file_entry := related_question + ":\n" + answer
    if stat.Size() > 0 { 
	file_entry = "\n" + file_entry
    }

    _, err = io.WriteString(file, file_entry)
    if err != nil {
        return err
    }

    return nil
}
