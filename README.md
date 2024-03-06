READ ME

Ignore test version.

This program is designed to accept .ppm image files and alter them based on a series of command line flags.

-file ppm_file_to_process	          Specify the name of the PPM image file to process. This argument is required.
-h	                                Flip image horizontally
-v	                                Flip image vertically
-g	                                Convert file to grayscale
-i	                                Invert image colors
-f [r][g][b]	                      Flatten the specified colors. Note the user can enter one or more colors (e.g., -f gb or -f b).
-x	                                Apply extreme contrast filter
