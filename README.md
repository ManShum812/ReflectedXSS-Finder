# Overview
The Reflected XSS Vulnerability Scanner is a tool designed to help security professionals and developers identify reflected XSS vulnerabilities in web applications. By scanning a list of URLs provided in a .txt file, this tool can efficiently detect potential security flaws that might be exploited by attackers.

# Features
Concurrency: Utilizes multiple workers to scan URLs concurrently, speeding up the scanning process.

Custom Payloads: Allows you to easily specify a payload to test for XSS vulnerabilities.

User-Agent Rotation: Rotates through a list of User-Agent strings to reduce detection.

Error Handling: Includes robust error handling and retries for HTTP/2 specific errors.

Logging and Output: Logs the scanning process and saves the results in a results.txt file.

# Installation
git clone https://github.com/ManShum812/ReflectedXSS-Finder.git

cd ReflectedXSS-Finder

go build main.go

# Running the Scanner
1. Prepare Input File: Create a file named input.txt with the URLs you want to scan, one per line.
   
2. Run the Script: ./main
   
3. Check the Output: The results will be saved in results.txt. The file will contain details about each URL scanned, including the status code and whether the URL is vulnerable to XSS.
   ![xss](https://github.com/ManShum812/ReflectedXSS-Finder/assets/43279996/508969f8-23be-4762-84ac-2bcee53ede25)
   ![xss2](https://github.com/ManShum812/ReflectedXSS-Finder/assets/43279996/66817218-5d1f-4da9-a772-1d8016153042)


# Notes
This script uses a custom payload ('"><12345) to test for XSS vulnerabilities. You can modify this payload as needed.
The script includes support for HTTP/2 and rotates through a list of User-Agent strings to mimic real browser requests.
