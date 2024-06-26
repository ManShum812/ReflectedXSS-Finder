# Overview
The Reflected XSS Vulnerability Scanner is a tool designed to help security professionals and developers identify reflected XSS vulnerabilities in web applications. By scanning a list of URLs provided in a .txt file, this tool can efficiently detect potential security flaws that might be exploited by attackers.

# Features
1. Efficiency: The tool scans multiple URLs concurrently, significantly reducing the time required to identify vulnerabilities.

2. Customization: Users can easily modify the payload to test for various types of XSS attacks, making the tool adaptable to different testing scenarios.

3. Reliability: With robust error handling and support for HTTP/2, the scanner ensures accurate results even in complex environments.

4. Stealth: By rotating through a list of User-Agent strings, the tool mimics real browser requests, reducing the likelihood of detection during scans.

This tool is an essential addition to any security professional's toolkit, providing a quick and effective way to improve the security posture of web applications.

# Why Use This Tool?
1. Identify Vulnerabilities: Quickly pinpoint reflected XSS vulnerabilities in your web applications before attackers do.

2. Improve Security: Enhance your application's security by addressing vulnerabilities identified by the scanner.

3. Save Time: Concurrent scanning and efficient HTTP request handling make the scanning process fast and effective.

The Reflected XSS Vulnerability Scanner is designed to be easy to use while providing powerful capabilities to help secure your web applications.

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
1. This script uses a custom payload ('"><12345) to test for XSS vulnerabilities. You can modify this payload as needed.
   
2. The script includes support for HTTP/2 and rotates through a list of User-Agent strings to mimic real browser requests.
   
3. You can also use this script to find other vulnerabilities such as SQL injection (SQLi), Server-Side Request Forgery (SSRF), and Server-Side Template Injection (SSTI) by changing the custom payload accordingly.

For example:
SQL Injection (SQLi): Replace all the parameter values in the URLs in the input.txt file with 1' and set the custom payload to "SQL syntax" to find SQL injection vulnerabilities.
