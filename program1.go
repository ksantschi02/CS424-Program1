// Author:	Kevin Santschi
// Date:	10 March, 2024
// Class:	CS424-01, Spring 2024
// What this program does:
// This program allows the user to input and alter a .ppm format image
// via a series of filters selected through command line arguments.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
)

// Main func
func main() {
	// Specify possible command line arguments and parse through them upon running program
	inputFileName := flag.String("file", "", "Specify the name of the PPM image file to process")
	flipHFlag := flag.Bool("h", false, "Flip image horizontally")
	flipVFlag := flag.Bool("v", false, "Flip image vertically")
	grayscaleFlag := flag.Bool("g", false, "Convert file to grayscale")
	invertFlag := flag.Bool("i", false, "Invert image colors")
	flattenFlag := flag.String("f", "", "Flatten the specified colors, can enter more than one")
	extremeFlag := flag.Bool("x", false, "Apply extreme contrast filter")
	flag.Parse()

	// Initialize the 3D pixel/RGB matrix for the given .ppm image file
	imageMatrix := readPPM(*inputFileName)

	// Check the command line flags, alter the image matrix accordingly:

	// Flip the image horizontally
	if *flipHFlag {
		imageMatrix = flipHorizontal(imageMatrix)
	} // Flip the image vertically
	if *flipVFlag {
		imageMatrix = flipVertical(imageMatrix)
	} // Convert the image to grayscale
	if *grayscaleFlag {
		imageMatrix = grayscaleMatrix(imageMatrix)
	} // Invert the image's colors
	if *invertFlag {
		imageMatrix = invertColors(imageMatrix)
	} // Flatten the specified colors
	if *flattenFlag != "" {
		imageMatrix = flattenColors(imageMatrix, *flattenFlag)
	} // Apply an extreme contrast filter
	if *extremeFlag {
		imageMatrix = extremeColors(imageMatrix)
	}

	// Write the resulting matrix to a new .ppm image file in an ASCII text format
	writePPM(*inputFileName, imageMatrix)
}

// This function takes the image matrix and flips it horizontally
func flipHorizontal(inputMatrix [][][]int) [][][]int {
	rows, cols := len(inputMatrix), len(inputMatrix[0])

	// Create a result matrix
	flippedMatrix := make([][][]int, rows)
	for i := range flippedMatrix {
		flippedMatrix[i] = make([][]int, cols)
		for j := range flippedMatrix[i] {
			flippedMatrix[i][j] = make([]int, 3)
		}
	}

	// Translate the input matrix from right to left into the result matrix
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			for k := 0; k < 3; k++ {
				flippedMatrix[i][j][k] = inputMatrix[i][cols-1-j][k]
			}
		}
	}

	// Return the flipped matrix
	return flippedMatrix
}

// This function takes the image matrix and flips it vertically
func flipVertical(inputMatrix [][][]int) [][][]int {
	rows, cols := len(inputMatrix), len(inputMatrix[0])

	// Create a result matrix
	flippedMatrix := make([][][]int, rows)
	for i := range flippedMatrix {
		flippedMatrix[i] = make([][]int, cols)
		for j := range flippedMatrix[i] {
			flippedMatrix[i][j] = make([]int, 3)
		}
	}

	// Translate the input matrix from bottom to top into the result matrix
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			for k := 0; k < 3; k++ {
				flippedMatrix[i][j][k] = inputMatrix[rows-1-i][j][k]
			}
		}
	}

	// Return the flipped matrix
	return flippedMatrix
}

// This function takes the image matrix and makes it grayscale
func grayscaleMatrix(inputMatrix [][][]int) [][][]int {
	for i := 0; i < len(inputMatrix); i++ {
		for j := 0; j < len(inputMatrix[0]); j++ {
			// For every pixel, average the RGB values and apply it to all three
			avgRGB := (inputMatrix[i][j][0] + inputMatrix[i][j][1] + inputMatrix[i][j][2]) / 3
			inputMatrix[i][j][0], inputMatrix[i][j][1], inputMatrix[i][j][2] = avgRGB, avgRGB, avgRGB
		}
	}

	// Return the altered matrix
	return inputMatrix
}

// This function takes the image matrix and inverts the color of each pixel
func invertColors(inputMatrix [][][]int) [][][]int {
	for i := 0; i < len(inputMatrix); i++ {
		for j := 0; j < len(inputMatrix[0]); j++ {
			for k := 0; k < 3; k++ {
				// Invert each RGB value by subtracting it from max value 255
				inputMatrix[i][j][k] = 255 - inputMatrix[i][j][k]
			}
		}
	}

	// Return the altered matrix
	return inputMatrix
}

// This function takes the image matrix and flattens the specified colors
func flattenColors(inputMatrix [][][]int, colChoices string) [][][]int {
	// Set bools to track selected colors
	boolR, boolG, boolB := false, false, false

	// Parse through the input string for the letters r, g, b and set bools accordingly
	for _, char := range colChoices {
		switch char {
		case 'R', 'r':
			boolR = true
		case 'G', 'g':
			boolG = true
		case 'B', 'b':
			boolB = true
		}
	}

	for i := 0; i < len(inputMatrix); i++ {
		for j := 0; j < len(inputMatrix[0]); j++ {
			// For each pixel, set RGB values to 0 depending on selected colors to flatten
			if boolR {
				inputMatrix[i][j][0] = 0
			}
			if boolG {
				inputMatrix[i][j][1] = 0
			}
			if boolB {
				inputMatrix[i][j][2] = 0
			}
		}
	}

	// Return the altered matrix
	return inputMatrix
}

// This function takes the image matrix and applies an extreme contrast filter
func extremeColors(inputMatrix [][][]int) [][][]int {
	for i := 0; i < len(inputMatrix); i++ {
		for j := 0; j < len(inputMatrix[0]); j++ {
			for k := 0; k < 3; k++ {
				// For each RGB value, either round to min value or max value
				if inputMatrix[i][j][k] < 128 {
					inputMatrix[i][j][k] = 0
				} else {
					inputMatrix[i][j][k] = 255
				}
			}
		}
	}

	// Return the altered matrix
	return inputMatrix
}

// This function reads a .ppm image file and parses through it, putting RGB values of all pixels
// into a 3D matrix for later alterations
func readPPM(inputFileName string) [][][]int {
	// Open the input file
	file, err := os.Open(inputFileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()

	// Create a scanner to read individual values from the file
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	// Create a string array to temporarily hold values to parse into matrix later
	var words []string
	for scanner.Scan() {
		word := scanner.Text()
		words = append(words, word)
	}

	// Pull the image size from the PPM header
	rows, _ := strconv.Atoi(words[2])
	cols, _ := strconv.Atoi(words[1])

	// Create a 3D matrix based on image size
	inputMatrix := make([][][]int, rows)
	for i := range inputMatrix {
		inputMatrix[i] = make([][]int, cols)
		for j := range inputMatrix[i] {
			inputMatrix[i][j] = make([]int, 3)
		}
	}

	// Create counter to hold the index of the values from the file
	// This will be used to parse the entire file and enter into the matrix
	valueCounter := 4

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			for k := 0; k < 3; k++ {
				// For each RGB value spot in the 3D matrix, accept one value from the file list
				inputMatrix[i][j][k], _ = strconv.Atoi(words[valueCounter])

				// and increment the index counter accordingly
				valueCounter += 1
			}
		}
	}

	// Return the resulting matrix with pixel/RGB values
	return inputMatrix
}

// This function takes the altered image matrix and writes it to a new PPM image file
func writePPM(inputFileName string, imageMatrix [][][]int) error {
	// Set the new file name to <inputFileName>_transformed.ppm
	filename := inputFileName[:len(inputFileName)-4] + "_transformed.ppm"

	// Create a new file by the name set above
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a writer to write to the new file
	writer := bufio.NewWriter(file)

	// Write the PPM header to the file
	header := fmt.Sprintf("P3\n%d %d\n255\n", len(imageMatrix[0]), len(imageMatrix))
	writer.WriteString(header)

	for i := 0; i < len(imageMatrix); i++ {
		for j := 0; j < len(imageMatrix[0]); j++ {
			for k := 0; k < 3; k++ {
				// Write every single RGB value to the file, creating a new line for rows in the image
				writer.WriteString(fmt.Sprintf("%d ", imageMatrix[i][j][k]))
			}
		}
		writer.WriteString("\n")
	}

	// Clear the buffer
	writer.Flush()

	return nil
}
