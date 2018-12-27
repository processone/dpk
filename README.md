# Data Portability Kit

This repository provides a set of tools to analyse / convert personal data from online services.

## Twitter

You can ask Twitter to download your archive here: [Your Twitter Data](https://twitter.com/settings/your_twitter_data).  
You will receive a link to download your archive when ready.

When you got it, unzip the file and convert your tweets to a Markdown directory structure with the command:

```bash
go run cmd/twitter-to-md/twitter-to-md.go ~/Downloads/twitter-2018-12-27-abcd121212 posts
``` 

It will create a directory with your data in a format you can reuse with your blogging tool platform.
