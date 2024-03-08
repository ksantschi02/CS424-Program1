// Author:	Kevin Santschi
// Date:	6 March, 2024
// Class:	CS424-01, Spring 2024
// What this program does:
//

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
)

func main() {
	inputFileName := flag.String("file", "", "Specify the name of the PPM image file to process")
	flipHFlag := flag.Bool("h", false, "Flip image horizontally")
	flipVFlag := flag.Bool("v", false, "Flip image vertically")
	grayscaleFlag := flag.Bool("g", false, "Convert file to grayscale")
	// invertFlag := flag.Bool("i", false, "Invert image colors")
	// flattenFlag := flag.String("f", "", "Flatten the specified colors, can enter more than one")
	// extremeFlag := flag.Bool("x", false, "Apply extreme contrast filter")
	flag.Parse()

	imageMatrix := readPPM(*inputFileName)

	if *flipHFlag {
		imageMatrix = flipHorizontal(imageMatrix)
	}
	if *flipVFlag {
		imageMatrix = flipVertical(imageMatrix)
	}
	if *grayscaleFlag {
		imageMatrix = grayscaleMatrix(imageMatrix)
	}

	// Print the populated 3D matrix
	fmt.Println("Populated 3D Matrix:")
	for i := range imageMatrix {
		for j := range imageMatrix[i] {
			fmt.Printf("%d ", imageMatrix[i][j])
		}
		fmt.Println()
	}
}

func flipHorizontal(inputMatrix [][][]int) [][][]int {
	rows, cols := len(inputMatrix), len(inputMatrix[0])

	// Create a 3D matrix
	flippedMatrix := make([][][]int, rows)
	for i := range flippedMatrix {
		flippedMatrix[i] = make([][]int, cols)
		for j := range flippedMatrix[i] {
			flippedMatrix[i][j] = make([]int, 3)
		}
	}

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			for k := 0; k < 3; k++ {
				flippedMatrix[i][j][k] = inputMatrix[i][cols-1-j][k]
			}
		}
	}

	return flippedMatrix
}

func flipVertical(inputMatrix [][][]int) [][][]int {
	rows, cols := len(inputMatrix), len(inputMatrix[0])

	// Create a 3D matrix
	flippedMatrix := make([][][]int, rows)
	for i := range flippedMatrix {
		flippedMatrix[i] = make([][]int, cols)
		for j := range flippedMatrix[i] {
			flippedMatrix[i][j] = make([]int, 3)
		}
	}

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			for k := 0; k < 3; k++ {
				flippedMatrix[i][j][k] = inputMatrix[rows-1-i][j][k]
			}
		}
	}

	return flippedMatrix
}

func grayscaleMatrix(inputMatrix [][][]int) [][][]int {
	for i := 0; i < len(inputMatrix); i++ {
		for j := 0; j < len(inputMatrix[0]); j++ {
			avgRGB := (inputMatrix[i][j][0] + inputMatrix[i][j][1] + inputMatrix[i][j][2]) / 3
			inputMatrix[i][j][0], inputMatrix[i][j][1], inputMatrix[i][j][2] = avgRGB, avgRGB, avgRGB
		}
	}

	return inputMatrix
}

func readPPM(inputFileName string) [][][]int {
	// Open your input file or use any other source of data
	file, err := os.Open(inputFileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()

	// Create a bufio scanner for reading individual values from the file
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	var words []string
	for scanner.Scan() {
		word := scanner.Text()
		words = append(words, word)
	}

	// Adjust these values based on your matrix dimensions
	rows, _ := strconv.Atoi(words[1])
	cols, _ := strconv.Atoi(words[2])

	// Create a 3D matrix
	inputMatrix := make([][][]int, rows)
	for i := range inputMatrix {
		inputMatrix[i] = make([][]int, cols)
		for j := range inputMatrix[i] {
			inputMatrix[i][j] = make([]int, 3)
		}
	}

	valueCounter := 4

	// Iterate over each line in the file
	for i := 0; i < rows; i++ {
		// Iterate over each value in the line and parse it into the matrix
		for j := 0; j < cols; j++ {
			// You can add error handling here if the conversion fails
			for k := 0; k < 3; k++ {
				inputMatrix[i][j][k], _ = strconv.Atoi(words[valueCounter])
				valueCounter += 1
			}
		}
	}
	return inputMatrix
}
