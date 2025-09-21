package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"unicode"
)

func main() {
	fmt.Println("Ну хеллоу :) Это моя первая программа, хоть и без гпт никуда.")
	fmt.Println("Как ты мог увидеть, странные числа в папке, попробуй сделать результат равный им")
	fmt.Println("Напиши 'помощь' для списка команд.")

	reader := bufio.NewReader(os.Stdin)
	history := []string{} // история вычислений

	for {
		fmt.Print("calc> ")
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)

		switch line {
		case "":
			fmt.Println("Введите выражение или команду. Для справки напишите 'помощь'.")
			continue
		case "выход":
			fmt.Println("Пока!")
			return
		case "помощь":
			fmt.Println("Доступные команды:")
			fmt.Println("  помощь   - показать эту справку")
			fmt.Println("  выход    - завершить работу калькулятора")
			fmt.Println("  очистить - очистить экран")
			fmt.Println("  история  - показать историю вычислений")
			fmt.Println("Примеры выражений: 2+2, 100-48, -4/2+2")
			continue
		case "очистить":
			clearScreen()
			continue
		case "история":
			if len(history) == 0 {
				fmt.Println("История пуста.")
			} else {
				fmt.Println("История вычислений:")
				for _, h := range history {
					fmt.Println("  " + h)
				}
			}
			continue
		}

		result, err := eval(line)
		if err != nil {
			fmt.Println("\033[31mОшибка:", err, "\033[0m") // красный цвет
			continue
		}

		fmt.Println("\033[32mРезультат:", result, "\033[0m") // зелёный цвет
		history = append(history, fmt.Sprintf("%s = %v", line, result))

		// 🎵 музыкальные бонусы
		if result == 52 {
			playMusic("52.mp3")
		} else if result == 4 {
			playMusic("4.mp3")
		} else if result == 0 {
			playMusic("0.mp3")
		}
	}
}

// ----------------- Основная логика -----------------

func eval(expr string) (float64, error) {
	tokens, err := tokenize(expr)
	if err != nil {
		return 0, err
	}

	rpn, err := shuntingYard(tokens)
	if err != nil {
		return 0, err
	}

	return evalRPN(rpn)
}

func tokenize(expr string) ([]string, error) {
	var tokens []string
	var number strings.Builder
	prev := "" // предыдущий токен

	for i, ch := range expr {
		if !unicode.IsDigit(ch) && ch != '.' && ch != '+' && ch != '-' && ch != '*' && ch != '/' && ch != '(' && ch != ')' {
			return nil, fmt.Errorf("недопустимый символ: %q", ch)
		}

		if unicode.IsDigit(ch) || ch == '.' {
			number.WriteRune(ch)
		} else {
			if number.Len() > 0 {
				tokens = append(tokens, number.String())
				number.Reset()
				prev = tokens[len(tokens)-1]
			}

			if strings.TrimSpace(string(ch)) == "" {
				continue
			}

			if ch == '-' {
				if i == 0 || prev == "" || prev == "(" || prev == "+" || prev == "-" || prev == "*" || prev == "/" {
					number.WriteRune(ch)
					continue
				}
			}

			tokens = append(tokens, string(ch))
			prev = string(ch)
		}
	}
	if number.Len() > 0 {
		tokens = append(tokens, number.String())
	}

	return tokens, nil
}

func shuntingYard(tokens []string) ([]string, error) {
	var output []string
	var stack []string
	precedence := map[string]int{"+": 1, "-": 1, "*": 2, "/": 2}

	for _, tok := range tokens {
		if (tok == "+" || tok == "-" || tok == "*" || tok == "/") &&
			(len(output) == 0 && (len(stack) == 0 || stack[len(stack)-1] == "(")) {
			return nil, fmt.Errorf("выражение не может начинаться с оператора %s", tok)
		}

		if isNumber(tok) {
			output = append(output, tok)
		} else if tok == "+" || tok == "-" || tok == "*" || tok == "/" {
			for len(stack) > 0 {
				top := stack[len(stack)-1]
				if top == "(" {
					break
				}
				if precedence[top] >= precedence[tok] {
					output = append(output, top)
					stack = stack[:len(stack)-1]
				} else {
					break
				}
			}
			stack = append(stack, tok)
		} else if tok == "(" {
			stack = append(stack, tok)
		} else if tok == ")" {
			foundParen := false
			for len(stack) > 0 {
				top := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				if top == "(" {
					foundParen = true
					break
				}
				output = append(output, top)
			}
			if !foundParen {
				return nil, fmt.Errorf("несоответствие скобок")
			}
		} else {
			return nil, fmt.Errorf("неизвестный токен: %s", tok)
		}
	}

	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if top == "(" || top == ")" {
			return nil, fmt.Errorf("несоответствие скобок")
		}
		output = append(output, top)
	}

	return output, nil
}

func evalRPN(tokens []string) (float64, error) {
	var stack []float64

	for _, tok := range tokens {
		if isNumber(tok) {
			num, err := strconv.ParseFloat(tok, 64)
			if err != nil {
				return 0, err
			}
			stack = append(stack, num)
		} else {
			if len(stack) < 2 {
				return 0, fmt.Errorf("недостаточно операндов")
			}
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			switch tok {
			case "+":
				stack = append(stack, a+b)
			case "-":
				stack = append(stack, a-b)
			case "*":
				stack = append(stack, a*b)
			case "/":
				if b == 0 {
					return 0, fmt.Errorf("деление на ноль")
				}
				stack = append(stack, a/b)
			}
		}
	}

	if len(stack) != 1 {
		return 0, fmt.Errorf("ошибка вычислений")
	}

	return stack[0], nil
}

func isNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// ----------------- Дополнительно -----------------

func clearScreen() {
	if runtime.GOOS == "windows" {
		// безопасная очистка для русской консоли
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		fmt.Print("\033[H\033[2J")
	}
}

func fileExists(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}

func playMusic(file string) {
	if !fileExists(file) {
		fmt.Println("Файл не найден:", file)
		return
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", file)
	case "darwin":
		cmd = exec.Command("open", file)
	case "linux":
		cmd = exec.Command("xdg-open", file)
	default:
		fmt.Println("Музыка не поддерживается на этой ОС")
		return
	}

	err := cmd.Start()
	if err != nil {
		fmt.Println("Ошибка при воспроизведении музыки:", err)
	}
}
